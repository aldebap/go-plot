////////////////////////////////////////////////////////////////////////////////
//	expression.go  -  Ago-25-2022  -  aldebap
//
//	Implementation of a simple expression parser
////////////////////////////////////////////////////////////////////////////////

package expression

import (
	"errors"
	"log"
	"strconv"
)

//	Parse parses a simple expression using a stack
func Parse(expression string) (int64, error) {

	postfix, err := infix2postfix(expression)
	if err != nil {
		log.Fatalf("error parsing expression: %s", err.Error())
	}

	return evaluatePolishReverse(postfix)
}

//	infix2postfix read the infix expression and create a stack with the postfix version of it
func infix2postfix(expression string) (Queue, error) {
	postfix := NewQueue()
	operand := ""
	operator := NewStack()
	parenthesis := 0

	for _, char := range expression {
		switch char {
		case '+', '-', '*', '/':
			operator.Push(string(char))

		case '(':
			parenthesis++

		case ')':
			if parenthesis == 0 {
				return nil, errors.New("expression with missing open parenthesis")
			}
			if !operator.IsEmpty() {
				postfix.Put(operator.Pop())
			}
			parenthesis--

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			operand += string(char)

		case ' ':
			if len(operand) > 0 {
				postfix.Put(operand)
				operand = ""
			}

		default:
			return nil, errors.New("syntax error: invalid character in expression: " + string(char))
		}
	}
	//	make sure the parenthesis are balanced till the end
	if parenthesis != 0 {
		return nil, errors.New("expression with unbalanced parenthesis")
	}

	//	make sure the last operand and operators are pushed to the stack
	if len(operand) > 0 {
		postfix.Put(operand)
	}
	if !operator.IsEmpty() {
		postfix.Put(operator.Pop())
	}

	return postfix, nil
}

//	evaluatePolishReverse evaluate the Polish reverse expression (postfix) and return a numerical result
func evaluatePolishReverse(postfix Queue) (int64, error) {

	operand := NewStack()

	for {
		item := postfix.Get()
		if item == nil {
			break
		}

		//	attempt to convert the item to a int number
		number, err := strconv.ParseInt(item.(string), 10, 64)
		if err == nil {
			operand.Push(number)
			continue
		}

		//	must be an operation
		var operand1 int64
		var operand2 int64

		if operand.IsEmpty() {
			return 0, errors.New("operation + requires two operands")
		}
		operand2 = operand.Pop().(int64)

		if operand.IsEmpty() {
			return 0, errors.New("operation + requires two operands")
		}
		operand1 = operand.Pop().(int64)

		switch item.(string) {
		case "+":
			operand.Push(operand1 + operand2)

		case "-":
			operand.Push(operand1 - operand2)

		case "*":
			operand.Push(operand1 * operand2)

		case "/":
			operand.Push(operand1 / operand2)
		}
	}

	if operand.IsEmpty() {
		return 0, nil
	}

	return operand.Pop().(int64), nil
}
