package repsonse

import (
	api_error "productanalyzer/api/errors"
	"productanalyzer/api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendFailureResponse(cn *gin.Context, err error) {
	var statusCode int
	var message string
	if err == nil {
		err = api_error.UnexpectedError(nil)
	}
	var data interface{}
	switch err := err.(type) {
	case *api_error.APIError:
		data = gin.H{"message": err.Message}
		statusCode = err.Code
		message = err.Title
	case validator.ValidationErrors:
		data = utils.FormatValidationErrors(err)
		message = "Validation Error"
		statusCode = 400
	default:
		data = gin.H{"message": err.Error()}
		statusCode = 500
		message = "Internal Server Error"
	}
	cn.JSON(statusCode, gin.H{
		"status":  "failed",
		"message": message,
		"data":    data,
	})

}

func SendSuccessResponse(cn *gin.Context, message string, data interface{}, statusCode *int) {
	if statusCode == nil {
		statusCode = new(int)
		*statusCode = 200
	}
	if data == nil {
		data = gin.H{}
	}
	cn.JSON(*statusCode, gin.H{
		"status":  "success",
		"data":    data,
		"message": message,
	})
}
