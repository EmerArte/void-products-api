package product

import (
	"context"
	"fmt"
)

// Service handles business logic for products
type Service struct {
	repo Repository
}

// NewService creates a new product service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateInput represents input for creating a product
type CreateInput struct {
	CompanyID        string
	SalePointID      string
	Name             string
	Description      string
	Category         string
	Photos           []string
	PriceVariations  []PriceVariation
	AvailableAddons  []Addon
	IsAddon          bool
	IsAvailable      bool
	IsUnlimitedStock bool
	Stock            *int
}

// UpdateInput represents input for updating a product
type UpdateInput struct {
	Name             *string
	Description      *string
	Category         *string
	Photos           *[]string
	PriceVariations  *[]PriceVariation
	AvailableAddons  *[]Addon
	IsAddon          *bool
	IsAvailable      *bool
	IsUnlimitedStock *bool
	Stock            **int // Pointer to pointer to allow setting to nil
}

// Create creates a new product
func (s *Service) Create(ctx context.Context, input CreateInput) (*Product, error) {
	// Create new product
	p := NewProduct(input.CompanyID, input.SalePointID, input.Name, input.Category, input.Description)

	// Set additional fields
	if len(input.Photos) > 0 {
		p.Photos = input.Photos
	}

	p.PriceVariations = input.PriceVariations
	p.AvailableAddons = input.AvailableAddons
	p.IsAddon = input.IsAddon
	p.IsAvailable = input.IsAvailable
	p.IsUnlimitedStock = input.IsUnlimitedStock
	p.Stock = input.Stock

	// Validate business rules
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Save to repository
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return p, nil
}

// GetByID retrieves a product by ID
func (s *Service) GetByID(ctx context.Context, id string) (*Product, error) {
	if id == "" {
		return nil, fmt.Errorf("product ID is required")
	}

	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetByCompanyID retrieves products by company ID with filters
func (s *Service) GetByCompanyID(ctx context.Context, companyID string, filters ProductFilters) ([]*Product, int64, error) {
	if companyID == "" {
		return nil, 0, fmt.Errorf("company ID is required")
	}

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
	total, err := s.repo.CountByCompanyID(ctx, companyID, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	products, err := s.repo.FindByCompanyID(ctx, companyID, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products: %w", err)
	}

	return products, total, nil
}

// GetBySalePointID retrieves products by sale point ID with filters
func (s *Service) GetBySalePointID(ctx context.Context, salePointID string, filters ProductFilters) ([]*Product, int64, error) {
	if salePointID == "" {
		return nil, 0, fmt.Errorf("sale point ID is required")
	}

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
	total, err := s.repo.CountBySalePointID(ctx, salePointID, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	products, err := s.repo.FindBySalePointID(ctx, salePointID, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products: %w", err)
	}

	return products, total, nil
}

// Update updates a product
func (s *Service) Update(ctx context.Context, id string, input UpdateInput) (*Product, error) {
	if id == "" {
		return nil, fmt.Errorf("product ID is required")
	}

	// Find existing product
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Category != nil {
		product.Category = *input.Category
	}
	if input.Photos != nil {
		product.Photos = *input.Photos
	}
	if input.PriceVariations != nil {
		product.PriceVariations = *input.PriceVariations
	}
	if input.AvailableAddons != nil {
		product.AvailableAddons = *input.AvailableAddons
	}
	if input.IsAddon != nil {
		product.IsAddon = *input.IsAddon
	}
	if input.IsAvailable != nil {
		product.IsAvailable = *input.IsAvailable
	}
	if input.IsUnlimitedStock != nil {
		product.IsUnlimitedStock = *input.IsUnlimitedStock
	}
	if input.Stock != nil {
		product.Stock = *input.Stock
	}

	// Validate business rules
	if err := product.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Update in repository
	if err := s.repo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

// Delete deletes a product
func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("product ID is required")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// GetCategoriesByCompanyID retrieves categories for a company
func (s *Service) GetCategoriesByCompanyID(ctx context.Context, companyID string) ([]string, error) {
	if companyID == "" {
		return nil, fmt.Errorf("company ID is required")
	}

	categories, err := s.repo.FindCategoriesByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

// GetCategoriesBySalePointID retrieves categories for a sale point
func (s *Service) GetCategoriesBySalePointID(ctx context.Context, salePointID string) ([]string, error) {
	if salePointID == "" {
		return nil, fmt.Errorf("sale point ID is required")
	}

	categories, err := s.repo.FindCategoriesBySalePointID(ctx, salePointID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}
