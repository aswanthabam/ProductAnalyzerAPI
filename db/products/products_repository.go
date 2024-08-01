package products_db

import (
	"context"
	"productanalyzer/api/db"
	api_error "productanalyzer/api/errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Inserts a new product into the database and returns the id of the product
func CreateProduct(product Product) (primitive.ObjectID, *api_error.APIError) {
	if _, err := GetProductByProductIDAUserID(product.ProductID, product.UserID); err == nil {
		return primitive.NilObjectID, api_error.NewAPIError("Product Already Exists", 409, "Product with the same product id already exists")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	product.CreatedAt = primitive.NewDateTimeFromTime(time.Now().UTC())
	product.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().UTC())
	result, err := db.Connection.Products.InsertOne(ctx, product)
	if err != nil {
		return primitive.NilObjectID, api_error.UnexpectedError(err)
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// Get a product by its object id
func GetProductByID(productId primitive.ObjectID) (*Product, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	product := Product{}
	if err := db.Connection.Products.FindOne(ctx, bson.M{"_id": productId}).Decode(&product); err != nil {
		return nil, api_error.NewAPIError("Product Not Found", 404, "Product not found")
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
