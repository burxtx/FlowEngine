package candidate

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	*gorm.DB
}

func (r *Repository) Save(cdt Candidate) (int64, error) {
	cdt.CreateTime = time.Now().Unix()
	cdt.UpdateTime = time.Now().Unix()
	err := r.Create(&cdt).Error
	if err != nil {
		return 0, err
	}
	return cdt.ID, nil
}

func (r *Repository) Get(id int64) (*Candidate, error) {
	cdt := Candidate{}
	err := r.First(&cdt, id).Error
	if err != nil {
		return nil, err
	}
	return &cdt, nil
}

func (r *Repository) List(query Candidate) ([]*Candidate, error) {
	var res []*Candidate
	err := r.Where(&query).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
