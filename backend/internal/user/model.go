package user

import "time"

type User struct {
	ID           string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"not null;default:user"`
	DailyLimit   int    `gorm:"not null;default:10"`
	SearchesUsed int    `gorm:"not null;default:0"`
	SearchesDate string `gorm:"not null;default:''"`
	CreatedAt    time.Time
}
