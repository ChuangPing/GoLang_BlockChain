package main

import "fmt"

func main() {
	// 初始化区块链 -- 在内部完成创世区块的创建，并返回区块链对象，可以添加区块
	bc := NewBlockChain()
	//	添加区块不用在指定前一个区块哈希，在内部已经实现，区块的产生工作量证明都在内部实现，对外只暴露简单的接口
	bc.AddBlock("Alice转账100BTC到Bob！")
	bc.AddBlock("Bob向Alice转账300BTC！")
	//	循环打印当前区块链中的区块
	for index, block := range bc.blocks {
		fmt.Printf("------当前区块高度：%d-------\n\n", index)
		fmt.Printf("前区块的哈希：%x\n", block.PreHash)
		fmt.Printf("当前区块的哈希：%x\n", block.Hash)
		fmt.Printf("当前区块的交易：%s\n", block.Data)
		fmt.Printf("当前区块的随机值：%d\n", block.Nonce)
		fmt.Printf("------当前区块结束：-------\n")
		fmt.Println()
	}
}
