package deletion_db

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	DELETION_REQUEST_TYPE_INITIAL = "initial"
	DELETION_REQUEST_TYPE_CONFIRM = "confirm"
	DELETION_TYPE_PRODUCT         = "product"
	DELETION_TYPE_ACCESS_KEY      = "access_key"
)

type DeletionList struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"` // primary key
	ObjectID primitive.ObjectID `bson:"object_id"`     // object id of the product
	Type     string             `bson:"type"`          // type of the object, product or access key
	Expires  primitive.DateTime `bson:"expires"`       // time at which the object will be deleted
}
