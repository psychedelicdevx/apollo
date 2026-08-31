package wallet

import (
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(s *WalletSnapshot) error {
	var existing WalletSnapshot
	err := r.db.Where("user_id = ? AND address = ?", s.UserID, s.Address).First(&existing).Error

	if err == nil {
		existing.BalanceWei = s.BalanceWei
		existing.BalanceEth = s.BalanceEth
		existing.RawResponse = s.RawResponse
		if err := r.db.Save(&existing).Error; err != nil {
			return err
		}
		*s = existing
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(s).Error
	}
	return err
}

func (r *Repository) ListByUser(userID string) ([]WalletSnapshot, error) {
	var snaps []WalletSnapshot
	err := r.db.Where("user_id = ?", userID).Order("updated_at DESC").Find(&snaps).Error
	return snaps, err
}

func (r *Repository) GetTokenCache(address string) (*TokenCache, error) {
	var c TokenCache
	err := r.db.Where("address = ?", address).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) SaveTokenCache(c *TokenCache) error {
	var existing TokenCache
	err := r.db.Where("address = ?", c.Address).First(&existing).Error
	if err == nil {

		existing.TotalUSD = c.TotalUSD
		existing.TokensJSON = c.TokensJSON
		existing.FetchedAt = c.FetchedAt
		return r.db.Save(&existing).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(c).Error
	}
	return err
}
