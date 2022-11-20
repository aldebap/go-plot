////////////////////////////////////////////////////////////////////////////////
//	symbol.go  -  Nov-19-2022  -  aldebap
//
//	Symbol manager for expression evaluator
////////////////////////////////////////////////////////////////////////////////

package expression

import "errors"

type SymbolTable interface {
	Exists(name string) bool
	SetValue(name string, _value float64)
	GetValue(name string) (float64, error)
}

type floatSymbolTable struct {
	symbol map[string]float64
}

//	New create a new float64 onlye symbol table
func NewFloatSymbolTable() SymbolTable {
	return &floatSymbolTable{
		symbol: make(map[string]float64),
	}
}

//	Exists returns true if the symbol exists on the table
func (f *floatSymbolTable) Exists(name string) bool {
	_, exists := f.symbol[name]

	return exists
}

//	SetValue set the value for a symbol on the table
func (f *floatSymbolTable) SetValue(name string, value float64) {
	f.symbol[name] = value
}

//	GetValue get the value for a symbol from the table
func (f *floatSymbolTable) GetValue(name string) (float64, error) {
	value, exists := f.symbol[name]

	if !exists {
		return 0, errors.New("unknown symbol name: " + name)
	}

	return value, nil
}
