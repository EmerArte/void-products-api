package app

import (
	"github.com/emerarteaga/products-api/internal/config"
	"github.com/emerarteaga/products-api/internal/handler"
	customhttp "github.com/emerarteaga/products-api/internal/infra/http"
	"github.com/gin-gonic/gin"
)

func SetupRouter(productHandler *handler.ProductHandler, orderHandler *handler.OrderHandler, cfg *config.Config) *gin.Engine {
	router := gin.New()
	router.Use(customhttp.Recovery())
	router.Use(customhttp.Logger())
	router.Use(customhttp.CORS(cfg.CORS))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Products API is running"})
	})

	v1 := router.Group("/api/v1")
	{
		// Product CRUD operations
		products := v1.Group("/products")
		{
			products.POST("", productHandler.Create)
			products.GET("/:id", productHandler.GetByID)
			products.PUT("/:id", productHandler.Update)
			products.DELETE("/:id", productHandler.Delete)

			// List products by company or sale point
			products.GET("/company/:company_id", productHandler.GetByCompanyID)
			products.GET("/sale-point/:sale_point_id", productHandler.GetBySalePointID)
		}

		// Categories endpoints
		categories := v1.Group("/categories")
		{
			categories.GET("/company/:company_id", productHandler.GetCategoriesByCompanyID)
			categories.GET("/sale-point/:sale_point_id", productHandler.GetCategoriesBySalePointID)
		}

		// Order endpoints
		orders := v1.Group("/orders")
		{
			// STAGE 1: Create order
			orders.POST("", orderHandler.Create)

			// STAGE 2: Public tracking (no auth required)
			orders.GET("/track/:code", orderHandler.Track)

			// STAGE 3: Partial update (PATCH - no products)
			orders.PATCH("", orderHandler.PartialUpdate)

			// STAGE 4: Modify order (PUT - products allowed)
			orders.PUT("", orderHandler.Modify)

			// STAGE 5: List orders with filters
			orders.GET("", orderHandler.GetAll)

			// STAGE 5: Get metrics and analytics
			orders.GET("/metrics", orderHandler.GetMetrics)

			// Get order by code (admin/internal)
			orders.GET("/:code", orderHandler.GetByCode)
		}
	}

	return router
}
