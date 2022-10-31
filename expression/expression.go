////////////////////////////////////////////////////////////////////////////////
//	expression.go  -  Ago-25-2022  -  aldebap
//
//	Expression parser and evaluator
////////////////////////////////////////////////////////////////////////////////

package expression

import (
	"errors"
	"fmt"
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

//	types of tokens used in expressions (terminal symbols)
const (
	LITERAL           uint8 = 1
	NAME              uint8 = 2
	ADD_OPERATOR      uint8 = 3
	SUB_OPERATOR      uint8 = 4
	TIMES_OPERATOR    uint8 = 5
	DIV_OPERATOR      uint8 = 6
	OPEN_PARENTHESIS  uint8 = 7
	CLOSE_PARENTHESIS uint8 = 8
	EMPTY             uint8 = 9
)

type token struct {
	category uint8
	value    string
}

//	lexicalAnalizer read the infix expression and create an array with all tokens
func lexicalAnalizer(expression string) ([]token, error) {

	var tokenList []token = make([]token, 0)
	var identifier string
	var literal string

	for _, char := range expression {
		switch char {
		//	a space after a valid token means the previous token have to be appended to the list
		case ' ':
			if len(identifier) > 0 {
				tokenList = append(tokenList, token{
					category: NAME,
					value:    identifier,
				})
				identifier = ""
			}
			if len(literal) > 0 {
				_, err := strconv.ParseFloat(literal, 64)
				if err != nil {
					return nil, errors.New("non numeric literal: " + err.Error())
				}
				tokenList = append(tokenList, token{
					category: LITERAL,
					value:    literal,
				})

				literal = ""
			}

		//	an operator can also means the previous token needs to be appended to the list
		case '+', '-', '*', '/':
			if len(identifier) > 0 {
				tokenList = append(tokenList, token{
					category: NAME,
					value:    identifier,
				})
				identifier = ""
			}
			if len(literal) > 0 {
				_, err := strconv.ParseFloat(literal, 64)
				if err != nil {
					return nil, errors.New("invalid numeric literal: " + err.Error())
				}
				tokenList = append(tokenList, token{
					category: LITERAL,
					value:    literal,
				})

				literal = ""
			}

			//	get operator's category
			var category uint8

			switch char {
			case '+':
				category = ADD_OPERATOR
			case '-':
				category = SUB_OPERATOR
			case '*':
				category = TIMES_OPERATOR
			case '/':
				category = DIV_OPERATOR
			}

			tokenList = append(tokenList, token{
				category: category,
				value:    string(char),
			})

		//	parenthesis can also means the previous token needs to be appended to the list
		case '(', ')':
			if len(identifier) > 0 {
				tokenList = append(tokenList, token{
					category: NAME,
					value:    identifier,
				})
				identifier = ""
			}
			if len(literal) > 0 {
				_, err := strconv.ParseFloat(literal, 64)
				if err != nil {
					return nil, errors.New("invalid numeric literal: " + err.Error())
				}
				tokenList = append(tokenList, token{
					category: LITERAL,
					value:    literal,
				})

				literal = ""
			}

			//	get parenthesis's category
			var category uint8

			if char == '(' {
				category = OPEN_PARENTHESIS
			} else {
				category = CLOSE_PARENTHESIS
			}
			tokenList = append(tokenList, token{
				category: category,
				value:    string(char),
			})

		//	a digit can be part of a literal or a name
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if len(identifier) > 0 {
				identifier += string(char)
			} else {
				literal += string(char)
			}

		//	a dot can only be part of a literal
		case '.':
			if len(identifier) > 0 {
				return nil, errors.New("invalid character in identifier: " + identifier + string(char))
			}
			literal += string(char)

		//	a letter can be part of a name
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '_':
			if len(literal) > 0 {
				return nil, errors.New("invalid numeric literal: " + literal + string(char))
			}
			identifier += string(char)
		}
	}

	//	add the last token when necessary
	if len(identifier) > 0 {
		tokenList = append(tokenList, token{
			category: NAME,
			value:    identifier,
		})
	} else if len(literal) > 0 {
		_, err := strconv.ParseFloat(literal, 64)
		if err != nil {
			return nil, errors.New("invalid numeric literal: " + err.Error())
		}
		tokenList = append(tokenList, token{
			category: LITERAL,
			value:    literal,
		})
	}

	return tokenList, nil
}

//	types of syntax elements used in expressions
const (
	TARGET          uint8 = 101
	EXPRESSION      uint8 = 102
	TERM            uint8 = 103
	EXPRESSION_LINE uint8 = 104
	FACTOR          uint8 = 105
	TERM_LINE       uint8 = 106
)

//	Context free grammar entry
type grammarEntry struct {
	symbol  uint8
	derives []uint8
}

//	all entries for mathematical expressions grammar
var expressionGrammar []grammarEntry = []grammarEntry{
	{symbol: TARGET, derives: []uint8{EXPRESSION}},
	{symbol: EXPRESSION, derives: []uint8{TERM, EXPRESSION_LINE}},
	{symbol: EXPRESSION_LINE, derives: []uint8{ADD_OPERATOR, TERM, EXPRESSION_LINE}},
	{symbol: EXPRESSION_LINE, derives: []uint8{SUB_OPERATOR, TERM, EXPRESSION_LINE}},
	{symbol: EXPRESSION_LINE, derives: []uint8{EMPTY}},
	{symbol: TERM, derives: []uint8{FACTOR, TERM_LINE}},
	{symbol: TERM_LINE, derives: []uint8{TIMES_OPERATOR, FACTOR, TERM_LINE}},
	{symbol: TERM_LINE, derives: []uint8{DIV_OPERATOR, FACTOR, TERM_LINE}},
	{symbol: TERM_LINE, derives: []uint8{EMPTY}},
	{symbol: FACTOR, derives: []uint8{OPEN_PARENTHESIS, EXPRESSION, CLOSE_PARENTHESIS}},
	{symbol: FACTOR, derives: []uint8{LITERAL}},
	{symbol: FACTOR, derives: []uint8{NAME}},
}

//	first and followe sets to avoid backtracking in the parser's search
var firstSet map[uint8][]uint8 = map[uint8][]uint8{
	EXPRESSION:      []uint8{OPEN_PARENTHESIS, LITERAL, NAME},
	EXPRESSION_LINE: []uint8{ADD_OPERATOR, SUB_OPERATOR, EMPTY},
	TERM:            []uint8{OPEN_PARENTHESIS, LITERAL, NAME},
	TERM_LINE:       []uint8{TIMES_OPERATOR, DIV_OPERATOR, EMPTY},
	FACTOR:          []uint8{OPEN_PARENTHESIS, LITERAL, NAME},
}

var followSet []grammarEntry = []grammarEntry{
	{symbol: EXPRESSION, derives: []uint8{EMPTY, CLOSE_PARENTHESIS}},
	{symbol: EXPRESSION_LINE, derives: []uint8{EMPTY, CLOSE_PARENTHESIS}},
	{symbol: TERM, derives: []uint8{EMPTY, ADD_OPERATOR, SUB_OPERATOR, CLOSE_PARENTHESIS}},
	{symbol: TERM_LINE, derives: []uint8{EMPTY, ADD_OPERATOR, SUB_OPERATOR, CLOSE_PARENTHESIS}},
	{symbol: FACTOR, derives: []uint8{EMPTY, ADD_OPERATOR, SUB_OPERATOR, TIMES_OPERATOR, DIV_OPERATOR, CLOSE_PARENTHESIS}},
}

//	syntax tree generated by the parser
type syntaxNode struct {
	grammarExpasion []uint8
	inputToken      *token
}

//	expressionParser parse the expression from an array of tokens
func expressionParser(tokenList []token) (Queue, error) {

	//	create the syntax tree from the input tokens (top-down)
	var syntaxTree []syntaxNode = make([]syntaxNode, 0)
	var currentToken = 0

	syntaxTree = append(syntaxTree, syntaxNode{
		grammarExpasion: []uint8{EXPRESSION},
		inputToken:      nil,
	})

	for {
		var terminalOnly = true
		var treeLevel = len(syntaxTree) - 1

		for i, symbol := range syntaxTree[treeLevel].grammarExpasion {
			var derives []int = make([]int, 0)

			//	fmt.Printf("[debug] [1] symbol: %d\n", symbol)
			for j, entry := range expressionGrammar {
				if symbol == entry.symbol {
					derives = append(derives, j)
				}
			}
			//	fmt.Printf("[debug] [2] derives: %d --> %v\n", len(derives), derives)

			if len(derives) == 0 {
				//	symbol is a terminal so, a token must be used
				if currentToken < len(tokenList) && currentToken == i {
					syntaxTree = append(syntaxTree, syntaxNode{
						grammarExpasion: make([]uint8, len(syntaxTree[treeLevel].grammarExpasion)),
						inputToken:      &tokenList[currentToken],
					})

					for j, symbol := range syntaxTree[treeLevel].grammarExpasion {
						syntaxTree[treeLevel+1].grammarExpasion[j] = symbol
					}
					//	fmt.Printf("[debug] [2.5] node #%d: symbols: %v - token: %v\n", treeLevel+1, syntaxTree[treeLevel+1].grammarExpasion, syntaxTree[treeLevel+1].inputToken)

					currentToken++
					terminalOnly = false
					break
				}

				continue
			}

			//	search the possible derived symbols for a match with the current token
			var chosenNode = 0
			var chosenNodeIsEmpty = false

			if len(derives) > 1 {
				chosenNode = -1

				if currentToken < len(tokenList) {
					for j, derivedSymbol := range derives {
						if expressionGrammar[derivedSymbol].derives[0] == tokenList[currentToken].category {
							chosenNode = j
							break
						}
					}
				}
				if chosenNode == -1 && expressionGrammar[derives[len(derives)-1]].derives[0] == EMPTY {
					chosenNodeIsEmpty = true
				}
			}
			//	fmt.Printf("[debug] [4] chosen node: %d\n", chosenNode)

			//	add a new level to the search tree
			if chosenNodeIsEmpty {
				syntaxTree = append(syntaxTree, syntaxNode{
					grammarExpasion: make([]uint8, len(syntaxTree[treeLevel].grammarExpasion)-1),
					inputToken:      nil,
				})
			} else {
				syntaxTree = append(syntaxTree, syntaxNode{
					grammarExpasion: make([]uint8, len(syntaxTree[treeLevel].grammarExpasion)-1+len(expressionGrammar[derives[chosenNode]].derives)),
					inputToken:      nil,
				})
			}

			//	concat the terminal symbols (left) + current symbol derivation + symbols to the right of current symbol
			var j = 0

			//	fmt.Printf("[debug] [5] tree level: %d / %d\n", treeLevel, treeLevel+1)
			for k, symbol := range syntaxTree[treeLevel].grammarExpasion {
				if k != i {
					syntaxTree[treeLevel+1].grammarExpasion[j] = symbol
					j++
				} else {
					if !chosenNodeIsEmpty {
						for _, newSymbol := range expressionGrammar[derives[chosenNode]].derives {
							syntaxTree[treeLevel+1].grammarExpasion[j] = newSymbol
							j++
						}
					}
				}
			}
			//	fmt.Printf("[debug] [6] node #%d: symbols: %v - token: %v\n", treeLevel+1, syntaxTree[treeLevel+1].grammarExpasion, syntaxTree[treeLevel+1].inputToken)

			terminalOnly = false
			break
		}

		//	if the there are only terminal symbols, the search is over
		if terminalOnly {
			break
		}
	}
	//	TODO: how to identify syntax errors ?

	fmt.Printf("[debug] syntax tree:\n")
	for i, node := range syntaxTree {
		fmt.Printf("[debug] node #%d: symbols: %v - token: %v\n", i+1, node.grammarExpasion, node.inputToken)
	}

	//	TODO: create the postfix from the syntax tree
	postfix := NewQueue()

	return postfix, nil
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
