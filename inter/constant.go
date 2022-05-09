package inter

import (
	"PyCompiler/lexer"
)

// Constant 常量类  ,代表编程语言里的常量,在编译阶段的代码生成中该怎么处理
type Constant struct {
	expr *Expr
}

// NewConstant 实例化
func NewConstant(line uint32, token *lexer.Token, expr_type *Type) *Constant {
	constant := &Constant{
		expr: NewExpr(line, token, expr_type),
	}
	return constant
}

// GetConstantTrue 定义两个常量 true和false
func GetConstantTrue() *Constant {
	tok := lexer.NewToken(lexer.TRUE)
	trueType := NewType("bool", lexer.TRUE, 1)
	return NewConstant(0, &tok, trueType)
}

func GetConstantFalse() *Constant {
	tok := lexer.NewToken(lexer.FALSE)
	falseType := NewType("bool", lexer.FALSE, 1)
	return NewConstant(0, &tok, falseType)
}

func (c *Constant) Errors(s string) error {
	return c.expr.Errors(s)
}

func (c *Constant) NewLabel() uint32 {
	return c.expr.NewLabel()
}

func (c *Constant) EmitLabel(l uint32) {
	c.expr.EmitLabel(l)
}

func (c *Constant) Emit(code string) {
	c.expr.Emit(code)
}

func (c *Constant) Gen() ExprInterface {
	return c
}

func (c *Constant) Reduce() ExprInterface {
	return c
}

func (c *Constant) Type() *Type {
	return c.expr.Type()
}

func (c *Constant) ToString() string {
	return c.expr.ToString()
}
