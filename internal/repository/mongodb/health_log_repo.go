package mongodb

import (
	"context"
	"time"

	"health-check-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const healthLogCollection = "health_logs"

type HealthLogRepository struct {
	collection *mongo.Collection
}

func NewHealthLogRepository() *HealthLogRepository {
	return &HealthLogRepository{
		collection: GetCollection(healthLogCollection),
	}
}

// Create creates a new health log
func (r *HealthLogRepository) Create(ctx context.Context, log *domain.HealthLog) error {
	log.ID = primitive.NewObjectID()
	log.CheckedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, log)
	return err
}

// GetByServiceID retrieves logs for a specific service
func (r *HealthLogRepository) GetByServiceID(ctx context.Context, serviceID string, limit int64) ([]domain.HealthLog, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "checked_at", Value: -1}}).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{"service_id": serviceID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []domain.HealthLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// GetLatestByServiceID retrieves the latest log for a specific service
func (r *HealthLogRepository) GetLatestByServiceID(ctx context.Context, serviceID string) (*domain.HealthLog, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "checked_at", Value: -1}})

	var log domain.HealthLog
	err := r.collection.FindOne(ctx, bson.M{"service_id": serviceID}, opts).Decode(&log)
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetRecent retrieves recent logs across all services
func (r *HealthLogRepository) GetRecent(ctx context.Context, limit int64) ([]domain.HealthLog, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "checked_at", Value: -1}}).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []domain.HealthLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// UpdateNotificationSent updates the notification_sent field
func (r *HealthLogRepository) UpdateNotificationSent(ctx context.Context, id primitive.ObjectID, sent bool) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"notification_sent": sent}},
	)
	return err
}

// GetDownServicesWithPendingNotification gets services that are down and haven't sent notification
func (r *HealthLogRepository) GetDownServicesWithPendingNotification(ctx context.Context) ([]domain.HealthLog, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"status":            domain.StatusDOWN,
		"notification_sent": false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []domain.HealthLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// EnsureIndexes creates necessary indexes
func (r *HealthLogRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "service_id", Value: 1}, {Key: "checked_at", Value: -1}},
			Options: options.Index().SetBackground(true),
		},
		{
			Keys:    bson.D{{Key: "checked_at", Value: -1}},
			Options: options.Index().SetBackground(true),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}, {Key: "notification_sent", Value: 1}},
			Options: options.Index().SetBackground(true),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
