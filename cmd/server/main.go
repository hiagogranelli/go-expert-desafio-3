package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"

	"desafio-cleanarchitecture/internal/application/usecase"
	"desafio-cleanarchitecture/internal/infra/database"
	"desafio-cleanarchitecture/internal/infra/graphql/graph"
	"desafio-cleanarchitecture/internal/infra/grpc/pb"
	grpcService "desafio-cleanarchitecture/internal/infra/grpc/service"
	webHandler "desafio-cleanarchitecture/internal/infra/web/handler"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	webPort := os.Getenv("WEB_SERVER_PORT")
	grpcPort := os.Getenv("GRPC_SERVER_PORT")
	graphqlPath := os.Getenv("GRAPHQL_SERVER_PATH")

	if dbURL == "" || webPort == "" || grpcPort == "" || graphqlPath == "" {
		log.Fatal("Environment variables DATABASE_URL, WEB_SERVER_PORT, GRPC_SERVER_PORT, GRAPHQL_SERVER_PATH must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection successful!")

	orderRepository := database.NewOrderRepositoryDb(db)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	webOrderHandler := webHandler.NewWebOrderHandler(createOrderUseCase, listOrdersUseCase)

	gqlSrv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			CreateOrderUseCase: createOrderUseCase,
			ListOrdersUseCase:  listOrdersUseCase,
		},
	}))

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Post("/order", webOrderHandler.CreateOrder)
	router.Get("/order", webOrderHandler.ListOrders)

	router.Handle("/", playground.Handler("GraphQL playground", graphqlPath))
	router.Handle(graphqlPath, gqlSrv)

	log.Printf("GraphQL playground available at http://localhost:%s/", webPort)
	log.Printf("GraphQL endpoint available at http://localhost:%s%s", webPort, graphqlPath)
	log.Printf("REST endpoint available at http://localhost:%s/order", webPort)

	go func() {
		log.Printf("Starting WEB server on port %s", webPort)
		if err := http.ListenAndServe(":"+webPort, router); err != nil {
			log.Fatalf("Failed to start WEB server: %v", err)
		}
	}()

	grpcOrderService := grpcService.NewOrderGrpcService(createOrderUseCase, listOrdersUseCase)
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, grpcOrderService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	log.Printf("Starting gRPC server on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
