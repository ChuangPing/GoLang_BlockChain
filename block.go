package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

//	定义区块结构体
type Block struct {
	//	1.版本号
	Version uint64
	//	2.前区块的哈希  -- 哈希都用byte类型存储
	PrevHash []byte
	//	3.merkel根  -- 是一个哈希值，在V4版本时在实现
	MerkelRoot []byte
	//	4.时间戳
	TimeStamp uint64
	//	5.难度值
	Difficulty uint64
	//	6.随机值  --  挖矿要找的数值
	Nonce uint64

	//	a.当前区块的哈希值
	Hash []byte
	//	b.块体的数据 -- 交易信息
	Data []byte
}

//辅助函数 -- 将uint64转换为byte类型  --  因为在求当前区块的哈希时，需要将区块的所有信息进行拼接在取哈希，转换成相同的类型便于拼接
func Uint64ToByte(num uint64) []byte {
	// func Write(w io.Writer, order ByteOrder, data interface{}) error
	//将data的binary编码格式写入w，data必须是定长值、定长值的切片、定长值的指针。order指定写入数据的字节序，写入结构体时，名字中有'_'的字段会置为0。

	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panicln(err)
	}
	return buffer.Bytes()
}

//创世纪快
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{}, //	暂时给一个空的
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0, //困难值谁便给的
		Nonce:      0, //同上
		Hash:       []byte{},
		Data:       []byte(data),
	}
	block.SetHash()
	return &block
}

//生成当前区块哈希方法
func (block *Block) SetHash() {
	//var blockInfo []byte
	//1. 拼装数据  -- 常规数据拼接的写法
	/*
		blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
		blockInfo = append(blockInfo, block.PrevHash...)
		blockInfo = append(blockInfo, block.MerkelRoot...)
		blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
		blockInfo = append(blockInfo, Uint64ToByte(block.Difficulty)...)
		blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
		blockInfo = append(blockInfo, block.Data...)
	*/
	// 使用Strings 包中的jion方法拼接
	//	定义二维切片
	temp := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(block.Nonce),
		block.Data,
	}
	//	将二维的切片数组链接起来，返回一个一维的切片
	blockInfo := bytes.Join(temp, []byte{})
	//	使用拼接好的区块，使用SHA256取哈希
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
