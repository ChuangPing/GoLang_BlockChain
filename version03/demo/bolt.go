package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

func main() {
	//1.打开数据库 如果数据库文件不存在会自动在当前目录创建， 0600 ：表示读写权限，固定的  -- 我们的操作都是基于这个数据库文件
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		fmt.Println("打开bolt数据库失败：", err)
	}
	defer db.Close() //关闭数据库
	//	对数据库进行操作
	db.Update(func(tx *bolt.Tx) error {
		//2. 找到抽屉bucket。如果在打开过程没有会自动创建
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil {
			//	没有抽屉，这时我们需要创建
			bucket, err = tx.CreateBucket([]byte("b1"))
			if err != nil {
				log.Panicln("创建抽屉Bucket失败：", err)
			}

		}
		// 写数据 -- 走到这一步说明已经存在、或者创建好了抽屉
		bucket.Put([]byte("11"), []byte("hello bolt"))

		return nil
	})
	//	读数据
	db.View(func(tx *bolt.Tx) error {
		//	1.找到抽屉 -- 没有直接退出或者报错    读写操作都是针对抽屉的
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil {
			log.Panicln("bucket b1 不应该为空，请检查")
		}
		res := bucket.Get([]byte("11"))
		fmt.Printf("键11对应的value：%s\n", res)
		return nil
	})
}
