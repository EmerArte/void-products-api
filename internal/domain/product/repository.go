package product

import "context"

// ProductFilters represents filters for querying products
type ProductFilters struct {
	CompanyID   *string
	SalePointID *string
	Category    *string
	IsAvailable *bool
	IsAddon     *bool
	Limit       int
	Offset      int
}

// Repository defines the contract for product data operations
type Repository interface {
	// Create creates a new product
	Create(ctx context.Context, product *Product) error

	// FindByID retrieves a product by its ID
	FindByID(ctx context.Context, id string) (*Product, error)

	// FindByCompanyID retrieves all products for a company with optional filters
	FindByCompanyID(ctx context.Context, companyID string, filters ProductFilters) ([]*Product, error)

	// FindBySalePointID retrieves all products for a sale point with optional filters
	FindBySalePointID(ctx context.Context, salePointID string, filters ProductFilters) ([]*Product, error)

	// FindAll retrieves all products with optional filters (deprecated, use FindByCompanyID or FindBySalePointID)
	FindAll(ctx context.Context, limit, offset int) ([]*Product, error)

	// Update updates an existing product
	Update(ctx context.Context, product *Product) error

	// Delete deletes a product by ID
	Delete(ctx context.Context, id string) error

	// FindCategories retrieves all unique categories
	FindCategoriesByCompanyID(ctx context.Context, companyID string) ([]string, error)
	FindCategoriesBySalePointID(ctx context.Context, salePointID string) ([]string, error)

	// Count returns the total number of products
	Count(ctx context.Context) (int64, error)

	// CountByCompanyID returns the total number of products for a company with filters
	CountByCompanyID(ctx context.Context, companyID string, filters ProductFilters) (int64, error)

	// CountBySalePointID returns the total number of products for a sale point with filters
	CountBySalePointID(ctx context.Context, salePointID string, filters ProductFilters) (int64, error)

	// Exists checks if a product exists
	Exists(ctx context.Context, id string) (bool, error)
}
