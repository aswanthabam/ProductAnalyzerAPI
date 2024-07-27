package products_db

import (
	api_error "productanalyzer/api/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProduct(product Product) *api_error.APIError {
	// TODO
	return nil
}

func GetProductByID(productId primitive.ObjectID) (*Product, *api_error.APIError) {
	// TODO
	return nil, nil
}

func GetProductsByUserID(userId primitive.ObjectID) (*[]Product, *api_error.APIError) {
	// TODO
	return nil, nil
}

func GetProductByProductIDAUserID(productId primitive.ObjectID, userId primitive.ObjectID) (*Product, *api_error.APIError) {
	// TODO
	return nil, nil
}
func UpdateProduct(product Product) *api_error.APIError {
	// TODO
	return nil
}

func DeleteProduct(productId primitive.ObjectID) *api_error.APIError {
	// TODO
	return nil
}
