package queue

import (
	"errors"
	"sync"

	"github.com/namelew/RPC/packages/messages"
)

type node struct {
	destiny uint64
	data    messages.Message
}

type Queue struct {
	elements []node
	mutex    *sync.Mutex
}

func (q *Queue) Check(id uint64) bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.elements[0].destiny == id
}

func (q *Queue) Dequeue() (messages.Message, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	n := len(q.elements)

	if n < 1 {
		return messages.Message{}, errors.New("empty queue")
	}

	data := q.elements[0].data

	if n > 1 {
		q.elements = q.elements[1:]
	} else {
		q.elements = []node{}
	}

	return data, nil
}

func (q *Queue) Enqueue(id uint64, data messages.Message) {
	q.mutex.Lock()

	q.elements = append(q.elements, node{id, data})

	q.mutex.Unlock()
}
