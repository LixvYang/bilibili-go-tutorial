package main

import (
	"errors"
	"fmt"
)

type Transition struct {
	from  string
	to    string
	event string
}

type StateMachine struct {
	state       string
	transitions []Transition
	handleEvent func(fromState string, toState string, args []interface{}) error
}

func NewStateMachine(init string, transitions []Transition, handleEvent func(fromState string, toState string, args []interface{}) error) *StateMachine {
	return &StateMachine{
		state:       init,
		transitions: transitions,
		handleEvent: handleEvent,
	}
}

func (m *StateMachine) changeState(state string) {
	m.state = state
}

func (m *StateMachine) findTransMatching(fromState string, event string) *Transition {
	for _, v := range m.transitions {
		if v.from == fromState && v.event == event {
			return &v
		}
	}
	return nil
}

func (m *StateMachine) Trigger(event string, args ...interface{}) (err error) {
	trans := m.findTransMatching(m.state, event)
	if trans == nil {
		err = errors.New("转换状态失败: 未找到trans")
		return
	}

	if trans.event == "" {
		err = errors.New("未找到具体的event")
		return
	}

	err = m.handleEvent(m.state, trans.to, args)
	if err != nil {
		err = errors.New("转换状态失败: 未找到handleEvent")
		return
	}

	m.changeState(trans.to)

	return
}

func main() {
	transitions := make([]Transition, 0)
	transitions = append(transitions, Transition{
		from:  "create",
		to:    "running",
		event: "start",
	})
	transitions = append(transitions, Transition{
		from:  "running",
		to:    "end",
		event: "work",
	})

	fsm := NewStateMachine("create", transitions, func(fromState string, toState string, args []interface{}) error {
		switch toState {
		case "end":
			fmt.Println("工作结束")
		}
		return nil
	})
	fsm.Trigger("start")
	fsm.Trigger("work")

	fmt.Println(fsm.state)
}
