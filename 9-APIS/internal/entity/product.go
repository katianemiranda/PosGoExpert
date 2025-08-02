package entity

import (
	"errors"
	"katianemiranda/PosGoExpert/9-APIS/pkg/entity"
	"time"
)

var (
	ErrIDisRequired    = errors.New("id is required")
	ErrInvalidID       = errors.New("invalid ID")
	ErrNameIsRequired  = errors.New("Name is required")
	ErrPriceIsRequired = errors.New("Price is required")
	ErrInvalidPrice    = errors.New("invalid price")
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func NewProduct(name string, price float64) (*Product, error) {
	product := &Product{
		ID:        entity.NewID(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
	}

	err := product.Validate()
	if err != nil {
		return nil, err
	}

	return product, nil

}

func (p *Product) Validate() error {
	if p.ID.String() == "" {
		return ErrIDisRequired
	}
	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return ErrInvalidID
	}
	if p.Name == "" {
		return ErrNameIsRequired
	}
	if p.Price < 0 {
		return ErrInvalidPrice
	}
	if p.Price == 0 {
		return ErrPriceIsRequired
	}
	return nil
}
