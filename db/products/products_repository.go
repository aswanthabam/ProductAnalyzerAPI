package products_db

import (
	"context"
	"log"
	"productanalyzer/api/db"
	api_error "productanalyzer/api/errors"
	"productanalyzer/api/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Inserts a new product into the database and returns the id of the product
func CreateProduct(product *Product) (primitive.ObjectID, *api_error.APIError) {
	if _, err := GetProductByProductIDAUserID(product.ProductID, product.UserID); err == nil {
		return primitive.NilObjectID, api_error.NewAPIError("Product Already Exists", 409, "Product with the same product id already exists")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	curTime := utils.GetCurrentTime()
	product.AccessKeys = []primitive.ObjectID{}
	product.CreatedAt = curTime
	product.UpdatedAt = curTime
	result, err2 := db.Connection.Products.InsertOne(ctx, product)
	if err2 != nil {
		return primitive.NilObjectID, api_error.UnexpectedError(err2)
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
func CreateProductAccessKey(productID primitive.ObjectID, scope string) (*ProductAccessKey, *api_error.APIError) {
	key, err := utils.GenerateAPIKey()
	if err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	accessKey := ProductAccessKey{
		ProductID: productID,
		AccessKey: key,
		Scope:     scope,
		CreatedAt: utils.GetCurrentTime(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := db.Connection.ProductAccessKey.InsertOne(ctx, accessKey)
	if err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	accessKey.ID = id.InsertedID.(primitive.ObjectID)
	return &accessKey, nil
}

// Get a product by its object id
func GetProductByID(productId primitive.ObjectID) (*Product, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	product := Product{}
	if err := db.Connection.Products.FindOne(ctx, bson.M{"_id": productId}).Decode(&product); err != nil {
		return nil, api_error.NewAPIError("Product Not Found", 404, "The requested product was not found")
	}
	return &product, nil
}

// Get all products created by a user
func GetProductsByUserID(userId primitive.ObjectID) (*[]Product, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	products := []Product{}
	cursor, err := db.Connection.Products.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	if err := cursor.All(ctx, &products); err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	return &products, nil
}

// Get a product by its object id and user id
func GetProductByProductIDAUserID(productId string, userId primitive.ObjectID) (*Product, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	product := Product{}
	if err := db.Connection.Products.FindOne(ctx, bson.M{"product_id": productId, "user_id": userId}).Decode(&product); err != nil {
		return nil, api_error.NewAPIError("Product Not Found", 404, "Product not found")
	}
	return &product, nil
}

// Update a product in the database
func UpdateProduct(product Product) *api_error.APIError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.Connection.Products.UpdateOne(ctx, bson.M{"_id": product.ID}, bson.M{"$set": product})
	if err != nil {
		return api_error.UnexpectedError(err)
	}
	return nil
}

// Delete a product from the database
func DeleteProduct(productId primitive.ObjectID) *api_error.APIError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.Connection.Products.DeleteOne(ctx, bson.M{"_id": productId})
	if err != nil {
		return api_error.UnexpectedError(err)
	}
	return nil
}

// VisitProduct visits a product and logs the activity,
// If a visit exists with the given session, it appends the activity to the existing visit,
// otherwise it creates a new visit with the activity.
func VisitProduct(productId primitive.ObjectID, sessionId primitive.ObjectID, activity ProductActivity) *api_error.APIError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	visit := ProductVisit{}
	err := db.Connection.ProductVisits.FindOne(ctx, bson.M{"session_id": sessionId}).Decode(&visit)
	if err != nil {
		visit = ProductVisit{
			ProductID: productId,
			SessionID: sessionId,
			Refferer:  "",
			Activities: []ProductActivity{
				activity,
			},
		}
		_, err = db.Connection.ProductVisits.InsertOne(ctx, visit)
	} else {
		visit.Activities = append(visit.Activities, activity)
		_, err = db.Connection.ProductVisits.UpdateOne(ctx, bson.M{"session_id": sessionId}, bson.M{"$set": bson.M{"activities": visit.Activities}})
	}
	if err != nil {
		return api_error.UnexpectedError(err)
	}
	return nil
}

// ValidateAPIKey validates the API Key and returns the ProductAccessKey if the key is valid
// and has the required scope, otherwise it returns an error. This method uses the hashed key for validation.
func ValidateAPIKey(productId primitive.ObjectID, apiKey, scope string) (*ProductAccessKey, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accessKey := ProductAccessKey{}
	if err := db.Connection.ProductAccessKey.FindOne(ctx, bson.M{
		"product_id": productId,
		"access_key": apiKey,
	}).Decode(&accessKey); err != nil {
		log.Print(err)
		return nil, api_error.NewAPIError("Invalid API Key", 401, "Invalid API Key")
	}
	if accessKey.Scope != scope && accessKey.Scope != PRODUCT_ACCESS_KEY_SCOPE_ALL {
		return nil, api_error.NewAPIError("Invalid API Key", 401, "Invalid Scope, the API Key does not have the required permissions")
	}
	return &accessKey, nil
}

// GetProductAccessKeys returns all the access keys for a product
func GetProductAccessKeys(productID primitive.ObjectID) (*[]ProductAccessKey, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	keys := []ProductAccessKey{}
	cursor, err := db.Connection.ProductAccessKey.Find(ctx, bson.M{"product_id": productID})
	if err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	if err := cursor.All(ctx, &keys); err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	return &keys, nil
}

// GetProductByAccessKeyAndProductID returns the product with the given product id and access key
func GetProductByAccessKeyAndProductID(apiKey string, productId string) (*Product, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	product := Product{}
	if err := db.Connection.Products.FindOne(ctx, bson.M{"product_id": productId, "access_keys.access_key": apiKey}).Decode(&product); err != nil {
		return nil, api_error.NewAPIError("Product Not Found", 404, "Product not found")
	}
	return &product, nil
}
