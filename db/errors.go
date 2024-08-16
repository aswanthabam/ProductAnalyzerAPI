package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func IsNotFoundError(err error) bool {
	return err == mongo.ErrNoDocuments
}
