package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
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
	serviceName = "service-b"
)

// --- Estruturas de Dados ---

type WeatherInput struct {
	CEP string `json:"cep"`
}

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherOutput struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

var tracer = otel.Tracer(serviceName)

// --- Funções Auxiliares de OTEL/Tracing ---

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

	// Define o Propagator para extrair o contexto do header
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)
	return tp
}

// --- Funções de Conversão ---

func convertTemperatures(celsius float64) (fahrenheit float64, kelvin float64) {
	// F = C * 1,8 + 32
	fahrenheit = celsius*1.8 + 32
	// K = C + 273
	kelvin = celsius + 273.0
	return fahrenheit, kelvin
}

// --- Handlers e Lógica de Negócio ---

func getWeatherHandler(w http.ResponseWriter, r *http.Request) {
	// O contexto é extraído automaticamente pelo chi/otel quando o Serviço A chama.
	// Usa o contexto da requisição para continuar o trace.
	ctx := r.Context()

	// Inicia um Span com o contexto propagado
	ctx, span := tracer.Start(ctx, "getWeatherHandler", trace.WithAttributes(attribute.String("component", "service-b-entry")))
	defer span.End()

	cep := r.URL.Query().Get("cep")

	// 1. Validação (redundante, mas para segurança)
	if len(cep) != 8 {
		respondWithError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		span.SetAttributes(attribute.Bool("error", true), attribute.String("status_code", "422"))
		return
	}

	// 2. Busca a cidade (ViaCEP)
	city, err := getCityFromCEP(ctx, cep)
	if err != nil {
		if err.Error() == "not found" {
			respondWithError(w, http.StatusNotFound, "can not find zipcode")
			span.SetAttributes(attribute.Bool("error", true), attribute.String("status_code", "404"))
			return
		}
		log.Printf("Erro ao buscar cidade: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		span.SetAttributes(attribute.Bool("error", true), attribute.String("status_code", "500"))
		return
	}

	// 3. Busca a temperatura (WeatherAPI)
	tempC, err := getTemperature(ctx, city)
	if err != nil {
		log.Printf("Erro ao buscar temperatura: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		span.SetAttributes(attribute.Bool("error", true), attribute.String("status_code", "500"))
		return
	}

	// 4. Converte e Formata
	tempF, tempK := convertTemperatures(tempC)

	output := WeatherOutput{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	span.SetAttributes(
		attribute.String("city", city),
		attribute.Float64("temp_c", tempC),
		attribute.String("status_code", "200"),
	)
	respondWithJSON(w, http.StatusOK, output)
}

// Função para buscar a cidade com Span
func getCityFromCEP(ctx context.Context, cep string) (string, error) {
	ctx, span := tracer.Start(ctx, "getCityFromCEP", trace.WithAttributes(attribute.String("http.endpoint", "viacep")))
	defer span.End()

	viaCEPURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	req, _ := http.NewRequestWithContext(ctx, "GET", viaCEPURL, nil)
	client := http.Client{Timeout: 5 * time.Second}

	span.SetAttributes(attribute.String("http.url", viaCEPURL))

	resp, err := client.Do(req)
	if err != nil {
		span.SetAttributes(attribute.Bool("http.error", true))
		return "", fmt.Errorf("erro na requisição ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		// ViaCEP retorna 400 para CEPs válidos no formato, mas inexistentes
		return "", fmt.Errorf("not found")
	}

	if resp.StatusCode != http.StatusOK {
		span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
		return "", fmt.Errorf("erro inesperado no ViaCEP, status: %d", resp.StatusCode)
	}

	var data ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("erro ao decodificar ViaCEP: %w", err)
	}

	if data.Erro || data.Localidade == "" {
		return "", fmt.Errorf("not found")
	}

	span.SetAttributes(attribute.String("city_result", data.Localidade))
	return data.Localidade, nil
}

// Função para buscar a temperatura com Span
func getTemperature(ctx context.Context, city string) (float64, error) {
	ctx, span := tracer.Start(ctx, "getTemperature", trace.WithAttributes(attribute.String("http.endpoint", "weatherapi")))
	defer span.End()

	weatherAPIKey := os.Getenv("WEATHERAPI_KEY")

	// MUDANÇA: Remove espaços e aspas duplas que podem ter sido lidas.
	weatherAPIKey = strings.TrimSpace(strings.ReplaceAll(weatherAPIKey, "\"", ""))

	if weatherAPIKey == "" {
		log.Println("ERRO CRÍTICO: Variável de ambiente WEATHERAPI_KEY está vazia após a limpeza.")
		span.SetAttributes(attribute.Bool("error", true), attribute.String("error.type", "API_KEY_MISSING"))
		return 0, fmt.Errorf("weatherapi key not found or is invalid")
	}

	// Codifica o nome da localidade para uso na URL
	encodedLocation := url.QueryEscape(city)

	weatherAPIURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", weatherAPIKey, encodedLocation)

	req, _ := http.NewRequestWithContext(ctx, "GET", weatherAPIURL, nil)
	//client := http.Client{Timeout: 5 * time.Second}
	client := http.Client{}

	span.SetAttributes(attribute.String("http.url_base", "http://api.weatherapi.com/v1/current.json"))

	// Adiciona o cabeçalho User-Agent para evitar bloqueio da API
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		span.SetAttributes(attribute.Bool("http.error", true))
		return 0, fmt.Errorf("erro na requisição WeatherAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode), attribute.Bool("http.error", true))
		// Loga o erro, mas retorna uma mensagem genérica no 500
		return 0, fmt.Errorf("erro na WeatherAPI: status %d, requisicao %s", resp.StatusCode, weatherAPIURL)
	}

	var data WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("erro ao decodificar WeatherAPI: %w", err)
	}

	span.SetAttributes(attribute.Float64("temp_c_raw", data.Current.TempC))
	return data.Current.TempC, nil
}

// --- Funções de Resposta HTTP ---

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

	// O Serviço B recebe a requisição do Serviço A via GET com query parameter
	r.Get("/weather", getWeatherHandler)

	log.Println("Serviço B rodando na porta 8081...")
	http.ListenAndServe(":8081", r)
}
