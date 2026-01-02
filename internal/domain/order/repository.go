package order

import "context"

// OrderFilters represents filters for querying orders
type OrderFilters struct {
	DateFrom    *string
	DateTo      *string
	Status      *OrderStatus
	SaleType    *SaleType
	ProductID   *string
	ProductName *string
	MinTotal    *int64
	MaxTotal    *int64
	Limit       int
	Offset      int
}

// OrderMetrics represents aggregated order metrics
type OrderMetrics struct {
	TotalSales     int64                 `json:"total_sales"`
	AvgTicket      int64                 `json:"avg_ticket"`
	OrdersByStatus map[OrderStatus]int   `json:"orders_by_status"`
	TopProducts    []ProductSalesSummary `json:"top_products"`
}

// ProductSalesSummary represents product sales aggregation
type ProductSalesSummary struct {
	ProductID     string `json:"product_id"`
	Name          string `json:"name"`
	TotalQuantity int    `json:"total_quantity"`
	TotalRevenue  int64  `json:"total_revenue"`
}

// Repository defines the contract for order data operations
type Repository interface {
	// Create creates a new order
	Create(ctx context.Context, order *Order) error

	// FindByID retrieves an order by its ID
	FindByID(ctx context.Context, id string) (*Order, error)

	// FindByCode retrieves an order by its tracking code
	FindByCode(ctx context.Context, code string) (*Order, error)

	// Update updates an existing order
	Update(ctx context.Context, order *Order) error

	// FindAll retrieves all orders with optional filters
	FindAll(ctx context.Context, filters OrderFilters) ([]*Order, error)

	// Count returns the total number of orders matching filters
	Count(ctx context.Context, filters OrderFilters) (int64, error)

	// ExistsByCode checks if an order exists with the given code
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// GetMetrics returns aggregated order metrics
	GetMetrics(ctx context.Context, filters OrderFilters) (*OrderMetrics, error)
}
