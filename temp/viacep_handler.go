package handler

import (
	"2-DESAFIO-MULTITHREADING/internal/entity"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

/*type ViaCepHandler struct {
	viacep entity.ViaCEPResponse
}*/

func buscaCepbyViaCEP(cep string) (*entity.ViaCEPResponse, error) {
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var viacep entity.ViaCEPResponse
	err = json.Unmarshal(body, &viacep)
	if err != nil {
		return nil, err
	}
	return &viacep, nil

}

func BuscaCEPbyViaCEPHandler(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()

	log.Println("Request iniciada By ViaCEP")
	defer log.Println("Request finalizada By ViaCEP")

	cep := chi.URLParam(r, "cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	viacep, err := BuscaCepbyViaCEP(cep)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(viacep)

}
