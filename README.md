# Products API

A REST API for product management built with Go, Gin, and MongoDB following Clean Architecture principles.

## 🚀 Features

- ✅ CRUD operations for products
- ✅ RESTful API design
- ✅ MongoDB integration
- ✅ Clean Architecture / Hexagonal Architecture
- ✅ Structured logging with slog
- ✅ Request validation
- ✅ Pagination support
- ✅ CORS enabled
- ✅ Graceful shutdown
- ✅ Environment-based configuration

## 📋 Prerequisites

- Go >= 1.21
- MongoDB >= 4.4
- Make (optional)

## 🛠️ Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd products-api
```

2. Install dependencies:
```bash
go mod download
```

3. Copy the example environment file:
```bash
cp .env.example .env
```

4. Update `.env` with your configuration (especially MongoDB URI)

## 🏃 Running the Application

### Using Make:
```bash
make run
```

### Using Go directly:
```bash
go run cmd/api/main.go
```

The API will start on `http://localhost:8080` by default.

## 📚 API Endpoints

### Health Check
- `GET /health` - Check API health

### Products
- `POST /api/v1/products` - Create a new product
- `GET /api/v1/products` - Get all products (with pagination)
- `GET /api/v1/products/:id` - Get a product by ID
- `PUT /api/v1/products/:id` - Update a product
- `DELETE /api/v1/products/:id` - Delete a product

### Example Request (Create Product):
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "stock": 50,
    "category": "Electronics"
  }'
```

### Example Request (Get Products with Pagination):
```bash
curl "http://localhost:8080/api/v1/products?limit=10&offset=0"
```

## 📁 Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── app/
│   │   ├── server.go              # Server initialization
│   │   └── router.go              # Route definitions
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── domain/
│   │   └── product/               # Domain layer
│   │       ├── entity.go          # Product entity
│   │       ├── repository.go      # Repository interface
│   │       └── service.go         # Business logic
│   ├── handler/
│   │   └── product_handler.go     # HTTP handlers
│   ├── repository/
│   │   └── product_mongo_repository.go  # MongoDB implementation
│   ├── infra/
│   │   ├── mongo/
│   │   │   └── client.go          # MongoDB client
│   │   ├── logger/
│   │   │   └── logger.go          # Structured logger
│   │   └── http/
│   │       └── middleware.go      # HTTP middlewares
│   ├── dto/
│   │   └── product_dto.go         # Data transfer objects
│   ├── response/
│   │   └── api_response.go        # API response formats
│   └── errors/
│       └── errors.go              # Custom errors
├── .env.example                    # Example environment file
├── .gitignore
├── Makefile
├── go.mod
└── README.md
```

## ��️ Architecture

This project follows **Clean Architecture** principles:

- **Domain Layer**: Contains business entities and logic (no external dependencies)
- **Application Layer**: Orchestrates use cases and coordinates data flow
- **Infrastructure Layer**: Implements external concerns (database, HTTP, logging)
- **Interface Layer**: Handles HTTP requests and responses

### Dependency Flow:
```
Handler → Service (Domain) → Repository Interface → Repository Implementation → Infrastructure
```

## 🔧 Configuration

Configuration is managed entirely through **environment variables**, following the [12-factor app methodology](https://12factor.net/config).

### Setup:
1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your values:
   ```bash
   # Server Configuration
   SERVER_PORT=8080
   SERVER_MODE=debug  # debug, release, or test
   
   # Database Configuration
   DATABASE_URI=mongodb://localhost:27017
   DATABASE_NAME=products_db
   
   # Logger Configuration
   LOGGER_LEVEL=info   # debug, info, warn, or error
   LOGGER_FORMAT=json  # json or text
   ```

### Environment Variables:
| Variable | Description | Default | Valid Values |
|----------|-------------|---------|--------------|
| `SERVER_PORT` | HTTP server port | `8080` | 1-65535 |
| `SERVER_MODE` | Gin mode | `debug` | `debug`, `release`, `test` |
| `DATABASE_URI` | MongoDB connection URI | `mongodb://localhost:27017` | Valid MongoDB URI |
| `DATABASE_NAME` | MongoDB database name | `products_db` | Non-empty string |
| `LOGGER_LEVEL` | Log level | `info` | `debug`, `info`, `warn`, `error` |
| `LOGGER_FORMAT` | Log output format | `json` | `json`, `text` |

### Configuration Priority:
1. **System environment variables** (highest priority - used in production)
2. **`.env` file** (loaded in development if present)
3. **Default values** (built into the application)

> **Note**: In production, set environment variables directly (via Docker, Kubernetes, cloud platform, etc.). The `.env` file is only for local development convenience.

## 🧪 Testing

Run tests:
```bash
make test
```

Or:
```bash
go test -v -cover ./...
```

## 📝 Code Quality

Format code:
```bash
make fmt
```

Run linter (requires golangci-lint):
```bash
make lint
```

## 🐳 Docker Support (Coming Soon)

A Dockerfile will be provided for containerized deployment.

## 📄 License

MIT License

## 👥 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
