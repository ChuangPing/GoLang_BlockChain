package main

import (
	//"GoLang_BlockChain/version1/GoLang_BlockChain/version01"
	"fmt"
)

func main() {
	// 初始化区块链 -- 创世区块已经在内部实现
	//bc := NewBlockChain()
	bc := NewBlockChain()
	bc.AddBlock("Alice转账100BTCBob！")
	bc.AddBlock("Bob转账300BTCAlice!")
	for index, block := range bc.blocks {
		fmt.Printf("------当前区块高度：%d-------\n\n", index)
		fmt.Printf("前区块哈希：%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希：%x\n", block.Hash)
		fmt.Printf("当前区块数据：%s\n", block.Data)
		fmt.Printf("------当前区块结束：-------\n")
		fmt.Println()
	}
}
