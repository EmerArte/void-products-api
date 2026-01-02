package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/emerarteaga/products-api/internal/domain/order"
	"github.com/emerarteaga/products-api/internal/dto"
	"github.com/emerarteaga/products-api/internal/infra/logger"
	"github.com/emerarteaga/products-api/internal/response"
	"github.com/gin-gonic/gin"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	service *order.Service
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(service *order.Service) *OrderHandler {
	return &OrderHandler{service: service}
}

// Create handles POST /api/v1/orders
func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
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

	o, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		// Map domain errors to HTTP status codes
		statusCode := h.mapErrorToStatusCode(err)
		logger.Error("failed to create order", "error", err)
		response.Error(c, statusCode, err, "Failed to create order")
		return
	}

	logger.Info("order created", "order_id", o.ID, "code", o.Code, "sale_type", o.SaleType)
	response.Success(c, http.StatusCreated, dto.ToCreatedResponse(o), "Order created successfully")
}

// Track handles GET /api/v1/orders/track/:code
func (h *OrderHandler) Track(c *gin.Context) {
	code := c.Param("code")

	o, err := h.service.GetByCode(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, order.ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, err, "Order not found")
			return
		}
		logger.Error("failed to track order", "error", err, "code", code)
		response.Error(c, http.StatusInternalServerError, err, "Failed to track order")
		return
	}

	// Return limited public information
	response.Success(c, http.StatusOK, dto.ToTrackResponse(o), "")
}

// PartialUpdate handles PATCH /api/v1/orders
func (h *OrderHandler) PartialUpdate(c *gin.Context) {
	var req dto.PartialUpdateOrderRequest
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

	// Check if products are being sent (not allowed in PATCH)
	var rawData map[string]interface{}
	if err := c.ShouldBindJSON(&rawData); err == nil {
		if _, hasProducts := rawData["products"]; hasProducts {
			logger.Warn("products not allowed in PATCH", "code", req.Code)
			response.Error(c, http.StatusBadRequest, order.ErrProductsNotAllowedInPatch, "Products cannot be updated via PATCH, use PUT instead")
			return
		}
	}

	// Convert DTO to service input
	input := req.ToPartialUpdateInput()

	o, err := h.service.PartialUpdate(c.Request.Context(), req.Code, input)
	if err != nil {
		statusCode := h.mapErrorToStatusCode(err)
		logger.Error("failed to partial update order", "error", err, "code", req.Code)
		response.Error(c, statusCode, err, "Failed to update order")
		return
	}

	logger.Info("order partially updated", "order_id", o.ID, "code", o.Code, "status", o.Status)
	response.Success(c, http.StatusOK, dto.ToOrderResponse(o), "Order updated successfully")
}

// Modify handles PUT /api/v1/orders
func (h *OrderHandler) Modify(c *gin.Context) {
	var req dto.ModifyOrderRequest
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
	input := req.ToModifyInput()

	o, err := h.service.Modify(c.Request.Context(), req.Code, input)
	if err != nil {
		statusCode := h.mapErrorToStatusCode(err)
		logger.Error("failed to modify order", "error", err, "code", req.Code)
		response.Error(c, statusCode, err, "Failed to modify order")
		return
	}

	logger.Info("order modified", "order_id", o.ID, "code", o.Code, "status", o.Status)
	response.Success(c, http.StatusOK, dto.ToOrderResponse(o), "Order modified successfully")
}

// GetMetrics handles GET /api/v1/orders/metrics
func (h *OrderHandler) GetMetrics(c *gin.Context) {
	filters := h.parseFilters(c)

	metrics, err := h.service.GetMetrics(c.Request.Context(), filters)
	if err != nil {
		logger.Error("failed to get metrics", "error", err)
		response.Error(c, http.StatusInternalServerError, err, "Failed to get metrics")
		return
	}

	response.Success(c, http.StatusOK, dto.ToMetricsResponse(metrics), "")
}

// GetAll handles GET /api/v1/orders
func (h *OrderHandler) GetAll(c *gin.Context) {
	filters := h.parseFilters(c)

	orders, total, err := h.service.GetAll(c.Request.Context(), filters)
	if err != nil {
		logger.Error("failed to get orders", "error", err)
		response.Error(c, http.StatusInternalServerError, err, "Failed to get orders")
		return
	}

	// Convert to response
	orderResponses := make([]dto.OrderResponse, len(orders))
	for i, o := range orders {
		orderResponses[i] = dto.ToOrderResponse(o)
	}

	response.Paginated(c, http.StatusOK, orderResponses, total, filters.Limit, filters.Offset)
}

// GetByCode handles GET /api/v1/orders/:code (internal/admin use)
func (h *OrderHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	o, err := h.service.GetByCode(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, order.ErrOrderNotFound) {
			response.Error(c, http.StatusNotFound, err, "Order not found")
			return
		}
		logger.Error("failed to get order", "error", err, "code", code)
		response.Error(c, http.StatusInternalServerError, err, "Failed to get order")
		return
	}

	response.Success(c, http.StatusOK, dto.ToOrderResponse(o), "")
}

// parseFilters parses query parameters into OrderFilters
func (h *OrderHandler) parseFilters(c *gin.Context) order.OrderFilters {
	filters := order.OrderFilters{}

	// Parse pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	filters.Limit = limit
	filters.Offset = offset

	// Parse date filters
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filters.DateFrom = &dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		filters.DateTo = &dateTo
	}

	// Parse status filter
	if statusStr := c.Query("status"); statusStr != "" {
		status := order.OrderStatus(statusStr)
		filters.Status = &status
	}

	// Parse sale type filter
	if saleTypeStr := c.Query("sale_type"); saleTypeStr != "" {
		saleType := order.SaleType(saleTypeStr)
		filters.SaleType = &saleType
	}

	// Parse product filters
	if productID := c.Query("product_id"); productID != "" {
		filters.ProductID = &productID
	}
	if productName := c.Query("product_name"); productName != "" {
		filters.ProductName = &productName
	}

	// Parse total filters
	if minTotalStr := c.Query("min_total"); minTotalStr != "" {
		if minTotal, err := strconv.ParseInt(minTotalStr, 10, 64); err == nil {
			filters.MinTotal = &minTotal
		}
	}
	if maxTotalStr := c.Query("max_total"); maxTotalStr != "" {
		if maxTotal, err := strconv.ParseInt(maxTotalStr, 10, 64); err == nil {
			filters.MaxTotal = &maxTotal
		}
	}

	return filters
}

// mapErrorToStatusCode maps domain errors to HTTP status codes
func (h *OrderHandler) mapErrorToStatusCode(err error) int {
	switch {
	case errors.Is(err, order.ErrOrderNotFound):
		return http.StatusNotFound
	case errors.Is(err, order.ErrInvalidStatusTransition):
		return http.StatusConflict
	case errors.Is(err, order.ErrOrderCannotBeModified):
		return http.StatusConflict
	case errors.Is(err, order.ErrOrderCodeAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, order.ErrNoProducts),
		errors.Is(err, order.ErrInvalidProductID),
		errors.Is(err, order.ErrInvalidProductName),
		errors.Is(err, order.ErrInvalidProductQuantity),
		errors.Is(err, order.ErrInvalidProductPrice),
		errors.Is(err, order.ErrDuplicateProduct),
		errors.Is(err, order.ErrCustomerRequiredForDelivery),
		errors.Is(err, order.ErrCustomerNameRequired),
		errors.Is(err, order.ErrCustomerPhoneRequired),
		errors.Is(err, order.ErrShippingAddressRequired),
		errors.Is(err, order.ErrTableNumberRequiredForOnSite),
		errors.Is(err, order.ErrInvalidSaleType),
		errors.Is(err, order.ErrInvalidStatus),
		errors.Is(err, order.ErrTotalMismatch):
		return http.StatusUnprocessableEntity
	case errors.Is(err, order.ErrProductsNotAllowedInPatch):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
