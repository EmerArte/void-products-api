package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/emerarteaga/products-api/internal/domain/product"
	"github.com/emerarteaga/products-api/internal/dto"
	"github.com/emerarteaga/products-api/internal/infra/logger"
	"github.com/emerarteaga/products-api/internal/response"
	"github.com/gin-gonic/gin"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	service *product.Service
}

// NewProductHandler creates a new product handler
func NewProductHandler(service *product.Service) *ProductHandler {
	return &ProductHandler{service: service}
}

// Create handles POST /api/v1/products
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("invalid request body", "error", err)
		// Format validation errors for user-friendly response
		errorMsg, details := FormatValidationErrors(err)
		if details != nil {
			// Convert to response format
			responseDetails := make([]response.ValidationErrorDetail, len(details))
			for i, d := range details {
				responseDetails[i] = response.ValidationErrorDetail{
					Field:   d.Field,
					Message: d.Message,
				}
			}
			response.ValidationError(c, http.StatusBadRequest, errorMsg, "Validation failed", responseDetails)
			return
		}
		response.Error(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	// Convert DTO to service input
	input := req.ToCreateInput()

	p, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		logger.Error("failed to create product", "error", err)
		response.Error(c, http.StatusInternalServerError, err, "Failed to create product")
		return
	}

	logger.Info("product created", "product_id", p.ID, "company_id", p.CompanyID, "sale_point_id", p.SalePointID)
	response.Success(c, http.StatusCreated, p, "Product created successfully")
}

// GetByID handles GET /api/v1/products/:id
func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	p, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, product.ErrProductNotFound) {
			response.Error(c, http.StatusNotFound, err, "Product not found")
			return
		}
		logger.Error("failed to get product", "error", err, "product_id", id)
		response.Error(c, http.StatusInternalServerError, err, "Failed to get product")
		return
	}

	response.Success(c, http.StatusOK, p, "")
}

// GetByCompanyID handles GET /api/v1/products/company/:company_id
func (h *ProductHandler) GetByCompanyID(c *gin.Context) {
	companyID := c.Param("company_id")

	// Parse filters from query parameters
	filters := h.parseFilters(c)

	products, total, err := h.service.GetByCompanyID(c.Request.Context(), companyID, filters)
	if err != nil {
		logger.Error("failed to get products", "error", err, "company_id", companyID)
		response.Error(c, http.StatusInternalServerError, err, "Failed to get products")
		return
	}

	// Convert to list responses (simplified view)
	listResponses := dto.ToListResponses(products)
	response.Paginated(c, http.StatusOK, listResponses, total, filters.Limit, filters.Offset)
}

// GetBySalePointID handles GET /api/v1/products/sale-point/:sale_point_id
func (h *ProductHandler) GetBySalePointID(c *gin.Context) {
	salePointID := c.Param("sale_point_id")

	// Parse filters from query parameters
	filters := h.parseFilters(c)

	products, total, err := h.service.GetBySalePointID(c.Request.Context(), salePointID, filters)
	if err != nil {
		logger.Error("failed to get products", "error", err, "sale_point_id", salePointID)
		response.Error(c, http.StatusInternalServerError, err, "Failed to get products")
		return
	}

	// Convert to list responses (simplified view)
	listResponses := dto.ToListResponses(products)
	response.Paginated(c, http.StatusOK, listResponses, total, filters.Limit, filters.Offset)
}

// Update handles PUT /api/v1/products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("invalid request body", "error", err)
		// Format validation errors for user-friendly response
		errorMsg, details := FormatValidationErrors(err)
		if details != nil {
			// Convert to response format
			responseDetails := make([]response.ValidationErrorDetail, len(details))
			for i, d := range details {
				responseDetails[i] = response.ValidationErrorDetail{
					Field:   d.Field,
					Message: d.Message,
				}
			}
			response.ValidationError(c, http.StatusBadRequest, errorMsg, "Validation failed", responseDetails)
			return
		}
		response.Error(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	// Convert DTO to service input
	input := req.ToUpdateInput()

	p, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, product.ErrProductNotFound) {
			response.Error(c, http.StatusNotFound, err, "Product not found")
			return
		}
		logger.Error("failed to update product", "error", err, "product_id", id)
		response.Error(c, http.StatusInternalServerError, err, "Failed to update product")
		return
	}

	logger.Info("product updated", "product_id", id)
	response.Success(c, http.StatusOK, p, "Product updated successfully")
}

// Delete handles DELETE /api/v1/products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, product.ErrProductNotFound) {
			response.Error(c, http.StatusNotFound, err, "Product not found")
			return
		}
		logger.Error("failed to delete product", "error", err, "product_id", id)
		response.Error(c, http.StatusInternalServerError, err, "Failed to delete product")
		return
	}

	logger.Info("product deleted", "product_id", id)
	response.Success(c, http.StatusOK, nil, "Product deleted successfully")
}

// GetCategoriesByCompanyID handles GET /api/v1/categories/company/:company_id
func (h *ProductHandler) GetCategoriesByCompanyID(c *gin.Context) {
	companyID := c.Param("company_id")

	categories, err := h.service.GetCategoriesByCompanyID(c.Request.Context(), companyID)
	if err != nil {
		logger.Error("failed to get categories", "error", err, "company_id", companyID)
		response.Error(c, http.StatusInternalServerError, err, "Failed to get categories")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"categories": categories}, "")
}

// GetCategoriesBySalePointID handles GET /api/v1/categories/sale-point/:sale_point_id
func (h *ProductHandler) GetCategoriesBySalePointID(c *gin.Context) {
	salePointID := c.Param("sale_point_id")

	categories, err := h.service.GetCategoriesBySalePointID(c.Request.Context(), salePointID)
	if err != nil {
		logger.Error("failed to get categories", "error", err, "sale_point_id", salePointID)
		response.Error(c, http.StatusInternalServerError, err, "Failed to get categories")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"categories": categories}, "")
}

// parseFilters parses query parameters into ProductFilters
func (h *ProductHandler) parseFilters(c *gin.Context) product.ProductFilters {
	filters := product.ProductFilters{}

	// Parse limit and offset
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	filters.Limit = limit
	filters.Offset = offset

	// Parse category filter
	if category := c.Query("category"); category != "" {
		filters.Category = &category
	}

	// Parse is_available filter
	if isAvailableStr := c.Query("is_available"); isAvailableStr != "" {
		isAvailable := isAvailableStr == "true"
		filters.IsAvailable = &isAvailable
	}

	// Parse is_addon filter
	if isAddonStr := c.Query("is_addon"); isAddonStr != "" {
		isAddon := isAddonStr == "true"
		filters.IsAddon = &isAddon
	}

	return filters
}
