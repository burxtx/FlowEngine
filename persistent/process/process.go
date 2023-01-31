package process

// Process 是持久化的流程实例
type Process struct {
	ID           int64
	BizID        int64
	Name         string `gorm:"column:service"`
	CurrentState string `gorm:"column:cur_state"`
	CurrentNode  int64  `gorm:"column:cur_node"`
	HeadNode     int64
	TailNode     int64
	CreateTime   int64 `gorm:"column:created_at"`
	UpdateTime   int64 `gorm:"column:updated_at"`
	CloseTime    int64 `gorm:"column:closed_at"`
}

func (Process) TableName() string {
	return "process"
}

func NewProcess(name, state string) Process {
	return Process{
		Name:         name,
		CurrentState: state,
	}
}

type ProcessRepository interface {
	Save(p Process) (int64, error)
	Update(p Process) error
	Get(id int64) (*Process, error)
}
