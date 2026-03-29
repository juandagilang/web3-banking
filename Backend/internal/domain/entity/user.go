package entity

import "time"

type User struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	WalletAddress string    `json:"wallet_address" gorm:"uniqueIndex;not null;size:42"`
	Nonce         string    `json:"nonce" gorm:"not null;size:64"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
