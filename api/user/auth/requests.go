package auth_route

type RegisterParams struct {
	Fullname string `form:"fullname" binding:"required" json:"fullname"`
	Email    string `form:"email" binding:"required,email" json:"email"`
	Password string `form:"password" binding:"required" json:"password"`
}

type VerifyEmailParams struct {
	OTP string `form:"otp" binding:"required" json:"otp"`
}

type ResendOTPParams struct {
	Scope string `form:"scope" binding:"required" json:"scope"`
}

type LoginParams struct {
	Email    string `form:"email" binding:"required,email" json:"email"`
	Password string `form:"password" binding:"required" json:"password"`
}

type GetAccessTokenParams struct {
	RefreshToken string `form:"refresh_token" binding:"required" json:"refresh_token"`
}

type LogoutParams struct {
	RefreshToken string `form:"refresh_token" binding:"required" json:"refresh_token"`
	All          bool   `form:"all" json:"all" default:"false"`
}
