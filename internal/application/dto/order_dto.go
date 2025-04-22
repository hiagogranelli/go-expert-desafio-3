package dto

// Input DTO para criação
type CreateOrderInputDTO struct {
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

// Output DTO genérico (pode ser usado por create e list)
type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

// Output DTO específico para listagem (se necessário diferenciar)
type ListOrdersOutputDTO struct {
	Orders []OrderOutputDTO `json:"orders"`
}
