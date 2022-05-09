package simple_parser

import "PyCompiler/inter"

// Symbol 符号表
type Symbol struct {
	id       *inter.ID
	exprType *inter.Type
}

// NewSymbol 实例化符号表
func NewSymbol(id *inter.ID, exprType *inter.Type) *Symbol {
	return &Symbol{
		id:       id,
		exprType: exprType,
	}
}
