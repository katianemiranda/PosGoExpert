package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katianemiranda00/posgoexpert/labsgo/1-cep/entity"
)

type Resposta struct {
	ViaCep entity.ViaCEPResponse `json:"viacep"`
}

type Temperatura struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

func BuscarCEP(cep string) (*Resposta, error) {
	// Lógica para buscar o CEP
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var viaCep entity.ViaCEPResponse
	err = json.Unmarshal(body, &viaCep)
	if err != nil {
		return nil, err
	}
	return &Resposta{ViaCep: viaCep}, nil
}

func BuscarCEPHandler(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	resposta, err := BuscarCEP(cep)
	if err != nil {
		http.Error(w, "Erro ao buscar o CEP", http.StatusInternalServerError)
		return
	}

	if resposta.ViaCep.Cep == "" {
		http.Error(w, "CEP não encontrado", http.StatusNotFound)
		return
	}

	if resposta.ViaCep.Cep == "erro" {
		http.Error(w, "CEP inválido", http.StatusUnprocessableEntity)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(resposta)

	clima, err := ConsultaClimaHandler(resposta.ViaCep.Localidade)
	if err != nil {
		http.Error(w, "Erro ao consultar o clima", http.StatusInternalServerError)
		return
	}

	//json.NewEncoder(w).Encode(clima)

	var temperatura Temperatura
	temperatura.TempC = clima.Current.TempC
	temperatura.TempF = clima.Current.TempC*1.8 + 32
	temperatura.TempK = clima.Current.TempC + 273.15

	json.NewEncoder(w).Encode(temperatura)

	//fmt.Printf("Clima em %s, %s: %.1f°C, %s\n", clima.Current.TempC, clima.Current.Condition.Text)

}
