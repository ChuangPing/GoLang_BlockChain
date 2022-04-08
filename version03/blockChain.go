package main

import (
	"github.com/boltdb/bolt"
	"log"
)

//定义区块链结构体
type BlockChain struct {
	//	在以前版本中，我们是将区块存储在数组中，这个版本使用bolt数据库存储
	db   *bolt.DB
	tail []byte // 储存最后一个区块的哈希
}

//定义常量，区块存储的数据库名称（数据库文件名称），存储的抽屉（数据库名）
const blockChaindb = "blockChain.db"
const blockBucket = "blockBucket"

//初始化区块链
func NewBlockChain() *BlockChain {
	//	创建创世区块
	genesisBlock := GenesisBlock()
	//	最后一个区块哈希 -- 从数据库中进行读取
	var lasthash []byte
	//	1.打开数据库 -- 一开始没有这个数据库文件，会自动创建blockChaindb.db 数据库文件
	db, err := bolt.Open(blockChaindb, 0600, nil)
	if err != nil {
		log.Panicln("打开数据库失败")
	}
	//	2. 操作数据库 -- 将生成的区块出入数据库中（存入特定的抽屉中，类似于表）
	db.Update(func(tx *bolt.Tx) error {
		//找到存放数据的抽屉，没有就创建 （---这个没有不会自动创建，需要判断后自动创建）
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//	抽屉没有，自动创建
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panicln("创建Bucket失败", err)
			}
			//	3. 将创建的创世区块存储在数据库中：key：使用区块的hash，value：当前区块的转换成的字节, 因此要写一个将结构体转成字节的方法（Block的方法） -- 序列化，v2版本使用的bianry，当前使用gob
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			//	4.将最后一个区块哈希存储在数据库中，需要记录最后一区块哈希  -- 即 当前区块的哈希
			bucket.Put([]byte("lastHashKey"), genesisBlock.Hash)
			lasthash = genesisBlock.Hash // 创世区块的哈希就是当前区块的最后一个哈希
		} else {
			//	说明区块链已经初始化，抽屉都已经创建好了，直接从数据库中拿最后一个区块哈希
			lasthash = bucket.Get([]byte("lastHashKey"))
		}

		return nil
	})
	return &BlockChain{
		db:   db, //当前区块链数据的连接对象
		tail: lasthash,
	}
}

//创世区块函数
func GenesisBlock() *Block {
	block := NewBlock("version03创世区块00", []byte{})
	return block
}

//添加区块方法
func (bc *BlockChain) AddBlock(data string) {
	//获取存储区块的数据库连接对象   -- 因为添加区块要存入数据库中，需要操作数据库
	db := bc.db
	//获取最后一个区块的哈希，添加完成后需要在数据库和内存中同时更新最后一个区块的哈希
	lastHash := bc.tail
	db.Update(func(tx *bolt.Tx) error {
		//找到抽屉，如果未找到直接报错 -- 要操作数据库先要找到抽屉 -- 因为这时区块链已经初始化完成，数据库已经初始化完成
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panicln("bucket不应该为空，请检查")
		}
		//创建新的区块
		block := NewBlock(data, lastHash)
		//将创建好的区块添加到数据库中  hash作为key， block的字节流作为value
		bucket.Put(block.Hash, block.Serialize())
		//当有新的区块添加进数据库中，需要更新最后一个区块hash，就是当前区块的hash
		bucket.Put([]byte("lastHashKey"), block.Hash)
		//更新一下当前运行时内存中的BlockChain结构体中最后一个区块哈希  -- 不跟新的话下一次调用AddBlock，读取的lastHash := bc.tail 没有发生变化，会导致调佣 NewBlock(data, lastHash)出错
		bc.tail = block.Hash
		return nil
	})
}
