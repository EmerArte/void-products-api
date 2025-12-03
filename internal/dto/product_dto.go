package dto

import "github.com/emerarteaga/products-api/internal/domain/product"

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	CompanyID        string                  `json:"company_id" binding:"required"`
	SalePointID      string                  `json:"sale_point_id" binding:"required"`
	Name             string                  `json:"name" binding:"required,min=2,max=200"`
	Description      string                  `json:"description" binding:"max=1000"`
	Category         string                  `json:"category" binding:"required,min=2,max=100"`
	Photos           []string                `json:"photos"`
	PriceVariations  []PriceVariationRequest `json:"price_variations" binding:"required,min=1,dive"`
	AvailableAddons  []AddonRequest          `json:"available_addons" binding:"dive"`
	IsAddon          bool                    `json:"is_addon"`
	IsAvailable      bool                    `json:"is_available"`
	IsUnlimitedStock bool                    `json:"is_unlimited_stock"`
	Stock            *int                    `json:"stock" binding:"omitempty,gte=0"`
}

// PriceVariationRequest represents a price variation in the request
type PriceVariationRequest struct {
	Type           string                `json:"type" binding:"required"`
	Price          int64                 `json:"price" binding:"required,gte=0"`
	IncludedAddons IncludedAddonsRequest `json:"included_addons"`
}

// IncludedAddonsRequest represents included addons in the request
type IncludedAddonsRequest struct {
	MaxSelections int            `json:"max_selections" binding:"gte=0"`
	Options       []AddonRequest `json:"options" binding:"dive"`
}

// AddonRequest represents an addon in the request
type AddonRequest struct {
	ID          string   `json:"id"`
	Name        string   `json:"name" binding:"required"`
	Price       int64    `json:"price" binding:"gte=0"`
	Photos      []string `json:"photos"`
	IsAvailable bool     `json:"is_available"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name             *string                  `json:"name" binding:"omitempty,min=2,max=200"`
	Description      *string                  `json:"description" binding:"omitempty,max=1000"`
	Category         *string                  `json:"category" binding:"omitempty,min=2,max=100"`
	Photos           *[]string                `json:"photos"`
	PriceVariations  *[]PriceVariationRequest `json:"price_variations" binding:"omitempty,min=1,dive"`
	AvailableAddons  *[]AddonRequest          `json:"available_addons" binding:"omitempty,dive"`
	IsAddon          *bool                    `json:"is_addon"`
	IsAvailable      *bool                    `json:"is_available"`
	IsUnlimitedStock *bool                    `json:"is_unlimited_stock"`
	Stock            **int                    `json:"stock" binding:"omitempty"`
}

// ToCreateInput converts DTO to service input
func (r *CreateProductRequest) ToCreateInput() product.CreateInput {
	// Convert price variations
	priceVariations := make([]product.PriceVariation, len(r.PriceVariations))
	for i, pv := range r.PriceVariations {
		options := make([]product.Addon, len(pv.IncludedAddons.Options))
		for j, opt := range pv.IncludedAddons.Options {
			options[j] = product.Addon{
				ID:          opt.ID,
				Name:        opt.Name,
				Price:       opt.Price,
				Photos:      opt.Photos,
				IsAvailable: opt.IsAvailable,
			}
		}

		priceVariations[i] = product.PriceVariation{
			Type:  pv.Type,
			Price: pv.Price,
			IncludedAddons: product.IncludedAddons{
				MaxSelections: pv.IncludedAddons.MaxSelections,
				Options:       options,
			},
		}
	}

	// Convert available addons
	availableAddons := make([]product.Addon, len(r.AvailableAddons))
	for i, addon := range r.AvailableAddons {
		availableAddons[i] = product.Addon{
			ID:          addon.ID,
			Name:        addon.Name,
			Price:       addon.Price,
			Photos:      addon.Photos,
			IsAvailable: addon.IsAvailable,
		}
	}

	return product.CreateInput{
		CompanyID:        r.CompanyID,
		SalePointID:      r.SalePointID,
		Name:             r.Name,
		Description:      r.Description,
		Category:         r.Category,
		Photos:           r.Photos,
		PriceVariations:  priceVariations,
		AvailableAddons:  availableAddons,
		IsAddon:          r.IsAddon,
		IsAvailable:      r.IsAvailable,
		IsUnlimitedStock: r.IsUnlimitedStock,
		Stock:            r.Stock,
	}
}

// ToUpdateInput converts DTO to service input
func (r *UpdateProductRequest) ToUpdateInput() product.UpdateInput {
	input := product.UpdateInput{
		Name:             r.Name,
		Description:      r.Description,
		Category:         r.Category,
		Photos:           r.Photos,
		IsAddon:          r.IsAddon,
		IsAvailable:      r.IsAvailable,
		IsUnlimitedStock: r.IsUnlimitedStock,
		Stock:            r.Stock,
	}

	// Convert price variations if provided
	if r.PriceVariations != nil {
		priceVariations := make([]product.PriceVariation, len(*r.PriceVariations))
		for i, pv := range *r.PriceVariations {
			options := make([]product.Addon, len(pv.IncludedAddons.Options))
			for j, opt := range pv.IncludedAddons.Options {
				options[j] = product.Addon{
					ID:          opt.ID,
					Name:        opt.Name,
					Price:       opt.Price,
					Photos:      opt.Photos,
					IsAvailable: opt.IsAvailable,
				}
			}

			priceVariations[i] = product.PriceVariation{
				Type:  pv.Type,
				Price: pv.Price,
				IncludedAddons: product.IncludedAddons{
					MaxSelections: pv.IncludedAddons.MaxSelections,
					Options:       options,
				},
			}
		}
		input.PriceVariations = &priceVariations
	}

	// Convert available addons if provided
	if r.AvailableAddons != nil {
		availableAddons := make([]product.Addon, len(*r.AvailableAddons))
		for i, addon := range *r.AvailableAddons {
			availableAddons[i] = product.Addon{
				ID:          addon.ID,
				Name:        addon.Name,
				Price:       addon.Price,
				Photos:      addon.Photos,
				IsAvailable: addon.IsAvailable,
			}
		}
		input.AvailableAddons = &availableAddons
	}

	return input
}

// ProductListResponse represents a simplified product for list views
type ProductListResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Photos      []string `json:"photos"`
	Category    string   `json:"category"`
	MinPrice    int64    `json:"min_price"` // Minimum price from variations
	IsAvailable bool     `json:"is_available"`
}

// ToListResponse converts a product to list response
func ToListResponse(p *product.Product) ProductListResponse {
	minPrice := int64(0)
	if len(p.PriceVariations) > 0 {
		minPrice = p.PriceVariations[0].Price
		for _, pv := range p.PriceVariations {
			if pv.Price < minPrice {
				minPrice = pv.Price
			}
		}
	}

	return ProductListResponse{
		ID:          p.ID,
		Name:        p.Name,
		Photos:      p.Photos,
		Category:    p.Category,
		MinPrice:    minPrice,
		IsAvailable: p.IsAvailable,
	}
}

// ToListResponses converts multiple products to list responses
func ToListResponses(products []*product.Product) []ProductListResponse {
	responses := make([]ProductListResponse, len(products))
	for i, p := range products {
		responses[i] = ToListResponse(p)
	}
	return responses
}
