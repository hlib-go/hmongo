package hmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Collection(db *mongo.Database, name string) *mongo.Collection {
	return db.Collection(name)
}

// 列表查询
func Find(ctx context.Context, c *mongo.Collection, filter interface{}, result interface{}, opts ...*options.FindOptions) (err error) {
	cursor, err := c.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	if cursor.Err() != nil {
		return cursor.Err()
	}
	err = cursor.All(ctx, result)
	if err != nil {
		return err
	}
	return
}

// 分页查询
func FindPage(ctx context.Context, c *mongo.Collection, filter interface{}, sort bson.M, pageSize, pageNum int64, result interface{}, fo ...*options.FindOptions) (total int64, err error) {
	if sort == nil {
		sort = bson.M{}
	}
	if sort["_id"] == nil {
		sort["_id"] = -1 // 默认根据_id倒序
	}
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	op := options.Find().SetSort(sort).SetSkip(pageSize * (pageNum - 1)).SetLimit(pageSize)
	fo = append(fo, op)
	err = Find(ctx, c, filter, result, fo...)
	if err != nil {
		return 0, err
	}
	return c.CountDocuments(ctx, filter)
}
