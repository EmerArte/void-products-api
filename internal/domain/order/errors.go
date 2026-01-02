package order

import "errors"

// General validation errors
var (
	ErrOrderNotFound          = errors.New("order not found")
	ErrInvalidOrderID         = errors.New("invalid order ID")
	ErrInvalidOrderCode       = errors.New("invalid order code")
	ErrOrderCodeAlreadyExists = errors.New("order code already exists")
)

// Product validation errors
var (
	ErrNoProducts                = errors.New("order must contain at least one product")
	ErrInvalidProductID          = errors.New("product ID is required")
	ErrInvalidProductName        = errors.New("product name is required")
	ErrInvalidProductQuantity    = errors.New("product quantity must be greater than 0")
	ErrInvalidProductPrice       = errors.New("product price must be greater than or equal to 0")
	ErrDuplicateProduct          = errors.New("duplicate product in order")
	ErrProductsNotAllowedInPatch = errors.New("products cannot be updated via PATCH, use PUT instead")
)

// Customer validation errors
var (
	ErrCustomerRequiredForDelivery    = errors.New("customer information is required for delivery orders")
	ErrCustomerNameRequired           = errors.New("customer name is required")
	ErrCustomerPhoneRequired          = errors.New("customer phone is required")
	ErrCustomerIdentificationRequired = errors.New("customer identification is required")
	ErrInvalidIDType                  = errors.New("invalid customer ID type")
)

// Shipping and table errors
var (
	ErrShippingAddressRequired            = errors.New("shipping address is required for delivery orders")
	ErrShippingAddressNotAllowedForOnSite = errors.New("shipping address is not allowed for on-site orders")
	ErrTableNumberRequiredForOnSite       = errors.New("table number is required for on-site orders")
	ErrTableNumberNotAllowedForDelivery   = errors.New("table number is not allowed for delivery orders")
	ErrInvalidTableNumber                 = errors.New("invalid table number")
)

// Sale type and status errors
var (
	ErrInvalidSaleType         = errors.New("invalid sale type")
	ErrInvalidStatus           = errors.New("invalid order status")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrOrderCannotBeModified   = errors.New("order cannot be modified in current status")
	ErrOrderAlreadyCancelled   = errors.New("order is already cancelled")
	ErrOrderAlreadyDelivered   = errors.New("order is already delivered")
)

// Payment errors
var (
	ErrInvalidPaymentAccountID  = errors.New("invalid payment account ID")
	ErrInvalidPaymentReceiptURL = errors.New("invalid payment receipt URL")
)

// Total validation errors
var (
	ErrTotalMismatch = errors.New("provided total does not match calculated total")
	ErrInvalidTotal  = errors.New("invalid total amount")
)
