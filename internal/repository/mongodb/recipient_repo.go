package mongodb

import (
	"context"
	"time"

	"health-check-service/internal/domain"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const recipientCollection = "recipients"

type RecipientRepository struct {
	collection *mongo.Collection
}

func NewRecipientRepository() *RecipientRepository {
	return &RecipientRepository{
		collection: GetCollection(recipientCollection),
	}
}

// Create creates a new recipient
func (r *RecipientRepository) Create(ctx context.Context, recipient *domain.Recipient) error {
	recipient.ID = primitive.NewObjectID()
	recipient.CreatedAt = time.Now()
	recipient.IsActive = true

	logrus.Info("Recipient created: ", recipient)
	_, err := r.collection.InsertOne(ctx, recipient)
	return err
}

// GetByID retrieves a recipient by ID
func (r *RecipientRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Recipient, error) {
	var recipient domain.Recipient
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&recipient)
	if err != nil {
		return nil, err
	}
	return &recipient, nil
}

// GetAll retrieves all recipients
func (r *RecipientRepository) GetAll(ctx context.Context) ([]domain.Recipient, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var recipients []domain.Recipient
	if err = cursor.All(ctx, &recipients); err != nil {
		return nil, err
	}
	return recipients, nil
}

// GetActive retrieves all active recipients
func (r *RecipientRepository) GetActive(ctx context.Context) ([]domain.Recipient, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var recipients []domain.Recipient
	if err = cursor.All(ctx, &recipients); err != nil {
		return nil, err
	}
	return recipients, nil
}

// Update updates a recipient
func (r *RecipientRepository) Update(ctx context.Context, id primitive.ObjectID, update *domain.UpdateRecipientRequest) error {
	updateDoc := bson.M{}

	if update.Name != "" {
		updateDoc["name"] = update.Name
	}
	if update.Phone != "" {
		updateDoc["phone"] = update.Phone
	}
	if update.IsActive != nil {
		updateDoc["is_active"] = *update.IsActive
	}

	if len(updateDoc) == 0 {
		return nil
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updateDoc},
	)
	return err
}

// Delete deletes a recipient
func (r *RecipientRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// EnsureIndexes creates necessary indexes
func (r *RecipientRepository) EnsureIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "is_active", Value: 1}},
		Options: options.Index().SetBackground(true),
	}
	_, err := r.collection.Indexes().CreateOne(ctx, indexModel)
	return err
}
