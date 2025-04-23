package usecase

import (
	"context"

	"desafio-cleanarchitecture/internal/application/dto"
	"desafio-cleanarchitecture/internal/domain/entity"
	"desafio-cleanarchitecture/internal/domain/repository"
)

type CreateOrderUseCase struct {
	OrderRepository repository.OrderRepository
}

func NewCreateOrderUseCase(repo repository.OrderRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: repo,
	}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, input dto.CreateOrderInputDTO) (*dto.OrderOutputDTO, error) {
	order, err := entity.NewOrder(input.Price, input.Tax)
	if err != nil {
		return nil, err
	}

	err = uc.OrderRepository.Save(ctx, order)
	if err != nil {
		return nil, err
	}

	output := &dto.OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}
	return output, nil
}
