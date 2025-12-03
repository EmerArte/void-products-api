package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/emerarteaga/products-api/internal/domain/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productMongoRepository struct {
	collection *mongo.Collection
}

// NewProductMongoRepository creates a new product repository
func NewProductMongoRepository(collection *mongo.Collection) product.Repository {
	return &productMongoRepository{collection: collection}
}

// CreateIndexes creates the necessary indexes for the products collection
func (r *productMongoRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "company_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "sale_point_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "company_id", Value: 1},
				{Key: "category", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "sale_point_id", Value: 1},
				{Key: "is_available", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "company_id", Value: 1},
				{Key: "is_available", Value: 1},
			},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// Create creates a new product
func (r *productMongoRepository) Create(ctx context.Context, p *product.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	return nil
}

// FindByID finds a product by ID
func (r *productMongoRepository) FindByID(ctx context.Context, id string) (*product.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var p product.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, product.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	return &p, nil
}

// FindByCompanyID retrieves all products for a company with optional filters
func (r *productMongoRepository) FindByCompanyID(ctx context.Context, companyID string, filters product.ProductFilters) ([]*product.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Build filter
	filter := bson.M{"company_id": companyID}
	r.applyFilters(filter, filters)

	// Set default pagination
	if filters.Limit <= 0 {
		filters.Limit = 50
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	opts := options.Find().
		SetLimit(int64(filters.Limit)).
		SetSkip(int64(filters.Offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find products: %w", err)
	}
	defer cursor.Close(ctx)

	return r.decodeProducts(ctx, cursor)
}

// FindBySalePointID retrieves all products for a sale point with optional filters
func (r *productMongoRepository) FindBySalePointID(ctx context.Context, salePointID string, filters product.ProductFilters) ([]*product.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Build filter
	filter := bson.M{"sale_point_id": salePointID}
	r.applyFilters(filter, filters)

	// Set default pagination
	if filters.Limit <= 0 {
		filters.Limit = 50
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	opts := options.Find().
		SetLimit(int64(filters.Limit)).
		SetSkip(int64(filters.Offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find products: %w", err)
	}
	defer cursor.Close(ctx)

	return r.decodeProducts(ctx, cursor)
}

// FindAll finds all products with pagination (deprecated)
func (r *productMongoRepository) FindAll(ctx context.Context, limit, offset int) ([]*product.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find products: %w", err)
	}
	defer cursor.Close(ctx)

	return r.decodeProducts(ctx, cursor)
}

// Update updates a product
func (r *productMongoRepository) Update(ctx context.Context, p *product.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	p.UpdatedAt = time.Now()

	update := bson.M{
		"$set": p,
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": p.ID}, update)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	if result.MatchedCount == 0 {
		return product.ErrProductNotFound
	}

	return nil
}

// Delete deletes a product (hard delete)
func (r *productMongoRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if result.DeletedCount == 0 {
		return product.ErrProductNotFound
	}

	return nil
}

// FindCategoriesByCompanyID retrieves all unique categories for a company
func (r *productMongoRepository) FindCategoriesByCompanyID(ctx context.Context, companyID string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"company_id": companyID}

	categories, err := r.collection.Distinct(ctx, "category", filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}

	result := make([]string, 0, len(categories))
	for _, cat := range categories {
		if str, ok := cat.(string); ok && str != "" {
			result = append(result, str)
		}
	}

	return result, nil
}

// FindCategoriesBySalePointID retrieves all unique categories for a sale point
func (r *productMongoRepository) FindCategoriesBySalePointID(ctx context.Context, salePointID string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"sale_point_id": salePointID}

	categories, err := r.collection.Distinct(ctx, "category", filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}

	result := make([]string, 0, len(categories))
	for _, cat := range categories {
		if str, ok := cat.(string); ok && str != "" {
			result = append(result, str)
		}
	}

	return result, nil
}

// Count returns the total number of products
func (r *productMongoRepository) Count(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}

// CountByCompanyID returns the total number of products for a company with filters
func (r *productMongoRepository) CountByCompanyID(ctx context.Context, companyID string, filters product.ProductFilters) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Build filter (same as FindByCompanyID but without pagination)
	filter := bson.M{"company_id": companyID}
	r.applyFilters(filter, filters)

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}

// CountBySalePointID returns the total number of products for a sale point with filters
func (r *productMongoRepository) CountBySalePointID(ctx context.Context, salePointID string, filters product.ProductFilters) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Build filter (same as FindBySalePointID but without pagination)
	filter := bson.M{"sale_point_id": salePointID}
	r.applyFilters(filter, filters)

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}

// Exists checks if a product exists
func (r *productMongoRepository) Exists(ctx context.Context, id string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, fmt.Errorf("failed to check product existence: %w", err)
	}

	return count > 0, nil
}

// applyFilters applies filters to the MongoDB filter document
func (r *productMongoRepository) applyFilters(filter bson.M, filters product.ProductFilters) {
	if filters.Category != nil {
		filter["category"] = *filters.Category
	}
	if filters.IsAvailable != nil {
		filter["is_available"] = *filters.IsAvailable
	}
	if filters.IsAddon != nil {
		filter["is_addon"] = *filters.IsAddon
	}
}

// decodeProducts decodes products from cursor
func (r *productMongoRepository) decodeProducts(ctx context.Context, cursor *mongo.Cursor) ([]*product.Product, error) {
	var products []*product.Product
	for cursor.Next(ctx) {
		var p product.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, fmt.Errorf("failed to decode product: %w", err)
		}
		products = append(products, &p)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return products, nil
}
