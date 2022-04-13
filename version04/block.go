package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

//	定义区块结构体
type Block struct {
	Version    uint64
	TimeStamp  uint64
	Difficult  uint64
	MerkelRoot []byte
	Nonce      uint64
	//Data       []byte
	//	version04 版本中，区块体中的数据保存真正的交易
	Transactions []*Transaction
	PreHash      []byte
	Hash         []byte //当前区块的哈希
}

//	工具函数 -- 将uint64转为byte
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panicln("将uint64转byte出错", err)
		return nil
	}
	return buffer.Bytes()
}

//	初始化函数 -- 初始化区块函数,即创建区块函数  -- 创建一个区块：前一个区块哈希，区块中的交易
func NewBlock(preHash []byte, txs []*Transaction) *Block {
	//1.根据传入的参数，初始化要产生的区块
	block := Block{
		Version:    04,
		PreHash:    preHash,
		MerkelRoot: []byte{},
		Difficult:  0,
		TimeStamp:  uint64(time.Now().Unix()),
		//Data:       []byte(data),	v4版本进行修改
		Transactions: txs,
		Hash:         []byte{},
	}
	//	设置区块melkerTree 根哈希
	block.MerkelRoot = block.MakeMelkerRoot()
	//	2。进行挖矿，确定当前区块的哈希，以及随机值
	//初始化工作量证明函数
	pow := NewProofOfWorlk(&block)
	//进行挖矿
	Nonce, Hash := pow.Run()
	block.Hash = Hash
	block.Nonce = Nonce
	return &block
}

//区块序列化函数 -- 区块存入bolt数据库，Key：区块哈希 value: 区块byte
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	//	1.定义编码器 -- 使用gob包
	encoder := gob.NewEncoder(&buffer)
	//	2.使用编码器进行编码  -- 编码的结果以经存放在buffer中，不需要返回值，因为一开始就是使用了&buffer
	err := encoder.Encode(&block)
	if err != nil {
		log.Panicln("序列化时，编码出错", err)
	}
	return buffer.Bytes()
}

//工具函数 -- 解码函数，从bolt数据库中获取到区块，使用解码函数还原成Block类型
func DeSerialize(blockByte []byte) *Block {
	//	定义解码类型
	var block *Block
	//	初始化解码器
	decoder := gob.NewDecoder(bytes.NewReader(blockByte))
	//	使用解码器进行解码
	err := decoder.Decode(&block)
	if err != nil {
		log.Panicln("解码器解码出错", err)
		return nil
	}
	return block
}

//将区块交易组织成melkerTree,并求根哈希 -- 返回melekerTree 的根哈希
func (block *Block) MakeMelkerRoot() []byte {
	//	TODO
	return []byte{}
}
