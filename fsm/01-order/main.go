package main

type OrderState int

const (
	CREATE OrderState = iota
	PAID
	DELIVERING
	RECEIVED
	DONE
	CANCELLING
	RETURNING
	CLOSED
)

type Order struct {
	state OrderState
}

func NewOrder() *Order {
	return &Order{
		state: CREATE,
	}
}

func (o *Order) can_pay() bool {
	return o.state == CREATE
}

func (o *Order) can_deliver() bool {
	return o.state == PAID
}

func (o *Order) can_cancel() bool {
	return o.state == CREATE || o.state == PAID
}

func (o *Order) can_receive() bool {
	return o.state == DELIVERING
}

func (o *Order) payment_service() bool {
	// 调用 RPC 接口完成支付
	return false
}

func (o *Order) pay() bool {
	if o.can_pay() {
		ok := o.payment_service()
		if ok {
			o.state = PAID
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
func (o *Order) cancel() bool {
	if o.can_cancel() {
		o.state = CANCELLING
		// 取消订单，申请审批和清理数据，如果顺利成功再——
		o.state = CLOSED
	} else {
		return false
	}
	return false
}
