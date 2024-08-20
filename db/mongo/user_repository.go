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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
	Connection *MongoConnection
}

type MongoUser = db_types.User[primitive.ObjectID]
type MongoOTP = db_types.OTP[primitive.ObjectID]

// Inserts a new user into the database
func (m *MongoUserRepository) CreateUser(usr *MongoUser) (primitive.ObjectID, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	usr.CreatedAt = utils.GetCurrentTime()
	result, err := m.Connection.User.InsertOne(ctx, usr)
	if err != nil {
		if m.Connection.IsDuplicateKeyError(err) {
			return primitive.NilObjectID, api_error.NewAPIError("User Already Exists", 409, "User already exists")
		}
		return primitive.NilObjectID, api_error.UnexpectedError(err)
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// Creates a new OTP for the user on the given scope
func (m *MongoUserRepository) CreateOTP(userId primitive.ObjectID, scope string) (string, *api_error.APIError) {
	code, err := utils.GenerateOTP(6)
	if err != nil {
		return "", api_error.UnexpectedError(err)
	}
	otp := MongoOTP{
		UserID:    userId,
		OTP:       code,
		Scope:     scope,
		Verified:  false,
		CreatedAt: utils.GetCurrentTime(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	previousOtp := MongoOTP{}
	opts := options.FindOne().SetSort(bson.D{{"created_at", -1}})
	if err := m.Connection.OTP.FindOne(ctx, bson.M{"user_id": userId, "scope": scope, "verified": false}, opts).Decode(&previousOtp); err == nil {
		if time.Now().UTC().Sub(previousOtp.CreatedAt) < 2*time.Minute {
			return "", api_error.NewAPIError("OTP Already Sent", 400, "Please wait atleast 2 minutes before trying again")
		}
	}

	_, err = m.Connection.OTP.InsertOne(ctx, otp)
	if err != nil {
		return "", api_error.UnexpectedError(err)
	}
	return code, nil
}

// Verifies the OTP for the user
func (m *MongoUserRepository) VerifyOTP(userId primitive.ObjectID, code, scope string) *api_error.APIError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	otp := MongoOTP{}
	m.Connection.OTP.FindOne(ctx, bson.M{"user_id": userId, "otp": code, "scope": scope}).Decode(&otp)
	if otp.ID.IsZero() {
		return api_error.NewAPIError("Invalid OTP", 400, "Invalid OTP")
	}
	if otp.Verified {
		return api_error.NewAPIError("OTP Already Verified", 400, "OTP already used")
	}
	if time.Now().UTC().Sub(otp.CreatedAt) > 10*time.Minute {
		return api_error.NewAPIError("OTP Expired", 400, "OTP expired")
	}
	_, err := m.Connection.OTP.UpdateOne(ctx, bson.M{"_id": otp.ID}, bson.M{"$set": bson.M{"verified": true}})
	if err != nil {
		return api_error.UnexpectedError(err)
	}
	return nil
}

// Fetches the user from the database by email
func (m *MongoUserRepository) GetUserByEmail(email string) (*MongoUser, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user := MongoUser{}
	err := m.Connection.User.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, api_error.NewAPIError("User not found", 404, "User not found")
	}
	return &user, nil
}

// Fetches the user from the database by ID
func (m *MongoUserRepository) GetUserByID(userId primitive.ObjectID) (*MongoUser, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user := MongoUser{}
	err := m.Connection.User.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		return nil, api_error.NewAPIError("User not found", 404, "User not found")
	}
	return &user, nil
}

// Sets the email verification status of the user
func (m *MongoUserRepository) SetEmailVerified(userId primitive.ObjectID, verified bool) *api_error.APIError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.Connection.User.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$set": bson.M{"email_verified": verified}})
	if err != nil {
		return api_error.UnexpectedError(err)
	}
	return nil
}

// Creates a new refresh token for the user
func (m *MongoUserRepository) CreateRefreshToken(userId primitive.ObjectID) (string, *api_error.APIError) {
	token, err := utils.CreateRefreshToken(userId.Hex())
	if err != nil {
		return "", api_error.UnexpectedError(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err2 := m.Connection.User.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$push": bson.M{"refresh_tokens": token}})
	if err2 != nil {
		return "", api_error.UnexpectedError(err)
	}
	return token, nil
}

// Removes the refresh token from the database
func (m *MongoUserRepository) RemoveRefreshToken(userId primitive.ObjectID, token string) *api_error.APIError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.Connection.User.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$pull": bson.M{"refresh_tokens": token}})
	if err != nil {
		return api_error.UnexpectedError(err)
	}
	return nil
}

// Checks if the refresh token is valid
func (m *MongoUserRepository) IsValidRefreshToken(userId primitive.ObjectID, token string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user := MongoUser{}
	err := m.Connection.User.FindOne(ctx, bson.M{"_id": userId, "refresh_tokens": bson.M{
		"$in": []string{token},
	}}).Decode(&user)
	return err == nil
}

// Removes all the old refresh tokens from the database
func (m *MongoUserRepository) RemoveOldTokens(userId primitive.ObjectID) *api_error.APIError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result map[string]interface{}
	err2 := m.Connection.User.FindOne(ctx, bson.M{"_id": userId}, options.FindOne().SetProjection(bson.M{"refresh_tokens": 1})).Decode(&result)
	if err2 != nil {
		return api_error.UnexpectedError(err2)
	}
	if refreshTokens, ok := result["refresh_tokens"].(primitive.A); ok {
		log.Println(len(refreshTokens))
		if len(refreshTokens) > db_types.MAX_USER_SESSIONS {
			_, err2 = m.Connection.User.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$pop": bson.M{"refresh_tokens": -1}})
			if err2 != nil {
				return api_error.UnexpectedError(err2)
			}
		}
	}
	return nil
}
