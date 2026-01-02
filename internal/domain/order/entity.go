package order

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	StatusCreated        OrderStatus = "CREATED"
	StatusVerified       OrderStatus = "VERIFIED"
	StatusInProgress     OrderStatus = "IN_PROGRESS"
	StatusOutForDelivery OrderStatus = "OUT_FOR_DELIVERY"
	StatusDelivered      OrderStatus = "DELIVERED"
	StatusCancelled      OrderStatus = "CANCELLED"
)

// SaleType represents the type of sale
type SaleType string

const (
	SaleTypeDelivery SaleType = "DELIVERY"
	SaleTypeOnSite   SaleType = "ON_SITE"
)

// IDType represents the type of identification
type IDType string

const (
	IDTypeCC       IDType = "CC" // Cédula de Ciudadanía
	IDTypeCE       IDType = "CE" // Cédula de Extranjería
	IDTypePassport IDType = "PASSPORT"
	IDTypeNIT      IDType = "NIT" // Tax ID for companies
)

// Order represents a sales order
type Order struct {
	ID                string         `json:"id" bson:"_id"`
	Code              string         `json:"code" bson:"code"`
	Status            OrderStatus    `json:"status" bson:"status"`
	SaleType          SaleType       `json:"sale_type" bson:"sale_type"`
	Products          []OrderProduct `json:"products" bson:"products"`
	Total             int64          `json:"total" bson:"total"` // In cents
	Note              *string        `json:"note,omitempty" bson:"note,omitempty"`
	Customer          *Customer      `json:"customer,omitempty" bson:"customer,omitempty"`
	ShippingAddress   *string        `json:"shipping_address,omitempty" bson:"shipping_address,omitempty"`
	TableNumber       *int           `json:"table_number,omitempty" bson:"table_number,omitempty"`
	PaymentReceiptURL *string        `json:"payment_receipt_url,omitempty" bson:"payment_receipt_url,omitempty"`
	PaymentAccountID  *string        `json:"payment_account_id,omitempty" bson:"payment_account_id,omitempty"`
	CreatedAt         time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at" bson:"updated_at"`
}

// OrderProduct represents a product in an order
type OrderProduct struct {
	ID          string  `json:"id" bson:"id"`
	Name        string  `json:"name" bson:"name"`
	Description *string `json:"description,omitempty" bson:"description,omitempty"`
	Observation *string `json:"observation,omitempty" bson:"observation,omitempty"`
	Price       int64   `json:"price" bson:"price"` // In cents
	Quantity    int     `json:"quantity" bson:"quantity"`
}

// Customer represents customer information
type Customer struct {
	Identification string `json:"identification" bson:"identification"`
	IDType         IDType `json:"id_type" bson:"id_type"`
	Name           string `json:"name" bson:"name"`
	Phone          string `json:"phone" bson:"phone"`
}

// NewOrder creates a new order
func NewOrder(saleType SaleType, products []OrderProduct) *Order {
	now := time.Now()
	order := &Order{
		ID:        uuid.New().String(),
		Code:      generateOrderCode(),
		Status:    StatusCreated,
		SaleType:  saleType,
		Products:  products,
		CreatedAt: now,
		UpdatedAt: now,
	}
	order.CalculateTotal()
	return order
}

// CalculateTotal calculates the total amount from products
func (o *Order) CalculateTotal() {
	total := int64(0)
	for _, product := range o.Products {
		total += product.Price * int64(product.Quantity)
	}
	o.Total = total
}

// Validate validates the order business rules
func (o *Order) Validate() error {
	// Validate products
	if len(o.Products) == 0 {
		return ErrNoProducts
	}

	for i, product := range o.Products {
		if product.ID == "" {
			return ErrInvalidProductID
		}
		if product.Name == "" {
			return ErrInvalidProductName
		}
		if product.Quantity <= 0 {
			return ErrInvalidProductQuantity
		}
		if product.Price < 0 {
			return ErrInvalidProductPrice
		}
		// Check for duplicate products
		for j := i + 1; j < len(o.Products); j++ {
			if o.Products[j].ID == product.ID {
				return ErrDuplicateProduct
			}
		}
	}

	// Validate sale type specific rules
	switch o.SaleType {
	case SaleTypeDelivery:
		if o.Customer == nil {
			return ErrCustomerRequiredForDelivery
		}
		if o.Customer.Name == "" {
			return ErrCustomerNameRequired
		}
		if o.Customer.Phone == "" {
			return ErrCustomerPhoneRequired
		}
		if o.ShippingAddress == nil || *o.ShippingAddress == "" {
			return ErrShippingAddressRequired
		}
		if o.TableNumber != nil {
			return ErrTableNumberNotAllowedForDelivery
		}
	case SaleTypeOnSite:
		if o.TableNumber == nil {
			return ErrTableNumberRequiredForOnSite
		}
		if *o.TableNumber <= 0 {
			return ErrInvalidTableNumber
		}
		if o.ShippingAddress != nil {
			return ErrShippingAddressNotAllowedForOnSite
		}
	default:
		return ErrInvalidSaleType
	}

	// Validate status
	if !o.IsValidStatus(o.Status) {
		return ErrInvalidStatus
	}

	return nil
}

// IsValidStatus checks if the status is valid
func (o *Order) IsValidStatus(status OrderStatus) bool {
	validStatuses := []OrderStatus{
		StatusCreated,
		StatusVerified,
		StatusInProgress,
		StatusOutForDelivery,
		StatusDelivered,
		StatusCancelled,
	}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// CanTransitionTo checks if the order can transition to the given status
func (o *Order) CanTransitionTo(newStatus OrderStatus) bool {
	// Define valid transitions
	validTransitions := map[OrderStatus][]OrderStatus{
		StatusCreated: {
			StatusVerified,
			StatusInProgress,
			StatusCancelled,
		},
		StatusVerified: {
			StatusInProgress,
			StatusCancelled,
		},
		StatusInProgress: {
			StatusOutForDelivery,
			StatusDelivered, // Direct delivery for ON_SITE
			StatusCancelled,
		},
		StatusOutForDelivery: {
			StatusDelivered,
			StatusCancelled,
		},
		StatusDelivered: {
			// Terminal state - no transitions
		},
		StatusCancelled: {
			// Terminal state - no transitions
		},
	}

	allowedStatuses, exists := validTransitions[o.Status]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}
	return false
}

// CanBeModified checks if the order can be modified (products, address, etc.)
func (o *Order) CanBeModified() bool {
	// Cannot modify if already out for delivery, delivered, or cancelled
	immutableStatuses := []OrderStatus{
		StatusOutForDelivery,
		StatusDelivered,
		StatusCancelled,
	}
	for _, status := range immutableStatuses {
		if o.Status == status {
			return false
		}
	}
	return true
}

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(newStatus OrderStatus) error {
	if !o.IsValidStatus(newStatus) {
		return ErrInvalidStatus
	}

	if !o.CanTransitionTo(newStatus) {
		return ErrInvalidStatusTransition
	}

	o.Status = newStatus
	o.UpdatedAt = time.Now()
	return nil
}

// UpdateProducts updates the order products and recalculates total
func (o *Order) UpdateProducts(products []OrderProduct) error {
	if !o.CanBeModified() {
		return ErrOrderCannotBeModified
	}

	if len(products) == 0 {
		return ErrNoProducts
	}

	o.Products = products
	o.CalculateTotal()
	o.Status = StatusVerified
	o.UpdatedAt = time.Now()
	return nil
}

// generateOrderCode generates a unique order code
func generateOrderCode() string {
	timestamp := time.Now().UnixNano()
	randomPart := uuid.New().String()[:8]
	return fmt.Sprintf("ORD-%d-%s", timestamp, randomPart)
}
