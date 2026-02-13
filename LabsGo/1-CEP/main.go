package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katianemiranda00/posgoexpert/labsgo/1-cep/infra/webserver/handler"
)

func main() {

	r := chi.NewRouter()

	r.Route("/cep", func(r chi.Router) {
		r.Get("/{cep}", handler.BuscarCEPHandler)
	})
	http.ListenAndServe(":8080", r)
}
