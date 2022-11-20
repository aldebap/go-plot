////////////////////////////////////////////////////////////////////////////////
//	symbol_test.go  -  Nov-19-2022  -  aldebap
//
//	Test cases for the symbol manager
////////////////////////////////////////////////////////////////////////////////

package expression

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

//	Test_lexicalAnalizer test cases for the lexical analizer
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
