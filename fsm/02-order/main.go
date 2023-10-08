package main

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"
)

func main() {
	order := fsm.NewFSM(
		"created",
		fsm.Events{
			{Name: "pay", Src: []string{"created"}, Dst: "paid"},
			{Name: "deliver", Src: []string{"paid"}, Dst: "delivering"},
			{Name: "receive", Src: []string{"delivering"}, Dst: "received"},
			{Name: "confirm", Src: []string{"received"}, Dst: "done"},

			{Name: "cancel", Src: []string{"received", "paid"}, Dst: "cancelling"},
			{Name: "return", Src: []string{"delivering", "received"}, Dst: "returning"},
			{Name: "close", Src: []string{"cancelling", "returning"}, Dst: "closed"},
		},
		fsm.Callbacks{
			"before_pay": func(_ context.Context, e *fsm.Event) {
				fmt.Println("支付服务申请中……")
				// 发送 payment_service 服务
			},
			"paid": func(_ context.Context, e *fsm.Event) {
				fmt.Println("支付成功")
			},
			"after_deliver": func(_ context.Context, e *fsm.Event) {
				fmt.Println("已通知用户：商品配送中")
			},
			"cancel": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("订单取消")
				e.Cancel()
			},
			"return": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("订单返回")
				e.Cancel()
			},
			"close": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("订单关闭")
				e.Cancel()
			},
		},
	)

	fmt.Println(order.Current())

	err := order.Event(context.Background(), "pay")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(order.Current())

	err = order.Event(context.Background(), "deliver")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(order.Current())

	err = order.Event(context.Background(), "receive")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(order.Current())

	err = order.Event(context.Background(), "confirm")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(order.Current())
}
