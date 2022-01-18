package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

func NewDatabase(dsn string, opts ...*options.ClientOptions) *mongo.Database {
	if dsn == "" {
		panic(errors.New("ERROR:MongoDB dsn is nil"))
	}

	var dbName = (strings.Split((strings.Split(dsn, "/"))[3], "?"))[0]
	opts = append(opts, options.Client().ApplyURI(dsn))

	client, err := mongo.Connect(context.Background(), opts...)
	if err != nil {
		panic(err)
	}
	return client.Database(dbName)
}
