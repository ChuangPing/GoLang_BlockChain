package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//定义工作量证明结构体
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

//初始化工作量证明函数
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
		//target:  -- 需要经过处理后才能赋值，不在这里赋值
	}
	targetStr := "000f0000000000000000000000000000000000000000000000000000000000"
	tempTarget := big.Int{}
	tempTarget.SetString(targetStr, 16)
	pow.target = &tempTarget
	return &pow
}

//不断计算当前区块的哈希方法 -- 挖矿,返回当前区块的哈希， 以及满足条件的随机值
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var nonce uint64
	//	获取当前挖矿的区块
	block := pow.block
	fmt.Println("---挖矿开始 ---")
	for {
		temp := [][]byte{
			block.MerkelRoot,
			block.PreHash,
			block.Hash,
			block.Data,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Version),
			Uint64ToByte(block.Difficult),
			Uint64ToByte(nonce), // 这个none会随着循环的过程变化
		}
		//	将二维byte切片数组连接起来，返回一个以为数组切片
		blockInfo := bytes.Join(temp, []byte{})
		//	将区块整体进行哈希运算
		hash := sha256.Sum256(blockInfo)
		//	将hsah  byte类型 转换为bigInt 类型 方便比较
		tempInt := big.Int{}
		tempInt.SetBytes(hash[:])
		//	比较当前的哈希与目标哈希值，如果当前的哈希值小于目标的哈希值，就说明找到了，否则继续找
		if tempInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功，hash：%x, nonce: %d\n", hash, nonce)
			return hash[:], nonce
		} else {
			nonce++
		}
	}
	//	将block区块进行拼接

}
