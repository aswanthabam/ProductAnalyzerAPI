package db_repositories

import (
	db_types "productanalyzer/api/db/types"
	api_error "productanalyzer/api/errors"
)

type DeletionListRepository[T db_types.PrimaryKey] interface {
	AddToDeletionList(objectID T, objectType string) (T, *api_error.APIError)
	GetFromDeletionList(objectId T) (*db_types.DeletionList[T], *api_error.APIError)
}
