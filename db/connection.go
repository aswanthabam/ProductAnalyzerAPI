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
	Client             *mongo.Client     // Connection to MongoDB
	Database           *mongo.Database   // Database instance
	User               *mongo.Collection // Collection of users
	Visits             *mongo.Collection // Collection of visits
	Products           *mongo.Collection // Collection of products
	OTP                *mongo.Collection // Collection of OTPs
	Location           *mongo.Collection // Collection of locations
	ProductUserSession *mongo.Collection // Collection of user sessions
	DeletionList       *mongo.Collection // Collection of deletion list
	ProductAccessKey   *mongo.Collection // Collection of product access keys
}

/*
Connect to MongoDB and initialize the connection, return error if anything went wrong.
The DBConnection will be initialized only after calling this method.

Requires Environment Variables: MONGODB_URI, MONGODB_DB set before calling this method.
*/
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

/*
Fetch all collections from the database and store them in the DBConnection struct.
*/
func (conn *DBConnection) FetchCollections() error {
	conn.User = conn.collection("users")
	conn.Visits = conn.collection("visits")
	conn.Products = conn.collection("products")
	conn.OTP = conn.collection("otp")
	conn.Location = conn.collection("locations")
	conn.ProductUserSession = conn.collection("product_user_sessions")
	conn.DeletionList = conn.collection("deletion_list")
	conn.ProductAccessKey = conn.collection("product_access_keys")
	return nil
}

/*
Close the connection to MongoDB.
*/
func (conn *DBConnection) Close() error {
	if conn.Client != nil {
		return conn.Client.Disconnect(context.TODO())
	}
	return nil
}

/*
Return a collection by name.
*/
func (conn *DBConnection) collection(name string) *mongo.Collection {
	return conn.Database.Collection(name)
}

/*
Initialize the database by creating indexes and other necessary operations.
*/
func (conn *DBConnection) Initialize() error {
	createUniqueIndex(conn.User, "email")
	// createUniqueIndex(conn.ProductUserSession, "hash")
	createUniqueIndex(conn.ProductAccessKey, "access_key")
	createUniqueIndex(conn.Location, "hash")
	return nil
}

/*
Create a unique index on the collection.
*/
func createUniqueIndex(collection *mongo.Collection, key string) {
	_, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    map[string]interface{}{key: 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Print(err)
	}
}
