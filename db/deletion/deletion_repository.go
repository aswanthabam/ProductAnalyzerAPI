package deletion_db

import (
	"context"
	"productanalyzer/api/db"
	api_error "productanalyzer/api/errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddToDeletionList(objectID primitive.ObjectID, objectType string) (primitive.ObjectID, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := db.Connection.DeletionList.InsertOne(ctx, bson.M{"object_id": objectID, "type": objectType,
		"expires": primitive.NewDateTimeFromTime(time.Now().Add(2 * time.Minute)).Time().UTC()})
	if err != nil {
		return primitive.NilObjectID, api_error.UnexpectedError(err)
	}
	if id == nil {
		return primitive.NilObjectID, api_error.NewAPIError("Failed to add to deletion list", 500, "Failed to add to deletion list")
	}
	return id.InsertedID.(primitive.ObjectID), nil
}

func GetFromDeletionList(objectId primitive.ObjectID) (*DeletionList, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	deletionList := DeletionList{}
	if err := db.Connection.DeletionList.FindOne(ctx, bson.M{"_id": objectId}).Decode(&deletionList); err != nil {
		return nil, api_error.NewAPIError("Unable to complete deletion.", 404, "Requested object not found in deletion list")
	}
	if deletionList.Expires.Time().UTC().Before(time.Now().UTC()) {
		return nil, api_error.NewAPIError("Unable to complete deletion", 404, "Time to delete the object has expired, try again")
	}
	return &deletionList, nil
}
