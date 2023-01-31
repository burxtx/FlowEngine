package candidate

type Candidate struct {
	ID         int64
	CreateTime int64 `gorm:"column:created_at"`
	UpdateTime int64 `gorm:"column:updated_at"`
	NodeID     int64
	PID        int64 `gorm:"column:flow_id"`
	// CandidateID   int64
	Approver string
}

func (Candidate) TableName() string {
	return "candidate"
}

func NewCandidate(user string, nodeID, pid int64) Candidate {
	return Candidate{
		Approver: user,
		NodeID:   nodeID,
		PID:      pid,
	}
}

type CandidateRepository interface {
	Save(c Candidate) (int64, error)
	Get(id int64) (*Candidate, error)
	List(query Candidate) ([]*Candidate, error)
}
