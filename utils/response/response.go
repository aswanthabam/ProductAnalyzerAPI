package repsonse

import (
	"fmt"

	utils "productanalyzer/api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendFailureResponse(cn *gin.Context, message string, err error, statusCode *int) {
	if statusCode == nil {
		statusCode = new(int)
		*statusCode = 400
	}
	if err == nil {
		err = fmt.Errorf("")
	}
	var data interface{}
	if _, ok := err.(validator.ValidationErrors); ok {
		data = utils.FormatValidationErrors(err)
	} else {
		data = gin.H{"message": err.Error()}
	}
	cn.JSON(*statusCode, gin.H{
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
		"message": message,
		"data":    data,
	})
}
