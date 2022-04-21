package main

import "fmt"

//根据命令行的命令执行相应的业务

//定义添加区块方法
func (cli *CLI) CommdAddBlock(data string) {
	//cli.bc.AddBlock(data)
	fmt.Printf("--- 添加区块命令执行成功 ---\n\n")
}

//打印区块链方法
func (cli *CLI) CommdPrintBlockChain() {
	//	初始化迭代器
	bcIterator := cli.bc.NewBlockChainIterator()
	for {
		block := bcIterator.Next()
		//	使用迭代器循环读取区块
		fmt.Printf("---当前区块 ---\n")
		fmt.Printf("当前区块哈希：%x\n", block.Hash)
		fmt.Printf("当前区块Nonce：%d\n", block.Nonce)
		fmt.Printf("当前区块时间戳：%v\n", block.TimeStamp)
		fmt.Printf("区块数据 :%s\n", block.Transactions[0].TXInputs[0].Sig)
		//fmt.Printf("当前区块交易：%s\n", block.Data) // data实际是byte类型，但是打印时可以选择以字符的形式打印
		fmt.Printf("前一个区块的哈希：%x\n\n", block.PreHash)
		//	退出循环条件
		if len(block.PreHash) == 0 { // 如果当前区块存储的前一个哈希为空，说明遍历结束已经到创世区块
			fmt.Printf("区块链遍历结束！\n")
			break
		}
	}
	fmt.Printf("--- 打印区块命令执行成功 ---\n\n")
}

//	查看用户余额的方法
func (cli *CLI) CommdGetBalance(address string) {
	total := 0.0
	//	查询地址的utxo -- 与账户相关的outpus ，就是别人向他账户转钱，但是账户没有花费
	utxos := cli.bc.FindUtxos(address)
	//	遍历utxos 中的，循环累加outputs中的value值 -- 即为余额
	for _, utxo := range utxos {
		fmt.Printf("utxo:%s", utxo)
		total += utxo.Value
	}
	fmt.Printf("%s账户余额为%f\n", address, total)
	fmt.Println("--- 执行获取账户余额命令成功 ---\n\n")
}

//	执行转账交易方法
func (cli *CLI) sendTransaction(from, to string, amount float64, miner, data string) {
	//1. 创建挖矿交易
	coinBase := NewCoinbaseTx(miner, data)
	//2. 创建普通交易
	tx := NewTransaction(from, to, amount, cli.bc)
	if tx == nil {
		fmt.Println("无效交易请检查")
		//阻止代码继续执行 -- 不然会进行挖矿，进而旷工得到假交易
		return
	}
	//3 旷工打包区块
	transaction := []*Transaction{coinBase, tx}
	//4.将打包好的交易添加进区块 -- 发布区块，进行挖矿
	cli.bc.AddBlock(transaction)
	fmt.Println("--- 执行转账交易命令成功 ---")
}
