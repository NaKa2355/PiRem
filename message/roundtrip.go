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
	var resp Message
	select {
	case resp, ok := <-m.returnChan:
		if !ok {
			return resp, ErrNoReply
		}
		return resp, nil
	case <-time.After(timeout):
		return resp, ErrTimeout
	}
}

func (m *RoundTrip) Close() {
	if m.returnChan == nil {
		return
	}

	close(m.returnChan)
	m.returnChan = nil
}
