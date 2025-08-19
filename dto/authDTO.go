package dto

type LoginRequest struct {
	Phone string `json:"phone" example:"09120000000"`
	OTP   string `json:"otp" example:"123456"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RequestOTPRequest struct {
	Phone string `json:"phone" example:"09120000000"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}
