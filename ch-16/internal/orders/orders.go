package orders

type OrderInfo struct {
	OrderCode   rune
	Amount      int
	OrderNumber uint16
	Items       []string
	IsReady     bool
}

type SmallOrderInfo struct {
	Items       []string
	Amount      int
	OrderCode   rune
	OrderNumber uint16
	IsReady     bool
}
