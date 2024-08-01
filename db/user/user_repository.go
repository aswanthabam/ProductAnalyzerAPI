package db

import (
	"context"
	"productanalyzer/api/db"
	api_error "productanalyzer/api/errors"
	"productanalyzer/api/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Inserts a new user into the database
func InsertUser(usr *User) (primitive.ObjectID, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	usr.CreatedAt = utils.GetCurrentTime()
	result, err := db.Connection.User.InsertOne(ctx, usr)
	if err != nil {
		if db.IsDuplicateKeyError(err) {
			return primitive.NilObjectID, api_error.NewAPIError("User Already Exists", 409, "User already exists")
		}
		return primitive.NilObjectID, api_error.UnexpectedError(err)
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// Creates a new OTP for the user on the given scope
func CreateOTP(userId primitive.ObjectID, scope string) (string, *api_error.APIError) {
	code, err := utils.GenerateOTP(6)
	if err != nil {
		return "", api_error.UnexpectedError(err)
	}
	otp := OTP{
		UserID:    userId,
		OTP:       code,
		Scope:     scope,
		Verified:  false,
		CreatedAt: utils.GetCurrentTime(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	previousOtp := OTP{}
	opts := options.FindOne().SetSort(bson.D{{"created_at", -1}})
	if err := db.Connection.OTP.FindOne(ctx, bson.M{"user_id": userId, "scope": scope, "verified": false}, opts).Decode(&previousOtp); err == nil {
		if time.Now().UTC().Sub(previousOtp.CreatedAt.Time()) < 2*time.Minute {
			return "", api_error.NewAPIError("OTP Already Sent", 400, "Please wait atleast 2 minutes before trying again")
		}
	}

	_, err = db.Connection.OTP.InsertOne(ctx, otp)
	if err != nil {
		return "", api_error.UnexpectedError(err)
	}
	return code, nil
}

// Verifies the OTP for the user
func VerifyOTP(userId primitive.ObjectID, code, scope string) *api_error.APIError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	otp := OTP{}
	db.Connection.OTP.FindOne(ctx, bson.M{"user_id": userId, "otp": code, "scope": scope}).Decode(&otp)
	if otp.ID.IsZero() {
		return api_error.NewAPIError("Invalid OTP", 400, "Invalid OTP")
	}
	if otp.Verified {
		return api_error.NewAPIError("OTP Already Verified", 400, "OTP already used")
	}
	if time.Now().UTC().Sub(otp.CreatedAt.Time()) > 10*time.Minute {
		return api_error.NewAPIError("OTP Expired", 400, "OTP expired")
	}
	_, err := db.Connection.OTP.UpdateOne(ctx, bson.M{"_id": otp.ID}, bson.M{"$set": bson.M{"verified": true}})
	if err != nil {
		return api_error.UnexpectedError(err)
	}
	return nil
}

// Fetches the user from the database by email
func GetUserByEmail(email string) (*User, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user := User{}
	err := db.Connection.User.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, api_error.NewAPIError("User not found", 404, "User not found")
	}
	return &user, nil
}

// Fetches the user from the database by ID
func GetUserByID(userId primitive.ObjectID) (*User, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user := User{}
	err := db.Connection.User.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		return nil, api_error.NewAPIError("User not found", 404, "User not found")
	}
	return &user, nil
}

// Sets the email verification status of the user
func SetEmailVerified(userId primitive.ObjectID, verified bool) *api_error.APIError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.Connection.User.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$set": bson.M{"email_verified": verified}})
	if err != nil {
		return api_error.UnexpectedError(err)
	}
	return nil
}
