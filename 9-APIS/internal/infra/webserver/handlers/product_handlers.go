package handlers

import (
	"encoding/json"
	"katianemiranda/PosGoExpert/9-APIS/internal/dto"
	"katianemiranda/PosGoExpert/9-APIS/internal/entity"
	"katianemiranda/PosGoExpert/9-APIS/internal/infra/database"
	"net/http"
)

type Producthandler struct {
	productDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *Producthandler {
	return &Producthandler{
		productDB: db,
	}
}

func (h *Producthandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.productDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
