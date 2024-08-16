package db_types

import "time"

const (
	OTP_SCOPE_EMAIL_VERIFICATION = "email_verification"
)

type User[T PrimaryKey] struct {
	ID            T         `bson:"_id,omitempty"`  // primary key
	Fullname      string    `bson:"fullname"`       // fullname of the user
	Email         string    `bson:"email"`          // email of the user
	EmailVerified bool      `bson:"email_verified"` // email verification status of the user
	Password      string    `bson:"password"`       // password of the user
	CreatedAt     time.Time `bson:"created_at"`     // time at which the user was created
}

type OTP[T PrimaryKey] struct {
	ID        T         `bson:"_id,omitempty"` // primary key
	UserID    T         `bson:"user_id"`       // Id of user collection
	OTP       string    `bson:"otp"`           // OTP generated for the user
	Scope     string    `bson:"scope"`         // Scope of the OTP
	Verified  bool      `bson:"verified"`      // Verification status of the OTP, whether it is used or not
	CreatedAt time.Time `bson:"created_at"`    // time at which the OTP was created
}

type UserPlan[T PrimaryKey] struct {
	ID         T         `bson:"_id,omitempty"`
	UserID     T         `bson:"user_id"`
	PlanID     T         `bson:"plan_id"`
	QuotaUsage Quota     `bson:"quota"`
	StartDate  time.Time `bson:"start_date"`
	EndDate    time.Time `bson:"end_date"`
}

/*
	INDIRECT TYPES
*/

type Quota struct {
	Hits     int `json:"hits"`
	Products int `json:"products"`
}
