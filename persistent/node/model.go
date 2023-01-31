package node

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	*gorm.DB
}

func (r *Repository) Save(n Node) (int64, error) {
	n.CreateTime = time.Now().Unix()
	n.UpdateTime = time.Now().Unix()
	err := r.Create(&n).Error
	if err != nil {
		return 0, err
	}
	return n.ID, nil
}

func struct2map(n Node) map[string]interface{} {
	return map[string]interface{}{
		"status":     n.Status,
		"value":      n.UserName,
		"remark":     n.Memo,
		"updated_at": n.UpdateTime,
	}
}

func (r *Repository) Update(n Node) error {
	n.UpdateTime = time.Now().Unix()
	err := r.Model(&n).Where("id = ?", n.ID).Updates(struct2map(n)).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Get(id int64) (*Node, error) {
	n := Node{}
	err := r.First(&n, id).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *Repository) AddLine(nl NodeLine) (int64, error) {
	nl.CreateTime = time.Now().Unix()
	nl.UpdateTime = time.Now().Unix()
	err := r.Create(&nl).Error
	if err != nil {
		return 0, err
	}
	return nl.ID, nil
}

func (r *Repository) List(query Node) ([]*Node, error) {
	var n []*Node
	err := r.Where(&query).Find(&n).Error
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (r *Repository) GetLine(query NodeLine) (*NodeLine, error) {
	nl := NodeLine{}
	err := r.Where(&query).First(&nl).Error
	if err != nil {
		return nil, err
	}
	return &nl, nil
}
