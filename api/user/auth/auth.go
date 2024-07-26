package auth

import (
	"productanalyzer/api/db"
	mailer "productanalyzer/api/mail"
	response "productanalyzer/api/utils/response"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var params RegisterParams
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	user := db.User{
		Fullname:      params.Fullname,
		Email:         params.Email,
		Password:      params.Password,
		EmailVerified: false,
	}
	userId, err := db.InsertUser(&user)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	otp, err := db.CreateOTP(userId, "email_verification")
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	err = mailer.SendHTMLEmail(params.Email, "Email Verification", "Your OTP is "+otp)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	response.SendSuccessResponse(c, "User registered successfully", nil, nil)
}
