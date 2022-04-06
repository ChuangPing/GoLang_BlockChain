package main

//定义区块链结构体
type BlockChain struct {
	//	定义block类型数组  -- Block结构体在同一main包，可以相互引用
	blocks []*Block
}

//初始化函数
func NewBlockChain() *BlockChain {
	//	初始化的过程：；将创世区块加入到区块链中 TODO
	genesisionBlock := GenesisionBlock()
	return &BlockChain{
		blocks: []*Block{genesisionBlock},
	}
}

//创世区块
func GenesisionBlock() *Block {
	return NewBlock("version2创世区块00", []byte{})
}

//添加区块的方法
func (bc *BlockChain) AddBlock(data string) {
	//	获取前一个区块 -- 为了获取前一个区块的哈希
	preBlock := bc.blocks[len(bc.blocks)-1]
	//	获取前一个区块的哈希
	preHash := preBlock.Hash
	//	创建区块
	newBlock := NewBlock(data, preHash)
	bc.blocks = append(bc.blocks, newBlock)

}
