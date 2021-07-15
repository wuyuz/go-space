package main

import (
    "fmt"
    "encoding/json"
    "time"
    "strconv"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
    "context"

    "github.com/Shopify/sarama"
)

var Client *mongo.Client

type KafkaConsumer struct {
    Node string
    Topic string
    DBUrl string
    DBName string
    Collection string
    MessageQueue chan string
    CurrentID float64
    CurrentData SaveData
    SaveQueue chan SaveData
}

func (this KafkaConsumer) Consume(){
    consumer, err := sarama.NewConsumer([]string{this.Node}, nil)
    if err != nil {
        fmt.Printf("kafka connnet failed, error[%v]", err.Error())
        return
    }
    defer consumer.Close()

    partitions, err := consumer.Partitions(this.Topic)
    if err != nil {
        fmt.Printf("get topic failed, error[%v]", err.Error())
        return
    }
    for _, p := range partitions {
        partitionConsumer, err := consumer.ConsumePartition(this.Topic, p, sarama.OffsetNewest)  // 最新的offset开始消费
        if err != nil {
            fmt.Printf("get partition consumer failed, error[%v]", err.Error())
            continue
        }

        for message := range partitionConsumer.Messages() {
            // fmt.Printf("message:[%v], key:[%v], offset:[%v]\n", string(message.Value), string(message.Key), string(message.Offset))
            this.MessageQueue <- string(message.Value)
        }
    }
}

type inner struct {
    BeginTime string `json:"beginTime"`
    Data []interface{} `json:"data"`
    ID string `json:"-"`
}

// 抓取数据结构
type  Message struct {
    Payload string `json:"payload"`
    Username string `json:"username"`
    Topic string `json:"topic"`
    Timestamp uint64 `json:"timestamp"`
    Qos int64 `json:"-"`
    Clientid string `json:"clientid"`
}

// 存储数据结构
type SaveData struct {
    ShotNumber float64 `bson: "yz_shoot_NO"`
    DeviceID string   `bson: "device_id"`
    TimeStamp time.Time  `bson: "create_time"`
    V1Speed float64 `bson: "yz_v1_speed"`
    V2Speed float64 `bson: "yz_v2_speed"`
    V3Speed float64 `bson: "yz_v3_speed"`
    V4Speed float64 `bson: "yz_v4_speed"`
    PreRiseTime float64 `bson: "yz_pre_rist_time"`
    UpSpeedTime float64 `bson:"yz_upspeed_time`
    DownSpeedTime float64 `bson:"yz_downspeed_time`
}

// 消费到mongodb中
func (this KafkaConsumer) consumerToMongo() {
    for value := range this.MessageQueue {
        var m Message
        err := json.Unmarshal([]byte(value), &m)
        if err != nil {
            fmt.Println("err: ",err)
        }
        var i inner
        err = json.Unmarshal([]byte(m.Payload), &i)
        if err != nil {
            fmt.Println("err: ",err)
        }

        t,_ := strconv.Atoi(i.BeginTime)
        currentId := i.Data[8].(map[string]interface{})["data"].([]interface{})[0].(float64)
        // 判断
        if currentId > this.CurrentID {
            this.SaveQueue <- this.CurrentData
        }

        this.CurrentID = currentId
        this.CurrentData = SaveData{
            ShotNumber: currentId,
            DeviceID: "40000009230",
            TimeStamp: time.Unix(int64(t/1000),0),
            V1Speed: i.Data[1].(map[string]interface{})["data"].([]interface{})[0].(float64)/100,
            V2Speed: i.Data[2].(map[string]interface{})["data"].([]interface{})[0].(float64)/100,
            V3Speed: i.Data[3].(map[string]interface{})["data"].([]interface{})[0].(float64)/100,
            V4Speed: i.Data[4].(map[string]interface{})["data"].([]interface{})[0].(float64)/100,
            PreRiseTime: i.Data[5].(map[string]interface{})["data"].([]interface{})[0].(float64)/10,
            UpSpeedTime: i.Data[6].(map[string]interface{})["data"].([]interface{})[0].(float64),
            DownSpeedTime: i.Data[7].(map[string]interface{})["data"].([]interface{})[0].(float64),
        }
       fmt.Println(this.CurrentData)
    }
}

func (this KafkaConsumer) saveToMongo() {
    var (
		err        error
		collection *mongo.Collection
		ctx        = context.Background()
		update     bson.M
		filter     bson.M
		op         options.UpdateOptions
	)

    for value := range this.SaveQueue {
        if value.ShotNumber == 0 {
            continue
        }

        collection = Client.Database(this.DBName).Collection(this.Collection)
        m, _ := time.ParseDuration("-10m")
        filter = bson.M{"shotnumber": value.ShotNumber, "deviceid": value.DeviceID, "timestamp": bson.M{
            "$gte": time.Now().Add(m),
        }}  // ?
        update = bson.M{"$set": value}
        op = options.UpdateOptions{}
        op.SetUpsert(true)

         // 插入或更新
        if _, err = collection.UpdateOne(ctx, filter, update, &op); err != nil {
            fmt.Println("[+] Mongo update error: %s", err.Error())
            return
        }
    }
}

func (this KafkaConsumer) GetMongoClientByURI(ctx context.Context) (*mongo.Client, error) {
	var (
		err    error
		client *mongo.Client
		opts   *options.ClientOptions
	)

	opts = options.Client()
	opts.ApplyURI(this.DBUrl)
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

func main(){
    var (
        err        error
		ctx        = context.Background()
        c = make(chan string, 500)
        s = make(chan SaveData, 500)
        kafkaConsumer = KafkaConsumer{
            Node: "192.168.0.86:9093",
            DBUrl: "mongodb://admin:123456@localhost/",
            Collection: "test",
            DBName:"yj",
            Topic: "yjtest",
            MessageQueue: c,
            CurrentID: 0,
            CurrentData: SaveData{},
            SaveQueue: s,
        }
	)

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