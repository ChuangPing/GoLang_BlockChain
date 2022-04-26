package main

import (
	"github.com/boltdb/bolt"
	"log"
)

//BlockChain迭代器 -- 按照区块哈希在数据库中对区块进行遍历

//定义迭代器结构体
type BlockchainIterator struct {
	//	获取数据库连接对象 -- 因为需要在数据库中读数据
	db *bolt.DB
	//	游标，在迭代时用于循环读取区块
	currentHashPointer []byte
}

//初始化迭代器方法 -- BlockChain的方法
func (bc *BlockChain) NewBlockChainIterator() *BlockchainIterator {
	bcIterator := BlockchainIterator{
		db:                 bc.db,
		currentHashPointer: bc.tail,
	}
	return &bcIterator
}

//定义迭代器迭代方法  -- Next
func (bcIterator *BlockchainIterator) Next() *Block {
	var block *Block
	//	获取数据库连接对象
	db := *bcIterator.db
	//	获取最后一个区块哈希
	lastBlockHash := bcIterator.currentHashPointer
	//	读取数据库
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("blockBucket"))
		if bucket == nil {
			log.Panicln("在迭代器遍历区块链时bucket不应该为空,请检查！！！")
			return nil
		}
		//	区块存入数据库是经过序列化 -- byte类型
		blockByte := bucket.Get(lastBlockHash)
		//	反序列化 -- 这也是为什么没有将反序列化函数写区块Block的方法的原因， 需要解码随时用
		block = DeSerialize(blockByte)
		//	将迭代器的游标向前移动
		bcIterator.currentHashPointer = block.PreHash
		return nil
	})
	return block
}
