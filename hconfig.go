package hmongo

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func NewHConfig(collection *mongo.Collection) *HConfig {
	return &HConfig{collection: collection}
}

// HConfig 使用Mongodb集合做项目配置
type HConfig struct {
	collection *mongo.Collection
}

// Get 获取配置
func (o *HConfig) Get(name string, v interface{}) (err error) {
	var result map[string]interface{}
	err = o.collection.FindOne(nil, bson.M{"name": name}).Decode(&result)
	if err != nil && err == mongo.ErrNoDocuments {
		err = errors.New("HConfig name=" + name + " no documents in result")
		return
	}
	value, err := json.Marshal(result["value"])
	if err != nil {
		return
	}
	err = json.Unmarshal(value, v)
	if err != nil {
		return
	}
	return
}

// Put 更新配置
func (o *HConfig) Put(name string, v interface{}, remark string) (err error) {
	_, err = o.collection.UpdateOne(nil, bson.M{"name": name}, bson.M{"$set": bson.M{"value": v, "remark": remark, "updatedAt": time.Now()}})
	return
}

// 待测试调整
// Watch 异步监听配置
func (o *HConfig) Watch(cbFunc func(name string, value []byte)) {
	go func() {
		stream, err := o.collection.Watch(nil, bson.A{})
		if err != nil {
			return
		}
		for stream.Next(context.Background()) {
			// "{\"_id\": {\"_data\": \"825FE42F0E000000012B022C0100296E5A10049ADEEABDEF114921BF5B60F6101755F446645F696400645FE42B0BD5BDC442049CB4BA0004\"},\"operationType\": \"update\",\"clusterTime\": {\"$timestamp\":{\"t\":\"1608789774\",\"i\":\"1\"}},\"ns\": {\"db\": \"himkt\",\"coll\": \"hm_config\"},\"documentKey\": {\"_id\": {\"$oid\":\"5fe42b0bd5bdc442049cb4ba\"}},\"updateDescription\": {\"updatedFields\": {\"value\": {\"accessKey\": \"OY6HO2TZMqSc_sBQUX3lgVRLznE4D4GoCIFecEwW\",\"secretKey\": \"OkW4hZXz9m7MYLptJv7bLFO_fkjdxd4D5FSXtt40\"}},\"removedFields\": []}}"
			log.Info(stream.Current.String())

			var result map[string]interface{}
			err = stream.Decode(&result)
			if err != nil {
				continue
			}
			value, err := json.Marshal(result["value"])
			if err != nil {
				return
			}
			cbFunc(result["name"].(string), value)
		}
		log.Info("HConfig Watch end.")
	}()
}
