package main

import (
	"fmt"
	"os"
)

//解析命令行输入的命令

//定义cli结构体
type CLI struct {
	// 由于这些命令主要添加区块，查看区块，因此需要操作数据库，因此需要当前区块bc拿到数据库的连接对象,
	bc *BlockChain
}

//定义提示常量
const Mesage = `
	addBlock --data DATA     "添加区块"
	printChain               "正向打印区块链"
`

//初始化CLI
func NewCLI(bc *BlockChain) *CLI {
	// 初始化CLi完成BlockC赋值 --拿到初始化BlockCHain的数据库连接对象等
	cli := CLI{
		bc: bc,
	}
	return &cli
}

//定义cli方法用于解析命令行参数，执行对应命令
func (cli *CLI) Run() {
	//	获取命令行输入
	args := os.Args
	//	校验数据
	if len(args) < 2 {
		fmt.Println(Mesage)
		return
	}
	//	分析命令 -- 由于上面的校验，使用数据组下标1 不存在越界
	commd := args[1]
	switch commd {
	case "addBlock":
		//	进行验证    ./block   addBlock    --data    "HelloWorld" -- 正确的命令格式
		if len(args) == 4 && args[2] == "--data" {
			//	添加区块命令正确
			fmt.Println("--- 执行添加区块命令 ---")
			data := args[3]
			cli.CommdAddBlock(data)
		}
	case "printBlock":
		//由于这一个命令没有参数，因此经过前面的校验，到这里肯定是输入了打印区块命令
		fmt.Println("--- 开始执行打印区块链命令 ---")
		cli.CommdPrintBlockChain()
	}
}
