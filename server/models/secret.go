package models

import "time"

type SecretInfo struct {
	Name        string `gorm:"primaryKey"`
	Value       string // This should be encrypted
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}