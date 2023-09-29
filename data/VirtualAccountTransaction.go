package data

import (
	"bootchamp-codeid/db/mangodb/config"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TRANSFER = "Transfer"
	TOPUP    = "TopUp"
	PAYMENT  = "Payment"
)

type VATransaction struct {
	BaseNoSql

	HistoryID          primitive.ObjectID `json:"history_id" bson:"history_id"`
	NoVirtualAccount   string             `json:"no_virtual_account" bson:"no_virtual_account"`
	NameVirtualAccount string             `json:"name_virtual_account" bson:"name_virtual_account"`
	TransactionAmount  float32            `json:"transaction_amount" bson:"transaction_amount"`
	TransactionType    string             `json:"transaction_type" bson:"transaction_type"`
	Description        string             `json:"description" bson:"description"`
	TransactionID      string             `json:"transaction_id" bson:"transaction_id"`
}

func CreateVaTransaction(client *mongo.Client) VATransaction {
	VAtransaction := VATransaction{}
	VAtransaction.Client = client
	return VAtransaction
}

func (VAtransaction VATransaction) Collection() *mongo.Collection {
	return VAtransaction.Client.Database(config.DATABASEMONGO).Collection("virtualaccount_transaction")
}

func (VAtransaction VATransaction) Truncate(ctx context.Context) error {
	return VAtransaction.Collection().Drop(ctx)

}
func (VAtransaction VATransaction) AddVirtualAccountTransaction(ctx context.Context,
	NoVirtualAccount string,
	NameVirtualAccount string,
	TransactionAmount float32,
	TransactionType string,
	TransactionID string,
	Description string) (*VATransaction, error) {

	VAtransaction.HistoryID = primitive.NewObjectID()
	VAtransaction.NoVirtualAccount = NoVirtualAccount
	VAtransaction.NameVirtualAccount = NameVirtualAccount
	VAtransaction.TransactionType = TransactionType
	VAtransaction.TransactionAmount = TransactionAmount
	VAtransaction.Description = Description
	VAtransaction.TransactionID = TransactionID

	_, err := VAtransaction.Collection().InsertOne(ctx, VAtransaction)
	if err != nil {
		return nil, err
	}
	return &VAtransaction, nil
}

func (VAtransaction VATransaction) FindOene(ctx context.Context, HistoryID string) (error, error) {
	id, err := primitive.ObjectIDFromHex(HistoryID)
	if err != nil {
		return err, nil
	}
	filter := bson.D{
		primitive.E{
			Key:   "history_id",
			Value: id,
		},
	}
	result := VAtransaction.Collection().FindOne(ctx, filter)
	if result.Err() != nil {
		log.Println(result.Err())
		if result.Err().Error() == config.MONGO_NO_DOCUMENT {
			return errors.New("HistoryID Not Found"), nil
		}
		return result.Err(), nil
	}
	err = result.Decode(&VAtransaction)
	if err != nil {
		return err, nil
	}
	return nil, nil
}

func (VAtransaction VATransaction) Update(ctx context.Context,
	HistoryID string,
	NoVirtualAccount string,
	NameVirtualAccount string,
	TransactionAmount float32,
	reference string) error {

	id, err := primitive.ObjectIDFromHex(HistoryID)
	if err != nil {
		return err
	}

	filter := bson.D{
		primitive.E{
			Key:   "history_id",
			Value: id,
		},
	}
	set := bson.D{
		primitive.E{
			Key:   "no_virtual_account",
			Value: NoVirtualAccount,
		},
		primitive.E{
			Key:   "name_virtual_account",
			Value: NameVirtualAccount,
		},
		primitive.E{
			Key:   "transaction_amount",
			Value: TransactionAmount,
		},
		primitive.E{
			Key:   "reference",
			Value: reference,
		},
	}

	optionsAfter := options.After
	updateAfter := &options.FindOneAndUpdateOptions{
		ReturnDocument: &optionsAfter,
	}
	result := VAtransaction.Collection().FindOneAndUpdate(ctx, filter, set, updateAfter)
	if result.Err() != nil {
		return result.Err()
	}

	return result.Decode(&VAtransaction)

}

func (VAtransaction VATransaction) Delete(ctx context.Context, HistoryID string) error {

	id, err := primitive.ObjectIDFromHex(HistoryID)
	if err != nil {
		return err
	}
	filter := bson.D{
		primitive.E{
			Key:   "history_id",
			Value: id,
		},
	}
	result, err := VAtransaction.Collection().DeleteOne(ctx, filter, nil)
	if err != nil {
		return err
	}
	log.Println(result)
	return nil
}

func (VAtrasaction VATransaction) ListVirtualAccountTransaction(ctx context.Context) ([]VATransaction, error) {
	VAtrasactions := []VATransaction{}
	filter := bson.D{}
	cursor, err := VAtrasaction.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &VAtrasactions)
	if err != nil {
		return nil, err
	}
	return VAtrasactions, nil

}
