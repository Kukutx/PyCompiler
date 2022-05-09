package lexer

/*关键词 token*/

// Word 单词结构体
type Word struct { //假设   int abc = 123;
	lexeme string // "abc"
	Tag    Token  // IDENTIFIER
}

// NewWordToken 返回一个WordToken对象
func NewWordToken(s string, tag Tag) Word {
	return Word{
		lexeme: s,
		Tag:    NewToken(tag),
	}
}

// ToString 返回lexeme内的字符串
func (w *Word) ToString() string {
	return w.lexeme
}

// GetKeyWords 提取关键字
func GetKeyWords() []Word {
	var keyWords []Word
	keyWords = append(keyWords, NewWordToken("&&", AND))
	keyWords = append(keyWords, NewWordToken("||", OR))
	keyWords = append(keyWords, NewWordToken("==", EQ))
	keyWords = append(keyWords, NewWordToken("!=", NE))
	keyWords = append(keyWords, NewWordToken("<=", LE))
	keyWords = append(keyWords, NewWordToken(">=", GE))
	keyWords = append(keyWords, NewWordToken("-", MINUS))
	keyWords = append(keyWords, NewWordToken("true", TRUE))
	keyWords = append(keyWords, NewWordToken("false", FALSE))
	keyWords = append(keyWords, NewWordToken("+", PLUS))
	keyWords = append(keyWords, NewWordToken("if", IF))
	keyWords = append(keyWords, NewWordToken("else", ELSE))

	//添加类型定义
	keyWords = append(keyWords, NewWordToken("int", BASIC))
	keyWords = append(keyWords, NewWordToken("float", BASIC))
	keyWords = append(keyWords, NewWordToken("bool", BASIC))
	keyWords = append(keyWords, NewWordToken("char", BASIC))
	return keyWords
}
