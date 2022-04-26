package main

import (
	"fmt"
	"os"
	"strconv"
)

//解析命令行输入的命令

//定义cli结构体
type CLI struct {
	// 由于这些命令主要添加区块，查看区块，因此需要操作数据库，因此需要当前区块bc拿到数据库的连接对象,
	bc *BlockChain
}

//定义提示常量
const Mesage = `
	请您输入以下命令实现对区块链进行操作：
	printChain               "正向打印区块链"
	getBalance --address ADDRESS "获取指定地址的余额"
	send FROM TO AMOUNT MINER DATA	"由from转mountBTC给to, 由旷工miner挖矿，同时miner旷工可以自定义写入data到区块"
	newWallet   "创建一个新的钱包(私钥公钥对)"
	listAddresses "列举所有的钱包地址"
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
	case "printChain":
		//由于这一个命令没有参数，因此经过前面的校验，到这里肯定是输入了打印区块命令
		fmt.Println("--- 开始执行打印区块链命令 ---")
		cli.CommdPrintBlockChain()
	case "getBalance":
		//getBalance --address ADDRESS "获取指定地址的余额"  -- 输入的命令：v4.exe getBalance --address ADDRESS "获取指定地址的余额"
		if len(args) == 4 && args[2] == "--address" {
			fmt.Println("--- 开始执行获取账户余额命令 ---")
			address := args[3]
			cli.CommdGetBalance(address)
		} else {
			fmt.Println("请检查您输入的命令", Mesage)
		}
	case "send":
		//	v4.exe send FROM TO AMOUNT MINER DATA
		if len(args) == 7 && args[1] == "send" {
			fmt.Println("开始执行创建转账交易命令")
			from := args[2]
			to := args[3]
			//	将string转为float64
			amount, _ := strconv.ParseFloat(args[4], 64)
			miner := args[5]
			data := args[6]
			cli.sendTransaction(from, to, amount, miner, data)
		} else {
			fmt.Println("请检查您输入的命令", Mesage)
		}
	case "newWallet":
		fmt.Println("--- 开始执行创建钱包命令 ---")
		cli.CommdNewWallet()
	case "listAddresses":
		fmt.Printf("列举所有地址...\n")
		cli.CommdListAddresses()
	default:
		fmt.Printf("无效的命令，请检查!\n")
		fmt.Printf(Mesage)
	}
}
