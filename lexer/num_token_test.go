package lexer

import (
	"github.com/stretchr/testify/require" //引入测试包
	"testing"
)

func TestNumToken(t *testing.T) {
	numToken := NewNumToken(123)
	require.Equal(t, numToken.value, 123)
	require.Equal(t, numToken.ToString(), "123")
	numTag := numToken.Tag
	require.Equal(t, numTag.ToString(), "NUM")
}

func TestRealToken(t *testing.T) {
	realToken := NewRealToken(3.1415926)
	require.Equal(t, realToken.value, 3.1415926)
	require.Equal(t, realToken.ToString(), "3.1415926")

	var realTag = realToken.Tag
	require.Equal(t, realTag.ToString(), "REAL")
}
