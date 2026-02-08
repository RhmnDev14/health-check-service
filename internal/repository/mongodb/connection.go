package mongodb

import (
	"context"
	"time"

	"health-check-service/config"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
)

// Connect establishes connection to MongoDB
func Connect(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	Client = client
	Database = client.Database(cfg.MongoDatabase)

	logrus.Infof("Connected to MongoDB: %s", cfg.MongoDatabase)
	return nil
}

// Disconnect closes the MongoDB connection
func Disconnect() error {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return Client.Disconnect(ctx)
	}
	return nil
}

// GetCollection returns a collection from the database
func GetCollection(name string) *mongo.Collection {
	return Database.Collection(name)
}
