# 🚀 Getting Started with Products API

## Quick Start Guide

### 1. Start MongoDB

#### Option A: Using Docker Compose (Recommended)
```bash
docker-compose up -d
```

#### Option B: Using Docker directly
```bash
docker run -d --name products-mongo -p 27017:27017 mongo:7.0
```

#### Option C: Local MongoDB
If you have MongoDB installed locally, ensure it's running on port 27017.

### 2. Verify MongoDB is Running

```bash
# Check Docker container status
docker ps | grep mongo

# Or check MongoDB connection
mongosh --eval "db.adminCommand('ping')"
```

### 3. Run the Application

#### Option A: Using Make
```bash
make run
```

#### Option B: Using Go directly
```bash
go run cmd/api/main.go
```

You should see output like:
```
{"time":"...","level":"INFO","msg":"application starting","version":"1.0.0"}
{"time":"...","level":"INFO","msg":"connecting to MongoDB","uri":"mongodb://localhost:27017","database":"products_db"}
{"time":"...","level":"INFO","msg":"successfully connected to MongoDB"}
{"time":"...","level":"INFO","msg":"starting HTTP server","port":8080,"mode":"debug"}
```

### 4. Test the API

#### Health Check
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "message": "Products API is running"
}
```

#### Create a Product
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

#### Get All Products
```bash
curl http://localhost:8080/api/v1/products
```

#### Get Product by ID
```bash
curl http://localhost:8080/api/v1/products/{product_id}
```

#### Update a Product
```bash
curl -X PUT http://localhost:8080/api/v1/products/{product_id} \
  -H "Content-Type: application/json" \
  -d '{
    "price": 899.99,
    "stock": 45
  }'
```

#### Delete a Product
```bash
curl -X DELETE http://localhost:8080/api/v1/products/{product_id}
```

### 5. Run Complete API Tests

Use the provided test script:
```bash
./test_api.sh
```

This will run through all CRUD operations automatically.

## 📝 Configuration

The application is configured entirely through **environment variables** (following the [12-factor app methodology](https://12factor.net/config)).

### Initial Setup

1. **Copy the example environment file**:
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your values**:
   ```bash
   nano .env  # or use your preferred editor
   ```

3. **Configure variables** (see table below for all options)

### Environment Variables

| Variable | Description | Default | Valid Values |
|----------|-------------|---------|--------------|
| `SERVER_PORT` | HTTP server port | `8080` | 1-65535 |
| `SERVER_MODE` | Gin mode | `debug` | `debug`, `release`, `test` |
| `DATABASE_URI` | MongoDB connection URI | `mongodb://localhost:27017` | Valid MongoDB URI |
| `DATABASE_NAME` | MongoDB database name | `products_db` | Non-empty string |
| `DATABASE_MAX_POOL_SIZE` | Max MongoDB connections | `100` | Positive integer |
| `DATABASE_TIMEOUT` | DB operation timeout (seconds) | `10` | Positive integer |
| `LOGGER_LEVEL` | Log level | `info` | `debug`, `info`, `warn`, `error` |
| `LOGGER_FORMAT` | Log output format | `json` | `json`, `text` |

### Example `.env` Configuration

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_MODE=debug  # Use 'release' in production

# Database Configuration
DATABASE_URI=mongodb://localhost:27017
DATABASE_NAME=products_db
DATABASE_MAX_POOL_SIZE=100
DATABASE_TIMEOUT=10

# Logger Configuration
LOGGER_LEVEL=debug  # Use 'info' or 'warn' in production
LOGGER_FORMAT=json
```

### Configuration Priority

1. **System environment variables** (highest - used in production)
2. **`.env` file** (local development only)
3. **Default values** (built-in fallbacks)

> **Production Note**: Set environment variables directly via your deployment platform (Docker, Kubernetes, cloud provider) instead of using `.env` files.

### Available Modes

- `debug`: Detailed logging and Gin debug mode (development)
- `release`: Optimized for production, minimal logging
- `test`: For running automated tests

### Logger Levels

- `debug`: Most verbose, includes all details (queries, internal state)
- `info`: General application flow and events
- `warn`: Warnings about potential issues
- `error`: Only errors that require attention

## 🛠️ Development Commands

```bash
# Run the application
make run

# Build the binary
make build

# Run tests
make test

# Format code
make fmt

# Clean build artifacts
make clean

# Show available commands
make help
```

## 🔍 Troubleshooting

### MongoDB Connection Error

**Problem**: `failed to connect to MongoDB`

**Solutions**:
1. Verify MongoDB is running: `docker ps | grep mongo`
2. Check MongoDB URI in `.env` file
3. Ensure port 27017 is not in use: `lsof -i :27017`
4. Restart MongoDB: `docker-compose restart`

### Port Already in Use

**Problem**: `bind: address already in use`

**Solutions**:
1. Change `SERVER_PORT` in `.env` to a different port
2. Kill the process using the port: `lsof -ti:8080 | xargs kill -9`

### Module Not Found Errors

**Problem**: `package xxx is not in GOROOT`

**Solutions**:
```bash
go mod tidy
go mod download
```

## 📊 Project Stats

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: MongoDB
- **Architecture**: Clean Architecture / Hexagonal
- **Lines of Code**: ~800+
- **Test Coverage**: Ready for unit & integration tests

## 🎯 Next Steps

1. ✅ Add unit tests for services
2. ✅ Add integration tests for repositories
3. ✅ Implement authentication (JWT)
4. ✅ Add API documentation (Swagger/OpenAPI)
5. ✅ Implement caching (Redis)
6. ✅ Add rate limiting
7. ✅ Deploy to production

## 📚 Additional Resources

- [Gin Documentation](https://gin-gonic.com/docs/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Best Practices](https://golang.org/doc/effective_go)

## 🤝 Support

If you encounter any issues:
1. Check the logs for detailed error messages
2. Verify all prerequisites are installed
3. Ensure MongoDB is running and accessible
4. Check the configuration in `.env`

Happy coding! 🎉
