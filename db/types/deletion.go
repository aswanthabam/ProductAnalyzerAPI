package db_types

import "time"

const (
	DELETION_REQUEST_TYPE_INITIAL = "initial"
	DELETION_REQUEST_TYPE_CONFIRM = "confirm"
	DELETION_TYPE_PRODUCT         = "product"
	DELETION_TYPE_ACCESS_KEY      = "access_key"
)

type DeletionList[T PrimaryKey] struct {
	ID       T         `bson:"_id,omitempty"` // primary key
	ObjectID T         `bson:"object_id"`     // object id of the product
	Type     string    `bson:"type"`          // type of the object, product or access key
	Expires  time.Time `bson:"expires"`       // time at which the object will be deleted
}
