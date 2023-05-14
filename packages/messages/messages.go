package messages

import "encoding/json"

type Action uint32

const (
	ERROR   Action = 0
	REFUSE  Action = 1
	ALLOW   Action = 2
	REQUEST Action = 3
	FREE    Action = 5
	ACKFREE Action = 6
)

type Message struct {
	Action   Action
	Lockback string
}

func (m *Message) Pack() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) Unpack(b []byte) error {
	return json.Unmarshal(b, m)
}
