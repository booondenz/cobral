package main

import "github.com/learn-go/cidr/cmd"

// main 只做一件事：把控制权交给 cmd 包。
// 真正的命令行逻辑（命令树、参数、帮助输出）全都在 cmd/ 目录里。
func main() {
	cmd.Execute()
}
