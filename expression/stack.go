////////////////////////////////////////////////////////////////////////////////
//	stack.go  -  Ago-24-2022  -  aldebap
//
//	Stack Data Structure
////////////////////////////////////////////////////////////////////////////////

package expression

//	originaly implemented in github.com/aldebap/algorithms_dataStructs/chapter_3/stack

type Stack interface {
	Push(value interface{})
	Pop() interface{}
	IsEmpty() bool
}

type StackAsArray struct {
	element []interface{}
}

//	New create a new stack as array
func NewStack() Stack {
	return &StackAsArray{
		element: make([]interface{}, 0),
	}
}

//	Push pushes a new item to the top of the stack
func (s *StackAsArray) Push(value interface{}) {
	//	in order to avoid relocating the slice when inserting a new item,
	//	the top is the last element of the slice
	s.element = append(s.element, value)
}

//	Pop return the item in the top of the stack after removing it
func (s *StackAsArray) Pop() interface{} {
	if len(s.element) == 0 {
		return nil
	}
	value := s.element[len(s.element)-1]
	s.element = s.element[:len(s.element)-1]

	return value
}

//	IsEmpty return true if the stack is empty
func (s *StackAsArray) IsEmpty() bool {
	return len(s.element) == 0
}
