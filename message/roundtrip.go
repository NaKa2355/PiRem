package message

import "time"

type RoundTrip struct {
	value      interface{}
	returnChan chan Message
}

func NewRoundTrip(value interface{}) Message {
	m := RoundTrip{}
	m.value = value
	m.returnChan = make(chan Message)
	return &m
}

func (m RoundTrip) GetValue() interface{} {
	return m.value
}

func (m *RoundTrip) SendBack(resp Message) {
	if m.returnChan == nil {
		return
	}

	m.returnChan <- resp
	close(m.returnChan)
	m.returnChan = nil
}

func (m RoundTrip) Receive(timeout time.Duration) (Message, error) {
	select {
	case resp := <-m.returnChan:
		return resp, ErrResponseTimeout
	case <-time.After(timeout):
		return &RoundTrip{}, nil
	}
}
