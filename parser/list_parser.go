package simple_parser

import (
	"PyCompiler/inter"
	"PyCompiler/lexer"
	"errors"
	"fmt"
)

type SimpleParser struct {
	lexer       lexer.Lexer
	top         *Env        //当前作用域的符号表
	saved       *Env        //进入下一个作用域时，他用来记录当前符号表
	curTok      lexer.Token //当前读到字符串对应标签
	usedStorage uint32      //用于存储变量的内存大小
}

func NewSimpleParser(lexer lexer.Lexer) *SimpleParser {
	return &SimpleParser{
		lexer: lexer,
		top:   nil,
		saved: nil,
	}
}

func (s *SimpleParser) Parse() {
	s.program()
}

func (s *SimpleParser) program() {
	// program -> block
	// block -> "{" stmts "}"
	//stmt 其实是seq所形成的队列的头结点
	s.top = nil
	stmt := s.block()
	begin := stmt.NewLabel()
	after := stmt.NewLabel()
	stmt.EmitLabel(begin)
	stmt.Gen(begin, after)
	stmt.EmitLabel(after)
}

//当前的字符串是否匹配
func (s *SimpleParser) matchLexeme(str string) error {
	if s.lexer.Lexeme == str {
		return nil
	}
	// 当出现错误返回错误信息
	errS := fmt.Sprintf("error token , expected:%s , got:%s", str, s.lexer.Lexeme)
	return errors.New(errS)
}

//判断代码标签
func (s *SimpleParser) matchTag(tag lexer.Tag) error {
	if s.curTok.Tag == tag { //判断标签
		return nil
	}

	errS := fmt.Sprintf("error tag, expected:%d, got %d", tag, s.curTok.Tag)
	return errors.New(errS)
}

//往前读取一个字符串
func (s *SimpleParser) moveBackward() {
	s.lexer.ReverseScan()
}

//往后读取一个字符串
func (s *SimpleParser) moveForward() error {
	var err error
	s.curTok, err = s.lexer.Scan()
	return err
}

//判断代码块
func (s *SimpleParser) block() inter.StmtInterface {
	// block -> "{" decls stmts "}"
	err := s.moveForward()
	if err != nil {
		panic(err)
	}
	//判断字符串是否匹配
	err = s.matchLexeme("{")
	if err != nil {
		panic(err)
	}

	err = s.moveForward()
	if err != nil {
		panic(err)
	}
	// 判断语义操作，作用域的判断
	s.saved = s.top //记录，开启环境
	s.top = NewEnv(s.top)
	err = s.decls() //进行左递归
	if err != nil {
		panic(err)
	}

	stmt := s.stmts()
	if err != nil {
		panic(err)
	}

	err = s.matchLexeme("}")
	if err != nil {
		panic(err)
	}
	// 记录 ，用于{} 的作用域，将环境返回后来
	s.top = s.saved
	return stmt
}

// 左递归过程
func (s *SimpleParser) decls() error {
	/*
		decls -> decls decl | ε
		decls 表示由零个或多个decl组成，decl对应语句为:
		int a; float b; char c;等，其中int, float, char对应的标号为BASIC,
		在进入到这里时我们并不知道要解析多少个decl,一个处理办法就是判断当前读到的字符串标号，
		如果当前读到了BASIC标号，那意味着我们遇到了一个decl对应的声明语句，于是就执行decl对应的语法
		解析，完成后我们再次判断接下来读到的是不是还是BASIC标号，如果是的话继续进行decl解析，
		由此我们可以破除左递归
	*/
	for s.curTok.Tag == lexer.BASIC {
		err := s.decl()
		if err != nil {
			return err
		}
	}
	return nil
}

//获取类型
func (s *SimpleParser) getType() (*inter.Type, error) {
	err := s.matchTag(lexer.BASIC)
	if err != nil {
		return nil, err
	}

	width := uint32(4) //赋予 内存大小 的属性
	switch s.lexer.Lexeme {
	case "int":
		width = 4 //4个字节
	case "float":
		width = 8
	case "char":
		width = 1
	case "bool":
		width = 1
	}
	//构造新的类型
	p := inter.NewType(s.lexer.Lexeme, lexer.BASIC, width)
	s.usedStorage = s.usedStorage + width
	return p, nil
}

func (s *SimpleParser) decl() error {
	p, err := s.getType()
	if err != nil {
		return err
	}

	err = s.moveForward()
	if err != nil {
		return err
	}
	//这里必须复制，因为s.curTok会不断变化因此不能直接传入s.curTok
	tok := lexer.NewTokenWithString(s.curTok.Tag, s.lexer.Lexeme)
	id := inter.NewID(s.lexer.Line, tok, p) //符号
	sym := NewSymbol(id, p)                 //符号表
	s.top.Put(s.lexer.Lexeme, sym)

	err = s.moveForward() //读取下一个字符
	if err != nil {
		return err
	}

	err = s.matchLexeme(";") //分号结尾
	if err != nil {
		return err
	}

	err = s.moveForward() //读取下一个字符
	return err
}

func (s *SimpleParser) stmts() inter.StmtInterface {
	if s.matchLexeme("}") == nil {
		return inter.NewStmt(s.lexer.Line) //返回个stmt对象
	}

	//注意这里，seq节点通过递归形成了一个链表
	return inter.NewSeq(s.lexer.Line, s.stmt(), s.stmts())
}

func (s *SimpleParser) stmt() inter.StmtInterface {
	//当前只解析算术表达式
	return s.expression()
}

//表达式
func (s *SimpleParser) expression() inter.StmtInterface {
	//当前读到的是分号，常量，或者变量名(ID)就会进入到算术表达式的解析中
	if s.matchTag(lexer.ID) == nil {
		s.moveForward() //读取符号决定是赋值还是运算
		//算术表达式有两种情况，一种是赋值，其特点是有一个等号 c = a,另一种是运算 a + b
		if s.matchTag(lexer.ASSIGN_OPERATOR) == nil { //如果等于等号的话
			//就进行赋值解析
			s.moveBackward()
			s.moveBackward() //回退到变量名 , 返回两次从 （c = a）到 （c）
			return s.assign()
		}
		s.moveBackward() //如果不是的话 可能是 (a + b),所以放回到 (a + )
	}

	expression := inter.NewExpression(s.lexer.Line, s.expr())
	return expression
}

//赋值解析
func (s *SimpleParser) assign() inter.StmtInterface {
	s.moveForward()
	sym := s.top.Get(s.lexer.Lexeme) //判断符号表，看看是否被定义了
	if sym == nil {
		//引用了未定义变量
		errS := fmt.Sprintf("undefined variable with name: %s", s.lexer.Lexeme)
		err := errors.New(errS)
		panic(err)
	}

	s.moveForward()  //读取=
	s.moveForward()  //读取 = 后面的字符串
	expr := s.expr() // “=” 后面的都是表达式 例如： c = 1, c = a, c = a + b 得到 = 右边的东西
	set, err := inter.NewSet(sym.id, expr)
	if err != nil {
		panic(err)
	}
	err = s.matchLexeme(";")
	if err != nil {
		panic(err)
	}
	s.moveForward()
	expression := inter.NewExpression(s.lexer.Line, set)
	return expression
}

//解析运算符
func (s *SimpleParser) expr() inter.ExprInterface {
	x := s.term()
	var err error

	for s.matchLexeme("+") == nil || s.matchLexeme("-") == nil {
		tok := lexer.NewTokenWithString(s.curTok.Tag, s.lexer.Lexeme)
		s.moveForward()
		x, err = inter.NewArith(s.lexer.Line, tok, x, s.term())
		if err != nil {
			panic(err)
		}

	}

	return x
}

func (s *SimpleParser) term() inter.ExprInterface {
	x := s.factor()
	return x
}

func (s *SimpleParser) factor() inter.ExprInterface {
	var x inter.ExprInterface
	tok := lexer.NewTokenWithString(s.curTok.Tag, s.lexer.Lexeme)
	if s.matchTag(lexer.NUM) == nil {
		t := inter.NewType("int", lexer.BASIC, 4)
		x = inter.NewConstant(s.lexer.Line, tok, t)
	} else if s.matchTag(lexer.REAL) == nil {
		t := inter.NewType("float", lexer.BASIC, 8)
		x = inter.NewConstant(s.lexer.Line, tok, t)
	} else {
		sym := s.top.Get(s.lexer.Lexeme)
		if sym == nil {
			errS := fmt.Sprintf("undefined variable with name: %s", s.lexer.Lexeme)
			err := errors.New(errS)
			panic(err)
		}

		x = sym.id
	}

	s.moveForward()
	return x
}
