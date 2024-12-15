package entities

import "time"

type OTP struct {
	Id string    `json:"id"`
	UserId string `json:"user_id"`
	Mobile string `json:"mobile"`
	OTP string `json:"otp"`
	Expiration time.Time `json:"expiration"`
	CreatedAt time.Time `json:"created_at"`
}

type OtpRequest struct {
	Mobile string `json:"mobile"`
}

type VerifyOtpRequest struct {
	Mobile string `json:"mobile"`
	OTP    string `json:"otp"`
}