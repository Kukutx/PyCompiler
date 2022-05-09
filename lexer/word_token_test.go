package lexer

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWordToken(t *testing.T) {
	word := NewWordToken("variable", ID)          //返回一个Token结构体对象，变量等标识符一般都是 ID (IDENTIFIER)
	require.Equal(t, "variable", word.ToString()) //判断字符串是否正确
	wordTag := word.Tag
	require.Equal(t, wordTag.ToString(), "IDENTIFIER") //判断标识号
}

func TestKeyWords(t *testing.T) {
	keyWords := GetKeyWords()
	require.Equal(t, len(keyWords), 12) //判断keyWords长度，也就是关键词绑定有多少个

	andKeyWord := keyWords[0]
	require.Equal(t, andKeyWord.ToString(), "&&")

	orKeyWord := keyWords[1]
	require.Equal(t, orKeyWord.ToString(), "||")
}
