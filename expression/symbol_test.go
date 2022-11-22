////////////////////////////////////////////////////////////////////////////////
//	symbol_test.go  -  Nov-19-2022  -  aldebap
//
//	Test cases for the symbol manager
////////////////////////////////////////////////////////////////////////////////

package expression

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
)

//	Test_SymbolTable test cases for the variable's symbol table
func Test_SymbolTable(t *testing.T) {

	//	a few test cases
	var testScenarios = []struct {
		scenario     string
		symbols      []string
		testedSymbol string
		result       bool
		resultValue  float64
	}{
		{scenario: "empty table", symbols: []string{}, testedSymbol: "x", result: false},
		{scenario: "single non existent symbol table", symbols: []string{"y=10"}, testedSymbol: "x", result: false},
		{scenario: "single existent symbol table", symbols: []string{"x=10"}, testedSymbol: "x", result: true, resultValue: 10},
		{scenario: "twin symbol table", symbols: []string{"x=10", "y=20"}, testedSymbol: "x", result: true, resultValue: 10},
		{scenario: "twin non existent symbol table", symbols: []string{"x=10", "y=20"}, testedSymbol: "z", result: false},
		{scenario: "multiple symbol table", symbols: []string{"x=10", "y=20", "z=30"}, testedSymbol: "z", result: true, resultValue: 30},
	}

	t.Run(">>> test the symbol table exists() func", func(t *testing.T) {

		for _, test := range testScenarios {

			fmt.Printf("scenario: %s\n", test.scenario)

			//	add symbols
			symbolTable := NewFloatSymbolTable()

			for _, symbolDefinition := range test.symbols {

				values := strings.Split(symbolDefinition, "=")
				number, err := strconv.ParseFloat(values[1], 64)
				if err == nil {
					symbolTable.SetValue(values[0], number)
				}
			}

			want := test.result
			got := symbolTable.Exists(test.testedSymbol)

			if got != want {
				t.Errorf("fail in symbol table exists: expected: %t result: %t", want, got)
			}
		}
	})

	t.Run(">>> test the symbol table getValue() func", func(t *testing.T) {

		for _, test := range testScenarios {

			if !test.result {
				continue
			}

			fmt.Printf("scenario: %s\n", test.scenario)

			//	add symbols
			symbolTable := NewFloatSymbolTable()

			for _, symbolDefinition := range test.symbols {

				values := strings.Split(symbolDefinition, "=")
				number, err := strconv.ParseFloat(values[1], 64)
				if err == nil {
					symbolTable.SetValue(values[0], number)
				}
			}

			want := test.resultValue
			got, _ := symbolTable.GetValue(test.testedSymbol)

			if got != want {
				t.Errorf("fail in symbol table exists: expected: %f result: %f", want, got)
			}
		}
	})
}

//	Test_FunctionTable test cases for the function's symbol table
func Test_FunctionTable(t *testing.T) {

	symbolTable := NewFloatSymbolTable()

	t.Run(">>> test the symbol table DefineFunc() func", func(t *testing.T) {

		//	add functions
		symbolTable.DefineFunc("sin", func(x ...float64) float64 {
			return math.Sin(x[0])
		}, 1)

		symbolTable.DefineFunc("cos", func(x ...float64) float64 {
			return math.Cos(x[0])
		}, 1)
	})

	t.Run(">>> test the symbol table InvokeFunc() func", func(t *testing.T) {

		//	function not found
		fmt.Printf("scenario: function not found\n")

		want := "unknown function name: tan"
		got := ""
		_, err := symbolTable.InvokeFunc("tan", 0)
		if err != nil {
			got = err.Error()
		}

		if got != want {
			t.Errorf("fail in symbol table InvokeFunc: expected: %s result: %s", want, got)
		}

		//	invalid number of parameters (less)
		fmt.Printf("scenario: invalid number of parameters - less\n")

		want = "invalid number or parameters invoking function: sin"
		got = ""
		_, err = symbolTable.InvokeFunc("sin")
		if err != nil {
			got = err.Error()
		}

		if got != want {
			t.Errorf("fail in symbol table InvokeFunc: expected: %s result: %s", want, got)
		}

		//	invalid number of parameters (more)
		fmt.Printf("scenario: invalid number of parameters - more\n")

		want = "invalid number or parameters invoking function: sin"
		got = ""
		_, err = symbolTable.InvokeFunc("sin", 0, 3.141592)
		if err != nil {
			got = err.Error()
		}

		if got != want {
			t.Errorf("fail in symbol table InvokeFunc: expected: %s result: %s", want, got)
		}

		//	valid function invokation
		fmt.Printf("scenario: valid function invokation\n")

		wantFloat := math.Sin(3.141592 / 2)
		gotFloat, err := symbolTable.InvokeFunc("sin", 3.141592/2)
		if err != nil {
			t.Errorf("fail in symbol table InvokeFunc: unexpected error: %s", err.Error())
		}

		if gotFloat != wantFloat {
			t.Errorf("fail in symbol table InvokeFunc: expected: %f result: %f", wantFloat, gotFloat)
		}
	})
}
