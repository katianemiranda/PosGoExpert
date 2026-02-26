package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// Constantes
const (
	zipkinURL   = "http://zipkin:9411/api/v2/spans"
	serviceName = "service-a"
	serviceBURL = "http://service-b:8081/weather" // Nome do serviço B no Docker Compose
)

// Estruturas de Dados
type CEPInput struct {
	CEP string `json:"cep"`
}

var cepRegex = regexp.MustCompile(`^\d{8}$`)
var tracer = otel.Tracer(serviceName)

// Função de Inicialização do OpenTelemetry/Zipkin
func initTracer() *sdktrace.TracerProvider {
	exporter, err := zipkin.New(zipkinURL)
	if err != nil {
		log.Fatalf("falha ao criar exportador zipkin: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	// Configura o Propagator (W3C Trace Context) para que o contexto seja injetado nas requisições de saída
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)

	return tp
}

// Handler Principal
func processCEPHandler(w http.ResponseWriter, r *http.Request) {
	// O span já foi iniciado pelo otelhttp.NewHandler e o contexto está em r.Context()
	ctx := r.Context()
	span := trace.SpanFromContext(ctx) // Obtém o span ativo

	// Decodifica o input
	var input CEPInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		// Trata erro de parsing JSON ou corpo vazio
		respondWithError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		span.SetAttributes(attribute.Bool("error", true), attribute.String("status_code", "422"), attribute.String("error.msg", "JSON decode error"))
		return
	}

	cep := input.CEP
	span.SetAttributes(attribute.String("cep", cep))

	// 2. Validação do CEP
	if !cepRegex.MatchString(cep) {
		respondWithError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		span.SetAttributes(attribute.Bool("error", true), attribute.String("status_code", "422"), attribute.String("error.msg", "Format invalid"))
		return
	}

	// 3. Encaminha para o Serviço B (com instrumentação OTEL)
	serviceBURLWithQuery := fmt.Sprintf("%s?cep=%s", serviceBURL, cep)

	// otelhttp.DefaultClient já utiliza o propagador configurado (TraceContext)
	client := otelhttp.DefaultClient
	req, err := http.NewRequestWithContext(ctx, "GET", serviceBURLWithQuery, nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		span.SetAttributes(attribute.Bool("error", true), attribute.String("status_code", "500"), attribute.String("error.msg", "Request creation error"))
		return
	}

	span.AddEvent("Chamando Service B com Tracing Context")

	resp, err := client.Do(req)
	if err != nil {
		// Erro de rede/DNS/conexão
		respondWithError(w, http.StatusInternalServerError, "internal server error: service B unreachable")
		span.SetAttributes(attribute.Bool("error", true), attribute.String("status_code", "500"), attribute.String("error.msg", fmt.Sprintf("Service B connection error: %v", err)))
		return
	}
	defer resp.Body.Close()

	// 4. Repassa a resposta do Serviço B

	// Copia o status code
	w.WriteHeader(resp.StatusCode)
	// Copia o Content-Type (deve ser application/json)
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	// Copia o body (resposta JSON)
	_, _ = io.Copy(w, resp.Body)

	// Adiciona o status da resposta ao span do Service A
	span.SetAttributes(attribute.Int("http.status_code_response", resp.StatusCode))
}

// Funções de Resposta HTTP
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Erro ao desligar o TracerProvider: %v", err)
		}
	}()

	r := chi.NewRouter()

	// Instrumentação para o handler de entrada do Service A (cria o span pai inicial)
	handler := otelhttp.NewHandler(http.HandlerFunc(processCEPHandler), "Service A Incoming Request")

	r.Method(http.MethodPost, "/", handler)

	log.Println("Serviço A rodando na porta 8080...")
	http.ListenAndServe(":8080", r)
}
