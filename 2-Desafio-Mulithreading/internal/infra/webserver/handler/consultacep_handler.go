package handler

import (
	"2-DESAFIO-MULTITHREADING/internal/entity"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Resposta struct {
	Msg       string
	ViaCep    entity.ViaCEPResponse
	BrasilAPI entity.BrasilAPIResponse
}

func BuscaCepbyBrasilAPI(cep string) (*entity.BrasilAPIResponse, error) {
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

func BuscaCepbyViaCEP(cep string) (*entity.ViaCEPResponse, error) {
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

// ConsultaCEP godoc
// @Sumary Consulta um cep
// @Description Consulta cep
// @Tags         cep
// @Accept       json
// @Produce      json
// @Param        cep  path     string  true
// @Success      200  {array}  handler.Resposta
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /cep [get]

func ConsultaCepHandler(w http.ResponseWriter, r *http.Request) {
	c1 := make(chan Resposta)
	c2 := make(chan Resposta)

	cep := chi.URLParam(r, "cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	go func() {
		viacep, err := BuscaCepbyViaCEP(cep)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		msg := Resposta{Msg: "ViaCEP", ViaCep: *viacep}
		c1 <- msg
	}()
	go func() {
		viacep, err := BuscaCepbyBrasilAPI(cep)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		msg := Resposta{Msg: "BrasilAPI", BrasilAPI: *viacep}
		c2 <- msg
	}()

	//for {
	select {
	case msg1 := <-c1:
		fmt.Printf("Received from:  Msg=%s\n", msg1.Msg)
		fmt.Println(msg1.ViaCep)

	case msg2 := <-c2:
		fmt.Printf("Received from: Msg=%s\n", msg2.Msg)
		fmt.Println(msg2.BrasilAPI)
	case <-time.After(time.Second * 1):
		fmt.Println("Timeout: No messages received within 1 seconds")

	}
	//}

}
