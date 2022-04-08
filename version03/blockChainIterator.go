package main

import (
	"github.com/boltdb/bolt"
	"log"
)

//定义迭代器结构体
type BlockChainIterator struct {
	//	获取数据库连接对象 -- 因为需要在数据库中读数据
	db *bolt.DB
	//	游标，用于不断索引
	currentHashPointer []byte
}

//初始化区块链迭代器  -- 主要完成获得当前区块存储的数据库连接对象和当前游标 ,因此这个初始化函数必须写成BlockChain的方法
func (bc *BlockChain) NewBlockChainIterator() *BlockChainIterator {
	return &BlockChainIterator{
		db:                 bc.db,
		currentHashPointer: bc.tail, // 迭代器的游标指向最后一个元素的hash
	}
}

//迭代器是属于区块链的
//Next方式是属于迭代器的
//1. 返回当前的区块
//2. 指针前移
func (this *BlockChainIterator) Next() *Block {
	var block Block
	db := this.db // 获取数据库连接对象
	//	对数据库进行读操作
	db.View(func(tx *bolt.Tx) error {
		//找到抽屉 -- 没有找到直接报错
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panicln("迭代器遍历时bucket不应该为空，请检查!")
		}
		//通过最后一区块的hash读取区块，注意这个区块是序列化后的数据 -- byte类型
		blockTmp := bucket.Get(this.currentHashPointer)
		//对数据进行反序列化 -- 即将对应byte转换为相应的Block类型
		block = DeSerialize(blockTmp)
		// 游标向前移，即向左移 因为是重最后一一个区块开始遍历
		this.currentHashPointer = block.PreHash //此时前一个区块的哈希为当前取出来的区块(最后一个区块)所记录的前一个区块的哈希，即最新的最后一个区块
		return nil
	})
	return &block
}
