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

//	定义账户需要的utxos方法 -- 找到满足转账要求的最低可用的utxo集合
func (bc *BlockChain) FindNeedUTXOs(from string, amount float64) (map[string][]uint64, float64) {
	//	定义满足需要的utxos
	utxos := make(map[string][]uint64)
	//	找到的余额
	var total float64
	//	调用FindUTXOTransactions 获取与账户相关包含utxo集合的交易
	txs := bc.FindUTXOTransactions(from)
	//	对账户所有包含utxo交易进行遍历
	for _, tx := range txs {
		for i, output := range tx.TXOutputs {
			if output.PubKeyHash == from {
				//	先判断，需要的余额（amount）与找到的余额（total）
				if total < amount {
					//	将outputs添加进账户可用的utxos  -- 这里可做改进，将交易Id tx.TXID也要返回，方便函数复用
					total += output.Value
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(i))
					if total > amount {
						//	已经找到足够的余额 -- 直接退出，提高效率
						return utxos, total
					}
				}
				if total < amount {
					//	正常情况走不到这一步，因为在里面就结束循环
					fmt.Printf("不满足转账金额,当前总额：%f， 目标金额: %f\n", total, amount)
					return utxos, total
				}

			}
		}
	}
	//// 调用FindUtxos 找到进行转账账户可用到的utxo -- 找自己可用的余额 -- 即找与自己相关的oupts交易
	//txs := bc.FindUtxos(from)  -- 这个函数不够好，没有返回交易的ID
	//// 遍历交易  对账户余额进行累加
	//for _, utxo := range txs {
	//	//	utxo的总余额小于转钱的总额 -- 钱不够继续寻找可用的utxo
	//	if total < amount {
	//		// 将使用的
	//	}
	//}
	return utxos, total
}

//	查找对应账户可用的utxo -- 未花费的交易. 注：为了方便并没有将整个交易返回，而是只返回TXOutPut数组，方便遍历计算
func (bc *BlockChain) FindUtxos(address string) []TXOutput {
	//	定义需要返回utxos数组  -- TXOutPut 数组
	var UTXO []TXOutput
	txs := bc.FindUTXOTransactions(address)
	for _, tx := range txs {
		for _, output := range tx.TXOutputs {
			if address == output.PubKeyHash {
				UTXO = append(UTXO, output)
			}
		}
	}

	return UTXO
}

//	定义寻找寻找账户utxo集合的方法
func (bc *BlockChain) FindUTXOTransactions(address string) []*Transaction {
	//	存储所有包含utxo交易集合
	var txs []*Transaction
	//	//我们定义一个map来保存消费过的output，key是这个output的交易id，value是这个交易中索引的数组  map[交易id][]int64
	spentOutputs := make(map[string][]int64)
	bcIterator := bc.NewBlockChainIterator()
	for {
		//	1,循环遍历整个区块链中的区块
		block := bcIterator.Next()
		//	2. 遍历每个区块中的交易中
		for _, tx := range block.Transactions {
			//
		OUTPUT:
			//	3. 遍历outPut，找到与账户相关的outPuts -- 找到自己账户的收入(--如果不做过滤，这样会找到所有转入账户的钱，也包括转入账户已经花过的钱)
			for i, output := range tx.TXOutputs {
				//	在这里做一个过滤，将所有消耗过的outputs和当前的所即将添加output对比一下 如果相同，则跳过，否则添加 如果当前的交易id存在于我们已经表示的map，那么说明这个交易里面有消耗过的output
				if spentOutputs[string(tx.TXID)] != nil {
					//	说明当前交易在花费的outputs里面 -- 花费过
					for _, spentOutputsIndex := range spentOutputs[string(tx.TXID)] {
						//	取当前花费过的交易的index 判断是否与output的index是否相等，若相等则说明这个outPut是当前账户消费过不进行utxo统计
						if int64(i) == spentOutputsIndex {
							continue OUTPUT
						}
					}
				}
				//	判断当前output是否和账户相关， 判断output中的PubKeyHash(地址) 是否和自己的相同
				if output.PubKeyHash == address {
					//	将与自己相关的utxo添加到UTXO
					//UTXO = append(UTXO, output)
					//	将包含账户的utxo的交易集合存放在交易切片中
					txs = append(txs, tx)
				}
			}
			//	4.遍历当前区块的input交易， --- 找出自己花过的钱 (即花过的outputs)
			//if !tx.IsCoinbase() { -- 虽然这样做会提高效率，但是这样会有BUG，因为每次对挖矿交易不做遍历，当挖矿交易做为下一个交易的输入，即挖矿交易的钱被花费，而有不对它进行遍历判断是否花费，导致获取当前账户余额时每次都会多加挖矿产生的钱
			//	判断当前交易是否为铸币交易，因为铸币交易只有一个TXInput 交易，而且还是认为自定义写的内容不用进行循环查找里面是否有花过的Output
			for _, input := range tx.TXInputs {
				//	判断当前input交易中是否有当前账户花过，因为花钱产生的交易是Inputs，且花钱人会进行签名证明这是自己的钱， --当前版本签名使用的是地址4，进行判断
				if input.Sig == address {
					//	将当前input交易添加到消费过的oupts交易 ，用于过滤。花费的outPuts就不用添加到可用的UTXO中
					//indexArray := spentOutputs[string(input.TXid)]
					//indexArray = append(indexArray, input.Index)
					spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index) // 与上面等价
				}
			}
			//} else {
			//	//fmt.Println("挖矿交易，不遍历Input交易")
			//}
		}
		if len(block.PreHash) == 0 {
			//	区块链遍历完成退出
			break
		}
	}
	return txs
}
