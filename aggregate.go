package hmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 聚合查询
func Aggregate(ctx context.Context, c *mongo.Collection, pipeline interface{}, result interface{}, opts ...*options.AggregateOptions) (err error) {
	cursor, err := c.Aggregate(ctx, pipeline, opts...)
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
