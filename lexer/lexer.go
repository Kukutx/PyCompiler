package lexer

import (
	"bufio"
	"strconv"
	"strings"
	"unicode"
)

//Lexer 词法分析器结构体
type Lexer struct {
	Lexeme      string   //用来存储当前扫描到的字符所存储的字符串
	lexemeStack []string //语法栈，用来存储临时的词法，用来进行压栈（比如运行语法树的运算，逆波兰表达式），在本代码里每次识别token或字符串时都需要压栈，然后压栈回滚
	tokenStack  []Token
	peek        byte          //读入的字符
	Line        uint32        //当前字符串所处的第几行
	reader      *bufio.Reader //用于读取字节流
	readPointer int
	keyWords    map[string]Token //存储关键词
}

// NewLexer 实例化词法分析器
func NewLexer(source string) Lexer {
	str := strings.NewReader(source)                      //构造reader对象，以便读取单个字符
	sourceReader := bufio.NewReaderSize(str, len(source)) //依靠这个对象读取单个字符
	lexer := Lexer{
		Line:     uint32(1),
		reader:   sourceReader,
		keyWords: make(map[string]Token),
	}
	lexer.reserve() //保留所有关键字
	return lexer    //返回 lexer结构体
}

// ReverseScan 当前读到的字符全部返回回去
func (l *Lexer) ReverseScan() {
	//backLen := len(l.Lexeme)
	//for i := 0; i < backLen; i++ { //我们用l.Lexeme 依次返回UnreadByte缓冲区里面去
	//	l.reader.UnreadByte() // 取消已读取的最后一个字节（即把字节重新放回读取缓冲区的前部）。只有最近一次读取的单个字节才能取消读取
	//}
	////赋值当前的栈顶元素 ，然后获取下一个元素
	//l.lexemeStack = l.lexemeStack[:(len(l.lexemeStack) - 1)] //弹出当前元素
	//l.Lexeme = l.lexemeStack[len(l.lexemeStack)-1]           //赋值当前的栈顶元素

	if l.readPointer > 0 {
		l.readPointer = l.readPointer - 1
	}
}

// 保存所有关键字
func (l *Lexer) reserve() {
	keyWords := GetKeyWords()          //存储所有关键字
	for _, keyWord := range keyWords { //遍历所有关键字
		l.keyWords[keyWord.ToString()] = keyWord.Tag //赋值所有关键字
	}
}

// ReadCh 读取字符，有错误就返回
func (l *Lexer) ReadCh() error {
	char, err := l.reader.ReadByte() //提前读取下一个字符，会清除后面的
	l.peek = char
	return err //有错误就返回
}

// ReadCharacter 判断当前读到的字符是否是给定的字符
func (l *Lexer) ReadCharacter(c byte) (bool, error) {
	chars, err := l.reader.Peek(1) //从缓冲区里面读取首个字符，读出给定数量的字符，但是不会从缓冲区里面清除掉
	//然后读入错误就返回
	if err != nil {
		return false, err
	}
	//如果读到的字符跟输入的字符不一样的话
	peekChar := chars[0]
	if peekChar != c {
		return false, nil
	}
	l.ReadCh() //越过当前peek的字符
	return true, nil
}

// UnRead 把读到的字符重新返回缓冲区里面
func (l *Lexer) UnRead() error {
	return l.reader.UnreadByte() //每次ReadCh()都会把字符从缓冲区里拿出来，现在我们又得放回去
}

// Scan 扫描  词法扫描就是依次读入源代码字符，然后看所读到的字符到底属于那种类型的字符串，然后将字符串与给定类型的token对象联系起来
func (l *Lexer) Scan() (Token, error) {
	if l.readPointer < len(l.lexemeStack) {
		l.Lexeme = l.lexemeStack[l.readPointer]
		token := l.tokenStack[l.readPointer]
		l.readPointer = l.readPointer + 1
		return token, nil
	} else {
		l.readPointer = l.readPointer + 1
	}

	for {
		err := l.ReadCh() //读取下一个字符
		if err != nil {
			return NewToken(ERROR), err
		}

		if l.peek == ' ' || l.peek == '\t' { //如果是空格或者是水平制表符
			continue //就跳过
		} else if l.peek == '\n' { //如果是换行
			l.Line = l.Line + 1 //行数 + 1
		} else {
			break //就跳出
		}
	}

	l.Lexeme = "" //置空

	// 判断符号
	switch l.peek {
	case ';':
		l.Lexeme = ";"
		l.lexemeStack = append(l.lexemeStack, l.Lexeme) //记录当前的符号
		token := NewToken(SEMICOLON)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case '{':
		l.Lexeme = "{"
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(LEFT_BRACE)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case '}':
		l.Lexeme = "}"
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(RIGHT_BRACE)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case '+':
		l.Lexeme = "+"
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(PLUS)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case '-':
		l.Lexeme = "-"
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(MINUS)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case '(':
		l.Lexeme = "("
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(LEFT_BRACKET)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case ')':
		l.Lexeme = ")"
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(RIGHT_BRACKET)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case '&': // 判断 & 是单个还是 && 双个
		l.Lexeme = "&"
		if ok, err := l.ReadCharacter('&'); ok { //读入下一个字符且判断下一个字符是否还是 &
			l.Lexeme = "&&"
			word := NewWordToken("&&", AND) //如果是 && 就返回对应的 Token
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, word.Tag)
			return word.Tag, err
		} else {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			token := NewToken(AND_OPERATOR)
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}
	case '|':
		l.Lexeme = "|"
		if ok, err := l.ReadCharacter('|'); ok {
			l.Lexeme = "||"
			word := NewWordToken("||", OR)
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, word.Tag)
			return word.Tag, err
		} else {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			token := NewToken(OR_OPERATOR)
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}

	case '=':
		l.Lexeme = "="
		if ok, err := l.ReadCharacter('='); ok {
			l.Lexeme = "=="
			word := NewWordToken("==", EQ)
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, word.Tag)
			return word.Tag, err
		} else {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			token := NewToken(ASSIGN_OPERATOR)
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}

	case '!':
		l.Lexeme = "!"
		if ok, err := l.ReadCharacter('='); ok {
			l.Lexeme = "!="
			word := NewWordToken("!=", NE)
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, word.Tag)
			return word.Tag, err
		} else {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			token := NewToken(NEGATE_OPERATOR)
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}

	case '<':
		l.Lexeme = "<"
		if ok, err := l.ReadCharacter('='); ok {
			l.Lexeme = "<="
			word := NewWordToken("<=", LE)
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, word.Tag)
			return word.Tag, err
		} else {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			token := NewToken(LESS_OPERATOR)
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}

	case '>':
		l.Lexeme = ">"
		if ok, err := l.ReadCharacter('='); ok {
			l.Lexeme = ">="
			word := NewWordToken(">=", GE)
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, word.Tag)
			return word.Tag, err
		} else {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			token := NewToken(GREATER_OPERATOR)
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}
	}

	//判断读入的字符串是否是数字
	if unicode.IsNumber(rune(l.peek)) {
		//如果读入数字，那么就一直读下去直到没有数字位置
		var v int
		var err error
		for {
			num, err := strconv.Atoi(string(l.peek)) //将数字的字符转换成对应的数字，当在peek读到的字符是数字的时候
			if err != nil {
				if l.peek != 0 { //l.peek == 0 意味着已经读完所有字符
					l.UnRead() //将字符放回以便下次扫描
				}
				break
			}
			v = 10*v + num             //转换成对应的十进制之数
			l.Lexeme += string(l.peek) //就把读到的数字记录下来
			l.ReadCh()
		}

		if l.peek != '.' { //如果没有点的话
			l.lexemeStack = append(l.lexemeStack, l.Lexeme) //记录到语法栈里
			token := NewToken(NUM)                          //那么返回一个整型
			token.lexeme = l.Lexeme
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}
		l.Lexeme += string(l.peek)
		l.ReadCh() //越过 "."
		//如果是浮点数
		x := float64(v)
		d := float64(10)
		for {
			l.ReadCh()
			num, err := strconv.Atoi(string(l.peek)) //转换数字
			if err != nil {
				if l.peek != 0 { //l.peek == 0 意味着已经读完所有字符
					l.UnRead() //将字符放回以便下次扫描
				}
				break
			}
			x = x + float64(num)/d //小数点后面的数字进行转换
			d = d * 10
			l.Lexeme += string(l.peek)
		}
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(REAL) //返回一个浮点数
		token.lexeme = l.Lexeme
		l.tokenStack = append(l.tokenStack, token)
		return token, err
	}

	//判断是否读到变量字符串，注意读到关键字要抽出来
	if unicode.IsLetter(rune(l.peek)) { //判断是否字符
		var buffer []byte //临时字符串缓存，用来存储读到的字符形成变量字符串或者关键词来进行判断
		for {
			buffer = append(buffer, l.peek) //追加字符
			l.Lexeme += string(l.peek)
			l.ReadCh()                           //读取下一个字符
			if !unicode.IsLetter(rune(l.peek)) { //如果读取到的不是字符
				if l.peek != 0 { //l.peek == 0 意味着已经读完所有字符
					l.UnRead() //将字符放回以便下次扫描
				}
				break
			}
		}
		s := string(buffer)        //将所有读到的字符合成字符串并赋值
		token, ok := l.keyWords[s] //看看字符串是否关键字
		if ok {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, token)
			return token, nil //是就返回token
		}
		l.lexemeStack = append(l.lexemeStack, l.Lexeme) //每次识别token或字符串时都需要压栈，然后压栈回归（语法解析需要）
		token = NewToken(ID)
		token.lexeme = l.Lexeme
		l.tokenStack = append(l.tokenStack, token)
		return NewToken(ID), nil //不是返回标识符
	}
	return NewToken(EOF), nil //返回 退出符
}
