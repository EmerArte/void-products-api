package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emerarteaga/products-api/internal/config"
	"github.com/emerarteaga/products-api/internal/domain/order"
	"github.com/emerarteaga/products-api/internal/domain/product"
	"github.com/emerarteaga/products-api/internal/handler"
	"github.com/emerarteaga/products-api/internal/infra/logger"
	"github.com/emerarteaga/products-api/internal/infra/mongo"
	"github.com/emerarteaga/products-api/internal/repository"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config      *config.Config
	httpServer  *http.Server
	mongoClient *mongo.Client
}

func NewServer(cfg *config.Config) *Server {
	return &Server{config: cfg}
}

func (s *Server) Start() error {
	ctx := context.Background()

	mongoClient, err := mongo.NewClient(ctx, &s.config.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	s.mongoClient = mongoClient
	logger.Info("MongoDB connected successfully")

	productsCollection := mongoClient.Database.Collection("products")
	productRepo := repository.NewProductMongoRepository(productsCollection)

	// Create indexes for products
	if mongoRepo, ok := productRepo.(interface{ CreateIndexes(context.Context) error }); ok {
		if err := mongoRepo.CreateIndexes(ctx); err != nil {
			logger.Warn("failed to create product indexes", "error", err)
		} else {
			logger.Info("product indexes created successfully")
		}
	}

	productService := product.NewService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// Initialize order module
	ordersCollection := mongoClient.Database.Collection("orders")
	orderRepo := repository.NewOrderMongoRepository(ordersCollection)

	// Create indexes for orders
	if mongoRepo, ok := orderRepo.(interface{ CreateIndexes(context.Context) error }); ok {
		if err := mongoRepo.CreateIndexes(ctx); err != nil {
			logger.Warn("failed to create order indexes", "error", err)
		} else {
			logger.Info("order indexes created successfully")
		}
	}

	orderService := order.NewService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	gin.SetMode(s.config.Server.Mode)
	router := SetupRouter(productHandler, orderHandler, s.config)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Server.Port),
		Handler: router,
	}

	go func() {
		logger.Info("starting HTTP server", "port", s.config.Server.Port, "mode", s.config.Server.Mode)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	s.waitForShutdown()
	return nil
}

func (s *Server) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	if err := s.mongoClient.Disconnect(ctx); err != nil {
		logger.Error("failed to disconnect MongoDB", "error", err)
	}

	logger.Info("server stopped gracefully")
}
