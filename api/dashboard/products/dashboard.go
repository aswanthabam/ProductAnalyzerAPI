package products_route

import (
	"log"
	deletion_db "productanalyzer/api/db/deletion"
	products_db "productanalyzer/api/db/products"
	user_db "productanalyzer/api/db/user"
	api_error "productanalyzer/api/errors"
	response "productanalyzer/api/utils/response"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create Product Request [POST]
func CreateProduct(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*user_db.User)
	var params CreateProductRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product := products_db.Product{
		Name:        params.Name,
		Description: params.Description,
		BaseUrl:     params.BaseUrl,
		ProductID:   params.ProductID,
		UserID:      user.ID,
	}
	productId, err := products_db.CreateProduct(&product)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	response.SendSuccessResponse(c, "Product created successfully", CreateProductResponse{
		ProductId: productId.Hex(),
	}, nil)
}

// Create Access Key Request [POST]
func CreateAccessKey(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*user_db.User)
	var params CreateAccessKeyRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product, err := products_db.GetProductByProductIDAUserID(params.ProductID, user.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	key, err := products_db.CreateProductAccessKey(product.ID, params.Scope)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product.AccessKeys = append(product.AccessKeys, key.ID)
	err = products_db.UpdateProduct(*product)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	response.SendSuccessResponse(c, "Access Key created successfully", CreateAccessKeyResponse{
		AccessKey: key.AccessKey,
		ProductID: product.ProductID,
		Scope:     params.Scope,
	}, nil)
}

// Product Info Request [POST]
func ProductInfo(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*user_db.User)
	var params ProductInfoRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product, err := products_db.GetProductByProductIDAUserID(params.ProductID, user.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	response.SendSuccessResponse(c, "Product information", ProductInfoResponse{
		ID:          product.ID.Hex(),
		Name:        product.Name,
		Description: product.Description,
		BaseUrl:     product.BaseUrl,
		ProductID:   product.ProductID,
	}, nil)
}

func ProductAccessKeys(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*user_db.User)
	var params ProductAccessKeysRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product, err := products_db.GetProductByProductIDAUserID(params.ProductID, user.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	accessKeys, err := products_db.GetProductAccessKeys(product.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	keys := []ProductAccessKeyResponse{}
	for _, key := range *accessKeys {
		keys = append(keys, ProductAccessKeyResponse{
			AccessKey: key.AccessKey,
			Scope:     key.Scope,
			CreatedAt: key.CreatedAt.Time().UTC().String(),
		})
	}
	response.SendSuccessResponse(c, "Product access keys", keys, nil)
}

// Delete Product Request [POST]
func DeleteProduct(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*user_db.User)
	var params DeleteProductRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	if params.Type == deletion_db.DELETION_REQUEST_TYPE_INITIAL {
		product, err := products_db.GetProductByProductIDAUserID(params.InstanceId, user.ID)
		if err != nil {
			response.SendFailureResponse(c, err)
			return
		}
		id, err := deletion_db.AddToDeletionList(product.ID, deletion_db.DELETION_TYPE_PRODUCT)
		if err != nil {
			response.SendFailureResponse(c, err)
			return
		}
		response.SendSuccessResponse(c, "Product deletion request initiated", bson.M{"instance_id": id.Hex()}, nil)
	} else if params.Type == deletion_db.DELETION_REQUEST_TYPE_CONFIRM {
		objectId, err2 := primitive.ObjectIDFromHex(params.InstanceId)
		if err2 != nil {
			response.SendFailureResponse(c, api_error.NewAPIError("Invalid Instance ID", 400, "The given instance id is invalid"))
			return
		}
		deletion, err := deletion_db.GetFromDeletionList(objectId)
		if err != nil {
			response.SendFailureResponse(c, err)
			return
		}
		if deletion.Type != deletion_db.DELETION_TYPE_PRODUCT {
			response.SendFailureResponse(c, api_error.NewAPIError("Invalid Instance ID", 400, "The given instance id is invalid"))
			return
		}
		log.Print(deletion.ObjectID)
		product, err := products_db.GetProductByID(deletion.ObjectID)
		if err != nil {
			response.SendFailureResponse(c, err)
			return
		}
		err = products_db.DeleteProduct(product.ID)
		if err != nil {
			response.SendFailureResponse(c, err)
			return
		}
		response.SendSuccessResponse(c, "Product deleted successfully", nil, nil)
	} else {
		response.SendFailureResponse(c, api_error.NewAPIError("Invalid Request Type", 400, "The given request type is invalid"))
	}
}

func VisitLog(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*user_db.User)
	var params VisitLogRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	if params.Unit != VISIT_LOG_ENTRY_UNIT_MINUTE && params.Unit != VISIT_LOG_ENTRY_UNIT_HOUR && params.Unit != VISIT_LOG_ENTRY_UNIT_DAY && params.Unit != VISIT_LOG_ENTRY_UNIT_MONTH {
		response.SendFailureResponse(c, api_error.NewAPIError("Invalid Unit", 400, "The given unit is invalid"))
		return
	}
	if params.ToDate.IsZero() {
		params.ToDate = time.Now().UTC()
	}
	product, err := products_db.GetProductByProductIDAUserID(params.ProductID, user.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	visitLogs, err := products_db.GetVisitLogs(product.ID, params.FromDate, params.ToDate)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	visits := []VisitLogEntry{}
	totalActivities := 0
	for _, log := range *visitLogs {
		totalActivities += log.ActivityCount
		if log.ActivityCount == 0 {
			continue
		}
		if params.Unit == VISIT_LOG_ENTRY_UNIT_MINUTE {
			log.CreatedAt = primitive.NewDateTimeFromTime(log.CreatedAt.Time().Truncate(time.Minute))
		} else if params.Unit == VISIT_LOG_ENTRY_UNIT_HOUR {
			log.CreatedAt = primitive.NewDateTimeFromTime(log.CreatedAt.Time().Truncate(time.Hour))
		} else if params.Unit == VISIT_LOG_ENTRY_UNIT_DAY {
			log.CreatedAt = primitive.NewDateTimeFromTime(log.CreatedAt.Time().Truncate(24 * time.Hour))
		} else if params.Unit == VISIT_LOG_ENTRY_UNIT_MONTH {
			log.CreatedAt = primitive.NewDateTimeFromTime(time.Date(log.CreatedAt.Time().Year(), log.CreatedAt.Time().Month(), 1, 0, 0, 0, 0, time.UTC))
		}
		if len(visits) > 0 && visits[len(visits)-1].CreatedAt == log.CreatedAt.Time().UTC().String() {
			visits[len(visits)-1].SessionCount += int64(log.ActivityCount)
			continue
		}
		visits = append(visits, VisitLogEntry{
			Type:          params.Unit,
			CreatedAt:     log.CreatedAt.Time().UTC().String(),
			UpdatedAt:     log.UpdatedAt.Time().UTC().String(),
			ActivityCount: int64(log.ActivityCount),
			SessionCount:  1,
			Referer:       log.Referer,
		})
	}
	response.SendSuccessResponse(c, "Visit logs", VisitLogResponse{
		Logs:            visits,
		TotalSessions:   int64(len(*visitLogs)),
		TotalActivities: int64(totalActivities),
	}, nil)
}
