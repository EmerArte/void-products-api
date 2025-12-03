package app

import (
	"github.com/emerarteaga/products-api/internal/config"
	"github.com/emerarteaga/products-api/internal/handler"
	customhttp "github.com/emerarteaga/products-api/internal/infra/http"
	"github.com/gin-gonic/gin"
)

func SetupRouter(productHandler *handler.ProductHandler, cfg *config.Config) *gin.Engine {
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
	}

	return router
}
