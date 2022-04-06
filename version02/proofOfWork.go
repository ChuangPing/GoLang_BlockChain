package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//定义工作量证明结构体
type ProofOfWork struct {
	block *Block
	// 目标值 -- 挖矿难度值,它是一个非常大的数，因此不用一般的int ,big.Int : 一个非常大的数，它存在非常多的方法，例如：比较，赋值等方法
	target *big.Int
}

//初始化工作量证明函数
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	//我们指定的难度值，现在是一个string类型，需要进行转换
	targetStr := "000f0000000000000000000000000000000000000000000000000000000000"
	// 引入赋值变量，目的将上面定义的字符串转换成bigInt
	tempInt := big.Int{}
	//	将难度阈值赋值给bigInt，指定进制为16
	tempInt.SetString(targetStr, 16)
	pow.target = &tempInt
	return &pow
}

//	计算Hash方法（工作量证明的方法）  -- 挖矿过程函数  -- 返回：当前区块的哈希，找到的随机值
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	//	1.拼接数据（区块的数据，还有不断变化的随机数） -- 只有这样，整个区块的哈希才能不断变化，最终达到目标阈值
	//	2.做哈希运算
	//	3.每一次的哈希结果与pow中的目标阈值进行比较
	//	a.找到，退出返回
	//	b.未找到，继续改变nonce继续找
	var nonce uint64 //初始化自动为0
	block := pow.block
	//	循环  -- 找到才退出
	fmt.Printf("---挖矿开始---\n")
	for {
		// 拼接数据
		temp := [][]byte{
			Uint64ToByte(block.Version),
			block.PreHash,
			block.MerkelRoot,
			block.Data,
			Uint64ToByte(block.Difficulty),
			Uint64ToByte(block.TimeStamp),
			//	block里的Nonce是永远不变
			//Uint64ToByte(block.Nonce),
			//	随着循环nonce一直在变
			Uint64ToByte(nonce),
		}
		//	将二维byte切片数组连接起来，返回一个以为数组切片
		blockInfo := bytes.Join(temp, []byte{})
		//	做哈希运算
		hash := sha256.Sum256(blockInfo)
		//	与pow中的目标阈值进行比较 -- 定义一个中间变量，将hash转换为bigInt
		tempInt := big.Int{}
		tempInt.SetBytes(hash[:])
		//比较当前的哈希与目标哈希值，如果当前的哈希值小于目标的哈希值，就说明找到了，否则继续找

		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if tempInt.Cmp(pow.target) == -1 {
			//	z找到
			fmt.Printf("挖矿成功！hash:%x, nonce:%d\n", hash, nonce)
			return hash[:], nonce
		} else {
			//fmt.Printf("%d\n", nonce)
			//fmt.Printf("%x\n", hash)
			nonce++
		}
	}
}
