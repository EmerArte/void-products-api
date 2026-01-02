package order

import (
	"context"
	"fmt"
)

// Service handles business logic for orders
type Service struct {
	repo Repository
}

// NewService creates a new order service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateInput represents input for creating an order
type CreateInput struct {
	SaleType          SaleType
	Products          []OrderProduct
	Note              *string
	Customer          *Customer
	ShippingAddress   *string
	TableNumber       *int
	PaymentReceiptURL *string
	PaymentAccountID  *string
}

// PartialUpdateInput represents input for partial update (PATCH)
type PartialUpdateInput struct {
	Status            *OrderStatus
	Note              *string
	PaymentReceiptURL *string
	PaymentAccountID  *string
}

// ModifyInput represents input for full modification (PUT)
type ModifyInput struct {
	Products        []OrderProduct
	ShippingAddress *string
	Customer        *Customer
	Note            *string
}

// Create creates a new order
func (s *Service) Create(ctx context.Context, input CreateInput) (*Order, error) {
	// Create new order
	o := NewOrder(input.SaleType, input.Products)

	// Set optional fields
	o.Note = input.Note
	o.Customer = input.Customer
	o.ShippingAddress = input.ShippingAddress
	o.TableNumber = input.TableNumber
	o.PaymentReceiptURL = input.PaymentReceiptURL
	o.PaymentAccountID = input.PaymentAccountID

	// Validate business rules
	if err := o.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Check if code already exists (very unlikely but possible)
	exists, err := s.repo.ExistsByCode(ctx, o.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to check code existence: %w", err)
	}
	if exists {
		// Regenerate code and try again
		o.Code = generateOrderCode()
	}

	// Save to repository
	if err := s.repo.Create(ctx, o); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return o, nil
}

// GetByID retrieves an order by ID
func (s *Service) GetByID(ctx context.Context, id string) (*Order, error) {
	if id == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetByCode retrieves an order by tracking code
func (s *Service) GetByCode(ctx context.Context, code string) (*Order, error) {
	if code == "" {
		return nil, ErrInvalidOrderCode
	}

	order, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// PartialUpdate updates an order partially (PATCH - no product changes)
func (s *Service) PartialUpdate(ctx context.Context, code string, input PartialUpdateInput) (*Order, error) {
	if code == "" {
		return nil, ErrInvalidOrderCode
	}

	// Find existing order
	order, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// Update allowed fields
	if input.Status != nil {
		if err := order.UpdateStatus(*input.Status); err != nil {
			return nil, err
		}
	}

	if input.Note != nil {
		order.Note = input.Note
	}

	if input.PaymentReceiptURL != nil {
		order.PaymentReceiptURL = input.PaymentReceiptURL
	}

	if input.PaymentAccountID != nil {
		order.PaymentAccountID = input.PaymentAccountID
	}

	// Update in repository
	if err := s.repo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return order, nil
}

// Modify modifies an order (PUT - products allowed, auto VERIFIED)
func (s *Service) Modify(ctx context.Context, code string, input ModifyInput) (*Order, error) {
	if code == "" {
		return nil, ErrInvalidOrderCode
	}

	// Find existing order
	order, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// Check if order can be modified
	if !order.CanBeModified() {
		return nil, ErrOrderCannotBeModified
	}

	// Update products if provided
	if len(input.Products) > 0 {
		if err := order.UpdateProducts(input.Products); err != nil {
			return nil, err
		}
	}

	// Update other fields
	if input.ShippingAddress != nil {
		order.ShippingAddress = input.ShippingAddress
	}

	if input.Customer != nil {
		order.Customer = input.Customer
	}

	if input.Note != nil {
		order.Note = input.Note
	}

	// Validate updated order
	if err := order.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Update in repository
	if err := s.repo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return order, nil
}

// GetAll retrieves all orders with filters
func (s *Service) GetAll(ctx context.Context, filters OrderFilters) ([]*Order, int64, error) {
	// Set default pagination
	if filters.Limit <= 0 {
		filters.Limit = 50
	}
	if filters.Limit > 100 {
		filters.Limit = 100 // Maximum limit
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	// Get total count with same filters
	total, err := s.repo.Count(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Get orders
	orders, err := s.repo.FindAll(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get orders: %w", err)
	}

	return orders, total, nil
}

// GetMetrics retrieves aggregated order metrics
func (s *Service) GetMetrics(ctx context.Context, filters OrderFilters) (*OrderMetrics, error) {
	metrics, err := s.repo.GetMetrics(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	return metrics, nil
}
