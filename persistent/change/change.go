package change

type Change struct {
	ID        int64
	BizID     int64
	ProcessID int64
	User      string
	Action    string
	Data      string
	AutoSkip  bool
}

func (Change) TableName() string {
	return "change_log"
}

type NodeRepository interface {
	Save(c Change) (int64, error)
	Update(c Change) error
	Get(id int64) (*Change, error)
	List(query Change) ([]*Change, error)
}
