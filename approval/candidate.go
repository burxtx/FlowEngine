package approval

import (
	"github.com/burxtx/FlowEngine/approval/domain"
	"github.com/burxtx/FlowEngine/persistent/candidate"
)

type NodeCandidate struct {
	candidate.CandidateRepository
}

// Prepare 绑定节点和审批人，并持久化
func (nc *NodeCandidate) Prepare(name string, nodeID, pid int64) error {
	// 序列化node
	newNC := candidate.NewCandidate(name, nodeID, pid)
	_, err := nc.CandidateRepository.Save(newNC)
	if err != nil {
		return err
	}
	return nil
}

func (nc *NodeCandidate) ListUserNodes(name string, pid ...int64) ([]*candidate.Candidate, error) {
	cdt := candidate.Candidate{
		Approver: name,
	}
	if len(pid) > 0 {
		cdt.PID = pid[0]
	}
	userNodes, err := nc.CandidateRepository.List(cdt)
	if err != nil {
		return nil, err
	}
	return userNodes, nil
}

func (nc *NodeCandidate) ListNodeUsers(nodeDOM *domain.Node) ([]string, error) {
	cdt := candidate.Candidate{
		NodeID: nodeDOM.ID,
	}
	userNodes, err := nc.CandidateRepository.List(cdt)
	if err != nil {
		return nil, err
	}
	users := make([]string, 0)
	for _, v := range userNodes {
		users = append(users, v.Approver)
	}
	return users, nil
}

func NewNodeCandidate(repo candidate.CandidateRepository) *NodeCandidate {
	return &NodeCandidate{
		CandidateRepository: repo,
	}
}
