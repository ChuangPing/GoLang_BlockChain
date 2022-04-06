package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"
)

type Block struct {
	Version    uint64
	PreHash    []byte
	TimeStamp  uint64
	MerkelRoot []byte
	Difficulty uint64
	Nonce      uint64
	Hash       []byte
	Data       []byte
}

//辅助函数 -- 将uint64转byte类型
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic("Uint64ToByte err :", err)
	}
	return buffer.Bytes()
}

//创建区块构造函数
func NewBlock(data string, preHash []byte) *Block {
	block := Block{
		Version:    11,
		PreHash:    preHash,
		TimeStamp:  uint64(time.Now().Unix()),
		MerkelRoot: []byte{},
		Difficulty: 0,
		Nonce:      0,
		Data:       []byte(data),
		Hash:       []byte{}, // 由于当前区块的哈希需要计算，因此在创建时先赋值为空byte切片，后面就行挖矿在进行赋值
	}
	// 此时还不能将block返回，还要进行挖矿。挖矿成功返回
	//	创建工作量证明对象
	pow := NewProofOfWork(&block)
	//	开始挖矿 -- 成后返回当前区块的哈希，以及找到的随机值
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}
