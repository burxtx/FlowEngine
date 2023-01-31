package domain

type BizInfo struct {
	Info   string
	Amount float64
}

func (b *BizInfo) SetInfo() {

}

func (b *BizInfo) GetInfo() string {
	return b.Info
}

func NewBizInfo(amount float64, info string) BizInfo {
	return BizInfo{
		Info:   info,
		Amount: amount,
	}
}
