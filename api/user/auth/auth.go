package auth

import (
	"productanalyzer/api/db"
	api_error "productanalyzer/api/errors"
	mailer "productanalyzer/api/mail"
	"productanalyzer/api/utils"
	response "productanalyzer/api/utils/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(c *gin.Context) {
	var params RegisterParams
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	if passwordHash, err := utils.HashPassword(params.Password); err != nil {
		response.SendFailureResponse(c, err)
		return
	} else {
		params.Password = passwordHash
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
	token, err := utils.CreateToken(userId.Hex())
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	tokenData := TokenData{
		AccessToken: token,
	}
	err = mailer.SendHTMLEmail(params.Email, "Email Verification", "Your OTP is "+otp)
	if err != nil {
		response.SendSuccessResponse(c, "User registered successfully", tokenData, nil)
		return
	}
	response.SendSuccessResponse(c, "User registered successfully", tokenData, nil)
}

func VerifyEmail(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(db.User)
	var params VerifyEmailParams
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	if err := db.VerifyOTP(user.ID, params.OTP, "email_verification"); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	_, err := db.Connection.User.UpdateOne(c, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"email_verified": true}})
	if err != nil {
		response.SendFailureResponse(c, api_error.UnexpectedError(err))
		return
	}
	response.SendSuccessResponse(c, "Email verified successfully", nil, nil)
}

func ResendOTP(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(db.User)
	var params ResendOTPParams
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	if params.Scope == "email_verification" && user.EmailVerified {
		response.SendFailureResponse(c, api_error.NewAPIError("Email already verified", 400, "Email is already verified"))
		return
	}
	otp, err := db.CreateOTP(user.ID, params.Scope)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	err = mailer.SendHTMLEmail(user.Email, "Email Verification", "Your OTP is "+otp)
	if err != nil {
		response.SendFailureResponse(c, api_error.NewAPIError("Couldn't send Email", 500, "We are unable to send email at the moment"))
		return
	}
	response.SendSuccessResponse(c, "OTP sent successfully", nil, nil)
}
