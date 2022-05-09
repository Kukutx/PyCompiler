package inter

type Expression struct {
	stmt *Stmt         //用来封装条件语句的父节点
	expr ExprInterface //用来封装运算符节点的父节点
}

func NewExpression(line uint32, expr ExprInterface) *Expression {
	return &Expression{
		stmt: NewStmt(line),
		expr: expr,
	}
}

func (e *Expression) Errors(str string) error {
	return e.stmt.Errors(str)
}

func (e *Expression) NewLabel() uint32 {
	return e.stmt.NewLabel()
}

func (e *Expression) EmitLabel(i uint32) {
	e.stmt.EmitLabel(i)
}

func (e *Expression) Emit(code string) {
	e.stmt.Emit(code)
}

func (e *Expression) Gen(start uint32, end uint32) {
	e.expr.Gen()
}
