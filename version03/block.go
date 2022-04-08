package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

//定义区块结构体
type Block struct {
	Version    uint64
	PreHash    []byte
	TimeStamp  uint64
	MerkelRoot []byte
	Difficult  uint64
	Nonce      uint64
	Data       []byte
	Hash       []byte //当前区块的哈希
}

//辅助函数 -- 将uint64转为byte切片
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panicln("Uint64ToByte err", err)
	}
	return buffer.Bytes()
}

//创建区块函数
func NewBlock(data string, preHash []byte) *Block {
	block := Block{
		Version:    03,
		PreHash:    preHash,
		TimeStamp:  uint64(time.Now().Unix()),
		Difficult:  0,
		MerkelRoot: []byte{},
		Data:       []byte(data),
		Hash:       []byte{},
	}
	//	进行挖矿产生区块，寻找符合难度要求的随机值
	// 1.初始化工作量证明函数
	pow := NewProofOfWork(&block)
	//	2.调用挖矿函数,成功返回当前区块的哈希和寻找的随机值
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

//将区块进行序列化的方法 -- 使用gob
func (block *Block) Serialize() []byte {
	// 注意这里是将整个当前区块结构体进行序列化，不是求哈希时需要将block所有内容使用字节的方式拼接在一起
	//	解码的数据放到buffer中
	var buffer bytes.Buffer
	//	1.定义编码器 -- 必须使用&buffer ，因为编码器在编码后没有返回值，结果字节存放在buffer地址中
	encoder := gob.NewEncoder(&buffer)
	//	2.使用编码器进行编码  -- 编码的结果以经存放在buffer中，不需要返回值，因为一开始就是使用了&buffer
	err := encoder.Encode(&block)
	if err != nil {
		log.Panicln("序列化时，编码出错", err)
	}
	return buffer.Bytes()
}

//解码函数，--不是Block的方法
func DeSerialize(data []byte) Block {
	//	定义需要解码的类型
	var block Block
	//	1.定义解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))
	//	2.使用解码器进行解码 -- 由于没有返回值，返回值只有err 因此要使用&block
	err := decoder.Decode(&block)
	if err != nil {
		log.Panicln("解码器解码出错", err)
	}
	return block
}
