package profile_route

import (
	"productanalyzer/api/db"
	api_error "productanalyzer/api/errors"
	response "productanalyzer/api/utils/response"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*db.User)
	response.SendSuccessResponse(c, "User information fetched successfully", UserInfoResponse{
		Fullname:      user.Fullname,
		EmailVerified: user.EmailVerified,
		ID:            user.ID.Hex(),
	}, nil)
}
