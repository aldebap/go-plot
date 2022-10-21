////////////////////////////////////////////////////////////////////////////////
//	expression.go  -  Ago-25-2022  -  aldebap
//
//	Expression parser and evaluator
////////////////////////////////////////////////////////////////////////////////

package expression

import (
	"errors"
	"regexp"
	"strconv"
)

//	based on original implemention from github.com/aldebap/algorithms_dataStructs/chapter_3/expression

type Expression interface {
	Evaluate(x_value float64) (float64, error)
}

type ParsedExpression struct {
	postfix Queue
}

//	New create a new parsed expression
func NewExpression(expressionStr string) (Expression, error) {
	postfix, err := infix2postfix(expressionStr)
	if err != nil {
		return nil, errors.New("error parsing expression: " + err.Error())
	}

	return &ParsedExpression{
		postfix: postfix,
	}, nil
}

//	infix2postfixV2 read the infix expression and create a stack with the postfix version of it
func infix2postfixV2(expression string) (Queue, error) {
	postfix := NewQueue()

	return postfix, nil
}

//	lexicalAnalizer read the infix expression and create an array with all tokens
func lexicalAnalizer(expression string) ([]string, error) {

	var token []string = make([]string, 0)
	var identifier string
	var numericLiteral string

	for _, char := range expression {
		switch char {
		case ' ':
			if len(identifier) > 0 {
				token = append(token, identifier)
				identifier = ""
			}
			if len(numericLiteral) > 0 {
				token = append(token, numericLiteral)
				numericLiteral = ""
			}

		case '+', '-', '*', '/':
			if len(identifier) > 0 {
				token = append(token, identifier)
				identifier = ""
			}
			if len(numericLiteral) > 0 {
				token = append(token, numericLiteral)
				numericLiteral = ""
			}
			token = append(token, string(char))

		case '(', ')':
			if len(identifier) > 0 {
				token = append(token, identifier)
				identifier = ""
			}
			if len(numericLiteral) > 0 {
				token = append(token, numericLiteral)
				numericLiteral = ""
			}
			token = append(token, string(char))

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if len(identifier) > 0 {
				identifier += string(char)
			} else {
				numericLiteral += string(char)
			}

		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '_':
			if len(numericLiteral) > 0 {
				return nil, errors.New("invalid numeric literal")
			}
			identifier += string(char)
		}
	}

	//	add the last token when necessary
	if len(identifier) > 0 {
		token = append(token, identifier)
	} else if len(numericLiteral) > 0 {
		token = append(token, numericLiteral)
	}

	return token, nil
}

//	infix2postfix read the infix expression and create a stack with the postfix version of it
func infix2postfix(expression string) (Queue, error) {
	postfix := NewQueue()
	previousToken := ""
	operand := ""
	operator := NewStack()
	operatorParenthesisLevel := NewStack()
	parenthesis := 0

	//	initialize some regexp needed for the parsing
	intnumRegEx, err := regexp.Compile(`^(\d+)$`)
	if err != nil {
		return nil, err
	}

	//	TODO: refactor the whole function
	for _, char := range expression {
		switch char {
		case '+', '-', '*', '/':
			match := intnumRegEx.FindAllStringSubmatch(previousToken, -1)
			if previousToken != ")" && len(match) != 1 && previousToken != "x" {
				return nil, errors.New("operator without a previous operand")
			}
			if !operator.IsEmpty() {
				operatorAux := operator.Pop()
				parenthesisLevelAux := operatorParenthesisLevel.Pop()

				if parenthesis == parenthesisLevelAux.(int) {
					postfix.Put(operatorAux)
				} else {
					operator.Push(operatorAux)
					operatorParenthesisLevel.Push(parenthesisLevelAux.(int))
				}
			}
			operator.Push(string(char))
			operatorParenthesisLevel.Push(parenthesis)
			previousToken = string(char)

		case '(':
			previousToken = string(char)
			parenthesis++

		case ')':
			if parenthesis == 0 {
				return nil, errors.New("expression with missing open parenthesis")
			}
			if !operator.IsEmpty() {
				operatorAux := operator.Pop()
				parenthesisLevelAux := operatorParenthesisLevel.Pop()

				if parenthesis == parenthesisLevelAux.(int) {
					postfix.Put(operatorAux)
				} else {
					operator.Push(operatorAux)
					operatorParenthesisLevel.Push(parenthesisLevelAux.(int))
				}
			}
			previousToken = string(char)
			parenthesis--

		//	TODO: need to improve this to allow signed numbers as well as decimal numbers
		//	TODO: temporary validation for x variable
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if len(operand) > 0 && operand == "x" {
				return nil, errors.New("mixed digits with x variable")
			}
			operand += string(char)

		//	TODO: temporary parsing for x variable
		case 'x':
			if len(operand) > 0 {
				return nil, errors.New("mixed digits with x variable")
			}
			operand = "x"

		case ' ':
			if len(operand) > 0 {
				postfix.Put(operand)
				previousToken = operand
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
func (p *ParsedExpression) Evaluate(x_value float64) (float64, error) {

	postfixAux := p.postfix.Copy()
	operand := NewStack()

	for {
		item := postfixAux.Get()
		if item == nil {
			break
		}

		//	check if current token is a reference to x variable
		token := item.(string)

		if token == "x" {
			operand.Push(x_value)
			continue
		}

		//	attempt to convert the item to a float64 number
		//	TODO: should improve this forced conversion
		number, err := strconv.ParseFloat(token, 64)
		if err == nil {
			operand.Push(number)
			continue
		}

		//	must be an operation
		var operand1 float64
		var operand2 float64

		//	TODO: this message is not necesarily correct for the current operation
		if operand.IsEmpty() {
			return 0, errors.New("operation + requires two operands")
		}
		operand2 = operand.Pop().(float64)

		if operand.IsEmpty() {
			return 0, errors.New("operation + requires two operands")
		}
		operand1 = operand.Pop().(float64)

		switch token {
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

	return operand.Pop().(float64), nil
}
