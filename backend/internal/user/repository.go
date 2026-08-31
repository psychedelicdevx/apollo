package user

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(u *User) error {
	return r.db.Create(u).Error
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	var u User

	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FindByID(id string) (*User, error) {
	var u User
	err := r.db.First(&u, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) ConsumeSearch(userID string) error {
	var u User
	if err := r.db.First(&u, "id = ?", userID).Error; err != nil {
		return err
	}

	today := time.Now().UTC().Format("2006-01-02")
	if u.SearchesDate != today {
		u.SearchesUsed = 0
		u.SearchesDate = today
	}

	if u.SearchesUsed >= u.DailyLimit {
		return ErrQuotaExceeded
	}

	u.SearchesUsed++
	return r.db.Save(&u).Error
}
