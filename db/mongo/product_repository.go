package db_mongo

import (
	"context"
	"log"
	db_types "productanalyzer/api/db/types"
	api_error "productanalyzer/api/errors"
	"productanalyzer/api/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoProductRepository struct {
	Connection *MongoConnection
}

type MongoProduct = db_types.Product[primitive.ObjectID]
type MongoProductUserSession = db_types.ProductUserSession[primitive.ObjectID]
type MongoProductAccessKey = db_types.ProductAccessKey[primitive.ObjectID]
type MongoLocation = db_types.Location[primitive.ObjectID]

// Inserts a new product into the database and returns the id of the product
func (m *MongoProductRepository) CreateProduct(product *MongoProduct) (primitive.ObjectID, *api_error.APIError) {
	if _, err := m.GetProductByProductIDAUserID(product.ProductID, product.UserID); err == nil {
		return primitive.NilObjectID, api_error.NewAPIError("Product Already Exists", 409, "Product with the same product id already exists")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	curTime := utils.GetCurrentTime()
	product.AccessKeys = []primitive.ObjectID{}
	product.CreatedAt = curTime
	product.UpdatedAt = curTime
	result, err2 := m.Connection.Products.InsertOne(ctx, product)
	if err2 != nil {
		return primitive.NilObjectID, api_error.UnexpectedError(err2)
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
func (m *MongoProductRepository) CreateProductAccessKey(productID primitive.ObjectID, scope string) (*MongoProductAccessKey, *api_error.APIError) {
	key, err := utils.GenerateAPIKey()
	if err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	accessKey := MongoProductAccessKey{
		ProductID: productID,
		AccessKey: key,
		Scope:     scope,
		CreatedAt: utils.GetCurrentTime(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := m.Connection.ProductAccessKey.InsertOne(ctx, accessKey)
	if err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	accessKey.ID = id.InsertedID.(primitive.ObjectID)
	return &accessKey, nil
}

// Get a product by its object id
func (m *MongoProductRepository) GetProductByID(productId primitive.ObjectID) (*MongoProduct, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	product := MongoProduct{}
	if err := m.Connection.Products.FindOne(ctx, bson.M{"_id": productId}).Decode(&product); err != nil {
		return nil, api_error.NewAPIError("Product Not Found", 404, "The requested product was not found")
	}
	return &product, nil
}

// Get all products created by a user
func (m *MongoProductRepository) GetProductsByUserID(userId primitive.ObjectID) (*[]MongoProduct, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	products := []MongoProduct{}
	cursor, err := m.Connection.Products.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	if err := cursor.All(ctx, &products); err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	return &products, nil
}

// Get a product by its object id and user id
func (m *MongoProductRepository) GetProductByProductIDAUserID(productId string, userId primitive.ObjectID) (*MongoProduct, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	product := MongoProduct{}
	if err := m.Connection.Products.FindOne(ctx, bson.M{"product_id": productId, "user_id": userId}).Decode(&product); err != nil {
		return nil, api_error.NewAPIError("Product Not Found", 404, "Product not found")
	}
	return &product, nil
}

// Update a product in the database
func (m *MongoProductRepository) UpdateProduct(product MongoProduct) *api_error.APIError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.Connection.Products.UpdateOne(ctx, bson.M{"_id": product.ID}, bson.M{"$set": product})
	if err != nil {
		return api_error.UnexpectedError(err)
	}
	return nil
}

// Delete a product from the database
func (m *MongoProductRepository) DeleteProduct(productId primitive.ObjectID) *api_error.APIError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.Connection.Products.DeleteOne(ctx, bson.M{"_id": productId})
	if err != nil {
		return api_error.UnexpectedError(err)
	}
	return nil
}

// VisitProduct visits a product and logs the activity,
// If a visit exists with the given session, it appends the activity to the existing visit,
// otherwise it creates a new visit with the activity.
func (m *MongoProductRepository) VisitProduct(ps *MongoProductUserSession, activity db_types.ProductActivity) *api_error.APIError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if ps.ID.IsZero() {
		return api_error.NewAPIError("Invalid Session", 400, "Invalid Session")
	}
	ps.Activities = append(ps.Activities, activity)
	ps.UpdatedAt = utils.GetCurrentTime()
	_, err := m.Connection.ProductUserSession.UpdateOne(ctx,
		bson.M{"_id": ps.ID},
		bson.M{
			"$set": bson.M{
				"activities": ps.Activities,
				"updated_at": ps.UpdatedAt,
			},
		})
	if err != nil {
		return api_error.UnexpectedError(err)
	}
	return nil
}

func (m *MongoProductRepository) GetVisitLogs(productId primitive.ObjectID, fromDate time.Time, toDate time.Time) (*[]db_types.VisitLogEntry, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sessions := []db_types.VisitLogEntry{}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{
			"product_id": productId,
			"updated_at": bson.M{"$gte": fromDate, "$lte": toDate},
		}}},
		bson.D{{Key: "$unwind", Value: "$activities"}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":            "$_id",
			"activity_count": bson.M{"$sum": 1},
			"created_at":     bson.M{"$first": "$created_at"},
			"updated_at":     bson.M{"$last": "$updated_at"},
			"referer":        bson.M{"$first": "$referer"},
		}}},
		bson.D{{Key: "$sort", Value: bson.M{"created_at": 1}}},
		bson.D{{Key: "$project", Value: bson.M{
			"created_at":     1,
			"updated_at":     1,
			"activity_count": 1,
			"referer":        1,
		}}},
	}
	cursor, err := m.Connection.ProductUserSession.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	return &sessions, nil
}

// ValidateAPIKey validates the API Key and returns the ProductAccessKey if the key is valid
// and has the required scope, otherwise it returns an error. This method uses the hashed key for validation.
func (m *MongoProductRepository) ValidateAPIKey(apiKey, scope string) (*MongoProductAccessKey, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accessKey := MongoProductAccessKey{}
	if err := m.Connection.ProductAccessKey.FindOne(ctx, bson.M{
		"access_key": apiKey,
	}).Decode(&accessKey); err != nil {
		log.Print(err)
		return nil, api_error.NewAPIError("Invalid API Key", 401, "Invalid API Key")
	}
	if accessKey.Scope != scope && accessKey.Scope != db_types.PRODUCT_ACCESS_KEY_SCOPE_ALL {
		return nil, api_error.NewAPIError("Invalid API Key", 401, "Invalid Scope, the API Key does not have the required permissions")
	}
	return &accessKey, nil
}

// GetProductAccessKeys returns all the access keys for a product
func (m *MongoProductRepository) GetProductAccessKeys(productID primitive.ObjectID) (*[]MongoProductAccessKey, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	keys := []MongoProductAccessKey{}
	cursor, err := m.Connection.ProductAccessKey.Find(ctx, bson.M{"product_id": productID})
	if err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	if err := cursor.All(ctx, &keys); err != nil {
		return nil, api_error.UnexpectedError(err)
	}
	return &keys, nil
}

// GetProductByAccessKeyAndProductID returns the product with the given product id and access key
func (m *MongoProductRepository) GetProductByAccessKeyAndProductID(apiKey string, productId string) (*MongoProduct, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	product := MongoProduct{}
	if err := m.Connection.Products.FindOne(ctx, bson.M{"product_id": productId, "access_keys.access_key": apiKey}).Decode(&product); err != nil {
		return nil, api_error.NewAPIError("Product Not Found", 404, "Product not found")
	}
	return &product, nil
}

func (m *MongoProductRepository) GetSessionById(sessionId primitive.ObjectID) (*MongoProductUserSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	session := MongoProductUserSession{}
	err := m.Connection.ProductUserSession.FindOne(ctx, bson.M{"_id": sessionId}).Decode(&session)
	return &session, err
}

/* SAVE METHODS */

func (m *MongoProductRepository) SaveLocation(lc *MongoLocation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if !lc.ID.IsZero() {
		_, err := m.Connection.Location.UpdateOne(ctx, bson.M{"_id": lc.ID}, bson.M{"$set": lc})
		return err
	}
	id, err := m.Connection.Location.InsertOne(ctx, lc)
	if err != nil {
		return err
	}
	lc.ID = id.InsertedID.(primitive.ObjectID)
	return nil
}

func (m *MongoProductRepository) SaveProductUserSession(ps *MongoProductUserSession) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if !ps.ID.IsZero() {
		_, err := m.Connection.ProductUserSession.UpdateOne(ctx, bson.M{"_id": ps.ID}, bson.M{"$set": ps})
		return err
	}
	ps.CreatedAt = utils.GetCurrentTime()
	id, err := m.Connection.ProductUserSession.InsertOne(ctx, ps)
	if err != nil {
		return err
	}
	ps.ID = id.InsertedID.(primitive.ObjectID)
	return nil
}

func (m *MongoProductRepository) ExistsLocationHash(lc MongoLocation) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := m.Connection.Location.CountDocuments(ctx, bson.M{"hash": lc.Hash})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *MongoProductRepository) ExistsProductUserSessionHash(ps MongoProductUserSession) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := m.Connection.ProductUserSession.CountDocuments(ctx, bson.M{"hash": ps.Hash})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *MongoProductRepository) GetLocationByHash(lc *MongoLocation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Connection.Location.FindOne(ctx, bson.M{"hash": lc.Hash}).Decode(lc)
}

func (m *MongoProductRepository) GetProductUserSessionByHash(ps *MongoProductUserSession) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.Connection.ProductUserSession.FindOne(
		ctx,
		bson.M{"hash": ps.Hash},
		options.FindOne().SetSort(bson.M{"created_at": -1}),
	).Decode(ps)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (m *MongoProductRepository) HashProductUserSession(ps *MongoProductUserSession) error {
	psCopy := MongoProductUserSession{
		ProductID: ps.ProductID,
		IPAddress: ps.IPAddress,
		Location:  ps.Location,
		Lat:       ps.Lat,
		Lon:       ps.Lon,
		UserAgent: ps.UserAgent,
		Proxy:     ps.Proxy,
		Isp:       ps.Isp,
		Device:    ps.Device,
		Os:        ps.Os,
		Browser:   ps.Browser,
		Bot:       ps.Bot,
		Referer:   ps.Referer,
	}
	hash, err := utils.HashStruct(psCopy)
	if err != nil {
		return err
	}
	ps.Hash = hash
	return nil
}
