package db_types

import "go.mongodb.org/mongo-driver/bson/primitive"

type PrimaryKey interface {
	primitive.ObjectID | string
}
