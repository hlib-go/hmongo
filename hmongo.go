package hmongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// 如果需要修改读取的系统环境变量，拷贝此文件到业务系统修改即可。

var (
	// 默认通过系统环境配置读取数据库连接
	DefaultDB = new(HMongo)
)

// Mongodb 连接
type HMongo struct {
	database *mongo.Database
	cmp      map[string]bool
}

// ConneStr 从系统环境读取Mongodb数据库连接字符串
func (o *HMongo) ConnStr() string {
	//**********从环境变量读取
	connStr := os.Getenv("HMONGO")
	if connStr != "" {
		return connStr
	}
	//**********从docker secrets 读取
	filename := "/run/secrets/hmongo"
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error hmongo file:", err.Error())
		return ""
	}
	return string(bs)
}

func (o *HMongo) Database() (db *mongo.Database) {
	db = o.database
	if db != nil {
		return
	}
	var (
		err error
	)
	defer func() {
		if err != nil {
			panic(err.Error())
		}
	}()

	connStr := o.ConnStr()
	if connStr == "" {
		err = errors.New("HMongo 连接字符串不能为空")
		return
	}
	dbName := (strings.Split((strings.Split(connStr, "/"))[3], "?"))[0]
	if dbName == "" {
		err = errors.New(fmt.Sprintf("Errror HMongo connStr %s", connStr))
		return
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connStr))
	if err != nil {
		err = errors.New(fmt.Sprintf("Errror HMongo Connect %s", err))
		return
	}
	o.database = client.Database(dbName)
	db = o.database
	return
}

func (o *HMongo) Client() (client *mongo.Client, err error) {
	db := o.Database()
	client = db.Client()
	return
}

func (o *HMongo) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	c := o.Database().Collection(name, opts...)
	if !o.cmp[name] {
		err := o.CreatedEmptyCollection(c)
		if err != nil {
			panic(err)
		}
	}
	return c
}

// CreatedEmptyCollection 因事务不支持不存在的集合，可在服务启动时执行检查
func (o *HMongo) CreatedEmptyCollection(c *mongo.Collection) (err error) {
	count, err := c.CountDocuments(nil, bson.M{}, options.Count().SetLimit(1))
	if err != nil {
		return
	}
	if count == 0 {
		_, err = c.InsertOne(nil, bson.M{"createdAt": time.Now()})
	}
	return
}
