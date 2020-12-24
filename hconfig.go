package hmongo

import (
	"context"
	"encoding/json"
	"errors"
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

// Watch 异步监听配置
func (o *HConfig) Watch(cbFunc func(name string, value []byte)) {
	go func() {
		stream, err := o.collection.Watch(nil, bson.A{bson.M{"$project": bson.M{"name": 1, "value": 1}}})
		if err != nil {
			return
		}
		for stream.TryNext(context.Background()) {
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
	}()
}
