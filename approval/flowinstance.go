package approval

import (
	"github.com/burxtx/FlowEngine/approval/domain"
	"github.com/burxtx/FlowEngine/persistent/process"
)

type FlowInstance struct {
	process.ProcessRepository
}

// func (f *FlowInstance) AddFirst(v *Node) {
// 	v.Parents = nil
// 	f.nodes[v.ID] = v

// }

// func (f *FlowInstance) AddLast(v *Node) {
// 	v.Children = nil
// 	f.nodes[v.ID] = v
// }

func (f *FlowInstance) AddBefore(bfID, nodeID int64) {

}

func (f *FlowInstance) AddAfter(afID, nodeID int64) {

}

func (f *FlowInstance) Remove(nodeID int64) error {
	return nil
}

func (f *FlowInstance) Append(nodeID int64) error {
	return nil
}

// func (f *FlowInstance) First() *Node {
// 	return f.HeadNode
// }

// func (f *FlowInstance) Last() *Node {
// 	return f.TailNode
// }

// func (f *FlowInstance) Name() string {
// 	return f.name
// }

func (f *FlowInstance) SetResult(fi *domain.FlowInstance, domResult domain.Result) error {
	p := process.Process{
		ID:           fi.ID,
		CurrentState: domResult.Name,
	}
	err := f.ProcessRepository.Update(p)
	if err != nil {
		return err
	}
	fi.Result = domResult
	return err
}

func (f *FlowInstance) SetCurrentNode(fi *domain.FlowInstance, v *domain.Node) error {
	p := process.Process{
		ID:          fi.ID,
		CurrentNode: v.ID,
	}
	err := f.ProcessRepository.Update(p)
	if err != nil {
		return err
	}
	fi.SetCurrentNode(v)
	return nil
}

func (f *FlowInstance) GetCurrentNode(fi *domain.FlowInstance) (int64, error) {
	p, err := f.ProcessRepository.Get(fi.ID)
	if err != nil {
		return 0, err
	}
	return p.CurrentNode, nil
}

func (f *FlowInstance) SetHeadNode(fi *domain.FlowInstance, v *domain.Node) error {
	p := process.Process{
		ID:       fi.ID,
		HeadNode: v.ID,
	}
	err := f.ProcessRepository.Update(p)
	if err != nil {
		return err
	}
	fi.HeadNode = v
	return nil
}

func (f *FlowInstance) SetTailNode(fi *domain.FlowInstance, v *domain.Node) error {
	p := process.Process{
		ID:       fi.ID,
		TailNode: v.ID,
	}
	err := f.ProcessRepository.Update(p)
	if err != nil {
		return err
	}
	fi.TailNode = v
	return nil
}

func (f *FlowInstance) Init(name string) (*domain.FlowInstance, error) {
	fi := domain.NewFlowInstance()
	p := process.NewProcess(name, domain.Initializing)
	pID, err := f.ProcessRepository.Save(p)
	if err != nil {
		return fi, err
	}
	fi.ID = pID
	return fi, nil
}

func newInstance(repo process.ProcessRepository) *FlowInstance {
	return &FlowInstance{
		ProcessRepository: repo,
	}
}
