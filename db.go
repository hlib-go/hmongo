package hmongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strings"
)

// 	Deprecated
// 	改为使用db.NewDatabase(dsn)
// 	解析Mongodb数据库连接字符串，返回数据库连接
func DB(connectionString string) (db *mongo.Database, err error) {
	if connectionString == "" {
		err = errors.New("Mongodb 连接字符串不能为空")
		return
	}
	dbName := (strings.Split((strings.Split(connectionString, "/"))[3], "?"))[0]
	if dbName == "" {
		err = errors.New(fmt.Sprintf("Errror Mongodb connectionString %s", connectionString))
		return
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString).SetReadPreference(readpref.SecondaryPreferred()))
	if err != nil {
		err = errors.New(fmt.Sprintf("Errror Connect mongodb exception %s", err))
		return
	}
	db = client.Database(dbName)
	return
}
