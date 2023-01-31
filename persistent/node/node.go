package node

// DAO
type Node struct {
	ID         int64
	CreateTime int64  `gorm:"column:created_at"`
	UpdateTime int64  `gorm:"column:updated_at"`
	FinishTime int64  `gorm:"column:finished_at"`
	Name       string `gorm:"column:role"`
	UserName   string `gorm:"column:value"`
	Memo       string `gorm:"column:remark"`
	Status     string
	// Type       string
	Result string `gorm:"column:result"`
	Pre    int64  `gorm:"column:pre_id"`
	Next   int64  `gorm:"column:next_id"`
	// BizID      int64
	// BizInfo    string
	ProcessID int64 `gorm:"column:flow_id"`
}

type NodeLine struct {
	ID         int64
	CreateTime int64 `gorm:"column:created_at"`
	UpdateTime int64 `gorm:"column:updated_at"`
	// UserID     int64
	// UserName   string
	// Memo      string
	// Status    string
	// Result    string
	Pre       int64 `gorm:"column:parent"`
	Next      int64 `gorm:"column:child"`
	ProcessID int64 `gorm:"column:flow_id"`
}

func (Node) TableName() string {
	return "task_node"
}

func (NodeLine) TableName() string {
	return "node_line"
}

func NewNode(name string, processID int64) Node {
	return Node{
		Name:      name,
		ProcessID: processID,
	}
}

func NewNodeLine(preID, nextID, processID int64) NodeLine {
	return NodeLine{
		Pre:       preID,
		Next:      nextID,
		ProcessID: processID,
	}
}

type NodeRepository interface {
	Save(n Node) (int64, error)
	Update(n Node) error
	Get(id int64) (*Node, error)
	AddLine(nl NodeLine) (int64, error)
	List(query Node) ([]*Node, error)
	GetLine(query NodeLine) (*NodeLine, error)
}
