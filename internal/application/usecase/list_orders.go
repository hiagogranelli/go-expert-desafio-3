package usecase

import (
	"context"

	"desafio-cleanarchitecture/internal/application/dto"
	"desafio-cleanarchitecture/internal/domain/repository"
)

type ListOrdersUseCase struct {
	OrderRepository repository.OrderRepository
}

func NewListOrdersUseCase(repo repository.OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: repo,
	}
}

func (uc *ListOrdersUseCase) Execute(ctx context.Context) (*dto.ListOrdersOutputDTO, error) {
	orders, err := uc.OrderRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var orderDTOs []dto.OrderOutputDTO
	for _, order := range orders {
		orderDTOs = append(orderDTOs, dto.OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return &dto.ListOrdersOutputDTO{Orders: orderDTOs}, nil
}
