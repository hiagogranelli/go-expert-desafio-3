# Desafio Go Expert - Listagem de Orders (REST, gRPC, GraphQL)

Este projeto implementa um serviço em Go para listar pedidos (`Orders`), expondo a funcionalidade através de três interfaces distintas:

1.  **REST API:** Endpoint `GET /order`
2.  **gRPC:** Serviço `OrderService` com o método `ListOrders`
3.  **GraphQL:** Query `listOrders`

O projeto utiliza Docker e Docker Compose para gerenciar o banco de dados PostgreSQL e a aplicação. As migrações de banco de dados são gerenciadas com `golang-migrate`.

## Pré-requisitos

*   Go (versão 1.22 ou superior)
*   Docker
*   Docker Compose
*   `make` (opcional, para facilitar comandos)
*   `protoc` (instalado localmente se não for usar o build multi-stage do Dockerfile para gerar protos)
*   Plugins `protoc-gen-go` e `protoc-gen-go-grpc` (instalados localmente se não for usar o build multi-stage)
*   `gqlgen` (instalado localmente se não for usar o build multi-stage)
*   `migrate` CLI (instalado localmente se não for usar o build multi-stage ou quiser rodar migrações manualmente)
*   Ferramentas de teste de API como `curl`, Postman, Insomnia, ou a extensão REST Client do VS Code (para usar `api/api.http`).
*   `grpcurl` (para testar o endpoint gRPC).

## Configuração e Execução

1.  **Clone o repositório:**
    ```bash
    git clone <url-do-seu-repositorio>
    cd <diretorio-do-projeto>
    ```

2.  **Construa e suba os containers:**
    Este comando irá construir a imagem Docker da aplicação Go (incluindo a geração de código gRPC e GraphQL e compilação) e iniciar os containers da aplicação e do banco de dados PostgreSQL. As migrações serão executadas automaticamente quando o container da aplicação iniciar.
    ```bash
    docker compose up --build -d
    ```
    *   `-d` executa em modo detached (background). Remova se quiser ver os logs diretamente.

3.  **Verifique os logs (opcional):**
    ```bash
    docker compose logs -f app # Logs da aplicação Go
    docker compose logs -f db  # Logs do PostgreSQL
    ```
    Procure por mensagens indicando que o banco de dados está pronto, as migrações foram aplicadas e os servidores (Web e gRPC) foram iniciados.

## Portas dos Serviços

*   **Servidor Web (REST / GraphQL):** `http://localhost:8080`
    *   Endpoint REST: `GET /order`, `POST /order`
    *   Endpoint GraphQL: `POST /query`
    *   GraphQL Playground: `http://localhost:8080/`
*   **Servidor gRPC:** `localhost:50051`

## Como Usar

### 1. Criar e Listar Orders via REST (Usando `api/api.http`)

*   Abra o arquivo `api/api.http` no VS Code com a extensão [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) instalada.
*   Clique em "Send Request" acima da requisição `@name createOrderRest` para criar um pedido.
*   Clique em "Send Request" acima da requisição `@name listOrdersRest` para listar os pedidos.

### 2. Criar e Listar Orders via GraphQL

*   **Usando o Playground:**
    *   Acesse `http://localhost:8080/` no seu navegador.
    *   No painel esquerdo, use a mutation `createOrder` para criar dados:
        ```graphql
        mutation CreateSampleOrder {
          createOrder(input: {price: 99.90, tax: 5.50}) {
            id
            price
            tax
            finalPrice
          }
        }
        ```
    *   Use a query `listOrders` para buscar os dados:
        ```graphql
        query ListAllOrders {
          listOrders {
            id
            price
            tax
            finalPrice
          }
        }
        ```
*   **Usando `api/api.http`:**
    *   Execute as requisições `@name createOrderGraphQL` e `@name listOrdersGraphQL` no arquivo `api/api.http`.

### 3. Listar Orders via gRPC (Usando `grpcurl`)

*   Certifique-se de ter o `grpcurl` instalado ([https://github.com/fullstorydev/grpcurl](https://github.com/fullstorydev/grpcurl)).
*   Liste os serviços disponíveis (opcional):
    ```bash
    grpcurl -plaintext localhost:50051 list
    # Deverá listar pb.OrderService
    ```
*   Liste os métodos do serviço (opcional):
    ```bash
    grpcurl -plaintext localhost:50051 list pb.OrderService
    # Deverá listar CreateOrder e ListOrders
    ```
*   Execute a chamada `ListOrders`:
    ```bash
    grpcurl -plaintext localhost:50051 pb.OrderService/ListOrders
    ```
*   Execute a chamada `CreateOrder` (exemplo):
    ```bash
    grpcurl -plaintext -d '{"price": 55.25, "tax": 4.75}' localhost:50051 pb.OrderService/CreateOrder
    ```

## Parando a Aplicação

```bash
docker compose down
