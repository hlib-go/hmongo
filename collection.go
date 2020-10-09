package hmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func Collection(db *mongo.Database, name string) *mongo.Collection {
	return db.Collection(name)
}

// 因事务不支持不存在的集合，
func CreatedEmptyCollection(c *mongo.Collection) (err error) {
	count, err := c.CountDocuments(nil, bson.M{}, options.Count().SetLimit(1))
	if err != nil {
		return
	}
	if count == 0 {
		_, err = c.InsertOne(nil, bson.M{"createdAt": time.Now()})
	}
	return
}

// Mongodb Update
func UpdateOne(ctx context.Context, col *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return col.UpdateOne(ctx, filter, update)
}
