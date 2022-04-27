package main

func main() {
	//初始化BlockChain
	bc := NewBlockChain("1uit8P6KgUZepKk7hCPMmSiCy8YNWGKFk")
	//	初始化命令行解析结构体处理命令行命令
	cli := NewCLI(bc)
	cli.Run()
}
