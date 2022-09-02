package message

import (
	"errors"
	"time"
)

type Message interface {
	GetValue() interface{}
	SendBack(Message)
	Receive(timeout time.Duration) (Message, error)
}

var ErrResponseTimeout = errors.New("response timeout")