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