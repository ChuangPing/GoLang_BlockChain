package main

import (
	"fmt"
	"os"
)

//对命令行的命令进行解析

type CLI struct {
	// 由于这些命令主要添加区块，查看区块，因此需要操作数据库，因此需要当前区块bc拿到数据库的连接对象
	bc *BlockChain
}

//定义提示常量
const Mesage = `
	addBlock --data DATA     "添加区块"
	printChain               "正向打印区块链"
`

//定义cli方法，用于处理接收到的参数
func (cli *CLI) Run() {
	//./block printChain  -- 输入block  printChain 实际commandLine读取到的数据
	//./block addBlock --data "HelloWorld"
	//	1. 获取到所有的命令
	args := os.Args
	//	2. 校验数据
	if len(args) < 2 {
		fmt.Println(Mesage)
		return
	}
	//	3. 分析命令
	cmd := args[1] // 使用下标1是安全的，因为前面我们已经对args的长度进行判断，不会出现越界
	switch cmd {
	case "addBlock":
		//	提示执行的命令
		//	进行验证    ./block   addBlock    --data    "HelloWorld"
		if len(args) == 4 && args[2] == "--data" {
			//说明按照 ./block   addBlock    --data    "HelloWorld"  指定格式正确输入
			//	获取添加区块的内容
			data := args[3] // 由于存在上面的验证不存在数组越界
			//	调用CLI的方法进行添加区块
			cli.AddBlocks(data)
		}
	case "printChain":
		//	走到这一步，经过前面的验证，说明至少输入了一个命令参数
		fmt.Println("执行打印区块命令")
		cli.PrintBlockChain()
	}
}
