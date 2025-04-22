package usecase

import (
	"context"

	"desafio-cleanarchitecture/internal/application/dto"
	"desafio-cleanarchitecture/internal/domain/entity"
	"desafio-cleanarchitecture/internal/domain/repository"
)

type CreateOrderUseCase struct {
	OrderRepository repository.OrderRepository
	// Uow             repository.Uow // Se usar Unit of Work
}

func NewCreateOrderUseCase(repo repository.OrderRepository /*, uow repository.Uow*/) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: repo,
		// Uow:             uow,
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

	// Se usar UoW:
	// err = uc.Uow.Do(ctx, func(uow *repository.Uow) error {
	//  repo, err := uow.GetRepository(ctx, "OrderRepository")
	//  if err != nil { return err }
	//  err = repo.(repository.OrderRepository).Save(ctx, order)
	//  return err
	// })
	// if err != nil { return nil, err }

	output := &dto.OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}
	return output, nil
}
