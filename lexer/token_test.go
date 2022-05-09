package lexer

// 测试
import (
	"github.com/stretchr/testify/require" //引入测试包
	"testing"
)

//对于github包需要 输入 go get github.com/stretchr/testify/tree/master/require 下载包权限
//然后运行 go test  //这里可参考 go语言的 test命令
func TestTokenName(t *testing.T) {
	//测试键值对的 键值和字符串是否正确
	indexToken := NewToken(REAL)
	require.Equal(t, "REAL", indexToken.ToString())

	realToken := NewToken(ID)
	require.Equal(t, "IDENTIFIER", realToken.ToString())
}
