package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

//定义一个自定义结构体
type Person struct {
	Name string
	Age  int
}

func main() {
	var alice Person
	alice.Age = 20
	alice.Name = "爱丽丝"
	//	将alice进行序列化  -- 即将结构体转成字节byte类型
	//	解码的数据放到buffer中
	var buffer bytes.Buffer
	//	使用gob进行序列化
	//1.定义一个编码器
	encoder := gob.NewEncoder(&buffer)
	//	2.使用解码器进行解码
	err := encoder.Encode(&alice)
	if err != nil {
		log.Panicln("编码失败", err)
	}
	fmt.Printf("编码后的alice：%v\n", buffer.Bytes()) // 存在buffer中，但是是byte类型，所有使用buffer的Bytes方法

	//	解码 -- 将字节转换为对应的类型
	//	1.定义一个解码器
	decoder := gob.NewDecoder(&buffer)
	//	2.使用解码器进行解码
	var bob Person // 定义一个需要解码的类型
	err = decoder.Decode(&bob)
	if err != nil {
		log.Panicln("解码器进行解码出错", err)
	}
	fmt.Printf("解码后的alice：%v\n", bob) // 解码后的alice：{爱丽丝 20} %v可以自己推导类型
}
