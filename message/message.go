package message

import (
	"errors"
)

type Message interface {
	GetValue() interface{}
	SendBack(Message)
	Receive() (Message, error)
	Close()
}

var ErrNoReply = errors.New("no reply from the destination")
