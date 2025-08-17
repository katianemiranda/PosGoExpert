package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func BuscaCepbyViaCEP(cep string) (*ViaCEPResponse, error) {
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var viacep ViaCEPResponse
	err = json.Unmarshal(body, &viacep)
	if err != nil {
		return nil, err
	}
	return &viacep, nil

}

func BuscaCEPbyViaCEPHandler(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()

	log.Println("Request iniciada")
	defer log.Println("Request finalizada")

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

func main() {
	r := chi.NewRouter()

	r.Route("/cep", func(r chi.Router) {
		r.Get("/{cep}", BuscaCEPbyViaCEPHandler)
	})

	//http.HandleFunc("/cep", BuscaCEPbyViaCEPHandler)
	http.ListenAndServe(":8000", r)
}
