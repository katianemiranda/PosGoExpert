package handler

import (
	"2-DESAFIO-MULTITHREADING/internal/entity"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func buscaCepbyBrasilAPI(cep string) (*entity.BrasilAPIResponse, error) {
	resp, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var brasilapi entity.BrasilAPIResponse
	err = json.Unmarshal(body, &brasilapi)
	if err != nil {
		return nil, err
	}
	return &brasilapi, nil

}

func BuscaCepByBrasilApiHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request iniciada By BrasilAPI")
	defer log.Println("Request finalizada By BrasilAPI")

	cep := chi.URLParam(r, "cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	viacep, err := BuscaCepbyBrasilAPI(cep)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(viacep)

}
