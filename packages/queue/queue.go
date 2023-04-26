package queue

import (
	"errors"
	"sync"
)

type node struct {
	destiny uint64
	data    interface{}
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

func (q *Queue) Dequeue() (interface{}, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	n := len(q.elements)

	if n < 1 {
		return nil, errors.New("empty queue")
	}

	data := q.elements[0].data

	if n > 1 {
		q.elements = q.elements[1:]
	} else {
		q.elements = []node{}
	}

	return data, nil
}

func (q *Queue) Enqueue(id uint64, data interface{}) {
	q.mutex.Lock()

	q.elements = append(q.elements, node{id, data})

	q.mutex.Unlock()
}
