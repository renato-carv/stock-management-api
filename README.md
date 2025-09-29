# API de Gerenciamento de Estoque

API RESTful desenvolvida em Go para gerenciamento de estoque com sistema de autenticação JWT.

## 📋 Descrição

Este projeto implementa uma API completa para controle de estoque com autenticação de usuários. A aplicação permite o cadastro e login de usuários, além de operações CRUD completas para produtos em estoque, com controle de quem criou cada item.

## 🚀 Funcionalidades

### Autenticação
- Registro de novos usuários
- Login com geração de token JWT
- Middleware de autenticação para rotas protegidas
- Hash seguro de senhas com bcrypt

### Gerenciamento de Estoque
- Criar produtos
- Listar todos os produtos
- Atualizar produtos por ID
- Deletar produtos por ID
- Rastreamento de quem criou cada produto

## 🛠️ Tecnologias Utilizadas

- **Go** - Linguagem de programação
- **PostgreSQL** - Banco de dados
- **JWT** - Autenticação via tokens
- **bcrypt** - Hash de senhas
- **UUID** - Identificadores únicos

### Dependências

```go
github.com/joho/godotenv
github.com/lib/pq
github.com/golang-jwt/jwt/v5
golang.org/x/crypto/bcrypt
github.com/google/uuid
```

## 📁 Estrutura do Projeto

```
auth-register-sistem/
├── cmd/
│   └── server/
│       └── main.go               # Ponto de entrada da aplicação
├── internal/
│   ├── config/
│   │   └── database.go           # Configuração do banco de dados
│   ├── handler/
│   │   ├── user_handler.go       # Handlers de usuário
│   │   └── stock_handler.go      # Handlers de estoque
│   ├── middleware/
│   │   └── auth.go               # Middleware de autenticação JWT
│   ├── model/
│   │   ├── user/
│   │   │   └── user.go           # Modelo de usuário
│   │   └── stock/
│   │       └── stock.go          # Modelo de produto
│   ├── repository/
│   │   ├── user_repository.go    # Repositório de usuários
│   │   └── stock_repository.go   # Repositório de estoque
│   └── routes/
│       └── routes.go             # Configuração de rotas
└── .env                          # Variáveis de ambiente
```

## ⚙️ Configuração

### 1. Pré-requisitos

- Go 1.19 ou superior
- PostgreSQL instalado e em execução

### 2. Instalação

Clone o repositório:
```bash
git clone <url-do-repositorio>
cd auth-register-sistem
```

Instale as dependências:
```bash
go mod download
```

### 3. Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```env
DB_HOST=localhost
DB_USERNAME=seu_usuario
DB_PASSWORD=sua_senha
DB_DATABASE=nome_do_banco
DB_PORT=5432
JWT_SECRET=sua_chave_secreta_jwt
```

### 4. Banco de Dados

A aplicação cria automaticamente as tabelas necessárias ao iniciar:

- **users**: Armazena informações dos usuários
- **stock**: Armazena informações dos produtos

## 🚀 Executando a Aplicação

```bash
go run cmd/main.go
```

O servidor iniciará na porta `8080`.

## 📡 Endpoints da API

### Autenticação

#### Registrar Usuário
```http
POST /register
Content-Type: application/json

{
  "name": "João Silva",
  "username": "joaosilva",
  "email": "joao@email.com",
  "password": "senha123"
}
```

**Resposta de Sucesso (201):**
```json
{
  "id": "uuid-do-usuario",
  "message": "User created successfully"
}
```

#### Login
```http
POST /login
Content-Type: application/json

{
  "username": "joaosilva",
  "password": "senha123"
}
```

**Resposta de Sucesso (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "Login successful"
}
```

### Gerenciamento de Estoque

> ⚠️ Todos os endpoints de estoque requerem autenticação via token JWT no header `Authorization: Bearer <token>`

#### Criar Produto
```http
POST /stock
Authorization: Bearer <seu-token>
Content-Type: application/json

{
  "name": "Notebook Dell",
  "quantity": 10
}
```

**Resposta de Sucesso (201):**
```json
{
  "id": "uuid-do-produto",
  "message": "Product created successfully"
}
```

#### Listar Todos os Produtos
```http
GET /stock
Authorization: Bearer <seu-token>
```

**Resposta de Sucesso (200):**
```json
[
  {
    "id": "uuid-do-produto",
    "name": "Notebook Dell",
    "quantity": 10,
    "created_at": "2025-09-29T10:00:00Z",
    "updated_at": "2025-09-29T10:00:00Z",
    "created_by": "uuid-do-usuario"
  }
]
```

#### Atualizar Produto
```http
PUT /stock?id=<uuid-do-produto>
Authorization: Bearer <seu-token>
Content-Type: application/json

{
  "name": "Notebook Dell Inspiron",
  "quantity": 15
}
```

**Resposta de Sucesso (200):**
```json
{
  "id": "uuid-do-produto",
  "message": "Product updated successfully"
}
```

#### Deletar Produto
```http
DELETE /stock?id=<uuid-do-produto>
Authorization: Bearer <seu-token>
```

**Resposta de Sucesso (200):**
```json
{
  "message": "Product deleted successfully"
}
```

## 🔒 Segurança

- Senhas são hasheadas com bcrypt antes de serem armazenadas
- Tokens JWT expiram após 24 horas
- Rotas de estoque protegidas por middleware de autenticação
- Validação de tokens em todas as requisições protegidas

## 📊 Modelos de Dados

### User
```go
{
  ID        uuid.UUID
  Name      string
  Username  string
  Email     string
  Password  string
  CreatedAt time.Time
  UpdatedAt time.Time
}
```

### Stock
```go
{
  ID        uuid.UUID
  Name      string
  Quantity  int
  CreatedAt time.Time
  UpdatedAt time.Time
  CreatedBy uuid.UUID
}
```

## 🧪 Testando a API

Você pode testar a API usando ferramentas como:
- [Postman](https://www.postman.com/)
- [Insomnia](https://insomnia.rest/)
- curl

### Exemplo com curl:

```bash
# Registrar usuário
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name":"João Silva","username":"joaosilva","email":"joao@email.com","password":"senha123"}'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"joaosilva","password":"senha123"}'

# Criar produto (substitua <TOKEN> pelo token recebido no login)
curl -X POST http://localhost:8080/stock \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Notebook Dell","quantity":10}'
```


## 👤 Autor

Renato Carvalho Assunção da Silva