package handler

import (
	"encoding/json"
	"net/http"

	"desafio-cleanarchitecture/internal/application/dto"
	"desafio-cleanarchitecture/internal/application/usecase"
)

type WebOrderHandler struct {
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

func NewWebOrderHandler(createUC *usecase.CreateOrderUseCase, listUC *usecase.ListOrdersUseCase) *WebOrderHandler {
	return &WebOrderHandler{
		CreateOrderUseCase: createUC,
		ListOrdersUseCase:  listUC,
	}
}

func (h *WebOrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateOrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.CreateOrderUseCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (h *WebOrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	output, err := h.ListOrdersUseCase.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
