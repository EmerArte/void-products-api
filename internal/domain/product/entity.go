package product

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product in the system with all its variations and addons
type Product struct {
	ID               string           `json:"id" bson:"_id"`
	CompanyID        string           `json:"company_id" bson:"company_id"`
	SalePointID      string           `json:"sale_point_id" bson:"sale_point_id"`
	Name             string           `json:"name" bson:"name"`
	Photos           []string         `json:"photos" bson:"photos"`
	PriceVariations  []PriceVariation `json:"price_variations" bson:"price_variations"`
	Category         string           `json:"category" bson:"category"`
	Description      string           `json:"description" bson:"description"`
	IsAddon          bool             `json:"is_addon" bson:"is_addon"`
	IsAvailable      bool             `json:"is_available" bson:"is_available"`
	IsUnlimitedStock bool             `json:"is_unlimited_stock" bson:"is_unlimited_stock"`
	Stock            *int             `json:"stock" bson:"stock"` // Pointer to allow null
	AvailableAddons  []Addon          `json:"available_addons" bson:"available_addons"`
	CreatedAt        time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" bson:"updated_at"`
}

// PriceVariation represents a variation of the product with different pricing
type PriceVariation struct {
	Type           string         `json:"type" bson:"type"`
	Price          int64          `json:"price" bson:"price"` // Price in cents/smallest currency unit
	IncludedAddons IncludedAddons `json:"included_addons" bson:"included_addons"`
}

// IncludedAddons represents addons that are included with a price variation
type IncludedAddons struct {
	MaxSelections int     `json:"max_selections" bson:"max_selections"`
	Options       []Addon `json:"options" bson:"options"`
}

// Addon represents an additional item that can be added to a product
type Addon struct {
	ID          string   `json:"id" bson:"_id"`
	Name        string   `json:"name" bson:"name"`
	Price       int64    `json:"price" bson:"price"` // Price in cents/smallest currency unit
	Photos      []string `json:"photos" bson:"photos"`
	IsAvailable bool     `json:"is_available" bson:"is_available"`
}

// NewProduct creates a new Product with generated UUID and timestamps
func NewProduct(companyID, salePointID, name, category, description string) *Product {
	now := time.Now()
	return &Product{
		ID:               uuid.New().String(),
		CompanyID:        companyID,
		SalePointID:      salePointID,
		Name:             name,
		Category:         category,
		Description:      description,
		Photos:           []string{},
		PriceVariations:  []PriceVariation{},
		AvailableAddons:  []Addon{},
		IsAddon:          false,
		IsAvailable:      true,
		IsUnlimitedStock: true,
		Stock:            nil,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// NewAddon creates a new Addon with generated UUID
func NewAddon(name string, price int64) *Addon {
	return &Addon{
		ID:          uuid.New().String(),
		Name:        name,
		Price:       price,
		Photos:      []string{},
		IsAvailable: true,
	}
}

// Validate performs business logic validation on the Product
func (p *Product) Validate() error {
	if p.CompanyID == "" {
		return ErrInvalidCompanyID
	}
	if p.SalePointID == "" {
		return ErrInvalidSalePointID
	}
	if p.Name == "" {
		return ErrInvalidName
	}
	if p.Category == "" {
		return ErrInvalidCategory
	}
	if len(p.PriceVariations) == 0 {
		return ErrNoPriceVariations
	}

	// Validate stock logic
	if !p.IsUnlimitedStock && p.Stock == nil {
		return ErrInvalidStock
	}
	if p.IsUnlimitedStock && p.Stock != nil {
		return ErrStockMustBeNullForUnlimited
	}
	if p.Stock != nil && *p.Stock < 0 {
		return ErrNegativeStock
	}

	// Validate price variations
	for i, pv := range p.PriceVariations {
		if pv.Type == "" {
			return ErrInvalidPriceVariationType
		}
		if pv.Price < 0 {
			return ErrNegativePrice
		}
		if pv.IncludedAddons.MaxSelections < 0 {
			return ErrInvalidMaxSelections
		}
		if pv.IncludedAddons.MaxSelections > 0 && len(pv.IncludedAddons.Options) == 0 {
			return ErrNoOptionsForMaxSelections
		}

		// Validate included addons
		for _, addon := range pv.IncludedAddons.Options {
			if addon.Name == "" {
				return ErrInvalidAddonName
			}
			if addon.Price < 0 {
				return ErrNegativeAddonPrice
			}
		}

		// Check for duplicate price variation types
		for j := i + 1; j < len(p.PriceVariations); j++ {
			if p.PriceVariations[j].Type == pv.Type {
				return ErrDuplicatePriceVariationType
			}
		}
	}

	// Validate available addons
	for _, addon := range p.AvailableAddons {
		if addon.Name == "" {
			return ErrInvalidAddonName
		}
		if addon.Price < 0 {
			return ErrNegativeAddonPrice
		}
	}

	return nil
}

// UpdateStock updates the product stock
func (p *Product) UpdateStock(newStock *int) error {
	if p.IsUnlimitedStock {
		return ErrCannotUpdateStockForUnlimited
	}
	if newStock != nil && *newStock < 0 {
		return ErrNegativeStock
	}
	p.Stock = newStock
	p.UpdatedAt = time.Now()
	return nil
}

// SetAvailability sets the availability status
func (p *Product) SetAvailability(available bool) {
	p.IsAvailable = available
	p.UpdatedAt = time.Now()
}

// AddPriceVariation adds a new price variation
func (p *Product) AddPriceVariation(variation PriceVariation) error {
	// Check for duplicates
	for _, pv := range p.PriceVariations {
		if pv.Type == variation.Type {
			return ErrDuplicatePriceVariationType
		}
	}

	if variation.Price < 0 {
		return ErrNegativePrice
	}

	p.PriceVariations = append(p.PriceVariations, variation)
	p.UpdatedAt = time.Now()
	return nil
}

// AddAvailableAddon adds a new available addon
func (p *Product) AddAvailableAddon(addon Addon) error {
	if addon.Name == "" {
		return ErrInvalidAddonName
	}
	if addon.Price < 0 {
		return ErrNegativeAddonPrice
	}

	// Check for duplicates by name
	for _, a := range p.AvailableAddons {
		if a.Name == addon.Name {
			return ErrDuplicateAddon
		}
	}

	p.AvailableAddons = append(p.AvailableAddons, addon)
	p.UpdatedAt = time.Now()
	return nil
}

// DecrementStock decrements the stock by the given amount
func (p *Product) DecrementStock(amount int) error {
	if p.IsUnlimitedStock {
		return nil // No stock management needed
	}

	if p.Stock == nil {
		return ErrInvalidStock
	}

	newStock := *p.Stock - amount
	if newStock < 0 {
		return ErrInsufficientStock
	}

	p.Stock = &newStock
	p.UpdatedAt = time.Now()
	return nil
}
