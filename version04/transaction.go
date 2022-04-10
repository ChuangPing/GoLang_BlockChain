package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

//	定义加一结构体 ,一个交易包含：输出交易：转出去的钱，输入交易：转入的钱
type Transaction struct {
	TXID      []byte     //	交易ID
	TXInputs  []TXInput  //	交易输入数组，因为一个交易有多个输入或者输出
	TXOutputs []TXOutput //	交易输出数组
}

//	定义交易输入类型结构体
type TXInput struct {
	//	引用的交易ID，相当于上一个向你地址转钱的交易ID -- 类似于区块哈希的做法
	TXid []byte
	//	引用的output的索引值，上一个向你地址进行转账的输出交易的索引
	Index int64
	//	解锁脚本的签名，version04中我们使用地址来模拟（转账人的地址） -- 相当于你自己用自己的余额，向别人转账前要验证一下这个钱是否是转给你的（地址是否是你）
	Sig string
}

//	定义交易输出类型结构体
type TXOutput struct {
	//	转账金额
	Value float64
	//	锁定脚本 -- version04使用地址进行模拟  你转账的地址
	PubKeyHash string
}

//	定义挖矿奖励  2022:6.25
const reward = 12.5

//	定义Transaction结构体方法，生成交易哈希
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	//	定义编码器
	encoder := gob.NewEncoder(&buffer)
	//	使用编码器进行编码
	err := encoder.Encode(tx)
	if err != nil {
		log.Panicln("编码器编码出错")
		return
	}
	//	获取到整个交易转成byte
	transactionByte := buffer.Bytes()
	//	取哈希
	hash := sha256.Sum256(transactionByte)
	tx.TXID = hash[:]
}

//	提供创建交易函数 -- 挖矿交易 简单：由于币（币是凭空产生），因此不用说明币的来源，此时交易的inputs可以自己指定
//	我们规定，input中：TXid为空，Index：-1，Sig：由于不用指定币的来源，因此旷工可以随意指定一些内容，列如写一下数据永久的保存在区块中
func NewCoinbaseTx(address string, data string) *Transaction {
	//挖矿交易的特点：--因此只需传入 旷工的地址，挖矿的币需要给到他账户，旷工指定写入区块的数据
	//1. 只有一个input
	//2. 无需引用交易id
	//3. 无需引用index
	//矿工由于挖矿时无需指定签名，所以这个sig字段可以由矿工自由填写数据，一般是填写矿池的名字
	input := TXInput{TXid: []byte{}, Index: -1, Sig: data}
	//	得到的钱是写入旷工的地址，因此是写在outputs中，表示账户收到的钱
	output := TXOutput{Value: reward, PubKeyHash: address}
	//	组装交易
	tx := Transaction{
		TXID:      []byte{},
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{output},
	}
	//	设置区块哈希
	tx.SetHash()
	return &tx
}
