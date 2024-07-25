package db

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func DuplicateError() error {
	return fmt.Errorf("DUPLICATE KEY ERROR")
}

func IsDuplicateKeyError(err error) bool {
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

func IsNotFoundError(err error) bool {
	return err == mongo.ErrNoDocuments
}
