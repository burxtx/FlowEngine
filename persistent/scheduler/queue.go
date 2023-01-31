package scheduler

import (
	"encoding/json"
	"time"

	"github.com/burxtx/FlowEngine/approval/domain"
	"gorm.io/gorm"
)

type TaskQueue interface {
	Pop() (curNode *domain.Node, err error)
	Push(res *ScheduleQueue) error
	List(query ScheduleQueue) ([]*ScheduleQueue, error)
}

type TaskQueueInstance struct {
	*gorm.DB
}

type ScheduleQueue struct {
	ID         int64
	NodeID     int64
	ProcessID  int64
	CreateTime int64 `gorm:"column:created_at"`
	User       string
	State      string
	Memo       string
	Name       string
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;null"`
	DomainNode string         `gorm:"column:data"`
}

func (ScheduleQueue) TableName() string {
	return "schedule_queue"
}

func (q *TaskQueueInstance) Pop() (curNode *domain.Node, err error) {
	sq := ScheduleQueue{}
	result := q.First(&sq)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if err = result.Error; err != nil {
		return
	}
	err = json.Unmarshal([]byte(sq.DomainNode), &curNode)
	if err != nil {
		return
	}
	err = q.Delete(&sq).Error
	if err != nil {
		return
	}
	return
}

func (q *TaskQueueInstance) Push(v *ScheduleQueue) (err error) {
	v.CreateTime = time.Now().Unix()
	err = q.Create(&v).Error
	return
}

func (q *TaskQueueInstance) List(query ScheduleQueue) (res []*ScheduleQueue, err error) {
	err = q.Where(&query).Find(&res).Error
	return
}

func NewTaskQueueInstance(db *gorm.DB) TaskQueue {
	return &TaskQueueInstance{
		db,
	}
}
