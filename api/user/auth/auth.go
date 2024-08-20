package auth_route

import (
	"productanalyzer/api/db"
	db_types "productanalyzer/api/db/types"
	api_error "productanalyzer/api/errors"
	mailer "productanalyzer/api/mail"
	"productanalyzer/api/utils"
	response "productanalyzer/api/utils/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	user := db.User{
		Fullname:      params.Fullname,
		Email:         params.Email,
		Password:      params.Password,
		EmailVerified: false,
	}
	userId, err := db.UserRepository.CreateUser(&user)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	otp, err := db.UserRepository.CreateOTP(userId, db_types.OTP_SCOPE_EMAIL_VERIFICATION)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	token, err := utils.CreateToken(userId.Hex())
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	refreshToken, err := db.UserRepository.CreateRefreshToken(userId)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	db.UserRepository.RemoveOldTokens(user.ID)
	tokenData := TokenData{
		AccessToken:  token,
		RefreshToken: refreshToken,
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
	user := usr.(*db.User)
	var params VerifyEmailParams
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	if err := db.UserRepository.VerifyOTP(user.ID, params.OTP, db_types.OTP_SCOPE_EMAIL_VERIFICATION); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	err := db.UserRepository.SetEmailVerified(user.ID, true)
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
	user := usr.(*db.User)
	var params ResendOTPParams
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	if params.Scope == db_types.OTP_SCOPE_EMAIL_VERIFICATION && user.EmailVerified {
		response.SendFailureResponse(c, api_error.NewAPIError("Email already verified", 400, "Email is already verified"))
		return
	}
	otp, err := db.UserRepository.CreateOTP(user.ID, params.Scope)
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
	user, err := db.UserRepository.GetUserByEmail(params.Email)
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
	}
	token, err := utils.CreateToken(user.ID.Hex())
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	refreshToken, err := db.UserRepository.CreateRefreshToken(user.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	db.UserRepository.RemoveOldTokens(user.ID)
	tokenData := TokenData{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}
	response.SendSuccessResponse(c, message, tokenData, nil)
}

func GetAccessToken(c *gin.Context) {
	var params GetAccessTokenParams
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	claims, tokenParseError := utils.ValidateToken(params.RefreshToken)
	if claims != nil {
		userObjectId, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			response.SendFailureResponse(c, api_error.NewAPIError("Invalid Token", 400, "Invalid User ID"))
			return
		}
		valid := db.UserRepository.IsValidRefreshToken(userObjectId, params.RefreshToken)
		if !valid {
			response.SendFailureResponse(c, api_error.NewAPIError("Invalid Token", 400, "Invalid refresh token"))
			return
		}
		if tokenParseError != nil {
			db.UserRepository.RemoveRefreshToken(userObjectId, params.RefreshToken)
			response.SendFailureResponse(c, tokenParseError)
			return
		}
	} else {
		response.SendFailureResponse(c, tokenParseError)
		return
	}
	if claims.TokenType != "refresh" {
		response.SendFailureResponse(c, api_error.NewAPIError("Invalid Token", 400, "Invalid token"))
		return
	}
	token, err := utils.CreateToken(claims.UserID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	tokenData := TokenData{
		AccessToken:  token,
		RefreshToken: params.RefreshToken,
	}
	response.SendSuccessResponse(c, "Access Token generated successfully", tokenData, nil)
}
