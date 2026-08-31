package label

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByAddress(address string) (*Label, error) {
	var l Label
	if err := r.db.First(&l, "address = ?", address).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *Repository) Count() (int64, error) {
	var n int64
	err := r.db.Model(&Label{}).Count(&n).Error
	return n, err
}

func (r *Repository) CreateMany(labels []Label) error {
	return r.db.Create(&labels).Error
}
