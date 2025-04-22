package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         string    `json:"id"`
	Price      float64   `json:"price"`
	Tax        float64   `json:"tax"`
	FinalPrice float64   `json:"final_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewOrder(price, tax float64) (*Order, error) {
	order := &Order{
		ID:    uuid.New().String(),
		Price: price,
		Tax:   tax,
		// CreatedAt e UpdatedAt serão definidos pelo DB ou na criação
	}
	err := order.Validate()
	if err != nil {
		return nil, err
	}
	order.CalculateFinalPrice()
	return order, nil
}

func (o *Order) Validate() error {
	if o.ID == "" {
		return errors.New("id is required")
	}
	if o.Price <= 0 {
		return errors.New("price must be positive")
	}
	if o.Tax < 0 { // Taxa pode ser zero
		return errors.New("tax cannot be negative")
	}
	return nil
}

func (o *Order) CalculateFinalPrice() {
	o.FinalPrice = o.Price + o.Tax
}
