package graph

import (
	"context"
	"fmt"

	"desafio-cleanarchitecture/internal/application/dto"
	"desafio-cleanarchitecture/internal/application/usecase"
)

type Resolver struct {
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateOrder(ctx context.Context, input CreateOrderInput) (*Order, error) {
	ucInput := dto.CreateOrderInputDTO{
		Price: input.Price,
		Tax:   input.Tax,
	}
	output, err := r.CreateOrderUseCase.Execute(ctx, ucInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return &Order{
		ID:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) ListOrders(ctx context.Context) ([]*Order, error) {
	output, err := r.ListOrdersUseCase.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	gqlOrders := make([]*Order, 0, len(output.Orders))
	for _, orderDTO := range output.Orders {
		gqlOrders = append(gqlOrders, &Order{
			ID:         orderDTO.ID,
			Price:      orderDTO.Price,
			Tax:        orderDTO.Tax,
			FinalPrice: orderDTO.FinalPrice,
		})
	}
	return gqlOrders, nil
}
