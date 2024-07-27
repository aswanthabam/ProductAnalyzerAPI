package auth_route

import (
	user_db "productanalyzer/api/db/user"
	api_error "productanalyzer/api/errors"
	mailer "productanalyzer/api/mail"
	"productanalyzer/api/utils"
	response "productanalyzer/api/utils/response"

	"github.com/gin-gonic/gin"
)

// Register User Endpoint [POST]
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
	user := user_db.User{
		Fullname:      params.Fullname,
		Email:         params.Email,
		Password:      params.Password,
		EmailVerified: false,
	}
	userId, err := user_db.InsertUser(&user)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	otp, err := user_db.CreateOTP(userId, user_db.OTP_SCOPE_EMAIL_VERIFICATION)
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

// Verify Email Endpoint [POST]
func VerifyEmail(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*user_db.User)
	var params VerifyEmailParams
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	if err := user_db.VerifyOTP(user.ID, params.OTP, user_db.OTP_SCOPE_EMAIL_VERIFICATION); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	err := user_db.SetEmailVerified(user.ID, true)
	if err != nil {
		response.SendFailureResponse(c, api_error.UnexpectedError(err))
		return
	}
	response.SendSuccessResponse(c, "Email verified successfully", nil, nil)
}

// Resend OTP Endpoint [POST]
func ResendOTP(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*user_db.User)
	var params ResendOTPParams
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	if params.Scope == user_db.OTP_SCOPE_EMAIL_VERIFICATION && user.EmailVerified {
		response.SendFailureResponse(c, api_error.NewAPIError("Email already verified", 400, "Email is already verified"))
		return
	}
	otp, err := user_db.CreateOTP(user.ID, params.Scope)
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

// Login User Endpoint [POST]
func Login(c *gin.Context) {
	var params LoginParams
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	user, err := user_db.GetUserByEmail(params.Email)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	if !utils.VerifyPassword(user.Password, params.Password) {
		response.SendFailureResponse(c, api_error.NewAPIError("Invalid credentials", 400, "Invalid email or password"))
		return
	}
	message := "Login successful"
	if !user.EmailVerified {
		message = "Login successful, Email not verified"
		return
	}
	token, err := utils.CreateToken(user.ID.Hex())
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	tokenData := TokenData{
		AccessToken: token,
	}
	response.SendSuccessResponse(c, message, tokenData, nil)
}
