package main

import "fmt"

//根据cli获取到的命令，执行相应的方法。从CLI文件抽离，使Cli文件只是解析命令，commandLie文件进行对应命令的执行

//添加区块方法
func (cli *CLI) AddBlocks(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("添加区块成功")
}

//打印区块方法
func (cli *CLI) PrintBlockChain() {
	//获取当前区块的实例对象 -- 迭代器需要操作数据库
	bc := cli.bc
	iterator := bc.NewBlockChainIterator()
	for {
		//	使用迭代器循环读取区块
		block := iterator.Next()
		fmt.Printf("---当前区块 ---\n")
		fmt.Printf("当前区块哈希：%x\n", block.Hash)
		fmt.Printf("当前区块Nonce：%d\n", block.Nonce)
		fmt.Printf("当前区块时间戳：%v\n", block.TimeStamp)
		fmt.Printf("当前区块交易：%s\n", block.Data) // data实际是byte类型，但是打印时可以选择以字符的形式打印
		fmt.Printf("前一个区块的哈希：%x\n\n", block.PreHash)
		//	退出循环条件
		if len(block.PreHash) == 0 { // 如果当前区块存储的前一个哈希为空，说明遍历结束已经到创世区块
			fmt.Printf("区块链遍历结束！\n")
			break
		}
	}
}
