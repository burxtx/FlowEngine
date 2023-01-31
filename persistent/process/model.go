package process

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	*gorm.DB
}

func (r *Repository) Save(p Process) (int64, error) {
	p.CreateTime = time.Now().Unix()
	p.UpdateTime = time.Now().Unix()
	err := r.Create(&p).Error
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (r *Repository) Update(p Process) error {
	p.UpdateTime = time.Now().Unix()
	err := r.Updates(&p).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Get(id int64) (*Process, error) {
	p := Process{}
	err := r.First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}
