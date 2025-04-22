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
	_ "github.com/lib/pq" // Driver PostgreSQL

	"desafio-cleanarchitecture/internal/application/usecase"
	"desafio-cleanarchitecture/internal/infra/database"
	"desafio-cleanarchitecture/internal/infra/graphql/graph" // Pacote gerado
	"desafio-cleanarchitecture/internal/infra/grpc/pb"       // Pacote gerado
	grpcService "desafio-cleanarchitecture/internal/infra/grpc/service"
	webHandler "desafio-cleanarchitecture/internal/infra/web/handler"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection" // Para grpcurl
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	webPort := os.Getenv("WEB_SERVER_PORT")
	grpcPort := os.Getenv("GRPC_SERVER_PORT")
	graphqlPath := os.Getenv("GRAPHQL_SERVER_PATH")

	if dbURL == "" || webPort == "" || grpcPort == "" || graphqlPath == "" {
		log.Fatal("Environment variables DATABASE_URL, WEB_SERVER_PORT, GRPC_SERVER_PORT, GRAPHQL_SERVER_PATH must be set")
	}

	// --- Conexão com Banco de Dados ---
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ping para verificar a conexão
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection successful!")

	// --- Inicialização de Dependências ---
	orderRepository := database.NewOrderRepositoryDb(db)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	// --- Configuração do Servidor Web (REST + GraphQL) ---
	webOrderHandler := webHandler.NewWebOrderHandler(createOrderUseCase, listOrdersUseCase)

	// Configuração do GraphQL
	gqlSrv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{ // Injete os use cases no resolver do GraphQL
			CreateOrderUseCase: createOrderUseCase,
			ListOrdersUseCase:  listOrdersUseCase,
		},
	}))

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second)) // Exemplo de middleware

	// Rotas REST
	router.Post("/order", webOrderHandler.CreateOrder)
	router.Get("/order", webOrderHandler.ListOrders)

	// Rota GraphQL
	router.Handle("/", playground.Handler("GraphQL playground", graphqlPath)) // Playground UI
	router.Handle(graphqlPath, gqlSrv)                                        // Endpoint GraphQL

	log.Printf("GraphQL playground available at http://localhost:%s/", webPort)
	log.Printf("GraphQL endpoint available at http://localhost:%s%s", webPort, graphqlPath)
	log.Printf("REST endpoint available at http://localhost:%s/order", webPort)

	// Iniciar servidor HTTP em uma goroutine
	go func() {
		log.Printf("Starting WEB server on port %s", webPort)
		if err := http.ListenAndServe(":"+webPort, router); err != nil {
			log.Fatalf("Failed to start WEB server: %v", err)
		}
	}()

	// --- Configuração do Servidor gRPC ---
	grpcOrderService := grpcService.NewOrderGrpcService(createOrderUseCase, listOrdersUseCase)
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, grpcOrderService)
	reflection.Register(grpcServer) // Habilita reflection para ferramentas como grpcurl

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	log.Printf("Starting gRPC server on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

	// O select {} manteria a main goroutine viva, mas como o grpcServer.Serve() bloqueia,
	// não é estritamente necessário aqui se o gRPC for o último a iniciar.
	// select {}
}
