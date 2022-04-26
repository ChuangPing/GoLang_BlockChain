package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallet.data"

//	定一个 Wallets结构，它保存所有的wallet以及它的地址
type Wallets struct {
	//	map[地址]钱包
	WalletMap map[string]*Wallet
}

//	创建函数返回当前所有钱包的实例
func NewWallets() *Wallets {
	var ws Wallets
	//	对Wallets结构体中的walletMap map类型开辟空间 -- map类型在使用前必须make，否则会报错
	ws.WalletMap = make(map[string]*Wallet)
	//	读取钱包文件，对Wallets进行初始化
	ws.loadFile()
	return &ws
}

//	定义重文件中加载钱包已有地址,返回当前所有钱包的实例
func (ws *Wallets) loadFile() {
	//	读取存储钱包地址文件
	_, err := os.Stat(walletFile) // 判断文件状态，是否文件为空
	if os.IsNotExist(err) {
		//log.Panicln("wallet.dat文件不存在") 坑：使用Panic会抛异常，并停止程序继续执行，因此第一次创建时，就会报错停止程序继续运行，导致出错
		return
	}
	//	读取文件
	contentByte, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panicln("读取文件出错")
		return
	}
	// 读出的文件是以byte形式，需要进行解码

	//由于存在interface类型，需要向gob进行注册。    panic: gob: type not registered for interface: elliptic.p256Curve
	gob.Register(elliptic.P256())
	//	初始化解码器
	decoder := gob.NewDecoder(bytes.NewReader(contentByte))
	//	定义需要解码后的类型
	wsLocal := Wallets{}
	err = decoder.Decode(&wsLocal)
	if err != nil {
		log.Panicln("解码钱包文件 内容出错：", err)
		return
	}
	//	将读取的文件内容给Wallets中的map赋值完成Wallets初始化
	ws.WalletMap = wsLocal.WalletMap
}

//	创建钱包方法,完成地址创建，并将创建后的地址保存在wallet.data中
func (ws *Wallets) CreateWallet() string {
	//	初始化钱包
	wallet := NewWallet()
	//	生成钱包地址
	address := wallet.NewAddress()
	//	将当前钱包加入到 wallets中
	ws.WalletMap[address] = wallet
	//	将生成的钱包存储在walle.data
	ws.saveToFile()
	return address
}
func (ws *Wallets) saveToFile() {
	var buffer bytes.Buffer
	//	存储文件前，需要将内容进行编码成字节
	//在gob中注册编码类型
	gob.Register(elliptic.P256())
	//	定义编码器
	encoder := gob.NewEncoder(&buffer)
	//	对内容进行编码
	err := encoder.Encode(ws)
	if err != nil {
		log.Panicln("文件编码出错：", err)
		return
	}
	//	将编码后的内容存入文件
	err = ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)
	if err != nil {
		log.Panicln("写入钱包文件出错：", err)
		return
	}
}

//	定义显示当前钱包所有地址的方法
func (ws Wallets) ListAllAddress() []string {
	var addresses []string
	//	遍历钱包，将所有的key（地址）加入addreses
	for address := range ws.WalletMap {
		addresses = append(addresses, address)

	}
	return addresses
}
