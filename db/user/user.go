package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	OTP_SCOPE_EMAIL_VERIFICATION = "email_verification"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`  // primary key
	Fullname      string             `bson:"fullname"`       // fullname of the user
	Email         string             `bson:"email"`          // email of the user
	EmailVerified bool               `bson:"email_verified"` // email verification status of the user
	Password      string             `bson:"password"`       // password of the user
	CreatedAt     primitive.DateTime `bson:"created_at"`     // time at which the user was created
}

type OTP struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // primary key
	UserID    primitive.ObjectID `bson:"user_id"`       // Id of user collection
	OTP       string             `bson:"otp"`           // OTP generated for the user
	Scope     string             `bson:"scope"`         // Scope of the OTP
	Verified  bool               `bson:"verified"`      // Verification status of the OTP, whether it is used or not
	CreatedAt primitive.DateTime `bson:"created_at"`    // time at which the OTP was created
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
