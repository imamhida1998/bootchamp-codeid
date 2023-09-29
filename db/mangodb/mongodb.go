package mangodb

import (
	"context"
	"fmt"

	"bootchamp-codeid/db/mangodb/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func OpenMongoDb(ctx context.Context) (*mongo.Client, error) {

	config := config.MONGO_CONFIGS[config.IMAM]

	connString := ""
	if config.Authentication {
		connString = fmt.Sprintf("mongodb://%s:%s@%s:%s",
			config.User, config.Pwd,
			config.Host, config.Port)

	} else {
		connString = fmt.Sprintf("mongodb://%s:%s",
			config.Host, config.Port)

	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}
func CloseMongoDb(ctx context.Context, client *mongo.Client) error {

	var err = client.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}
