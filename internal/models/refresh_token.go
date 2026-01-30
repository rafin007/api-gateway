package models

import "time"

type RefreshToken struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	TokenHash  string    `json:"token_hash"`
	DeviceInfo string    `json:"device_info"`
	IPAddress  string    `json:"ip_address"`
	ExpiresAt  time.Time `json:"expires_at"`
	Revoked    bool      `json:"revoked"`
	CreatedAt  time.Time `json:"created_at"`
}
