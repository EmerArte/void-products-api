package product

import "errors"

// Domain errors for Product entity
var (
	// General validation errors
	ErrInvalidCompanyID   = errors.New("company_id is required")
	ErrInvalidSalePointID = errors.New("sale_point_id is required")
	ErrInvalidName        = errors.New("product name is required")
	ErrInvalidCategory    = errors.New("category is required")

	// Stock errors
	ErrInvalidStock                  = errors.New("stock must be set when is_unlimited_stock is false")
	ErrStockMustBeNullForUnlimited   = errors.New("stock must be null when is_unlimited_stock is true")
	ErrNegativeStock                 = errors.New("stock cannot be negative")
	ErrInsufficientStock             = errors.New("insufficient stock available")
	ErrCannotUpdateStockForUnlimited = errors.New("cannot update stock for unlimited stock products")

	// Price variation errors
	ErrNoPriceVariations           = errors.New("at least one price variation is required")
	ErrInvalidPriceVariationType   = errors.New("price variation type is required")
	ErrNegativePrice               = errors.New("price cannot be negative")
	ErrDuplicatePriceVariationType = errors.New("duplicate price variation type")
	ErrInvalidMaxSelections        = errors.New("max_selections cannot be negative")
	ErrNoOptionsForMaxSelections   = errors.New("options must be provided when max_selections > 0")

	// Addon errors
	ErrInvalidAddonName   = errors.New("addon name is required")
	ErrNegativeAddonPrice = errors.New("addon price cannot be negative")
	ErrDuplicateAddon     = errors.New("addon already exists")

	// Not found error
	ErrProductNotFound = errors.New("product not found")
)
