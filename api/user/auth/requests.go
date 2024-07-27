package auth

type RegisterParams struct {
	Fullname string `form:"fullname" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type VerifyEmailParams struct {
	OTP string `form:"otp" binding:"required"`
}

type ResendOTPParams struct {
	Scope string `form:"scope" binding:"required"`
}

type LoginParams struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}
