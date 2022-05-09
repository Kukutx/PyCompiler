package inter

import (
	"errors"
	"fmt"
	"strconv"
)

var labels uint32 //用于实现跳转的标号

type Node struct {
	lexLine uint32 // 当前解析的行数
}

func NewNode(line uint32) *Node {
	labels = 0
	return &Node{
		lexLine: line,
	}
}

// Errors 语法错误
func (n *Node) Errors(s string) error {
	errS := "\nnear line " + strconv.FormatUint(uint64(n.lexLine), 10) + s //如果有错误打印代码出错误位置的位置，在那一行
	return errors.New(errS)
}

func (n *Node) NewLabel() uint32 {
	labels = labels + 1
	return labels
}

// EmitLabel 跳转
func (n *Node) EmitLabel(i uint32) {
	fmt.Print("\nL" + strconv.FormatUint(uint64(i), 10) + ":\n")
}

func (n *Node) Emit(s string) {
	fmt.Print("\t" + s)
}
