package db

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	api_error "productanalyzer/api/errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Fullname      string             `bson:"fullname"`
	Email         string             `bson:"email"`
	EmailVerified bool               `bson:"email_verified"`
	Password      string             `bson:"password"`
	CreatedAt     primitive.DateTime `bson:"created_at"`
}

type OTP struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	OTP       string             `bson:"otp"`
	Scope     string             `bson:"scope"`
	Verified  bool               `bson:"verified"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

type UserPlan struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
	PlanID     primitive.ObjectID `bson:"plan_id"`
	QuotaUsage Quota              `bson:"quota"`
	StartDate  primitive.DateTime `bson:"start_date"`
	EndDate    primitive.DateTime `bson:"end_date"`
}

/*
	INDIRECT TYPES
*/

type Quota struct {
	Hits     int `json:"hits"`
	Products int `json:"products"`
}

func InsertUser(usr *User) (primitive.ObjectID, *api_error.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	usr.CreatedAt = primitive.NewDateTimeFromTime(time.Now().UTC())
	result, err := Connection.User.InsertOne(ctx, usr)
	if err != nil {
		if IsDuplicateKeyError(err) {
			return primitive.NilObjectID, api_error.NewAPIError("User Already Exists", 409, "User already exists")
		}
		return primitive.NilObjectID, api_error.UnexpectedError(err)
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func CreateOTP(userId primitive.ObjectID, scope string) (string, *api_error.APIError) {
	code, err := generateOTP(6)
	if err != nil {
		return "", api_error.UnexpectedError(err)
	}
	otp := OTP{
		UserID:    userId,
		OTP:       code,
		Scope:     scope,
		Verified:  false,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now().UTC()),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = Connection.OTP.InsertOne(ctx, otp)
	if err != nil {
		return "", api_error.UnexpectedError(err)
	}
	return code, nil
}

func generateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %v", err)
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}
