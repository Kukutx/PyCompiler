package inter

import (
	"PyCompiler/lexer"
)

type Type struct {
	width  uint32 //用多少字节存储该类型变量，列如 int,用4个字节来存储，float用8个字节来存储
	tag    lexer.Tag
	Lexeme string // int float bool 字符串
}

// NewType 实例化
func NewType(lexeme string, tag lexer.Tag, w uint32) *Type {
	return &Type{
		width:  w,
		tag:    tag,
		Lexeme: lexeme,
	}
}

// Numberic 类型转换只发送在int,float,char他们都对应数值类型
func Numberic(p *Type) bool {
	//查看给定类型是否属于数值类，判断类型是否为数值类型
	numberic := false
	switch p.Lexeme {
	case "int":
		numberic = true
	case "float":
		numberic = true
	case "char":
		numberic = true
	}

	return numberic
}

// MaxType 类型提升
func MaxType(p1 *Type, p2 *Type) *Type {
	/*
		float > int > char,例如p1是int，p2是float, 那么就提升为float ,类型提升必须对数值类型才有效,
		a + b 这里指令需要两个变量都是相同类型才能进行操作
	*/

	if Numberic(p1) == false && Numberic(p2) == false {
		return nil
	}
	//如果两者有其一是float类型，那么就提升为float，要不然就是int
	if p1.Lexeme == "float" || p2.Lexeme == "float" {
		return NewType("float", lexer.BASIC, 8)
	} else if p1.Lexeme == "int" || p2.Lexeme == "int" {
		return NewType("int", lexer.BASIC, 4)
	}

	return NewType("char", lexer.BASIC, 1)
}
