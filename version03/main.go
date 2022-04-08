package main

func main() {
	//初始化区块链
	bc := NewBlockChain()
	//初始化命令行结构体 -- 使用命令行完成  下面的操作
	cli := CLI{
		bc: bc,
	}
	// 运行cli的读取命令行函数
	cli.Run()

	//弊端：每次运行时都会添加相同的区块，添加多次
	//向区块链中添加区块
	//bc.AddBlock("Alice向Bob转账100BTC")
	//bc.AddBlock("Bob向Alice装转200BTC")
	////初始化迭代器
	//iterator := bc.NewBlockChainIterator()
	//for {
	//	block := iterator.Next()
	//	fmt.Printf("---当前区块 ---\n")
	//	fmt.Printf("当前区块哈希：%x\n", block.Hash)
	//	fmt.Printf("当前区块Nonce：%d\n", block.Nonce)
	//	fmt.Printf("当前区块时间戳：%v\n", block.TimeStamp)
	//	fmt.Printf("当前区块交易：%s\n", block.Data) // data实际是byte类型，但是打印时可以选择以字符的形式打印
	//	fmt.Printf("前一个区块的哈希：%x\n\n", block.PreHash)
	//	//	退出循环条件
	//	if len(iterator.currentHashPointer) == 0 {
	//		break
	//	}
	//}
	////	由于我们将区块链存储在数据库中，因此我们在查看区块链中的每一个区块时，就不能使用循环，需要自己写一个迭代器，完成对区块链的遍历

}
