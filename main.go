package main

import (
	"PyCompiler/lexer"
	"PyCompiler/parser"
)

func main() {

	////语法树实现代码生成 示例：e = (a + b) - (c + d)
	//exprType := inter.NewType("int", lexer.BASIC, 4)
	//idA := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "a"), exprType)
	//idB := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "b"), exprType)
	//// a + b -> Arith
	//arith1, _ := inter.NewArith(1, lexer.NewTokenWithString(lexer.PLUS, "+"), idA, idB)
	////c + d
	//idC := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "c"), exprType)
	//idD := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "d"), exprType)
	////c + d -> Arith
	//arith2, _ := inter.NewArith(1, lexer.NewTokenWithString(lexer.PLUS, "+"), idC, idD)
	//// Arith (a + b) - (c + d)
	//arith3, _ := inter.NewArith(1, lexer.NewTokenWithString(lexer.PLUS, "-"), arith1, arith2)
	/////e 节点
	//idE := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "e"), exprType)
	////e = (a + b) - (c + d)
	//set, _ := inter.NewSet(idE, arith3)
	//set.Gen()

	source := `{int x; float y ; float c; float d;
	              x = 1; y = 3.14;
	              c = x + y;
	              d = x + y + c;
	              }`
	myLexer := lexer.NewLexer(source)
	parser := simple_parser.NewSimpleParser(myLexer)
	parser.Parse()

}
