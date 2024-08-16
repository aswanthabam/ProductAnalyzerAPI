package db

import (
	db_repositories "productanalyzer/api/db/repositories"
	db_types "productanalyzer/api/db/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectionInterface[T db_types.PrimaryKey] interface {
	Connect() error
	Close() error
	Initialize() error
	UserRepository() db_repositories.UserRepository[T]
	DeletionListRepository() db_repositories.DeletionListRepository[T]
	ProductRepository() db_repositories.ProductRepository[T]
}

var Connection ConnectionInterface[primitive.ObjectID]

var UserRepository db_repositories.UserRepository[primitive.ObjectID]
var DeletionListRepository db_repositories.DeletionListRepository[primitive.ObjectID]
var ProductRepository db_repositories.ProductRepository[primitive.ObjectID]

type User = db_types.User[primitive.ObjectID]
type DeletionList = db_types.DeletionList[primitive.ObjectID]
type Product = db_types.Product[primitive.ObjectID]
type ProductAccessKey = db_types.ProductAccessKey[primitive.ObjectID]
type ProductUserSession = db_types.ProductUserSession[primitive.ObjectID]
type Location = db_types.Location[primitive.ObjectID]

func InitializeRepositories() {
	UserRepository = Connection.UserRepository()
	DeletionListRepository = Connection.DeletionListRepository()
	ProductRepository = Connection.ProductRepository()
}
