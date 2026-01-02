package dto

import "github.com/emerarteaga/products-api/internal/domain/order"

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	SaleType          order.SaleType        `json:"sale_type" binding:"required,oneof=DELIVERY ON_SITE"`
	Products          []OrderProductRequest `json:"products" binding:"required,min=1,dive"`
	Note              *string               `json:"note" binding:"omitempty,max=500"`
	Customer          *CustomerRequest      `json:"customer" binding:"omitempty"`
	ShippingAddress   *string               `json:"shipping_address" binding:"omitempty,max=500"`
	TableNumber       *int                  `json:"table_number" binding:"omitempty,gte=1"`
	PaymentReceiptURL *string               `json:"payment_receipt_url" binding:"omitempty,url"`
	PaymentAccountID  *string               `json:"payment_account_id" binding:"omitempty"`
}

// OrderProductRequest represents a product in the request
type OrderProductRequest struct {
	ID          string  `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required,min=1,max=200"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Observation *string `json:"observation" binding:"omitempty,max=500"`
	Price       int64   `json:"price" binding:"required,gte=0"`
	Quantity    int     `json:"quantity" binding:"required,gte=1"`
}

// CustomerRequest represents customer information in the request
type CustomerRequest struct {
	Identification string       `json:"identification" binding:"required"`
	IDType         order.IDType `json:"id_type" binding:"required,oneof=CC CE PASSPORT NIT"`
	Name           string       `json:"name" binding:"required,min=2,max=200"`
	Phone          string       `json:"phone" binding:"required,min=7,max=20"`
}

// ToCreateInput converts DTO to service input
func (r *CreateOrderRequest) ToCreateInput() order.CreateInput {
	// Convert products
	products := make([]order.OrderProduct, len(r.Products))
	for i, p := range r.Products {
		products[i] = order.OrderProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Observation: p.Observation,
			Price:       p.Price,
			Quantity:    p.Quantity,
		}
	}

	// Convert customer if present
	var customer *order.Customer
	if r.Customer != nil {
		customer = &order.Customer{
			Identification: r.Customer.Identification,
			IDType:         r.Customer.IDType,
			Name:           r.Customer.Name,
			Phone:          r.Customer.Phone,
		}
	}

	return order.CreateInput{
		SaleType:          r.SaleType,
		Products:          products,
		Note:              r.Note,
		Customer:          customer,
		ShippingAddress:   r.ShippingAddress,
		TableNumber:       r.TableNumber,
		PaymentReceiptURL: r.PaymentReceiptURL,
		PaymentAccountID:  r.PaymentAccountID,
	}
}

// OrderCreatedResponse represents the response after creating an order
type OrderCreatedResponse struct {
	ID        string            `json:"id"`
	Code      string            `json:"code"`
	Status    order.OrderStatus `json:"status"`
	SaleType  order.SaleType    `json:"sale_type"`
	Total     int64             `json:"total"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

// ToCreatedResponse converts order to created response
func ToCreatedResponse(o *order.Order) OrderCreatedResponse {
	return OrderCreatedResponse{
		ID:        o.ID,
		Code:      o.Code,
		Status:    o.Status,
		SaleType:  o.SaleType,
		Total:     o.Total,
		CreatedAt: o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: o.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// OrderTrackResponse represents the public tracking response
type OrderTrackResponse struct {
	Code         string            `json:"code"`
	Status       order.OrderStatus `json:"status"`
	CustomerName string            `json:"customer_name"`
	UpdatedAt    string            `json:"updated_at"`
}

// ToTrackResponse converts order to tracking response (public, limited data)
func ToTrackResponse(o *order.Order) OrderTrackResponse {
	customerName := "Guest"
	if o.Customer != nil {
		customerName = o.Customer.Name
	}

	return OrderTrackResponse{
		Code:         o.Code,
		Status:       o.Status,
		CustomerName: customerName,
		UpdatedAt:    o.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ===================================
// STAGE 3: PARTIAL UPDATE (PATCH)
// ===================================

// PartialUpdateOrderRequest represents the request for partial update
type PartialUpdateOrderRequest struct {
	Code              string             `json:"code" binding:"required"`
	Status            *order.OrderStatus `json:"status" binding:"omitempty,oneof=CREATED VERIFIED IN_PROGRESS OUT_FOR_DELIVERY DELIVERED CANCELLED"`
	Note              *string            `json:"note" binding:"omitempty,max=500"`
	PaymentReceiptURL *string            `json:"payment_receipt_url" binding:"omitempty,url"`
	PaymentAccountID  *string            `json:"payment_account_id" binding:"omitempty"`
	// Products explicitly NOT allowed in PATCH
}

// ToPartialUpdateInput converts DTO to service input
func (r *PartialUpdateOrderRequest) ToPartialUpdateInput() order.PartialUpdateInput {
	return order.PartialUpdateInput{
		Status:            r.Status,
		Note:              r.Note,
		PaymentReceiptURL: r.PaymentReceiptURL,
		PaymentAccountID:  r.PaymentAccountID,
	}
}

// ===================================
// STAGE 4: MODIFY (PUT)
// ===================================

// ModifyOrderRequest represents the request to modify an order
type ModifyOrderRequest struct {
	Code            string                `json:"code" binding:"required"`
	Products        []OrderProductRequest `json:"products" binding:"omitempty,dive"`
	ShippingAddress *string               `json:"shipping_address" binding:"omitempty,max=500"`
	Customer        *CustomerRequest      `json:"customer" binding:"omitempty"`
	Note            *string               `json:"note" binding:"omitempty,max=500"`
}

// ToModifyInput converts DTO to service input
func (r *ModifyOrderRequest) ToModifyInput() order.ModifyInput {
	// Convert products if provided
	var products []order.OrderProduct
	if len(r.Products) > 0 {
		products = make([]order.OrderProduct, len(r.Products))
		for i, p := range r.Products {
			products[i] = order.OrderProduct{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Observation: p.Observation,
				Price:       p.Price,
				Quantity:    p.Quantity,
			}
		}
	}

	// Convert customer if present
	var customer *order.Customer
	if r.Customer != nil {
		customer = &order.Customer{
			Identification: r.Customer.Identification,
			IDType:         r.Customer.IDType,
			Name:           r.Customer.Name,
			Phone:          r.Customer.Phone,
		}
	}

	return order.ModifyInput{
		Products:        products,
		ShippingAddress: r.ShippingAddress,
		Customer:        customer,
		Note:            r.Note,
	}
}

// ===================================
// COMMON RESPONSES
// ===================================

// OrderResponse represents a complete order response
type OrderResponse struct {
	ID                string                 `json:"id"`
	Code              string                 `json:"code"`
	Status            order.OrderStatus      `json:"status"`
	SaleType          order.SaleType         `json:"sale_type"`
	Products          []OrderProductResponse `json:"products"`
	Total             int64                  `json:"total"`
	Note              *string                `json:"note,omitempty"`
	Customer          *CustomerResponse      `json:"customer,omitempty"`
	ShippingAddress   *string                `json:"shipping_address,omitempty"`
	TableNumber       *int                   `json:"table_number,omitempty"`
	PaymentReceiptURL *string                `json:"payment_receipt_url,omitempty"`
	PaymentAccountID  *string                `json:"payment_account_id,omitempty"`
	CreatedAt         string                 `json:"created_at"`
	UpdatedAt         string                 `json:"updated_at"`
}

// OrderProductResponse represents a product in the response
type OrderProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Observation *string `json:"observation,omitempty"`
	Price       int64   `json:"price"`
	Quantity    int     `json:"quantity"`
}

// CustomerResponse represents customer information in the response
type CustomerResponse struct {
	Identification string       `json:"identification"`
	IDType         order.IDType `json:"id_type"`
	Name           string       `json:"name"`
	Phone          string       `json:"phone"`
}

// ToOrderResponse converts order to full response
func ToOrderResponse(o *order.Order) OrderResponse {
	// Convert products
	products := make([]OrderProductResponse, len(o.Products))
	for i, p := range o.Products {
		products[i] = OrderProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Observation: p.Observation,
			Price:       p.Price,
			Quantity:    p.Quantity,
		}
	}

	// Convert customer if present
	var customer *CustomerResponse
	if o.Customer != nil {
		customer = &CustomerResponse{
			Identification: o.Customer.Identification,
			IDType:         o.Customer.IDType,
			Name:           o.Customer.Name,
			Phone:          o.Customer.Phone,
		}
	}

	return OrderResponse{
		ID:                o.ID,
		Code:              o.Code,
		Status:            o.Status,
		SaleType:          o.SaleType,
		Products:          products,
		Total:             o.Total,
		Note:              o.Note,
		Customer:          customer,
		ShippingAddress:   o.ShippingAddress,
		TableNumber:       o.TableNumber,
		PaymentReceiptURL: o.PaymentReceiptURL,
		PaymentAccountID:  o.PaymentAccountID,
		CreatedAt:         o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         o.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ===================================
// STAGE 5: FILTERS AND ANALYTICS
// ===================================

// OrderMetricsResponse represents the metrics response
type OrderMetricsResponse struct {
	Metrics     MetricsData                 `json:"metrics"`
	TopProducts []order.ProductSalesSummary `json:"top_products"`
}

// MetricsData represents aggregated metrics
type MetricsData struct {
	TotalSales     int64                     `json:"total_sales"`
	AvgTicket      int64                     `json:"avg_ticket"`
	OrdersByStatus map[order.OrderStatus]int `json:"orders_by_status"`
}

// ToMetricsResponse converts order metrics to response
func ToMetricsResponse(m *order.OrderMetrics) OrderMetricsResponse {
	return OrderMetricsResponse{
		Metrics: MetricsData{
			TotalSales:     m.TotalSales,
			AvgTicket:      m.AvgTicket,
			OrdersByStatus: m.OrdersByStatus,
		},
		TopProducts: m.TopProducts,
	}
}
