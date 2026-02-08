package mongodb

import (
	"context"
	"time"

	"health-check-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const serviceCollection = "services"

type ServiceRepository struct {
	collection *mongo.Collection
}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{
		collection: GetCollection(serviceCollection),
	}
}

// Create creates a new service
func (r *ServiceRepository) Create(ctx context.Context, service *domain.Service) error {
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	service.IsActive = true

	_, err := r.collection.InsertOne(ctx, service)
	return err
}

// GetByID retrieves a service by ID
func (r *ServiceRepository) GetByID(ctx context.Context, id string) (*domain.Service, error) {
	var service domain.Service
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&service)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// GetAll retrieves all services
func (r *ServiceRepository) GetAll(ctx context.Context) ([]domain.Service, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var services []domain.Service
	if err = cursor.All(ctx, &services); err != nil {
		return nil, err
	}
	return services, nil
}

// GetActive retrieves all active services
func (r *ServiceRepository) GetActive(ctx context.Context) ([]domain.Service, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var services []domain.Service
	if err = cursor.All(ctx, &services); err != nil {
		return nil, err
	}
	return services, nil
}

// Update updates a service
func (r *ServiceRepository) Update(ctx context.Context, id string, update *domain.UpdateServiceRequest) error {
	updateDoc := bson.M{"updated_at": time.Now()}

	if update.Name != "" {
		updateDoc["name"] = update.Name
	}
	if update.URL != "" {
		updateDoc["url"] = update.URL
	}
	if update.CheckInterval > 0 {
		updateDoc["check_interval"] = update.CheckInterval
	}
	if update.RetryInterval > 0 {
		updateDoc["retry_interval"] = update.RetryInterval
	}
	if update.IsActive != nil {
		updateDoc["is_active"] = *update.IsActive
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updateDoc},
	)
	return err
}

// Delete deletes a service
func (r *ServiceRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// EnsureIndexes creates necessary indexes
func (r *ServiceRepository) EnsureIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "is_active", Value: 1}},
		Options: options.Index().SetBackground(true),
	}
	_, err := r.collection.Indexes().CreateOne(ctx, indexModel)
	return err
}
