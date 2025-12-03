---
applyTo: "**"
---

# Instrucciones para Copilot: Backend en Go + Gin + MongoDB

## 1. Contexto del proyecto

Este repositorio implementa un **backend REST** en **Go (>= 1.21)** usando:

- **Framework HTTP**: `github.com/gin-gonic/gin`
- **Base de datos**: **MongoDB** usando el **driver oficial** `go.mongodb.org/mongo-driver/mongo`
- **Arquitectura**: Clean Architecture / Hexagonal

### Principios arquitectónicos obligatorios:

1. **Separación de capas estricta**: `Handler → Service → Repository → Infrastructure`
2. **Inyección de dependencias**: Todas las dependencias deben inyectarse vía constructores
3. **Inversión de dependencias**: Las capas externas dependen de abstracciones (interfaces) definidas en el dominio
4. **Sin dependencias circulares**: Nunca importar capas superiores desde capas inferiores
5. **Código idiomático Go**: Seguir convenciones estándar de Go (naming, error handling, etc.)

### Objetivos de calidad:

- Alto rendimiento y baja latencia
- Testabilidad completa (unit + integration tests)
- Mantenibilidad y escalabilidad
- Documentación clara y concisa

## 2. Stack tecnológico y librerías

### 2.1. Framework HTTP (obligatorio)

**Librería**: `github.com/gin-gonic/gin`

**Responsabilidades**:

- Configuración del servidor HTTP y grupos de rutas
- Implementación de middlewares (CORS, logging, recovery, etc.)
- Binding y validación de requests
- Serialización/deserialización JSON automática
- Manejo de errores HTTP con códigos apropiados

**Reglas**:

- ❌ **NUNCA** importar Gin en `internal/domain`
- ✅ Usar solo en `internal/handler` y `internal/app`
- ✅ Aprovechar `c.ShouldBindJSON()` para validación automática
- ✅ Usar grupos de rutas para versioning: `v1 := router.Group("/api/v1")`

### 2.2. MongoDB (obligatorio)

**Librerías**:

- `go.mongodb.org/mongo-driver/mongo`
- `go.mongodb.org/mongo-driver/mongo/options`
- `go.mongodb.org/mongo-driver/bson`
- `go.mongodb.org/mongo-driver/bson/primitive` (para ObjectID)

**Reglas críticas**:

- ✅ **Siempre** usar `context.WithTimeout` (5-10s para operaciones normales)
- ✅ **Cliente singleton**: Inicializar UNA VEZ en `main.go` y pasar a repositorios
- ✅ **Pool de conexiones**: Configurar MaxPoolSize (ej: 100) y timeouts
- ✅ **Índices**: Crear índices en inicialización, no por request
- ✅ **Proyecciones**: Usar `options.Find().SetProjection()` para optimizar queries
- ❌ **NUNCA** devolver `bson.M` o `bson.D` fuera de `internal/repository`
- ❌ **NUNCA** exponer `*mongo.Collection` al dominio
- ✅ Mapear documentos BSON a entidades del dominio en el repositorio

### 2.3. Logging (obligatorio)

**Librería recomendada**: `log/slog` (Go >= 1.21) - logging estructurado nativo

**Alternativas opcionales**:

- `go.uber.org/zap` - alto rendimiento, producción pesada
- `github.com/rs/zerolog` - logging JSON ultra-rápido

**Reglas estrictas**:

- ❌ **PROHIBIDO**: `fmt.Println`, `log.Printf`, `panic()` en lógica de negocio
- ✅ **Centralizado**: Todo logging DEBE pasar por `internal/infra/logger`
- ✅ **Estructurado**: Usar campos clave-valor: `logger.Info("user created", "user_id", id, "email", email)`
- ✅ **Niveles apropiados**:
  - `Debug`: Información de desarrollo (queries, datos internos)
  - `Info`: Eventos normales del sistema (startup, requests exitosos)
  - `Warn`: Situaciones anormales pero manejables (retry, caché miss)
  - `Error`: Errores que requieren atención (DB down, validación fallida)
- ✅ **Contexto**: Incluir request_id, user_id cuando sea relevante
- ✅ **No loggear**: Passwords, tokens, datos sensibles (PII)

### 2.4. Configuración y variables de entorno

**Librería**:

- `github.com/joho/godotenv` - cargar .env (solo desarrollo)

**Patrón de configuración**:

1. Definir structs tipados en `internal/config/config.go`
2. Cargar variables de entorno con `os.Getenv()` usando helpers para conversión de tipos
3. Proporcionar valores por defecto para desarrollo
4. Validar configuración al inicio (fail-fast)

**Reglas**:

- ❌ **NUNCA** llamar `os.Getenv()` directamente desde handlers/services
- ✅ **Inyectar** config struct a través de constructores
- ✅ **Validación**: Implementar método `Validate()` en struct Config
- ✅ **Secretos**: Cargar desde variables de entorno, no hardcodear
- ✅ **12-factor app**: Toda configuración DEBE venir de environment variables
- ✅ **Ejemplo de estructura**:

```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Logger   LoggerConfig
}

func LoadConfig() (*Config, error) {
    config := &Config{
        Server: ServerConfig{
            Port: getEnvAsInt("SERVER_PORT", 8080),
            Mode: getEnv("SERVER_MODE", "debug"),
        },
        // ... más campos
    }

    if err := config.Validate(); err != nil {
        return nil, err
    }

    return config, nil
}

// Helper para strings con default
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

// Helper para integers con default
func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}
```

**Uso en main.go**:

```go
// Cargar .env si existe (desarrollo)
if err := godotenv.Load(); err != nil {
    log.Println("No .env file found, using system environment variables")
}

// Cargar configuración desde environment variables
cfg, err := config.LoadConfig()
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
```

### 2.5. Validación

**Librería**: `github.com/go-playground/validator/v10` (incluida en Gin)

**Estrategia de validación en capas**:

1. **Capa HTTP (Handler)**: Validación de formato y tipos con binding tags
2. **Capa Dominio (Service)**: Validación de reglas de negocio

**Tags comunes**:

```go
type CreateUserDTO struct {
    Email    string `json:"email" binding:"required,email"`
    Name     string `json:"name" binding:"required,min=3,max=50"`
    Age      int    `json:"age" binding:"required,gte=18,lte=120"`
    Password string `json:"password" binding:"required,min=8"`
}
```

**Reglas**:

- ✅ Usar `c.ShouldBindJSON()` en handlers (no detiene ejecución)
- ✅ Manejar errores de validación centralmente con respuesta 400
- ✅ Devolver mensajes claros: `{"field": "email", "error": "invalid format"}`
- ✅ Validaciones custom para reglas complejas
- ❌ No mezclar validación HTTP con lógica de negocio

### 2.6. Testing

Usar:

- testing (nativo)
- github.com/stretchr/testify/assert
- github.com/stretchr/testify/require

### 2.7. Utilidades adicionales

- UUID: github.com/google/uuid
- HTTP resiliente: github.com/hashicorp/go-retryablehttp (opcional)
- Middlewares:

  - CORS
  - Logging
  - Request ID
  - Recovery

## 3\. Estructura de carpetas del proyecto

La estructura sigue patrones recomendados para Go + Gin + Clean Architecture:
.
├── cmd/
│ └── api/
│ └── main.go
├── internal/
│ ├── app/
│ │ ├── server.go
│ │ └── router.go
│ ├── config/
│ │ └── config.go
│ ├── domain/
│ │ └── user/
│ │ ├── entity.go
│ │ ├── repository.go
│ │ └── service.go
│ ├── handler/
│ │ └── user_handler.go
│ ├── repository/
│ │ └── user_mongo_repository.go
│ ├── infra/
│ │ ├── mongo/
│ │ │ └── client.go
│ │ ├── logger/
│ │ │ └── logger.go
│ │ └── http/
│ │ └── middleware.go
│ ├── dto/
│ │ └── user_dto.go
│ ├── response/
│ │ └── api_response.go
│ ├── errors/
│ │ └── errors.go
│ ├── util/
│ │ └── time.go
│ └── tests/
│ └── user_handler_test.go
├── pkg/
│ └── (utilidades reutilizables entre proyectos, si se requieren)
├── docs/
│ └── openapi.yaml
├── scripts/
│ └── dev.sh
├── .env.example
├── Dockerfile
├── Makefile
├── go.mod
└── go.sum
Esta estructura organiza el código en capas, siguiendo principios de Clean Architecture y buenas prácticas de Go.

## 4\. Propósito de cada carpeta

### 4.1. cmd/api

**Punto de entrada del servidor**.

Responsable de:

- Cargar configuración
- Inicializar logger
- Crear cliente MongoDB
- Construir servidor HTTP
- Manejar graceful shutdown

### 4.2. internal/app

Contiene:

- server.go: inicializa Gin y dependencias
- router.go: define rutas y middlewares

Reglas:

- No escribir rutas en main.go
- Inyección de dependencias aquí

### 4.3. internal/config

Contiene structs de configuración y LoadConfig().

### 4.4. internal/domain

Capa de negocio (Clean Architecture).

Contiene:

- **Entidades**
- **Interfaces** (UserRepository)
- **Servicios (casos de uso)**

Reglas importantes:

- No importar Gin
- No importar MongoDB
- No importar librerías de configuración (godotenv)
- No retornar tipos BSON
- Solo lógica de negocio

### 4.6. internal/handler

Controladores HTTP responsables de:

- Recibir requests
- Convertir requests -> DTOs
- Llamar a servicios
- Manejar errores
- Devolver respuestas JSON

### 4.7. internal/repository

Implementaciones concretas para MongoDB.

Reglas:

- Usar mongo.Collection singleton
- Mapear documentos BSON a entidades
- Usar contextos con timeout

### 4.8. internal/infra/mongo

Inicialización del cliente MongoDB.

Incluye:

- Cliente global (singleton)
- Configuración de pool
- Timeouts
- Indices

### 4.9. internal/infra/logger

Logger global de la aplicación.

### 4.10. internal/infra/http

Middlewares como:

- Logging
- CORS
- Request ID
- Recovery

### 4.11. internal/dto

Estructuras de entrada/salida para los handlers.

### 4.12. internal/response

Formato estandarizado de respuestas:

- SuccessResponse
- ErrorResponse
- Paginación

### 4.13. internal/errors

Errores tipados:

- NotFoundError
- ValidationError
- DomainError

### 4.14. internal/util

Utilidades genéricas.

### 4.15. internal/tests

Pruebas unitarias y de integración.

## 5. Reglas arquitectónicas (OBLIGATORIAS)

### 5.1. Flujo de dependencias

**Dirección correcta** (solo hacia adentro):

```
Handler → Service (Domain) → Repository Interface → Repository Implementation → Infrastructure
```

**Prohibiciones**:

- ❌ Domain NO puede importar: Gin, MongoDB, godotenv, HTTP
- ❌ Service NO puede importar: Handler, Repository implementation
- ❌ Handler NO puede importar: Repository, Infrastructure
- ✅ Handler SÍ puede importar: Service interfaces, DTOs, Errors

**Inyección de dependencias**:

```go
// ✅ Correcto: Inyectar dependencias en constructor
func NewUserService(repo domain.UserRepository) *UserService {
    return &UserService{repo: repo}
}

// ❌ Incorrecto: Crear dependencias dentro del servicio
func NewUserService() *UserService {
    repo := mongo.NewUserRepository() // ❌ Acoplamiento directo
    return &UserService{repo: repo}
}
```

### 5.2. Manejo de contextos

**Reglas**:

- ✅ **Siempre** pasar `context.Context` como primer parámetro
- ✅ **Siempre** usar `context.WithTimeout()` para operaciones I/O
- ✅ Propagar contexto a través de todas las capas
- ✅ Cancelar contexto para operaciones largas
- ❌ **NUNCA** usar `context.Background()` en handlers (usar `c.Request.Context()`)

**Timeouts recomendados**:

- HTTP handlers: 30s
- DB queries: 5-10s
- External APIs: 10-15s

### 5.3. Manejo de errores

**Principios**:

- ❌ **NUNCA** ignorar errores con `_`
- ✅ **Siempre** propagar errores hacia arriba o manejarlos
- ✅ Usar errores tipados en `internal/errors/`
- ✅ Loggear errores antes de propagarlos
- ✅ Wrappear errores con contexto: `fmt.Errorf("failed to create user: %w", err)`

**Errores HTTP**:

```go
// ✅ Mapeo correcto de errores
if errors.Is(err, domain.ErrNotFound) {
    return c.JSON(404, ErrorResponse{Message: "User not found"})
}
if errors.Is(err, domain.ErrValidation) {
    return c.JSON(400, ErrorResponse{Message: err.Error()})
}
// Error genérico
return c.JSON(500, ErrorResponse{Message: "Internal server error"})
```

### 5.4. Logging estratégico

**Qué loggear**:

- ✅ Inicio/parada del servidor
- ✅ Conexión/desconexión a DB
- ✅ Errores y excepciones
- ✅ Request/Response en endpoints críticos
- ✅ Operaciones de negocio importantes (creación de usuario, pago, etc.)

**Qué NO loggear**:

- ❌ Passwords, tokens, API keys
- ❌ Datos personales sensibles (a menos que sea necesario y esté hasheado)
- ❌ Loops o ejecuciones de alta frecuencia en Debug (afecta rendimiento)

### 5.5. MongoDB - Mejores prácticas

**Inicialización (una sola vez)**:

```go
// En main.go
client, err := mongo.Connect(ctx, options.Client().
    ApplyURI(config.MongoURI).
    SetMaxPoolSize(100).
    SetTimeout(10*time.Second))

db := client.Database("mydb")
usersCollection := db.Collection("users")

// Crear índices
userRepo := repository.NewUserMongoRepository(usersCollection)
if err := userRepo.CreateIndexes(ctx); err != nil {
    log.Fatal(err)
}
```

**Operaciones**:

- ✅ Usar proyecciones para queries grandes
- ✅ Implementar paginación con `Skip()` y `Limit()`
- ✅ Usar índices compuestos para queries complejas
- ✅ Transactions para operaciones atómicas multi-documento
- ❌ Evitar `cursor.All()` con colecciones grandes (usar streaming)

### 5.6. Testing (obligatorio)

**Tipos de tests**:

1. **Unit tests**: Lógica de servicios (mock de repositorios)
2. **Integration tests**: Repositorios con MongoDB (testcontainers)
3. **E2E tests**: Endpoints completos (opcional)

**Herramientas**:

- `testing` (nativo)
- `github.com/stretchr/testify` (assertions)
- `github.com/testcontainers/testcontainers-go` (MongoDB en tests)

**Cobertura mínima esperada**:

- Services: 80%+
- Repositories: 70%+
- Handlers: 60%+

**Estructura de test**:

```go
func TestUserService_Create(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    // Act
    user, err := service.Create(ctx, input)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    mockRepo.AssertExpectations(t)
}
```

## 6\. Ejemplo de flujo esperado

Handler → DTO → Servicio → Repositorio → Infraestructura

Handler:

- Recibe request
- Valida
- Llama servicio
- Maneja errores
- Devuelve JSON

Servicio:

- Procesa lógica de negocio

Repositorio:

- Interactúa con MongoDB

## 7. Estándares de código Go

### 7.1. Convenciones de nombres

**Variables y funciones**:

- ✅ `camelCase` para privados: `userService`, `findByID`
- ✅ `PascalCase` para públicos: `UserService`, `CreateUser`
- ✅ Acrónimos en mayúscula: `ID`, `HTTP`, `URL`, `JSON`
- ❌ No usar snake_case (excepto en tags JSON/DB)

**Paquetes**:

- ✅ Nombres cortos, singulares: `user`, `auth`, `config`
- ❌ Evitar: `users`, `user_service`, `userService`

**Interfaces**:

- ✅ Sufijo `-er` para interfaces de una función: `Reader`, `Writer`, `Validator`
- ✅ Para múltiples funciones: `UserRepository`, `UserService`

### 7.2. Código idiomático

**Error handling**:

```go
// ✅ Correcto
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// ❌ Incorrecto
if err == nil {
    // happy path
} else {
    return err
}
```

**Composición sobre herencia**:

```go
// ✅ Correcto: Composición
type UserService struct {
    repo   UserRepository
    logger *slog.Logger
}

// ❌ Go no tiene herencia clásica
```

**Interfaces pequeñas**:

```go
// ✅ Correcto: Interface segregation
type UserReader interface {
    FindByID(ctx context.Context, id string) (*User, error)
}

type UserWriter interface {
    Create(ctx context.Context, user *User) error
}

// ❌ Evitar: Interfaces gigantes con 10+ métodos
```

### 7.3. Organización de código

**Orden en archivos**:

1. Package declaration
2. Imports (agrupados: stdlib, external, internal)
3. Constants
4. Types (structs, interfaces)
5. Constructors
6. Methods (públicos primero, privados después)
7. Helper functions

**Imports**:

```go
import (
    // Standard library
    "context"
    "fmt"

    // External packages
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"

    // Internal packages
    "myapp/internal/domain/user"
    "myapp/internal/response"
)
```

### 7.4. Performance y optimización

**Buenas prácticas**:

- ✅ Usar punteros para structs grandes (>64 bytes)
- ✅ Preallocar slices cuando se conoce el tamaño: `make([]User, 0, 100)`
- ✅ Usar `strings.Builder` para concatenación de strings
- ✅ Cerrar recursos: `defer response.Body.Close()`
- ✅ Usar pools para objetos reutilizables: `sync.Pool`
- ❌ Evitar goroutines sin control (usar worker pools)

### 7.5. Documentación

**Comentarios de documentación**:

```go
// UserService handles business logic for user management.
// It coordinates between handlers and repositories.
type UserService struct {
    repo UserRepository
}

// Create creates a new user after validating business rules.
// Returns ErrValidation if input is invalid.
// Returns ErrDuplicate if user already exists.
func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*User, error) {
    // Implementation
}
```

**Cuándo documentar**:

- ✅ Funciones y tipos públicos (exportados)
- ✅ Comportamientos no obvios
- ✅ Errores que puede retornar
- ❌ No documentar lo obvio: `// GetID returns the ID`

## 8. Checklist de calidad

Antes de considerar código como completo, verificar:

- [ ] ✅ Sin errores de compilación o linting
- [ ] ✅ Tests escritos y pasando
- [ ] ✅ Errores manejados correctamente (sin `_` innecesarios)
- [ ] ✅ Contextos propagados con timeouts
- [ ] ✅ Logging apropiado en puntos clave
- [ ] ✅ Sin dependencias circulares o hacia arriba
- [ ] ✅ Interfaces usadas para abstracción
- [ ] ✅ Configuración inyectada, no hardcodeada
- [ ] ✅ Código formateado con `gofmt`
- [ ] ✅ Imports organizados con `goimports`
- [ ] ✅ Sin datos sensibles en logs o código
