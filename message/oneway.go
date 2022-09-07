package message

import "time"

type OneWay struct {
	value interface{}
}

func NewOneWay(value interface{}) Message {
	m := OneWay{}
	m.value = value
	return m
}

func (m OneWay) GetValue() interface{} {
	return m.value
}

func (m OneWay) SendBack(Message) {
}

func (m OneWay) Receive(timeout time.Duration) (Message, error) {
	return OneWay{}, nil
}

func (m OneWay) Close() {}
