package main

//区块链结构体
type BlockChain struct {
	//	定义block类型数组  -- Block结构体在同一main包，可以相互引用
	blocks []*Block
}

//初始化函数
func NewBlockChain() *BlockChain {
	//初始化创世区块并加入区块链
	genesisBlock := GenesisBlock()
	return &BlockChain{
		blocks: []*Block{genesisBlock},
	}
}

//定义创世区块
func GenesisBlock() *Block {
	return NewBlock("创世区块00", []byte{})
}

//添加区块方法
func (bc *BlockChain) AddBlock(data string) {
	//如何获取前区块的哈希呢？？
	// 1.获取前一个区块即最后一个区块(相对于当前加入区块) 取数组最后一个下标元素
	lastBlock := bc.blocks[len(bc.blocks)-1]
	preHash := lastBlock.PrevHash
	//	创建区块
	block := NewBlock(data, preHash)
	//	将创建的区块添加到区块链中 -- 即添加到数组中
	bc.blocks = append(bc.blocks, block)
}
