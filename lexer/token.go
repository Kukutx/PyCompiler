package lexer

// Tag uint32 取别名 Tag
type Tag uint32

//Token 表    为Token设置不用的标志号 如 123 + 456   => NUM PLUS NUM  ，token，它其实是将源代码中的字符串进行分类，并给每种分类赋予一个整数值
const (
	AND Tag = iota + 256
	BREAK
	DO
	EQ
	FALSE
	GE
	ID
	IF
	ELSE
	INDEX
	LE
	INT
	FLOAT
	MINUS
	PLUS
	NE
	NUM
	OR
	REAL
	STRING
	TRUE
	WHILE
	LEFT_BRACE    // "{"
	RIGHT_BRACE   // "}"
	LEFT_BRACKET  //"("
	RIGHT_BRACKET //")"
	AND_OPERATOR
	OR_OPERATOR
	ASSIGN_OPERATOR
	NEGATE_OPERATOR
	LESS_OPERATOR
	GREATER_OPERATOR
	BASIC //对应int , float, bool, char 等类型定义
	TEMP  //对应中间代码的临时寄存器变量
	SEMICOLON

	EOF   //end of file  结束
	ERROR //错误
)

//token_map集合 初始化map
var tokenMap = make(map[Tag]string)

//token初始化
func init() {
	tokenMap[AND] = "&&"
	tokenMap[BREAK] = "break"
	tokenMap[BASIC] = "BASIC"
	tokenMap[DO] = "do"
	tokenMap[ELSE] = "else"
	tokenMap[EQ] = "=="
	tokenMap[FALSE] = "FALSE"
	tokenMap[GE] = "GE"
	tokenMap[ID] = "IDENTIFIER"
	tokenMap[IF] = "if"
	tokenMap[INT] = "int"
	tokenMap[FLOAT] = "float"
	tokenMap[INDEX] = "INDEX"
	tokenMap[LE] = "<="
	tokenMap[MINUS] = "-"
	tokenMap[PLUS] = "+"
	tokenMap[NE] = "!="
	tokenMap[NUM] = "NUM"
	tokenMap[OR] = "OR"
	tokenMap[REAL] = "REAL"
	tokenMap[STRING] = "String"
	tokenMap[TEMP] = "t"
	tokenMap[TRUE] = "TRUE"
	tokenMap[WHILE] = "while"
	tokenMap[AND_OPERATOR] = "&"
	tokenMap[OR_OPERATOR] = "|"
	tokenMap[ASSIGN_OPERATOR] = "="
	tokenMap[NEGATE_OPERATOR] = "!"
	tokenMap[LESS_OPERATOR] = "<"
	tokenMap[GREATER_OPERATOR] = ">"
	tokenMap[LEFT_BRACE] = "{"
	tokenMap[RIGHT_BRACE] = "}"
	tokenMap[LEFT_BRACKET] = "("
	tokenMap[RIGHT_BRACKET] = ")"
	tokenMap[EOF] = "EOF"
	tokenMap[ERROR] = "ERROR"
	tokenMap[SEMICOLON] = ";"
}

type Token struct {
	lexeme string
	Tag    Tag
}

// ToString 利用 map的键值返回对应Tag的字符串
func (t *Token) ToString() string {
	if t.lexeme == "" {
		return tokenMap[t.Tag]
	}
	return t.lexeme
}

// NewToken 实例化Token
func NewToken(tag Tag) Token {
	return Token{
		lexeme: "",
		Tag:    tag,
	}
}

func NewTokenWithString(tag Tag, lexeme string) *Token {
	return &Token{
		lexeme: lexeme,
		Tag:    tag,
	}
}
