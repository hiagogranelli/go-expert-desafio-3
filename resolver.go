package graph

import (
	"context"
	"desafio-cleanarchitecture/internal/infra/graphql/graph"
)

type Resolver struct{}

func (r *mutationResolver) CreateOrder(ctx context.Context, input graph.CreateOrderInput) (*graph.Order, error) {
	panic("not implemented")
}

func (r *queryResolver) ListOrders(ctx context.Context) ([]*graph.Order, error) {
	panic("not implemented")
}

func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
