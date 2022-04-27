package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

//这里的钱包时一结构，每一个钱包保存了公钥,私钥对
type Wallet struct {
	//	私钥
	Private *ecdsa.PrivateKey
	//PubKey *ecdsa.PublicKey  -- 正常情况
	//约定，这里的PubKey不存储原始的公钥，而是存储X和Y拼接的字符串，在校验端重新拆分（参考r,s传递）
	PubKey []byte
}

//	创建钱包
func NewWallet() *Wallet {
	//	c创建生成公私钥需要的椭圆曲线
	curve := elliptic.P256()
	//	生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panicln("生成私钥出错：", err)
		return nil
	}
	//	有私钥生成公钥
	pubKeyOrig := privateKey.PublicKey
	//	拼接x, y
	pubKey := append(pubKeyOrig.X.Bytes(), pubKeyOrig.Y.Bytes()...)
	return &Wallet{
		privateKey,
		pubKey,
	}
}

//	钱包生成地址方法 比特币地址由三部分组成：Version           Public key hash                       Checksum
//                                   00         62E907B15CBF27D5425399EBF6F0FB50EBB88F18   C29B7D93 3
func (wallet *Wallet) NewAddress() string {
	//	在比特币中地址，是由用户公钥导出的
	publicKey := wallet.PubKey
	//	对公钥取rip160哈希
	rip160HashValue := PubKeyHash(publicKey)
	version := byte(00)
	//	拼接version
	payLoad := append([]byte{version}, rip160HashValue...)
	//	获四位校取验码
	checkCode := CheckSum(payLoad)
	//	拼接组合    -- 25字节数据
	payLoad = append(payLoad, checkCode...)
	//go语言有一个库，叫做btcd,这个是go语言实现的比特币全节点源码
	// 将最终的payLoad进行base58编码  --- 若没有base58 使用命令解决：go get github.com/btcsuite/btcutil/base58
	address := base58.Encode(payLoad)
	return address
}

//	对公钥进行rip160哈希
func PubKeyHash(data []byte) []byte {
	//	先对公钥进行取sha256哈希
	hash := sha256.Sum256(data)
	//理解为编码器     下面可能会报错使用命令解决：   go get -v github.com/btcsuite/btcd
	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])
	if err != nil {
		log.Panicln("对公钥取ripemd160出错：", err)
		return nil
	}
	//返回rip160的哈希结果
	rip160HashValue := rip160hasher.Sum(nil)
	return rip160HashValue
}

//	公钥验证码函数 -- 在后面会用到，来验证地址是否合法，是否是用这种方式生成的地址
func CheckSum(data []byte) []byte {
	//	对公钥进行rip160哈希后进行两次sha256
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	//	获取前四个字节的校验码
	checkCode := hash2[:4]
	return checkCode
}

//	验证地址是否合法
func IsValidAddress(address string) bool {
	//	将address进行解码拿到payLoad 与 checkSum,然后再由payLoa调用CheckSum 生成checksum00,判断生成checksum00 == checkSum
	//	1.解码
	addressBytes := base58.Decode(address)
	if len(addressBytes) < 4 {
		fmt.Printf("地址：%s非法！！\n", address)
		return false
	}
	//	2.取数据
	payLoad := addressBytes[:len(addressBytes)-4] //checkSum长度为4个字节
	checkSum := addressBytes[len(addressBytes)-4:]
	//	3.由payLoad生成checkSum1
	checkSum1 := CheckSum(payLoad)
	//	4.返回比较结果
	return bytes.Equal(checkSum, checkSum1)
}
