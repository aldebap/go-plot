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

//	originaly implemented in github.com/aldebap/algorithms_dataStructs/chapter_3/expression

//	Test_infix2postfix test cases for the conversion from infix -> postfix
func Test_infix2postfix(t *testing.T) {

	//	a few test cases
	var testScenarios = []struct {
		scenario string
		input    string
		output   string
	}{
		{scenario: "addition", input: "2 + 5", output: "2 5 +"},
		{scenario: "subtraction", input: "5 - 2", output: "5 2 -"},
		{scenario: "multiplication", input: "2 * 5", output: "2 5 *"},
		{scenario: "division", input: "10 / 2", output: "10 2 /"},
		{scenario: "one parenthesis", input: "( 4 + 6 ) / 2", output: "4 6 + 2 /"},
		{scenario: "two parenthesis", input: "( 4 + ( 2 * 3 ) ) / 2", output: "4 2 3 * + 2 /"},
		{scenario: "unbalanced parenthesis", input: "( 4 + 6 / 2", output: "expression with unbalanced parenthesis"},
		{scenario: "fix to parenthesis parsing", input: "( x * x ) - ( 3 * x ) + 2", output: "x x * 3 x * - 2 +"},
		{scenario: "fix to parenthesis parsing", input: "x * x - ( 3 * x ) + 2", output: "x x * 3 x * - 2 +"},
	}

	t.Run(">>> test conversion from infix -> postfix", func(t *testing.T) {

		for _, test := range testScenarios {

			fmt.Printf("scenario: %s\n", test.scenario)

			//	execute conversion from infix -> postfix
			postfix, err := infix2postfix(test.input)
			if err != nil {
				if err.Error() != test.output {
					t.Errorf("unexpected error converting from infix -> postfix: %s", err)
				}
				continue
			}

			want := test.output
			got := ""

			for {
				item := postfix.Get()
				if item == nil {
					break
				}
				got += " " + item.(string)
			}
			got = strings.TrimLeft(got, " ")
			fmt.Printf("[debug] postfix result: %s\n", got)

			//	check the result
			if want != got {
				t.Errorf("fail converting from infix -> postfix: expected: %s result: %s", want, got)
			}
		}
	})
}

//	Test_infix2postfixV2 test cases for the conversion from infix -> postfix
func Test_infix2postfixV2(t *testing.T) {

	//	a few test cases
	var testScenarios = []struct {
		scenario string
		input    string
		output   string
	}{
		{scenario: "addition", input: "2 + 5", output: "2 5 +"},
		{scenario: "subtraction", input: "5 - 2", output: "5 2 -"},
		{scenario: "multiplication", input: "2 * 5", output: "2 5 *"},
		{scenario: "division", input: "10 / 2", output: "10 2 /"},
		{scenario: "one parenthesis", input: "( 4 + 6 ) / 2", output: "4 6 + 2 /"},
		{scenario: "two parenthesis", input: "( 4 + ( 2 * 3 ) ) / 2", output: "4 2 3 * + 2 /"},
		{scenario: "unbalanced parenthesis", input: "( 4 + 6 / 2", output: "expression with unbalanced parenthesis"},
		{scenario: "fix to parenthesis parsing", input: "( x * x ) - ( 3 * x ) + 2", output: "x x * 3 x * - 2 +"},
		{scenario: "fix to parenthesis parsing", input: "x * x - ( 3 * x ) + 2", output: "x x * 3 x * - 2 +"},
	}

	t.Run(">>> test conversion from infix -> postfix", func(t *testing.T) {

		_ = testScenarios
		/*
			for _, test := range testScenarios {

				fmt.Printf("scenario: %s\n", test.scenario)

				//	execute conversion from infix -> postfix
				postfix, err := infix2postfixV2(test.input)
				if err != nil {
					if err.Error() != test.output {
						t.Errorf("unexpected error converting from infix -> postfix: %s", err)
					}
					continue
				}

				want := test.output
				got := ""

				for {
					item := postfix.Get()
					if item == nil {
						break
					}
					got += " " + item.(string)
				}
				got = strings.TrimLeft(got, " ")
				fmt.Printf("[debug] postfix result: %s\n", got)

				//	check the result
				if want != got {
					t.Errorf("fail converting from infix -> postfix: expected: %s result: %s", want, got)
				}
			}
		*/
	})
}

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
		{scenario: "using a function call", input: "x*sin(2*x)", output: []string{"x", "*", "sin", "(", "2", "*", "x", ")"}},
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
			{category: ADD_OPERATOR, value: "*"},
			{category: LITERAL, value: "5"},
			{category: ADD_OPERATOR, value: "*"},
			{category: LITERAL, value: "7"},
		}, output: "2 5 7 * *"},
		{scenario: "using precedence", input: []token{
			{category: LITERAL, value: "2"},
			{category: ADD_OPERATOR, value: "*"},
			{category: LITERAL, value: "5"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "8"},
		}, output: "2 5 * 8 +"},
		{scenario: "more precedence", input: []token{
			{category: LITERAL, value: "2"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "5"},
			{category: ADD_OPERATOR, value: "*"},
			{category: LITERAL, value: "8"},
		}, output: "2 5 8 * +"},
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
			{category: LITERAL, value: "x"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "x"},
			{category: CLOSE_PARENTHESIS, value: ")"},
			{category: SUB_OPERATOR, value: "-"},
			{category: OPEN_PARENTHESIS, value: "("},
			{category: LITERAL, value: "3"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "x"},
			{category: CLOSE_PARENTHESIS, value: ")"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "2"},
		}, output: "x x * 3 x * - 2 +"},
		{scenario: "using precedence", input: []token{
			{category: LITERAL, value: "x"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "x"},
			{category: SUB_OPERATOR, value: "-"},
			{category: LITERAL, value: "3"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "x"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "2"},
		}, output: "x x * 3 x * - 2 +"},
		{scenario: "using a variable name", input: []token{
			{category: LITERAL, value: "x"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "4"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "y"},
			{category: SUB_OPERATOR, value: "-"},
			{category: LITERAL, value: "z"},
			{category: TIMES_OPERATOR, value: "*"},
			{category: LITERAL, value: "z"},
		}, output: "x 4 y * + z z * -"},
		{scenario: "using a variable name with underscore", input: []token{
			{category: LITERAL, value: "var_x"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "2"},
		}, output: "var_x 2 +"},

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
		}, output: "syntax error: expected token 7"},
		{scenario: "unbalanced parenthesis", input: []token{
			{category: OPEN_PARENTHESIS, value: "("},
			{category: LITERAL, value: "4"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "6"},
			{category: DIV_OPERATOR, value: "/"},
			{category: LITERAL, value: "2"},
		}, output: "syntax error: expected token 8"},
		{scenario: "close parenthesis before opening it", input: []token{
			{category: LITERAL, value: "4"},
			{category: ADD_OPERATOR, value: "+"},
			{category: LITERAL, value: "6"},
			{category: OPEN_PARENTHESIS, value: ")"},
			{category: DIV_OPERATOR, value: "/"},
			{category: LITERAL, value: "2"},
		}, output: "expression with unbalanced parenthesis"},
	}

	t.Run(">>> test tokens found by lexical analizer", func(t *testing.T) {

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
					got += " " + item.(string)
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

//	Test_evaluatePolishReverse test cases for the Polish Reverse evaluation function
func Test_evaluatePolishReverse(t *testing.T) {

	//	a few test cases
	var testScenarios = []struct {
		scenario string
		input    string
		x_value  float64
		output   float64
	}{
		{scenario: "addition", input: "2 5 +", x_value: 0, output: 7},
		{scenario: "subtraction", input: "5 2 -", x_value: 0, output: 3},
		{scenario: "multiplication", input: "2 5 *", x_value: 0, output: 10},
		{scenario: "division", input: "10 2 /", x_value: 0, output: 5},
		{scenario: "one parenthesis", input: "4 6 + 2 /", x_value: 0, output: 5},
		{scenario: "two parenthesis", input: "4 2 3 * + 2 /", x_value: 0, output: 5},
		{scenario: "addition with x", input: "2 x +", x_value: 5, output: 7},
		{scenario: "subtraction with x", input: "5 x -", x_value: 2, output: 3},
		{scenario: "multiplication by x", input: "2 x *", x_value: 5, output: 10},
		{scenario: "division by x", input: "10 x /", x_value: 2, output: 5},
	}

	t.Run(">>> test Polish Reverse evaluation", func(t *testing.T) {

		for _, test := range testScenarios {

			fmt.Printf("scenario: %s\n", test.scenario)

			//	conversion from input string to the postfix queue
			postfix := NewQueue()

			for _, item := range strings.Split(test.input, " ") {
				if len(item) > 0 {
					postfix.Put(item)
				}
			}

			expr := &ParsedExpression{
				postfix: postfix,
			}

			//	Polish Reverse evaluation of postfix expression
			want := test.output
			got, err := expr.Evaluate(test.x_value)
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
