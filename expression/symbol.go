////////////////////////////////////////////////////////////////////////////////
//	symbol.go  -  Nov-19-2022  -  aldebap
//
//	Symbol manager for expression evaluator
////////////////////////////////////////////////////////////////////////////////

package expression

import "errors"

//	types of symbols
const (
	UNKNOWN  uint8 = 0
	CONSTANT uint8 = 1
	VARIABLE uint8 = 2
	FUNCTION uint8 = 3
)

type SymbolTable interface {
	Exists(name string) bool
	SetValue(name string, value float64)
	GetValue(name string) (float64, error)

	DefineFunc(name string, function func(parameter ...float64) float64, params int)
	InvokeFunc(name string, parameter ...float64) (float64, error)
}

type floatSymbolTable struct {
	variable       map[string]float64
	function       map[string]func(parameter ...float64) float64
	functionParams map[string]int
}

//	New create a new float64 onlye symbol table
func NewFloatSymbolTable() SymbolTable {
	return &floatSymbolTable{
		variable:       make(map[string]float64),
		function:       make(map[string]func(parameter ...float64) float64),
		functionParams: make(map[string]int),
	}
}

//	Exists returns true if the symbol exists on the table
func (f *floatSymbolTable) Exists(name string) bool {
	_, exists := f.variable[name]

	return exists
}

//	SetValue set the value for a variable on the table
func (f *floatSymbolTable) SetValue(name string, value float64) {
	f.variable[name] = value
}

//	GetValue get the value for a variable from the table
func (f *floatSymbolTable) GetValue(name string) (float64, error) {
	value, exists := f.variable[name]

	if !exists {
		return 0, errors.New("unknown symbol name: " + name)
	}

	return value, nil
}

//	DefineFunc set the function associated to a symbol name on the table
func (f *floatSymbolTable) DefineFunc(name string, function func(parameter ...float64) float64, params int) {
	f.function[name] = function
	f.functionParams[name] = params
}

//	InvokeFunc invoke the function associated to a symbol name on the table
func (f *floatSymbolTable) InvokeFunc(name string, parameter ...float64) (float64, error) {
	function, exists := f.function[name]

	if !exists {
		return 0, errors.New("unknown function name: " + name)
	}

	if len(parameter) != f.functionParams[name] {
		return 0, errors.New("invalid number or parameters invoking function: " + name)
	}

	return function(parameter...), nil
}
