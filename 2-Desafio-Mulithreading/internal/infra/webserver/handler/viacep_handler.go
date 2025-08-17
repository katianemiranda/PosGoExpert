package webserver

import (
	"net/http"

	_ " github.com/katianemiranda/0.2-Desafio-Multithreading/internal/entity"
	"github.com/go-chi/chi/v5"
)

type ViaCepHandler struct {
	viacep ViaCEPResponse
}

func (v *ViaCepHandler) BuscaCEPbyViaCEPHandler1(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//resp, _ := http.Get("https://viacep.com.br/ws/" + cep + "/json/")

	//	defer resp.Body.Close()

	//body, _ := io.ReadAll(resp.Body)

	//var viacep ViaCEPResponse
	//err := json.Unmarshal(body, &viacep)

	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(address)
}
