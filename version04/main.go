package main

func main() {
	//初始化BlockChain
	bc := NewBlockChain("创世旷工00")
	//	初始化命令行解析结构体处理命令行命令
	cli := NewCLI(bc)
	cli.Run()
}
