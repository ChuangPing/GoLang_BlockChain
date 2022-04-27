package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
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
	//Sig string
	//真正的数字签名，由r，s拼成的[]byte  -- version05
	Signature []byte
	//约定，这里的PubKey不存储原始的公钥，而是存储X和Y拼接的字符串，在校验端重新拆分（参考r,s传递）  -- 注意，是公钥，不是哈希，也不是地址
	PubKey []byte
}

//	定义交易输出类型结构体
type TXOutput struct {
	//	转账金额
	Value float64
	//	锁定脚本 -- version04使用地址进行模拟  你转账的地址
	//PubKeyHash string
	//收款方的公钥的哈希，注意，是哈希而不是公钥，也不是地址 -- version05
	PubKeyHash []byte
}

//	定义挖矿奖励  2022:6.25
const reward = 12.5

//	由于现在存储的字段是地址的公钥哈希，所以无法直接创建TXOutput，  为了能够得到公钥哈希，我们需要处理一下，写一个Lock函数
func (outPut *TXOutput) Lock(address string) {
	//1. 解码
	//2. 截取出公钥哈希：去除version（1字节），去除校验码（4字节）
	//真正的锁定动作！！！！！
	outPut.PubKeyHash = GetPubKeyHashFromAddress(address)

}

//	TXOutput初始化函数
func NewTXOutput(value float64, address string) *TXOutput {
	output := TXOutput{
		Value: value,
	}
	output.Lock(address)
	return &output
}

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
	//矿工由于挖矿时无需指定签名，所以这个PubKey字段可以由矿工自由填写数据(data)，一般是填写矿池的名字 -- version05
	//签名先填写为空，后面创建完整交易后，最后做一次签名即可 -- version05
	input := TXInput{TXid: []byte{}, Index: -1, Signature: nil, PubKey: []byte(data)}
	//	得到的钱是写入旷工的地址，因此是写在outputs中，表示账户收到的钱
	//output := TXOutput{Value: reward, PubKeyHash: address}
	//	初始化output
	output := NewTXOutput(reward, address)
	//	组装交易
	tx := Transaction{
		TXID:     []byte{},
		TXInputs: []TXInput{input},
		//	output初始化后是指针
		TXOutputs: []TXOutput{*output},
	}
	//	设置区块哈希
	tx.SetHash()
	return &tx
}

//	判断当前交易是否为挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	//	1.交易的input只有一个 ，因为币是“凭空造的”，不是转账来的  这个input都是人为写的一些东西
	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].TXid) == 0 && tx.TXInputs[0].Index == -1 {
		return true
	}
	// 这是一个bug 解决
	return false
}

//创建普通交易
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	//00. 创建交易之后要进行数字签名->所以需要私钥->打开钱包"NewWallets()" -- version05
	ws := NewWallets()
	//01.根据自己的地址找到自己的钱包wallet， wallet里面存有当前地址的公私钥  -- 转账的账户
	wallet := ws.WalletMap[from]
	if wallet == nil {
		fmt.Printf("没有找到地址为：%s钱包，创建钱包失败\n", from)
		return nil // 阻止代码向下运行
	}
	//02.得到当前转账账户的公私钥
	pubKey := wallet.PubKey
	//privateKey := wallet.Private
	//1.	获取自己账户(from)中满足当前交易可用的最合理utxo  -- 最合理：只要找到当前账户可用的余额大于等于要花费的钱就不用再继续找，可以进行转账
	//	传递公钥的哈希，而不是传递地址  -- TXOutput中使用的收款人的pubKeyHash
	pubKeyHash := PubKeyHash(pubKey)
	utxos, resValue := bc.FindNeedUTXOs(pubKeyHash, amount)
	if resValue < amount {
		//	说明当前账户余额不足
		fmt.Println("当前交易失败，余额不足")
		return nil
	}
	//	定义交易组成 -- 一个交易：输出， 输入
	inputs := []TXInput{}   // 输出交易类型切片
	outputs := []TXOutput{} //	输入交易类型切片
	//	2. 创建交易输出, 将这些UTXO逐一转成inputs  -- 将自己的余额花出去，即将自己的outputs变为inputs -- 自己花出去
	//map[2222] = []int64{0}    utxos中存储的示例
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TXInput{
				TXid:  []byte(id),
				Index: int64(i),
				// 花自己的钱，要自己进行签名
				//Sig: from,
				//	花自己的钱自己要签名，稍后实现 -- version05
				Signature: nil,
				//	花钱人的公钥哈希
				PubKey: pubKeyHash,
			}
			inputs = append(inputs, input)
		}
	}

	//	创建输出交易 -- 钱转给谁
	//output := TXOutput{
	//	Value: amount,
	//	//	收款人的地址-- 签名
	//	PubKeyHash: to,
	//}
	//	在version05 中TXOutput有自己的初始化方法
	output := NewTXOutput(amount, to)
	outputs = append(outputs, *output)

	//	找零  -- 剩下的钱需要转给自己账户
	if amount < resValue {
		//	找到的合理的utxos 中的余额大于本次转账的钱 --- 有剩余
		//outputs = append(outputs, TXOutput{resValue - amount, from})
		output := NewTXOutput(resValue-amount, from)
		outputs = append(outputs, *output)
	}
	//	初始化交易，将输出交易和输入交易打包
	transaction := Transaction{
		//	先赋空的byte类型 在下面调用Transaction的方式设置哈希
		TXID:      []byte{},
		TXInputs:  inputs,
		TXOutputs: outputs,
	}
	//	设置当前交易的哈希  -- 对当前区块进行取哈希
	transaction.SetHash()
	return &transaction
}
