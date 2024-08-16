package db_repositories

import (
	db_types "productanalyzer/api/db/types"
	api_error "productanalyzer/api/errors"
	"time"
)

type ProductRepository[T db_types.PrimaryKey] interface {
	CreateProduct(product *db_types.Product[T]) (T, *api_error.APIError)
	CreateProductAccessKey(productID T, scope string) (*db_types.ProductAccessKey[T], *api_error.APIError)
	GetProductByID(productId T) (*db_types.Product[T], *api_error.APIError)
	GetProductsByUserID(userId T) (*[]db_types.Product[T], *api_error.APIError)
	GetProductByProductIDAUserID(productId string, userId T) (*db_types.Product[T], *api_error.APIError)
	UpdateProduct(product db_types.Product[T]) *api_error.APIError
	DeleteProduct(productId T) *api_error.APIError
	VisitProduct(ps *db_types.ProductUserSession[T], activity db_types.ProductActivity) *api_error.APIError
	GetVisitLogs(productId T, fromDate time.Time, toDate time.Time) (*[]db_types.VisitLogEntry, *api_error.APIError)
	ValidateAPIKey(apiKey, scope string) (*db_types.ProductAccessKey[T], *api_error.APIError)
	GetProductAccessKeys(productID T) (*[]db_types.ProductAccessKey[T], *api_error.APIError)
	GetProductByAccessKeyAndProductID(apiKey string, productId string) (*db_types.Product[T], *api_error.APIError)
	GetSessionById(sessionId T) (*db_types.ProductUserSession[T], error)
	SaveLocation(lc *db_types.Location[T]) error
	SaveProductUserSession(ps *db_types.ProductUserSession[T]) error
	ExistsLocationHash(lc db_types.Location[T]) (bool, error)
	ExistsProductUserSessionHash(ps db_types.ProductUserSession[T]) (bool, error)
	GetLocationByHash(lc *db_types.Location[T]) error
	GetProductUserSessionByHash(ps *db_types.ProductUserSession[T]) (bool, error)
	HashProductUserSession(ps *db_types.ProductUserSession[T]) error
}
