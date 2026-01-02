package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/emerarteaga/products-api/internal/domain/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderMongoRepository struct {
	collection *mongo.Collection
}

// NewOrderMongoRepository creates a new order repository
func NewOrderMongoRepository(collection *mongo.Collection) order.Repository {
	return &orderMongoRepository{collection: collection}
}

// CreateIndexes creates the necessary indexes for the orders collection
func (r *orderMongoRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "code", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "sale_type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "updated_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "products.id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "products.name", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "sale_type", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// Create creates a new order
func (r *orderMongoRepository) Create(ctx context.Context, o *order.Order) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, o)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return order.ErrOrderCodeAlreadyExists
		}
		return fmt.Errorf("failed to insert order: %w", err)
	}

	return nil
}

// FindByID finds an order by ID
func (r *orderMongoRepository) FindByID(ctx context.Context, id string) (*order.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var o order.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&o)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, order.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	return &o, nil
}

// FindByCode finds an order by tracking code
func (r *orderMongoRepository) FindByCode(ctx context.Context, code string) (*order.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var o order.Order
	err := r.collection.FindOne(ctx, bson.M{"code": code}).Decode(&o)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, order.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	return &o, nil
}

// Update updates an order
func (r *orderMongoRepository) Update(ctx context.Context, o *order.Order) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	o.UpdatedAt = time.Now()

	update := bson.M{
		"$set": o,
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": o.ID}, update)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	if result.MatchedCount == 0 {
		return order.ErrOrderNotFound
	}

	return nil
}

// FindAll retrieves all orders with optional filters
func (r *orderMongoRepository) FindAll(ctx context.Context, filters order.OrderFilters) ([]*order.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Build filter
	filter := bson.M{}
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
		return nil, fmt.Errorf("failed to find orders: %w", err)
	}
	defer cursor.Close(ctx)

	var orders []*order.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, fmt.Errorf("failed to decode orders: %w", err)
	}

	return orders, nil
}

// Count returns the total number of orders matching filters
func (r *orderMongoRepository) Count(ctx context.Context, filters order.OrderFilters) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{}
	r.applyFilters(filter, filters)

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count orders: %w", err)
	}

	return count, nil
}

// ExistsByCode checks if an order exists with the given code
func (r *orderMongoRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{"code": code})
	if err != nil {
		return false, fmt.Errorf("failed to check order existence: %w", err)
	}

	return count > 0, nil
}

// GetMetrics returns aggregated order metrics
func (r *orderMongoRepository) GetMetrics(ctx context.Context, filters order.OrderFilters) (*order.OrderMetrics, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Build base filter
	matchFilter := bson.M{}
	r.applyFilters(matchFilter, filters)

	// Aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchFilter}},
		{{Key: "$facet", Value: bson.M{
			"metrics": []bson.M{
				{
					"$group": bson.M{
						"_id":         nil,
						"total_sales": bson.M{"$sum": "$total"},
						"avg_ticket":  bson.M{"$avg": "$total"},
						"count":       bson.M{"$sum": 1},
					},
				},
			},
			"by_status": []bson.M{
				{
					"$group": bson.M{
						"_id":   "$status",
						"count": bson.M{"$sum": 1},
					},
				},
			},
			"top_products": []bson.M{
				{"$unwind": "$products"},
				{
					"$group": bson.M{
						"_id": bson.M{
							"id":   "$products.id",
							"name": "$products.name",
						},
						"total_quantity": bson.M{"$sum": "$products.quantity"},
						"total_revenue": bson.M{
							"$sum": bson.M{
								"$multiply": []interface{}{
									"$products.price",
									"$products.quantity",
								},
							},
						},
					},
				},
				{"$sort": bson.M{"total_quantity": -1}},
				{"$limit": 10},
				{
					"$project": bson.M{
						"product_id":     "$_id.id",
						"name":           "$_id.name",
						"total_quantity": 1,
						"total_revenue":  1,
						"_id":            0,
					},
				},
			},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate metrics: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		Metrics []struct {
			TotalSales int64 `bson:"total_sales"`
			AvgTicket  int64 `bson:"avg_ticket"`
		} `bson:"metrics"`
		ByStatus []struct {
			Status order.OrderStatus `bson:"_id"`
			Count  int               `bson:"count"`
		} `bson:"by_status"`
		TopProducts []order.ProductSalesSummary `bson:"top_products"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode metrics: %w", err)
	}

	if len(results) == 0 {
		return &order.OrderMetrics{
			OrdersByStatus: make(map[order.OrderStatus]int),
			TopProducts:    []order.ProductSalesSummary{},
		}, nil
	}

	result := results[0]
	metrics := &order.OrderMetrics{
		OrdersByStatus: make(map[order.OrderStatus]int),
		TopProducts:    result.TopProducts,
	}

	if len(result.Metrics) > 0 {
		metrics.TotalSales = result.Metrics[0].TotalSales
		metrics.AvgTicket = result.Metrics[0].AvgTicket
	}

	for _, statusCount := range result.ByStatus {
		metrics.OrdersByStatus[statusCount.Status] = statusCount.Count
	}

	return metrics, nil
}

// applyFilters applies filters to the query
func (r *orderMongoRepository) applyFilters(filter bson.M, filters order.OrderFilters) {
	if filters.Status != nil {
		filter["status"] = *filters.Status
	}

	if filters.SaleType != nil {
		filter["sale_type"] = *filters.SaleType
	}

	if filters.ProductID != nil {
		filter["products.id"] = *filters.ProductID
	}

	if filters.ProductName != nil {
		filter["products.name"] = bson.M{"$regex": *filters.ProductName, "$options": "i"}
	}

	if filters.MinTotal != nil || filters.MaxTotal != nil {
		totalFilter := bson.M{}
		if filters.MinTotal != nil {
			totalFilter["$gte"] = *filters.MinTotal
		}
		if filters.MaxTotal != nil {
			totalFilter["$lte"] = *filters.MaxTotal
		}
		filter["total"] = totalFilter
	}

	if filters.DateFrom != nil || filters.DateTo != nil {
		dateFilter := bson.M{}
		if filters.DateFrom != nil {
			if t, err := time.Parse(time.RFC3339, *filters.DateFrom); err == nil {
				dateFilter["$gte"] = t
			}
		}
		if filters.DateTo != nil {
			if t, err := time.Parse(time.RFC3339, *filters.DateTo); err == nil {
				dateFilter["$lte"] = t
			}
		}
		if len(dateFilter) > 0 {
			filter["created_at"] = dateFilter
		}
	}
}
