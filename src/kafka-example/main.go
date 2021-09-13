package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Shopify/sarama"
)

// 计算函数
type Process func(float64) float64

var Client *mongo.Client
var tmpMap map[string]interface{}
var ProMap map[string]Process

type KafkaConsumer struct {
	Node         string
	Topic        string
	DBUrl        string
	DBName       string
	Collection   string
	MessageQueue chan string
	CurrentID    float64
	CurrentData  map[string]interface{}
	SaveQueue    chan map[string]interface{}
}

// 获取数据发送到chnnal中
func (thi KafkaConsumer) Consume() {
	consumer, err := sarama.NewConsumer([]string{thi.Node}, nil)
	if err != nil {
		fmt.Printf("kafka connnet failed, error[%v]", err.Error())
		return
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(thi.Topic)
	if err != nil {
		fmt.Printf("get topic failed, error[%v]", err.Error())
		return
	}

	for _, p := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(thi.Topic, p, sarama.OffsetNewest) // 最新的offset开始消费
		if err != nil {
			fmt.Printf("get partition consumer failed, error[%v]", err.Error())
			continue
		}

		for message := range partitionConsumer.Messages() {
			// fmt.Printf("message:[%v], key:[%v], offset:[%v]\n", string(message.Value), string(message.Key), string(message.Offset))
			thi.MessageQueue <- string(message.Value)
		}
	}
}

// 抓取数据结构
type Message struct {
	Payload   string `json:"payload"`
	Username  string `json:"username"`
	Topic     string `json:"topic"`
	Timestamp uint64 `json:"timestamp"`
	Qos       int64  `json:"-"`
	Clientid  string `json:"clientid"`
}

// 消费到mongodb中
func (thi KafkaConsumer) consumerToMongo() {
	for value := range thi.MessageQueue {
		var m Message
		var b bool
		err := json.Unmarshal([]byte(value), &m)
		if err != nil {
			fmt.Println("err: ", err)
		}

		err = json.Unmarshal([]byte(m.Payload), &tmpMap)
		if err != nil {
			fmt.Println("err: ", err)
		}
		// fmt.Println(tmpMap)
		for _, d := range tmpMap["d"].([]interface{}) {
			var (
				val float64
			)
			d_1, _ := d.(map[string]interface{})
			tag := d_1["tag"].(string)

			// 异常数据判断
			if strings.Contains(tag, "#") {
				b = true
				break
			}

			// 过滤参数
			if strings.Contains(tag, "未知点") {
				continue
			}

			// 计算处理
			if f, ok := ProMap[tag]; ok {
				val = f(d_1["value"].(float64))
			} else {
				val = d_1["value"].(float64)
			}

			tmpMap[tag] = val
		}

		if b {
			// 跳过异常数据
			continue
		}

		// 删除旧文件
		delete(tmpMap, "d")

		currentId := tmpMap["压射号"].(float64)
		// 判断
		if currentId > thi.CurrentID {
			thi.SaveQueue <- thi.CurrentData
		}

		thi.CurrentID = currentId
		thi.CurrentData = tmpMap
		// fmt.Println("ok")
	}
}

func (thi KafkaConsumer) saveToMongo() {
	var (
		err        error
		collection *mongo.Collection
		ctx        = context.Background()
		update     bson.M

		filter bson.M
		op     options.UpdateOptions
	)

	for value := range thi.SaveQueue {
		if len(value) == 0 {
			continue
		}

		y := value["压射号"].(float64)
		if y == 0 {
			continue
		}

		collection = Client.Database(thi.DBName).Collection(thi.Collection)

		filter = bson.M{"压射号": y} // ?
		update = bson.M{"$set": value}
		op = options.UpdateOptions{}
		op.SetUpsert(true)

		// 插入或更新
		if _, err = collection.UpdateOne(ctx, filter, update, &op); err != nil {
			fmt.Printf("[+] Mongo update error: %s", err.Error())
			return
		}
	}
}

func (thi KafkaConsumer) GetMongoClientByURI(ctx context.Context) (*mongo.Client, error) {
	var (
		err    error
		client *mongo.Client
		opts   *options.ClientOptions
	)

	opts = options.Client()
	opts.ApplyURI(thi.DBUrl)
	opts.SetConnectTimeout(3 * time.Second) // mongo连接异常
	opts.SetMaxPoolSize(20)

	if client, err = mongo.Connect(ctx, opts); err != nil {
		return client, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}

func mov(f float64) float64 {
	return f / 100
}

func main() {
	var (
		err           error
		ctx           = context.Background()
		c             = make(chan string, 500)
		s             = make(chan map[string]interface{}, 500)
		kafkaConsumer = KafkaConsumer{
			Node: "192.168.0.69:9092",
			// DBUrl:        "mongodb://yj:dXMxoD2dk@192.168.4.177:27017/yj?authSource=admin&directConnection=true&ssl=false",
			DBUrl:        "mongodb://admin:123456@localhost:27017/?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false",
			Collection:   "yj_2",
			DBName:       "dev",
			Topic:        "yjtest",
			MessageQueue: c,
			CurrentID:    0,
			CurrentData:  tmpMap,
			SaveQueue:    s,
		}
	)

	ProMap = map[string]Process{
		"吹气时间": mov,
	}

	// 初始化DB
	ctx = context.Background()
	Client, err = kafkaConsumer.GetMongoClientByURI(ctx)
	if err != nil {
		fmt.Println(err)
	}

	go kafkaConsumer.saveToMongo()
	go kafkaConsumer.consumerToMongo()
	kafkaConsumer.Consume()

}
