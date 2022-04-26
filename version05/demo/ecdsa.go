package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
)

//ecdsa包实现了椭圆曲线数字签名算法
//1. 演示如何使用ecdsa生成公私钥
//2.签名校验

func main() {
	//	1.发送端进行签名
	//	创建椭圆曲线
	curve := elliptic.P256()
	//	生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panicln("生成私钥出错：", err)
	}
	//	生成公钥
	pubkey := privateKey.PublicKey
	//	需要签名的数据
	data := "hello world"
	//	对原始数据取哈希
	hash := sha256.Sum256([]byte(data))

	//	对数据进行签名  rand.Reader:随机数
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Panicln("签名出错：", err)
	}
	//	打印输出r , s  -- 从打印可以看出， r, s 均为长度32，等长的
	fmt.Println("publick:%v", pubkey)
	fmt.Printf("r: %v, rLen: %d\n", r, len(r.Bytes()))
	fmt.Println("s:%v, sLen:%d\n", s, len(s.Bytes()))
	//	一般签名者签名后，需要将签名的 公钥，签名产生的r s 在网络中发送，在发送时，我们一般会把这些进行转换为字节流，在网络是上传输

	//	把 s, r 进行序列化  --- r,s都是big.Int类型，bigInt有转换为byte的方法
	signature := append(r.Bytes(), s.Bytes()...)

	//	2.接收端验证签名
	//	准备两个bigInt类型辅助变量，将接收到的byte字节流转换为r,s -- 验证签名需要使用
	r1 := big.Int{}
	s1 := big.Int{}
	//	将在网络中接收到的签名，反序列化 -- 平均分字节流，因为s和r等长，前面的为r后面的部分为s
	r1.SetBytes(signature[:len(signature)/2]) // 这里不用担心是否能除尽，因为r与s等长，一定能整除
	s1.SetBytes(signature[len(signature)/2:])
	//	校验时需要：原始数据 签名，公钥
	data1 := "hello world"
	hash1 := sha256.Sum256([]byte(data1))
	//func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool {
	res := ecdsa.Verify(&pubkey, hash1[:], &r1, &s1)
	if res {
		fmt.Println("校验成功：", res)
	} else {
		fmt.Println("校验失败", res)
	}

}
