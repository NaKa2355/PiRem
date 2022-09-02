package message

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

func (m RoundTrip) Receive() (Message, error) {
	resp, ok := <-m.returnChan
	if !ok {
		return resp, ErrNoReply
	}
	return resp, nil
}

func (m *RoundTrip) Close() {
	if m.returnChan == nil {
		return
	}

	close(m.returnChan)
	m.returnChan = nil
}
