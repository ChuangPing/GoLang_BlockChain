package main

import (
	"fmt"
	"os"
)

//go 读取命令行参数test

func main() {
	len1 := len(os.Args) // os.Arg 能获取到命令行输入的参数，是以数组形式
	fmt.Println("命令行输入的参数长度：", len1)
	for index, value := range os.Args {
		fmt.Printf("arg[%d]=%s\n", index, value)
	}
}
