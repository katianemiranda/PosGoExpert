package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/katianemiranda00/posgoexpert/labsgo/1-cep/entity"
)

func ConsultaClimaHandler(location string) (*entity.WeatherAPIResponse, error) {
	// Lógica para consultar a API de clima

	baseURL := "http://api.weatherapi.com/v1/current.json"
	apiKey := "6355c1c4be774137adb225023260102" // Substitua pela sua chave de API real

	city := location

	params := url.Values{}
	params.Add("key", apiKey)
	params.Add("q", city)
	params.Add("lang", "pt")

	url := baseURL + "?" + params.Encode()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao consultar a API de clima: %s", resp.Status)
	}

	var weather entity.WeatherAPIResponse

	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil

}
