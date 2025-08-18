package main

import (
	"2-DESAFIO-MULTITHREADING/internal/infra/webserver/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	r := chi.NewRouter()

	r.Route("/cep", func(r chi.Router) {
		r.Get("/{cep}", handler.ConsultaCepHandler)
	})
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))
	http.ListenAndServe(":8000", r)

}
