////////////////////////////////////////////////////////////////////////////////
//	expression_test.go  -  Ago-25-2022  -  aldebap
//
//	Test cases for the simple expression parser
////////////////////////////////////////////////////////////////////////////////

package expression

import (
	"fmt"
	"strings"
	"testing"
)

//	Test_lexicalAnalizer test cases for the lexical analizer
func Test_lexicalAnalizer(t *testing.T) {

	//	a few test cases
	var testScenarios = []struct {
		scenario string
		input    string
		output   []string
	}{
		{scenario: "addition", input: "2 + 5", output: []string{"2", "+", "5"}},
		{scenario: "subtraction", input: "5 - 2", output: []string{"5", "-", "2"}},
		{scenario: "multiplication", input: "2 * 5", output: []string{"2", "*", "5"}},
		{scenario: "division", input: "10 / 2", output: []string{"10", "/", "2"}},
		{scenario: "one parenthesis", input: "( 4 + 6 ) / 2", output: []string{"(", "4", "+", "6", ")", "/", "2"}},
		{scenario: "two parenthesis", input: "(4+(2*3))/2", output: []string{"(", "4", "+", "(", "2", "*", "3", ")", ")", "/", "2"}},
		{scenario: "unbalanced parenthesis", input: "(4+6/2", output: []string{"(", "4", "+", "6", "/", "2"}},
		{scenario: "fix to parenthesis parsing", input: "(x*x)-(3*x)+2", output: []string{"(", "x", "*", "x", ")", "-", "(", "3", "*", "x", ")", "+", "2"}},
		{scenario: "using precedence", input: "x*x-3*x+2", output: []string{"x", "*", "x", "-", "3", "*", "x", "+", "2"}},
		{scenario: "using a variable name", input: "x+4*y-z*z", output: []string{"x", "+", "4", "*", "y", "-", "z", "*", "z"}},
		{scenario: "using a variable name with underscore", input: "var_x+2", output: []string{"var_x", "+", "2"}},
		{scenario: "using a function call", input: "sin(x)", output: []string{"sin", "(", "x", ")"}},
		{scenario: "expression with a function call", input: "x*sin(2*x)", output: []string{"x", "*", "sin", "(", "2", "*", "x", ")"}},
	}

	t.Run(">>> test tokens found by lexical analizer", func(t *testing.T) {

		for _, test := range testScenarios {

			fmt.Printf("scenario: %s\n", test.scenario)

			//	execute lexical analizer
			tokenList, err := lexicalAnalizer(test.input)
			if err != nil {
				t.Errorf("unexpected error in lexical parser: %s", err)
				continue
			}

			want := test.output

			for i, tokenAux := range tokenList {
				if tokenAux.value != want[i] {
					t.Errorf("fail in lexical analizer: expected token: %s result: %v", want[i], tokenAux)
				}
			}
		}
	})
}

//	Test_expressionParser test cases for the expression parser
func Test_expressionParser(t *testing.T) {

	//	a few test cases
	var testScenarios = []struct {
		scenario string
		input    []token
		output   string
	}{
		{scenario: "addition", input: []token{
			{category: LITERAL, value: "2"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "5"},
		}, output: "2 5 +"},
		{scenario: "subtraction", input: []token{
			{category: LITERAL, value: "5"},
			{category: SUB_OPERATOR, value: "-"},
			{category: LITERAL, value: "2"},
		}, output: "5 2 -"},
		{scenario: "multiplication", input: []token{
			{category: LITERAL, value: "2"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "5"},
		}, output: "2 5 *"},
		{scenario: "division", input: []token{
			{category: LITERAL, value: "10"},
			{category: DIV_OPERATOR, value: "/"},
			{category: LITERAL, value: "2"},
		}, output: "10 2 /"},
		{scenario: "grouped addition", input: []token{
			{category: LITERAL, value: "2"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "5"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "7"},
		}, output: "2 5 7 + +"},
		{scenario: "grouped multiplication", input: []token{
			{category: LITERAL, value: "2"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "5"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "7"},
		}, output: "2 5 7 * *"},
		{scenario: "using precedence", input: []token{
			{category: LITERAL, value: "2"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "5"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "8"},
		}, output: "2 5 * 8 +"},
		{scenario: "more precedence", input: []token{
			{category: LITERAL, value: "2"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "5"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "8"},
		}, output: "2 5 8 * +"},
		{scenario: "multiples precedences", input: []token{
			{category: NAME, value: "x"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: NAME, value: "x"},
			{category: SUB_OPERATOR, value: "-"},
			{category: LITERAL, value: "3"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: NAME, value: "x"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "2"},
		}, output: "x x * 3 x * 2 + -"},
		{scenario: "one parenthesis", input: []token{
			{category: OPEN_PARENTHESIS, value: "("},
			{category: LITERAL, value: "4"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "6"},
			{category: CLOSE_PARENTHESIS, value: ")"},
			{category: DIV_OPERATOR, value: "/"},
			{category: LITERAL, value: "2"},
		}, output: "4 6 + 2 /"},
		{scenario: "two parenthesis", input: []token{
			{category: OPEN_PARENTHESIS, value: "("},
			{category: LITERAL, value: "4"},
			{category: ADD_OPERATOR, value: "+"},
			{category: OPEN_PARENTHESIS, value: "("},
			{category: LITERAL, value: "2"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "3"},
			{category: CLOSE_PARENTHESIS, value: ")"},
			{category: CLOSE_PARENTHESIS, value: ")"},
			{category: DIV_OPERATOR, value: "/"},
			{category: LITERAL, value: "2"},
		}, output: "4 2 3 * + 2 /"},
		{scenario: "fix to parenthesis parsing", input: []token{
			{category: OPEN_PARENTHESIS, value: "("},
			{category: NAME, value: "x"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: NAME, value: "x"},
			{category: CLOSE_PARENTHESIS, value: ")"},
			{category: SUB_OPERATOR, value: "-"},
			{category: OPEN_PARENTHESIS, value: "("},
			{category: LITERAL, value: "3"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: NAME, value: "x"},
			{category: CLOSE_PARENTHESIS, value: ")"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "2"},
		}, output: "x x * 3 x * 2 + -"},
		{scenario: "using a variable name", input: []token{
			{category: NAME, value: "x"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "4"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: NAME, value: "y"},
			{category: SUB_OPERATOR, value: "-"},
			{category: NAME, value: "z"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: NAME, value: "z"},
		}, output: "x 4 y * z z * - +"},
		{scenario: "using a variable name with underscore", input: []token{
			{category: NAME, value: "var_x"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "2"},
		}, output: "var_x 2 +"},
		{scenario: "using a function call", input: []token{
			{category: NAME, value: "sin"},
			{category: OPEN_PARENTHESIS, value: "("},
			{category: NAME, value: "x"},
			{category: CLOSE_PARENTHESIS, value: ")"},
		}, output: "x sin"},

		//	syntax error expressions scenarios
		{scenario: "operation without an operator", input: []token{
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "2"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "5"},
		}, output: "syntax error: unexpected token +"},
		{scenario: "operator without an operation", input: []token{
			{category: LITERAL, value: "2"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "5"},
			{category: ADD_OPERATOR, value: "+"},
		}, output: "syntax error: expected token 2"},
		{scenario: "unbalanced parenthesis", input: []token{
			{category: OPEN_PARENTHESIS, value: "("},
			{category: LITERAL, value: "4"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "6"},
			{category: DIV_OPERATOR, value: "/"},
			{category: LITERAL, value: "2"},
		}, output: "syntax error: expected token 9"},
		{scenario: "close parenthesis before opening it", input: []token{
			{category: LITERAL, value: "4"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "6"},
			{category: CLOSE_PARENTHESIS, value: ")"},
			{category: DIV_OPERATOR, value: "/"},
			{category: LITERAL, value: "2"},
		}, output: "syntax error: unexpected token )"},
	}

	t.Run(">>> test postfix expressions given by expressionParser()", func(t *testing.T) {

		for _, test := range testScenarios {

			fmt.Printf("scenario: %s\n", test.scenario)

			//	execute lexical analizer
			var want = test.output
			var got string

			postfix, err := expressionParser(test.input)
			if err != nil {
				got = err.Error()
			} else {
				for {
					item := postfix.Get()
					if item == nil {
						break
					}
					got += " " + item.(*token).value
				}
				got = strings.TrimLeft(got, " ")
				fmt.Printf("[debug] postfix result: %s\n", got)
			}

			//	check the result
			if want != got {
				t.Errorf("fail parsing the expression: postfix expected: %s result: %s", want, got)
			}
		}
	})
}

//	printSyntaxTree print a parsing/syntax tree structure
func printSyntaxTree(syntaxTree *syntaxNode) {

	type printNode struct {
		level      uint8
		syntaxTree *syntaxNode
	}

	var treeNodeDebug = NewStack()

	treeNodeDebug.Push(&printNode{
		level:      1,
		syntaxTree: syntaxTree,
	})
	for {
		if treeNodeDebug.IsEmpty() {
			break
		}

		node := treeNodeDebug.Pop().(*printNode)
		fmt.Printf("[debug] node #%d: grammar item: %d - #childs: %d - token: %v\n", node.level, node.syntaxTree.grammarItem,
			len(node.syntaxTree.childNodes), node.syntaxTree.inputToken)

		for i := len(node.syntaxTree.childNodes) - 1; i >= 0; i-- {
			treeNodeDebug.Push(&printNode{
				level:      node.level + 1,
				syntaxTree: node.syntaxTree.childNodes[i],
			})
		}
	}
}

//	printSyntaxTree print a parsing/syntax tree structure
func syntaxTree2string(syntaxTree *syntaxNode) string {

	type outputNode struct {
		level      uint8
		syntaxTree *syntaxNode
	}

	//	post-order traverse the syntax tree to create a string for testing the tree structure
	var treeOutput string
	var treeSearch = NewStack()

	treeSearch.Push(&outputNode{
		level:      1,
		syntaxTree: syntaxTree,
	})

	for {
		if treeSearch.IsEmpty() {
			break
		}

		searchNode := treeSearch.Pop().(*outputNode)
		if searchNode.syntaxTree.childNodes == nil {
			if searchNode.syntaxTree.inputToken == nil {
				continue
			}

			if searchNode.syntaxTree.inputToken.category == LITERAL ||
				searchNode.syntaxTree.inputToken.category == NAME {

				treeOutput += fmt.Sprintf("[%d] operand: %s; ", searchNode.level, searchNode.syntaxTree.inputToken.value)
			}

			if searchNode.syntaxTree.inputToken.category == FUNCTION_NAME {

				treeOutput += fmt.Sprintf("[%d] function call: %s; ", searchNode.level, searchNode.syntaxTree.inputToken.value)
			}

			if searchNode.syntaxTree.inputToken.category == ADD_OPERATOR ||
				searchNode.syntaxTree.inputToken.category == SUB_OPERATOR ||
				searchNode.syntaxTree.inputToken.category == TIMES_OPERATOR ||
				searchNode.syntaxTree.inputToken.category == DIV_OPERATOR {

				treeOutput += fmt.Sprintf("[%d] operator %s; ", searchNode.level, searchNode.syntaxTree.inputToken.value)
			}
			continue
		}

		//	insert left node, itself, and the right
		if len(searchNode.syntaxTree.childNodes) == 1 {
			treeSearch.Push(&outputNode{
				level:      searchNode.level + 1,
				syntaxTree: searchNode.syntaxTree.childNodes[0],
			})
		} else {
			if len(searchNode.syntaxTree.childNodes) == 2 {
				treeSearch.Push(&outputNode{
					level:      searchNode.level + 1,
					syntaxTree: searchNode.syntaxTree.childNodes[1],
				})
				treeSearch.Push(&outputNode{
					level:      searchNode.level + 1,
					syntaxTree: searchNode.syntaxTree.childNodes[0],
				})
			} else {
				if len(searchNode.syntaxTree.childNodes) == 3 {
					treeSearch.Push(&outputNode{
						level:      searchNode.level + 1,
						syntaxTree: searchNode.syntaxTree.childNodes[1],
					})
					treeSearch.Push(&outputNode{
						level:      searchNode.level + 1,
						syntaxTree: searchNode.syntaxTree.childNodes[2],
					})
					treeSearch.Push(&outputNode{
						level:      searchNode.level + 1,
						syntaxTree: searchNode.syntaxTree.childNodes[0],
					})
				} else {
					treeSearch.Push(&outputNode{
						level:      searchNode.level + 1,
						syntaxTree: searchNode.syntaxTree.childNodes[0],
					})
					treeSearch.Push(&outputNode{
						level:      searchNode.level + 1,
						syntaxTree: searchNode.syntaxTree.childNodes[2],
					})
					treeSearch.Push(&outputNode{
						level:      searchNode.level + 1,
						syntaxTree: searchNode.syntaxTree.childNodes[3],
					})
					treeSearch.Push(&outputNode{
						level:      searchNode.level + 1,
						syntaxTree: searchNode.syntaxTree.childNodes[1],
					})
				}
			}
		}
	}

	return treeOutput
}

//	Test_createParsingTree test cases for the createParsingTree function
func Test_createParsingTree(t *testing.T) {

	//	a few test cases
	var testScenarios = []struct {
		scenario string
		input    []token
		output   string
	}{
		{scenario: "constant", input: []token{
			{category: LITERAL, value: "2"},
		}, output: "[4] operand: 2;"},
		{scenario: "just a variable", input: []token{
			{category: NAME, value: "x"},
		}, output: "[4] operand: x;"},
		{scenario: "addition", input: []token{
			{category: LITERAL, value: "2"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "5"},
		}, output: "[4] operand: 2; [3] operator +; [5] operand: 5;"},
		{scenario: "subtraction", input: []token{
			{category: LITERAL, value: "5"},
			{category: SUB_OPERATOR, value: "-"},
			{category: LITERAL, value: "2"},
		}, output: "[4] operand: 5; [3] operator -; [5] operand: 2;"},
		{scenario: "multiplication", input: []token{
			{category: LITERAL, value: "2"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "5"},
		}, output: "[4] operand: 2; [4] operator *; [5] operand: 5;"},
		{scenario: "division", input: []token{
			{category: LITERAL, value: "10"},
			{category: DIV_OPERATOR, value: "/"},
			{category: LITERAL, value: "2"},
		}, output: "[4] operand: 10; [4] operator /; [5] operand: 2;"},
		{scenario: "parenthesis", input: []token{
			{category: OPEN_PARENTHESIS, value: "("},
			{category: LITERAL, value: "2"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "5"},
			{category: CLOSE_PARENTHESIS, value: ")"},
		}, output: "[7] operand: 2; [6] operator +; [8] operand: 5;"},
		{scenario: "function call", input: []token{
			{category: NAME, value: "sin"},
			{category: OPEN_PARENTHESIS, value: "("},
			{category: NAME, value: "x"},
			{category: CLOSE_PARENTHESIS, value: ")"},
		}, output: "[8] operand: x; [4] operand: sin;"},
		{scenario: "expression with a function call", input: []token{
			{category: LITERAL, value: "2"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: NAME, value: "x"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: NAME, value: "sin"},
			{category: OPEN_PARENTHESIS, value: "("},
			{category: LITERAL, value: "3"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: NAME, value: "x"},
			{category: CLOSE_PARENTHESIS, value: ")"},
		}, output: "[4] operand: 2; [4] operator *; [5] operator *; [10] operand: 3; [10] operator *; [11] operand: x; [6] operand: sin; [5] operand: x;"},

		//	syntax error expressions scenarios
		{scenario: "operands without an operator", input: []token{
			{category: LITERAL, value: "2"},
			{category: LITERAL, value: "5"},
		}, output: "syntax error: unexpected token 5"},
		{scenario: "operation before an operand", input: []token{
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "2"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "5"},
		}, output: "syntax error: unexpected token +"},
		{scenario: "operator without an operation", input: []token{
			{category: LITERAL, value: "2"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "5"},
			{category: ADD_OPERATOR, value: "+"},
		}, output: "syntax error: expected token 2"},
		{scenario: "unbalanced parenthesis", input: []token{
			{category: OPEN_PARENTHESIS, value: "("},
			{category: LITERAL, value: "4"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "6"},
			{category: DIV_OPERATOR, value: "/"},
			{category: LITERAL, value: "2"},
		}, output: "syntax error: expected token 9"},
		{scenario: "close parenthesis before opening it", input: []token{
			{category: LITERAL, value: "4"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "6"},
			{category: CLOSE_PARENTHESIS, value: ")"},
			{category: DIV_OPERATOR, value: "/"},
			{category: LITERAL, value: "2"},
		}, output: "syntax error: unexpected token )"},
	}

	t.Run(">>> test trees generated by parser", func(t *testing.T) {

		for _, test := range testScenarios {

			fmt.Printf("scenario: %s\n", test.scenario)

			//	execute lexical analizer
			var want = test.output
			var got string

			parsingTree, err := createParsingTree(test.input)
			if err != nil {
				got = err.Error()
			} else {
				got = syntaxTree2string(parsingTree)
				got = strings.TrimRight(got, " ")

				fmt.Printf("[debug] parsing tree: %s\n", got)
			}

			//	check the result
			if want != got {
				printSyntaxTree(parsingTree)
				t.Errorf("fail creating parsing tree: expected: %s result: %s", want, got)
			}
		}
	})
}

//	Test_createSyntaxTree test cases for the createSyntaxTree function
func Test_createSyntaxTree(t *testing.T) {

	//	a few test cases
	var testScenarios = []struct {
		scenario string
		input    []token
		output   string
	}{
		{scenario: "constant", input: []token{
			{category: LITERAL, value: "2"},
		}, output: "[4] operand: 2;"},
		{scenario: "just a variable", input: []token{
			{category: NAME, value: "x"},
		}, output: "[4] operand: x;"},
		{scenario: "addition", input: []token{
			{category: LITERAL, value: "2"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "5"},
		}, output: "[4] operand: 2; [5] operand: 5; [2] operator +;"},
		{scenario: "subtraction", input: []token{
			{category: LITERAL, value: "5"},
			{category: SUB_OPERATOR, value: "-"},
			{category: LITERAL, value: "2"},
		}, output: "[4] operand: 5; [5] operand: 2; [2] operator -;"},
		{scenario: "multiplication", input: []token{
			{category: LITERAL, value: "2"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "5"},
		}, output: "[4] operand: 2; [5] operand: 5; [3] operator *;"},
		{scenario: "division", input: []token{
			{category: LITERAL, value: "10"},
			{category: DIV_OPERATOR, value: "/"},
			{category: LITERAL, value: "2"},
		}, output: "[4] operand: 10; [5] operand: 2; [3] operator /;"},
		{scenario: "parenthesis", input: []token{
			{category: OPEN_PARENTHESIS, value: "("},
			{category: LITERAL, value: "2"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "5"},
			{category: CLOSE_PARENTHESIS, value: ")"},
		}, output: "[7] operand: 2; [8] operand: 5; [5] operator +;"},
		{scenario: "function call", input: []token{
			{category: NAME, value: "sin"},
			{category: OPEN_PARENTHESIS, value: "("},
			{category: NAME, value: "x"},
			{category: CLOSE_PARENTHESIS, value: ")"},
		}, output: "[8] operand: x; [4] function call: sin;"},
		{scenario: "expression with a function call", input: []token{
			{category: LITERAL, value: "2"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: NAME, value: "x"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: NAME, value: "sin"},
			{category: OPEN_PARENTHESIS, value: "("},
			{category: LITERAL, value: "3"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: NAME, value: "x"},
			{category: CLOSE_PARENTHESIS, value: ")"},
		}, output: "[4] operand: 2; [5] operand: x; [10] operand: 3; [11] operand: x; [9] operator *; [6] function call: sin; [4] operator *; [3] operator *;"},
	}

	t.Run(">>> test parser tree conversion to syntax tree", func(t *testing.T) {

		for _, test := range testScenarios {

			fmt.Printf("scenario: %s\n", test.scenario)

			//	execute lexical analizer
			var syntaxTree *syntaxNode
			var want = test.output
			var got string

			parsingTree, err := createParsingTree(test.input)
			if err != nil {
				got = err.Error()
			} else {
				syntaxTree, err = createSyntaxTree(parsingTree)
				if err != nil {
					got = err.Error()
				} else {

					got = syntaxTree2string(syntaxTree)
					got = strings.TrimRight(got, " ")

					fmt.Printf("[debug] syntax tree: %s\n", got)
				}
			}

			//	check the result
			if want != got {
				printSyntaxTree(syntaxTree)
				t.Errorf("fail creating syntax tree: expected: %s result: %s", want, got)
			}
		}
	})
}

//	Test_evaluatePolishReverse test cases for the Polish Reverse evaluation function
func Test_evaluatePolishReverse(t *testing.T) {

	//	a few test cases
	var testScenarios = []struct {
		scenario string
		input    []token
		x_value  float64
		output   float64
	}{
		{scenario: "addition", input: []token{
			{category: LITERAL, value: "2"},
			{category: LITERAL, value: "5"},
			{category: ADD_OPERATOR, value: "+"},
		}, x_value: 0, output: 7},
		{scenario: "subtraction", input: []token{
			{category: LITERAL, value: "5"},
			{category: LITERAL, value: "2"},
			{category: SUB_OPERATOR, value: "-"},
		}, x_value: 0, output: 3},
		{scenario: "multiplication", input: []token{
			{category: LITERAL, value: "2"},
			{category: LITERAL, value: "5"},
			{category: TIMES_OPERATOR, value: "*"},
		}, x_value: 0, output: 10},
		{scenario: "division", input: []token{
			{category: LITERAL, value: "10"},
			{category: LITERAL, value: "2"},
			{category: DIV_OPERATOR, value: "/"},
		}, x_value: 0, output: 5},
		{scenario: "one parenthesis", input: []token{
			{category: LITERAL, value: "4"},
			{category: LITERAL, value: "6"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "2"},
			{category: DIV_OPERATOR, value: "/"},
		}, x_value: 0, output: 5},
		{scenario: "two parenthesis", input: []token{
			{category: LITERAL, value: "4"},
			{category: LITERAL, value: "2"},
			{category: LITERAL, value: "3"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "2"},
			{category: DIV_OPERATOR, value: "/"},
		}, x_value: 0, output: 5},
		{scenario: "addition with x", input: []token{
			{category: LITERAL, value: "2"},
			{category: NAME, value: "x"},
			{category: ADD_OPERATOR, value: "+"},
		}, x_value: 5, output: 7},
		{scenario: "subtraction with x", input: []token{
			{category: LITERAL, value: "5"},
			{category: NAME, value: "x"},
			{category: SUB_OPERATOR, value: "-"},
		}, x_value: 2, output: 3},
		{scenario: "multiplication by x", input: []token{
			{category: LITERAL, value: "2"},
			{category: NAME, value: "x"},
			{category: TIMES_OPERATOR, value: "*"},
		}, x_value: 5, output: 10},
		{scenario: "division by x", input: []token{
			{category: LITERAL, value: "10"},
			{category: NAME, value: "x"},
			{category: DIV_OPERATOR, value: "/"},
		}, x_value: 2, output: 5},
	}

	t.Run(">>> test Polish Reverse evaluation", func(t *testing.T) {

		for _, test := range testScenarios {

			fmt.Printf("scenario: %s\n", test.scenario)

			//	conversion from input string to the postfix queue
			postfix := NewQueue()

			for i, _ := range test.input {
				postfix.Put(&test.input[i])
			}

			expr := &ParsedExpression{
				postfix: postfix,
			}

			//	create the symbol table
			symbolTable := NewFloatSymbolTable()

			symbolTable.SetValue("x", test.x_value)

			//	Polish Reverse evaluation of postfix expression
			want := test.output
			got, err := expr.Evaluate(symbolTable)
			if err != nil {
				t.Errorf("unexpected error converting from string -> postfix: %s", err)
				continue
			}
			fmt.Printf("[debug] postfix evaluation result: %f\n", got)

			//	check the result
			if want != got {
				t.Errorf("fail evaluating the postfix expression: expected: %f result: %f", want, got)
			}
		}
	})
}
