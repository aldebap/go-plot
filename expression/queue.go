////////////////////////////////////////////////////////////////////////////////
//	queue.go  -  Ago-24-2022  -  aldebap
//
//	Queue Data Structure
////////////////////////////////////////////////////////////////////////////////

package expression

//	originaly implemented in github.com/aldebap/algorithms_dataStructs/chapter_3/queue

type Queue interface {
	Put(value interface{})
	Get() interface{}
	IsEmpty() bool
	Copy() Queue
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

//	Create a copy of queue as array
func (q *QueueAsArray) Copy() Queue {

	element := make([]interface{}, len(q.element))

	for i, _ := range q.element {
		element[i] = q.element[i]
	}

	return &QueueAsArray{
		element: element,
	}
}
