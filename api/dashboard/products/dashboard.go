package products_route

import (
	"log"
	"productanalyzer/api/db"
	db_types "productanalyzer/api/db/types"
	api_error "productanalyzer/api/errors"
	response "productanalyzer/api/utils/response"
	"productanalyzer/api/websockets"
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
	user := usr.(*db.User)
	var params CreateProductRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product := db.Product{
		Name:        params.Name,
		Description: params.Description,
		BaseUrl:     params.BaseUrl,
		ProductID:   params.ProductID,
		UserID:      user.ID,
	}
	productId, err := db.ProductRepository.CreateProduct(&product)
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
	user := usr.(*db.User)
	var params CreateAccessKeyRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product, err := db.ProductRepository.GetProductByProductIDAUserID(params.ProductID, user.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	key, err := db.ProductRepository.CreateProductAccessKey(product.ID, params.Scope)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product.AccessKeys = append(product.AccessKeys, key.ID)
	err = db.ProductRepository.UpdateProduct(*product)
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
	user := usr.(*db.User)
	var params ProductInfoRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product, err := db.ProductRepository.GetProductByProductIDAUserID(params.ProductID, user.ID)
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

// Product Access Keys Request [GET]
func ProductAccessKeys(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*db.User)
	var params ProductAccessKeysRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product, err := db.ProductRepository.GetProductByProductIDAUserID(params.ProductID, user.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	accessKeys, err := db.ProductRepository.GetProductAccessKeys(product.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	keys := []ProductAccessKeyResponse{}
	for _, key := range *accessKeys {
		keys = append(keys, ProductAccessKeyResponse{
			AccessKey: key.AccessKey,
			Scope:     key.Scope,
			CreatedAt: key.CreatedAt.UTC().String(),
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
	user := usr.(*db.User)
	var params DeleteProductRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	if params.Type == db_types.DELETION_REQUEST_TYPE_INITIAL {
		product, err := db.ProductRepository.GetProductByProductIDAUserID(params.InstanceId, user.ID)
		if err != nil {
			response.SendFailureResponse(c, err)
			return
		}
		id, err := db.DeletionListRepository.AddToDeletionList(product.ID, db_types.DELETION_TYPE_PRODUCT)
		if err != nil {
			response.SendFailureResponse(c, err)
			return
		}
		response.SendSuccessResponse(c, "Product deletion request initiated", bson.M{"instance_id": id.Hex()}, nil)
	} else if params.Type == db_types.DELETION_REQUEST_TYPE_CONFIRM {
		objectId, err2 := primitive.ObjectIDFromHex(params.InstanceId)
		if err2 != nil {
			response.SendFailureResponse(c, api_error.NewAPIError("Invalid Instance ID", 400, "The given instance id is invalid"))
			return
		}
		deletion, err := db.DeletionListRepository.GetFromDeletionList(objectId)
		if err != nil {
			response.SendFailureResponse(c, err)
			return
		}
		if deletion.Type != db_types.DELETION_TYPE_PRODUCT {
			response.SendFailureResponse(c, api_error.NewAPIError("Invalid Instance ID", 400, "The given instance id is invalid"))
			return
		}
		log.Print(deletion.ObjectID)
		product, err := db.ProductRepository.GetProductByID(deletion.ObjectID)
		if err != nil {
			response.SendFailureResponse(c, err)
			return
		}
		err = db.ProductRepository.DeleteProduct(product.ID)
		if err != nil {
			response.SendFailureResponse(c, err)
			return
		}
		response.SendSuccessResponse(c, "Product deleted successfully", nil, nil)
	} else {
		response.SendFailureResponse(c, api_error.NewAPIError("Invalid Request Type", 400, "The given request type is invalid"))
	}
}

// Get Log of a product, filtered by date range and unit(month, day ..) [GET]
func VisitLog(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*db.User)
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
	product, err := db.ProductRepository.GetProductByProductIDAUserID(params.ProductID, user.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	visitLogs, err := db.ProductRepository.GetVisitLogs(product.ID, params.FromDate, params.ToDate)
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
			log.CreatedAt = log.CreatedAt.Truncate(time.Minute)
		} else if params.Unit == VISIT_LOG_ENTRY_UNIT_HOUR {
			log.CreatedAt = log.CreatedAt.Truncate(time.Hour)
		} else if params.Unit == VISIT_LOG_ENTRY_UNIT_DAY {
			log.CreatedAt = log.CreatedAt.Truncate(24 * time.Hour)
		} else if params.Unit == VISIT_LOG_ENTRY_UNIT_MONTH {
			log.CreatedAt = time.Date(log.CreatedAt.Year(), log.CreatedAt.Month(), 1, 0, 0, 0, 0, time.UTC)
		}
		if len(visits) > 0 && visits[len(visits)-1].CreatedAt == log.CreatedAt.UTC().String() {
			visits[len(visits)-1].SessionCount += int64(log.ActivityCount)
			continue
		}
		visits = append(visits, VisitLogEntry{
			CreatedAt:     log.CreatedAt.UTC().String(),
			UpdatedAt:     log.UpdatedAt.UTC().String(),
			ActivityCount: int64(log.ActivityCount),
			SessionCount:  1,
			Referer:       log.Referer,
		})
	}
	response.SendSuccessResponse(c, "Visit logs", VisitLogResponse{
		Unit:            params.Unit,
		Logs:            visits,
		TotalSessions:   int64(len(*visitLogs)),
		TotalActivities: int64(totalActivities),
	}, nil)
}

// Websocket for log connection [GET]
func HandleLogWebsocket(c *gin.Context) {
	conn, exists := c.Get("websocket")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	u, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := u.(*db.User)
	ws := conn.(*websockets.WebsocketConnection)
	productId := c.Param("product_id")
	product, err := db.ProductRepository.GetProductByProductIDAUserID(productId, user.ID)
	if err != nil {
		ws.SendErrorMesage(err.Message)
		ws.Close()
		return
	}
	err2 := websockets.OpenLogConnection(ws, product.ID.Hex())
	if err2 != nil {
		ws.SendErrorMesage(err2.Error())
		ws.Close()
		return
	}
	ws.SendData(gin.H{"message": "Connection established"})
}

// List all products [GET]
func ListProducts(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*db.User)
	products, err := db.ProductRepository.GetProductsByUserID(user.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	productsResponse := []ProductInfoResponse{}
	for _, product := range *products {
		productsResponse = append(productsResponse, ProductInfoResponse{
			ID:          product.ID.Hex(),
			Name:        product.Name,
			Description: product.Description,
			BaseUrl:     product.BaseUrl,
			ProductID:   product.ProductID,
		})
	}
	response.SendSuccessResponse(c, "Products", productsResponse, nil)
}
