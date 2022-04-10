package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//定义区块链结构体
type BlockChain struct {
	// bolt 数据库的连接对象
	db *bolt.DB
	//	最后一个区块哈希
	tail []byte
}

//定义常量，区块存储的数据库名称（数据库文件名称），存储的抽屉（数据库名）
const blockChaindb = "blockChain.db"
const blockBucket = "blockBucket"

//初始化BlockChain 函数
func NewBlockChain(address string) *BlockChain {
	//处理初始化需要的数据
	var lastBlockHash []byte
	//	1.创建数据库 -- 存储区块的数据库
	//	a.打开数据库，数据库文件名不存在会在当前目录创建数据库文件
	db, err := bolt.Open(blockChaindb, 0600, nil)
	if err != nil {
		log.Panicln("打开/创建数据库失败", err)
		return nil
	}
	//	b.使用连接对象打开抽屉（bucket）
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//	说明区块链刚初始化，没有添加任何区块，抽屉都没有创建，也没有创世区块
			//	1）添加抽屉
			bucket, err := tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panicln("创建抽屉失败", err)
				return nil
			}
			//	2)添加创世区块
			genesisBlock := GenesisBlock(address)
			//	3）将创世区块存入数据库抽屉中
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			//	4) 将最后区块的哈希存入数据库 -- 这样才可以按照哈希链的方式访问数据库  -- 重要
			bucket.Put([]byte("lastBlockHash"), genesisBlock.Hash)
			//	5） 更新内存中的最后一个区块哈希，因为BlockChain结构体初始化需要赋值
			lastBlockHash = genesisBlock.Hash
		} else {
			//	说明不是一次初始化区块链，直接获取一下此时区块链数据库存储的而最后一个哈希
			lastBlockHash = bucket.Get([]byte("lastBlockHash"))
		}
		return nil
	})
	// 初始化结构体字段
	blockChain := BlockChain{
		db:   db,
		tail: lastBlockHash,
	}
	return &blockChain
}

//初始化创世区块函数
func GenesisBlock(address string) *Block {
	//	调用产生挖矿交易方法，并指定矿工地址，因为挖矿交易产生的币会转到这个地址上
	conbaseTx := NewCoinbaseTx(address, "version04创世区块00")
	block := NewBlock([]byte{}, []*Transaction{conbaseTx})
	fmt.Println("%v", block)
	return block
}

//添加区块方法
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	//	获取当前BolckChain数据库的连接对象
	db := bc.db
	//	获取当前区块链最后一个区块哈希
	lastBlockHash := bc.tail
	//	添加区块 -- 不能写在这里，因为这里就会挖矿，如果数据库出现异常，就会导致挖矿没用
	//	block := NewBlock(lastBlockHash, data)  -- 坑
	//	将组装好的区块存储在区块链数据库中
	db.Update(func(tx *bolt.Tx) error {
		//	拿到要操作的抽屉
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panicln("添加区块时bucket不应该为空，请检查！！")
			return nil
		}
		block := NewBlock(lastBlockHash, txs)
		err := bucket.Put(block.Hash, block.Serialize())
		err = bucket.Put([]byte("lastBlockHash"), block.Hash)
		if err != nil {
			log.Panicln("向数据库中插入区块数据时报错", err)
		}
		//更新一下当前运行时内存中的BlockChain结构体中最后一个区块哈希  -- 不跟新的话下一次调用AddBlock，读取的lastHash := bc.tail 没有发生变化，会导致调佣 NewBlock(data, lastHash)出错
		bc.tail = block.Hash
		return nil
	})
}
