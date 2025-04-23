package service

import (
	"context"

	"desafio-cleanarchitecture/internal/application/dto"
	"desafio-cleanarchitecture/internal/application/usecase"
	"desafio-cleanarchitecture/internal/infra/grpc/pb"
)

type OrderGrpcService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

func NewOrderGrpcService(createUC *usecase.CreateOrderUseCase, listUC *usecase.ListOrdersUseCase) *OrderGrpcService {
	return &OrderGrpcService{
		CreateOrderUseCase: createUC,
		ListOrdersUseCase:  listUC,
	}
}

func (s *OrderGrpcService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	input := dto.CreateOrderInputDTO{
		Price: float64(req.Price),
		Tax:   float64(req.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Order: &pb.Order{
			Id:         output.ID,
			Price:      float32(output.Price),
			Tax:        float32(output.Tax),
			FinalPrice: float32(output.FinalPrice),
		},
	}, nil
}

func (s *OrderGrpcService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	output, err := s.ListOrdersUseCase.Execute(ctx)
	if err != nil {
		return nil, err
	}

	var grpcOrders []*pb.Order
	for _, orderDTO := range output.Orders {
		grpcOrders = append(grpcOrders, &pb.Order{
			Id:         orderDTO.ID,
			Price:      float32(orderDTO.Price),
			Tax:        float32(orderDTO.Tax),
			FinalPrice: float32(orderDTO.FinalPrice),
		})
	}

	return &pb.ListOrdersResponse{Orders: grpcOrders}, nil
}
