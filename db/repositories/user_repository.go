package db_repositories

import (
	db_types "productanalyzer/api/db/types"
	api_error "productanalyzer/api/errors"
)

type UserRepository[T db_types.PrimaryKey] interface {
	CreateUser(user *db_types.User[T]) (T, *api_error.APIError)           // Inserts a new user into the database
	CreateOTP(userId T, scope string) (string, *api_error.APIError)       // Creates a new OTP for the user on the given scope
	VerifyOTP(userId T, code, scope string) *api_error.APIError           // Verifies the OTP
	GetUserByEmail(email string) (*db_types.User[T], *api_error.APIError) // Gets the user by email
	GetUserByID(userId T) (*db_types.User[T], *api_error.APIError)        // Gets the user by ID
	SetEmailVerified(userId T, verified bool) *api_error.APIError         // Sets the email verification status of the user
	CreateRefreshToken(userId T) (string, *api_error.APIError)            // Creates a new refresh token for the user
	RemoveRefreshToken(userId T, token string) *api_error.APIError        // Removes the refresh token from the database
	RemoveAllRefreshTokens(userId T) *api_error.APIError                  // Removes all refresh tokens from the database
	IsValidRefreshToken(userId T, token string) bool                      // Checks if the refresh token is valid
	RemoveOldTokens(userId T) *api_error.APIError                         // Removes old refresh tokens, and only maintain 5 latest refresh tokens
}
