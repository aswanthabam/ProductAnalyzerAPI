package profile_route

type UserInfoResponse struct {
	ID            string `json:"id"`
	Fullname      string `json:"fullname"`
	EmailVerified bool   `json:"email_verified"`
}
