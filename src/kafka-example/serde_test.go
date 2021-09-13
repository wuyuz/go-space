package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func JsonToMap(str string) map[string]interface{} {
	var (
		tmpMap map[string]interface{}
	)

	err := json.Unmarshal([]byte(str), &tmpMap)
	if err != nil {
		panic(err)
	}
	return tmpMap
}

func TestSerde(t *testing.T) {
	var msg = "{\"d\":[{\"tag\":\"#DISABLE_DEVICE_压铸机三菱Q6\",\"value\":0},{\"tag\":\"#BATCH_WRITE_压铸机三菱Q6\",\"value\":0},{\"tag\":\"#DEVICE_ERROR_压铸机三菱Q6\",\"value\":-32764.00},{\"tag\":\"压射号\",\"value\":18072},{\"tag\":\"V1\",\"value\":0.00},{\"tag\":\"V2\",\"value\":0.00},{\"tag\":\"V3\",\"value\":0.00},{\"tag\":\"V4\",\"value\":0.00},{\"tag\":\"加速位置\",\"value\":0.00},{\"tag\":\"减速位置\",\"value\":0.00},{\"tag\":\"铸造压力\",\"value\":0.00},{\"tag\":\"升压时间\",\"value\":0.00},{\"tag\":\"料饼厚度\",\"value\":0.00},{\"tag\":\"锁模力\",\"value\":0.00},{\"tag\":\"锁模力MN\",\"value\":0.00},{\"tag\":\"高速区间\",\"value\":0.00},{\"tag\":\"循环时间\",\"value\":0.00},{\"tag\":\"浇注时间\",\"value\":0.00},{\"tag\":\"产品冷却时间\",\"value\":0.00},{\"tag\":\"顶出时间\",\"value\":0.00},{\"tag\":\"喷淋时间\",\"value\":0.00},{\"tag\":\"#DISABLE_DEVICE_点冷机S7_200MB\",\"value\":0},{\"tag\":\"#BATCH_WRITE_点冷机S7_200MB\",\"value\":0},{\"tag\":\"#DEVICE_ERROR_点冷机S7_200MB\",\"value\":0.00},{\"tag\":\"压射延时\",\"value\":0.00},{\"tag\":\"压射时间\",\"value\":3.00},{\"tag\":\"吹气延时\",\"value\":128.00},{\"tag\":\"吹气时间\",\"value\":2560.00},{\"tag\":\"未知点1\",\"value\":0},{\"tag\":\"未知点2\",\"value\":0},{\"tag\":\"未知点3\",\"value\":0},{\"tag\":\"未知点4\",\"value\":0}],\"ts\":\"2021-08-26T05:02:45+0000\"}"

	serde := JsonToMap(msg)

	fmt.Println(serde)
}
