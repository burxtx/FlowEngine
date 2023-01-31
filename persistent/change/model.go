package change

import (
	"gorm.io/gorm"
)

type Repository struct {
	*gorm.DB
}

func (r *Repository) Save(n Change) (int64, error) {
	err := r.Create(&n).Error
	if err != nil {
		return 0, err
	}
	return n.ID, nil
}

func (r *Repository) Update(n Change) error {
	err := r.Updates(&n).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Get(id int64) (*Change, error) {
	n := Change{}
	err := r.First(&n, id).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *Repository) List(query Change) ([]*Change, error) {
	var n []*Change
	err := r.Where(&query).Find(&n).Error
	if err != nil {
		return nil, err
	}
	return n, nil
}
