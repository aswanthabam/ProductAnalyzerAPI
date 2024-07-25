package db

import (
	"context"
	"fmt"
	"log"
	"productanalyzer/api/config"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Connection DBConnection

var mu sync.Mutex

type DBConnection struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func (conn *DBConnection) Connect() error {
	mu.Lock()
	defer mu.Unlock()

	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(config.Config.MONGODB_URI))
	if err != nil {
		return err
	}
	database := client.Database(config.Config.MONGODB_DB)
	if database == nil {
		return fmt.Errorf("DATABASE '%s' NOT FOUND", config.Config.MONGODB_DB)
	}
	conn.Client = client
	conn.Database = database
	return nil
}

func (conn *DBConnection) Close() error {
	if conn.Client != nil {
		return conn.Client.Disconnect(context.TODO())
	}
	return nil
}

func (conn *DBConnection) Collection(name string) *mongo.Collection {
	return conn.Database.Collection(name)
}

func (conn *DBConnection) Initialize() error {
	user := conn.Collection("users")
	createUniqueIndex(user, "email")
	return nil
}

func createUniqueIndex(collection *mongo.Collection, key string) {
	_, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    map[string]interface{}{key: 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Print(err)
	}
}
