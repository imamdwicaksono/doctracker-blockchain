package models

type OtpRequest struct {
	Email string `json:"email"`
}

type VerifyOtpRequest struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}
