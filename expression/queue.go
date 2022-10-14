////////////////////////////////////////////////////////////////////////////////
//	queue.go  -  Ago-24-2022  -  aldebap
//
//	Queue Data Structure
////////////////////////////////////////////////////////////////////////////////

package expression

type Queue interface {
	Put(value interface{})
	Get() interface{}
	IsEmpty() bool
}

type QueueAsArray struct {
	element []interface{}
}

//	New create a new queue as array
func NewQueue() Queue {
	return &QueueAsArray{
		element: make([]interface{}, 0),
	}
}

func (q *QueueAsArray) Put(value interface{}) {
	q.element = append(q.element, value)
}

func (q *QueueAsArray) Get() interface{} {
	if len(q.element) == 0 {
		return nil
	}
	value := q.element[0]
	q.element = q.element[1:]

	return value
}

func (q *QueueAsArray) IsEmpty() bool {
	return len(q.element) == 0
}
