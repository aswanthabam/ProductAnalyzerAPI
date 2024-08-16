package db_mongo

import (
	"context"
	"fmt"
	"log"
	"productanalyzer/api/config"
	db_repositories "productanalyzer/api/db/repositories"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mu sync.Mutex

type MongoConnection struct {
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
The MongoConnection will be initialized only after calling this method.

Requires Environment Variables: MONGODB_URI, MONGODB_DB set before calling this method.
*/
func (conn *MongoConnection) Connect() error {
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

	err = conn.FetchCollections()
	if err != nil {
		log.Panic(err)
	}
	return nil
}

/*
Fetch all collections from the database and store them in the MongoConnection struct.
*/
func (conn *MongoConnection) FetchCollections() error {
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
func (conn *MongoConnection) Close() error {
	if conn.Client != nil {
		return conn.Client.Disconnect(context.TODO())
	}
	return nil
}

/*
Return a collection by name.
*/
func (conn *MongoConnection) collection(name string) *mongo.Collection {
	return conn.Database.Collection(name)
}

/*
Initialize the database by creating indexes and other necessary operations.
*/
func (conn *MongoConnection) Initialize() error {
	createUniqueIndex(conn.User, "email")
	// createUniqueIndex(conn.ProductUserSession, "hash")
	createUniqueIndex(conn.ProductAccessKey, "access_key")
	createUniqueIndex(conn.Location, "hash")
	return nil
}

func (conn *MongoConnection) UserRepository() db_repositories.UserRepository[primitive.ObjectID] {
	return &MongoUserRepository{Connection: conn}
}

func (conn *MongoConnection) ProductRepository() db_repositories.ProductRepository[primitive.ObjectID] {
	return &MongoProductRepository{Connection: conn}
}

func (conn *MongoConnection) DeletionListRepository() db_repositories.DeletionListRepository[primitive.ObjectID] {
	return &MongoDeletionListRepository{Connection: conn}
}

/* ERRORS */

func (conn *MongoConnection) IsDuplicateKeyError(err error) bool {
	if mongoErr, ok := err.(mongo.WriteException); ok {
		for _, writeError := range mongoErr.WriteErrors {
			log.Println(writeError.Code)
			if writeError.Code == 11000 {
				return true
			}
		}
	}
	return false
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
