package auth

type RegisterParams struct {
	Fullname string `form:"fullname" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}
