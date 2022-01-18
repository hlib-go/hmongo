package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hlib-go/hmongo/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

/*
	新建配置对象
	dsn Mongodb数据库连接uri
	name 集合名
	keys 读取的配置key名字数字，为空读取所有
*/
func New(dsn, name string, keys ...string) *Config {
	collection := db.NewDatabase(dsn).Collection(name)
	c := &Config{
		cache:      sync.Map{},
		collection: collection,
		keys:       keys,
	}
	c.find()
	go c.watch()
	return c
}

type Config struct {
	cache      sync.Map
	collection *mongo.Collection
	keys       []string
}

// Get 读取指定配置
func (c *Config) Get(key string, value interface{}) {
	vBytes, ok := c.cache.Load(key)
	if !ok || vBytes == nil {
		panic(errors.New(fmt.Sprintf("ERROR:未读取到【%s】配置，请检查项目配置", key)))
	}

	err := json.Unmarshal(vBytes.([]byte), value)
	if err != nil {
		panic(errors.New(fmt.Sprintf("ERROR:无法解析【%s】配置，请检查类型是否正确", key)))
	}
	return
}

// GetCache 读取所有配置项
func (c *Config) GetCache() map[string]string {
	m := make(map[string]string)
	c.cache.Range(func(key, value interface{}) bool {
		m[key.(string)] = string(value.([]byte))
		return true
	})
	return m
}

// GetItem 读取其中一项配置
func (c *Config) GetItem(name string) (v []byte) {
	value, ok := c.cache.Load(name)
	if ok {
		v = value.([]byte)
	}
	return
}

//  查询配置
func (c *Config) find() {
	filter := bson.M{}
	if c.keys != nil && len(c.keys) > 0 {
		bsonA := bson.A{}
		for _, key := range c.keys {
			bsonA = append(bsonA, key)
		}
		filter["key"] = bson.M{"$in": bsonA}
	}

	cur, err := c.collection.Find(nil, filter)
	if err != nil {
		panic(err)
	}
	var configs []*Object
	err = cur.All(nil, &configs)
	if err != nil {
		panic(err)
	}
	for _, conf := range configs {
		value, _ := json.Marshal(conf.Value)
		if value == nil {
			continue
		}
		c.cache.Store(conf.Key, value)
	}
}

//  监听配置更新
func (c *Config) watch() {
	ctx := context.TODO()
	cs, err := c.collection.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		fmt.Println(err.Error())
		time.Sleep(5 * time.Minute)
		c.watch() // 异常情况，间隔5分钟重试
	}
	defer cs.Close(ctx)
	for cs.Next(ctx) {
		fmt.Println("Watch config ", cs.Current)
		c.find()
	}
}

type Object struct {
	Key    string                 `bson:"key" json:"key"`       // 配置Key
	Value  map[string]interface{} `bson:"value" json:"value"`   // 配置对象
	Remark string                 `bson:"remark" json:"remark"` // 备注说明
}
