package data

import (
	"bootchamp-codeid/db/mangodb/config"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Transfer struct {
	BaseNoSql

	TransferID               primitive.ObjectID `json:"transfer_id" bson:"transfer_id"`
	SourceNoVirtualAccount   string             `json:"source_no_virtual_account" bson:"source_no_virtual_account"`
	SourceNameVirtualAccount string             `json:"source_name_virtual_account" bson:"source_name_virtual_account"`
	DesNoVirtualAccount      string             `json:"des_no_virtual_account" bson:"des_no_virtual_account"`
	DescNameVirtualAccount   string             `json:"desc_name_virtual_account" bson:"desc_name_virtual_account"`
	TransferAmount           float32            `json:"transfer_amount" bson:"transfer_amount"`
	CreateAt                 string             `json:"create_at" bson:"create_at"`
}

func CreateTransfer(client *mongo.Client) Transfer {
	tf := Transfer{}
	tf.Client = client
	return tf

}
func (tf Transfer) Collection() *mongo.Collection {
	return tf.Client.Database(config.DATABASEMONGO).Collection("transfer")
}
func (tf Transfer) Truncate(ctx context.Context) error {
	return tf.Collection().Drop(ctx)
}
func (tf Transfer) AddTransfer(ctx context.Context,
	SourceNoVirtualAccount string,
	SourceNameVirtualAccount string,
	DesNoVirtualAccount string,
	DescNameVirtualAccount string,
	TransferAmount float32) (*Transfer, error) {

	tf.TransferID = primitive.NewObjectID()
	tf.SourceNoVirtualAccount = SourceNoVirtualAccount
	tf.SourceNameVirtualAccount = SourceNameVirtualAccount
	tf.DesNoVirtualAccount = DesNoVirtualAccount
	tf.DescNameVirtualAccount = DescNameVirtualAccount
	tf.TransferAmount = TransferAmount
	tf.CreateAt = time.Now().Format(time.RFC1123)

	_, err := tf.Collection().InsertOne(ctx, tf)
	if err != nil {
		return nil, err
	}
	return &tf, nil

}

func (tf Transfer) FindOne(ctx context.Context, TransferID string) error {
	var id, err = primitive.ObjectIDFromHex(TransferID)
	if err != nil {
		return err
	}
	filter := bson.D{
		primitive.E{
			Key:   "transfer_id",
			Value: id,
		},
	}
	result := tf.Collection().FindOne(ctx, filter)
	if result.Err() != nil {
		return result.Err()
	}
	log.Println(result.Err())
	if result.Err().Error() == config.MONGO_NO_DOCUMENT {
		return errors.New("topup_not_found")
	}

	err = result.Decode(&tf)
	if err != nil {
		return err
	}

	return nil

}

func (tf Transfer) Delete(ctx context.Context, TransferID string) (error, error) {
	id, err := primitive.ObjectIDFromHex(TransferID)
	if err != nil {
		return err, nil

	}
	filter := bson.D{
		primitive.E{
			Key:   "transfer_id",
			Value: id,
		},
	}
	result, err := tf.Collection().DeleteOne(ctx, filter, nil)
	if err != nil {
		return err, nil
	}
	log.Println("File Deleted", result)
	return nil, nil

}

func (tf Transfer) ListTransfer(ctx context.Context) ([]Transfer, error) {
	transfers := []Transfer{}
	filter := bson.D{}
	cursor, err := tf.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &transfers)
	if err != nil {
		return nil, err
	}
	return transfers, nil

}
