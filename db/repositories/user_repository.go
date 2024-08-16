package db_repositories

import (
	db_types "productanalyzer/api/db/types"
	api_error "productanalyzer/api/errors"
)

type UserRepository[T db_types.PrimaryKey] interface {
	CreateUser(user *db_types.User[T]) (T, *api_error.APIError)
	CreateOTP(userId T, scope string) (string, *api_error.APIError)
	VerifyOTP(userId T, code, scope string) *api_error.APIError
	GetUserByEmail(email string) (*db_types.User[T], *api_error.APIError)
	GetUserByID(userId T) (*db_types.User[T], *api_error.APIError)
	SetEmailVerified(userId T, verified bool) *api_error.APIError
}
