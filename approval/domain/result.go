package domain

type Result struct {
	Name   string
	CnName string
}

func (r *Result) String() string {
	return r.Name
}

func NewResult(name, cnName string) Result {
	return Result{
		Name:   name,
		CnName: cnName,
	}
}
