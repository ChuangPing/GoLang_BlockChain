package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//	工作量证明结构体-- 挖矿
type ProofOfWorlk struct {
	block  *Block   // 当前需要进行挖矿的区块
	target *big.Int // 挖矿难度值
}

func NewProofOfWorlk(block *Block) *ProofOfWorlk {
	//根据初始化参数对ProofOfWorlk进行赋值
	pow := ProofOfWorlk{
		block: block,
	}
	targetStr := "000f0000000000000000000000000000000000000000000000000000000000"
	var targetTemp big.Int
	targetTemp.SetString(targetStr, 16)
	pow.target = &targetTemp
	return &pow
}

//挖矿方法
func (pow *ProofOfWorlk) Run() (uint64, []byte) {
	//	获取当前需要挖矿的区块
	block := pow.block
	var Nonce uint64 //随机值，默认为0
	fmt.Println("--- 挖矿开始 ---")
	for {
		//	使用join方法对挖矿区块进行拼接 -- 将区块拼接成byte不断取哈希
		blockTemp := [][]byte{
			block.MerkelRoot,
			//	在v4版本中，我们块体是交易，因此我们并不需要将区块交易添加取哈希，我们只需要将区块的交易两两去哈希，构建merkel根。把merkelRoot
			//添加进区块，当区块交易发送改变，会影响到melkerRoot因此我们可以用merkelRoot来取代块体的交易，因为交易发生变化会影响到merkelRoot从而影响区块哈希
			//block.Data,
			block.PreHash,
			block.Hash,
			Uint64ToByte(block.Version),
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.DIfficult),
			// 这个Nonce随着挖矿的进行不断调整
			Uint64ToByte(Nonce),
		}
		//当前区块信息组成的byte切片
		blockInfo := bytes.Join(blockTemp, []byte{})
		hash := sha256.Sum256(blockInfo)
		//中间变量 -- 将找到的哈希值转换为特殊的big.Int 方便与pow中目标值进行比较
		hashTemp := big.Int{}
		hashTemp.SetBytes(hash[:])
		if hashTemp.Cmp(pow.target) == -1 {
			fmt.Println("--- 挖矿成功 ---")
			//	挖矿成功将，找到的随机值，当前区块哈希返回
			return Nonce, hash[:]
		} else {
			//	没有找到，改变Nonce
			Nonce++
		}
	}
}
