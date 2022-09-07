package message

import (
	"errors"
	"time"
)

type Message interface {
	GetValue() interface{}
	SendBack(Message)
	Receive(timeout time.Duration) (Message, error)
	Close()
}

var (
	ErrNoReply = errors.New("no reply from the destination")
	ErrTimeout = errors.New("connection timeout")
)
