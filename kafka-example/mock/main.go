package main

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	// // 等待服务器所有副本都保存成功后的响应
	// config.Producer.RequiredAcks = sarama.WaitForAll
	// // 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	// config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应
	config.Producer.Return.Successes = true

	// 使用给定代理地址和配置创建一个同步生产者
	producer, err := sarama.NewSyncProducer([]string{"192.168.0.86:9093"}, config)
	if err != nil {
		panic(err)
	}

	defer producer.Close()

	//构建发送的消息，
	msg := &sarama.ProducerMessage{
		Topic:     "yjtest",                    //包含了消息的主题
		Partition: int32(10),                   //
		Key:       sarama.StringEncoder("key"), //
	}

	for {
		msg.Value = sarama.ByteEncoder("xxxx")
		//fmt.Println(value)
		//SendMessage：该方法是生产者生产给定的消息
		//生产成功的时候返回该消息的分区和所在的偏移量, 生产失败的时候返回error
		partition, offset, err := producer.SendMessage(msg)

		if err != nil {
			fmt.Println("Send message Fail", err)
		}

		fmt.Printf("Partition = %d, offset=%d\n", partition, offset)
		time.Sleep(1 * time.Second)
	}
}
