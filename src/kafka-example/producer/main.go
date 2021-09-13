package main

import (
	"fmt"
	"strconv"
	"strings"
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
	producer, err := sarama.NewSyncProducer([]string{"192.168.0.69:9092"}, config)
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
	var count = 3825

	var v = ` 
	{
	   "username": "yj",
	   "topic": "data/device_id",
	   "timestamp": 1629954185743,
	   "qos": 1,
	   "payload": "{\"d\":[{\"tag\":\"压射号\",\"value\":18072},{\"tag\":\"V1\",\"value\":0.00},{\"tag\":\"V2\",\"value\":0.00},{\"tag\":\"V3\",\"value\":0.00},{\"tag\":\"V4\",\"value\":0.00},{\"tag\":\"加速位置\",\"value\":0.00},{\"tag\":\"减速位置\",\"value\":0.00},{\"tag\":\"铸造压力\",\"value\":0.00},{\"tag\":\"升压时间\",\"value\":0.00},{\"tag\":\"料饼厚度\",\"value\":0.00},{\"tag\":\"锁模力\",\"value\":0.00},{\"tag\":\"锁模力MN\",\"value\":0.00},{\"tag\":\"高速区间\",\"value\":0.00},{\"tag\":\"循环时间\",\"value\":0.00},{\"tag\":\"浇注时间\",\"value\":0.00},{\"tag\":\"产品冷却时间\",\"value\":0.00},{\"tag\":\"顶出时间\",\"value\":0.00},{\"tag\":\"喷淋时间\",\"value\":0.00},{\"tag\":\"压射延时\",\"value\":0.00},{\"tag\":\"压射时间\",\"value\":3.00},{\"tag\":\"吹气延时\",\"value\":128.00},{\"tag\":\"吹气时间\",\"value\":2560.00},{\"tag\":\"未知点1\",\"value\":0},{\"tag\":\"未知点2\",\"value\":0},{\"tag\":\"未知点3\",\"value\":0},{\"tag\":\"未知点4\",\"value\":0}],\"ts\":\"2021-08-26T05:02:45+0000\"}",
	   "clientid": "edgelink20210825153942"
	}`
	for {
		for i := 0; i < 10; i++ {
			msg.Value = sarama.ByteEncoder(strings.Replace(v, "18072", strconv.Itoa(count), -1))
			//fmt.Println(value)
			//SendMessage：该方法是生产者生产给定的消息
			//生产成功的时候返回该消息的分区和所在的偏移量, 生产失败的时候返回error
			partition, offset, err := producer.SendMessage(msg)

			if err != nil {
				fmt.Println("Send message Fail", err)
			}
			fmt.Printf("Partition = %d, offset=%d, count=%d \n", partition, offset, count)

		}
		count += 1
		time.Sleep(8 * time.Second)
	}
}
