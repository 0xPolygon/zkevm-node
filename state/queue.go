package state

import (
	"errors"
	"sync"
)

// ErrQueueEmpty is returned when a queue operation
// depends on the queue to not be empty, but it is empty
var ErrQueueEmpty = errors.New("queue is empty")

// Queue is a generic queue implementation that implements FIFO
type Queue[T any] struct {
	items []T
	mutex *sync.Mutex
}

// NewQueue creates a new instance of queue and initializes it
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0),
		mutex: &sync.Mutex{},
	}
}

// Push enqueue an item
func (q *Queue[T]) Push(item T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.items = append(q.items, item)
}

// Top returns the top level item without removing it
func (q *Queue[T]) Top() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	var v T
	if len(q.items) == 0 {
		return v, ErrQueueEmpty
	}
	return q.items[0], nil
}

// Pop returns the top level item and unqueues it
func (q *Queue[T]) Pop() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	var v T
	if len(q.items) == 0 {
		return v, ErrQueueEmpty
	}
	v = q.items[0]
	q.items = q.items[1:]
	return v, nil
}

// Len returns the size of the queue
func (q *Queue[T]) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.items)
}

// IsEmpty returns false if the queue has itens, otherwise true
func (q *Queue[T]) IsEmpty() bool {
	return q.Len() == 0
}
