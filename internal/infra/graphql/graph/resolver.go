package graph

// Este arquivo NÃO SERÁ sobrescrito pelo gqlgen após a primeira geração.

import (
	"context"
	"fmt" // Exemplo

	// --- CORRIGIR IMPORTS ---
	// Importe seus DTOs corretamente
	"desafio-cleanarchitecture/internal/application/dto"
	// Importe seus UseCases corretamente
	"desafio-cleanarchitecture/internal/application/usecase"
	// Não importe mais ".../graph/model"
)

// Resolver struct injeta as dependências (use cases)
type Resolver struct {
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

// --- AS FUNÇÕES ABAIXO DEVEM CORRESPONDER AO QUE gqlgen ESPERA ---

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

// CreateOrder is the resolver for the createOrder field.
// Use os tipos gerados pelo gqlgen (que estarão no pacote 'graph')
func (r *mutationResolver) CreateOrder(ctx context.Context, input CreateOrderInput) (*Order, error) {
	// Converter input do GraphQL para DTO do use case
	ucInput := dto.CreateOrderInputDTO{
		Price: input.Price,
		Tax:   input.Tax,
	}
	// Chamar o use case
	output, err := r.CreateOrderUseCase.Execute(ctx, ucInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	// Converter output do use case para o model do GraphQL (do pacote 'graph')
	return &Order{ // Note: Não precisa de 'graph.Order' aqui pois já estamos no pacote 'graph'
		ID:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice, // Nome do campo como definido no schema.graphqls
	}, nil
}

type queryResolver struct{ *Resolver }

// ListOrders is the resolver for the listOrders field.
// Use os tipos gerados pelo gqlgen (que estarão no pacote 'graph')
func (r *queryResolver) ListOrders(ctx context.Context) ([]*Order, error) {
	// Chamar o use case de listagem
	output, err := r.ListOrdersUseCase.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	// Converter a lista de DTOs para a lista de models do GraphQL (do pacote 'graph')
	gqlOrders := make([]*Order, 0, len(output.Orders)) // Melhor prática de inicialização de slice
	for _, orderDTO := range output.Orders {
		gqlOrders = append(gqlOrders, &Order{ // Note: Não precisa de 'graph.Order'
			ID:         orderDTO.ID,
			Price:      orderDTO.Price,
			Tax:        orderDTO.Tax,
			FinalPrice: orderDTO.FinalPrice,
		})
	}
	return gqlOrders, nil
}
