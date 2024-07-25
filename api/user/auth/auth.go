package auth

import (
	"productanalyzer/api/db"
	response "productanalyzer/api/utils/response"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var params RegisterParams
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, "Invalid Request", err, nil)
		return
	}
	users := db.Connection.Collection("users")
	_, err := users.InsertOne(c, gin.H{
		"email":    params.Email,
		"password": params.Password,
	})
	if err != nil {
		if db.IsDuplicateKeyError(err) {
			response.SendFailureResponse(c, "User Already Exists", db.DuplicateError(), nil)
			return
		}
		response.SendFailureResponse(c, "Failed to Register user", err, nil)
		return
	}
	response.SendSuccessResponse(c, "User registered successfully", nil, nil)
}
