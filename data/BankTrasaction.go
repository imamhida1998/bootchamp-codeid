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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BankTransaction struct {
	BaseNoSql

	TransactionID     primitive.ObjectID `json:"transaction_id" bson:"transaction_id"`
	BankAccount_No    string             `json:"NoBankAccount" bson:"NoBankAccount"`
	BankAccount_Name  string             `json:"NameBankAccount" bson:"NameBankAccount"`
	TransactionAmount float32            `json:"transaction_amount" bson:"transaction_amount"`
	Reference         string             `json:"reference" bson:"reference"`
	UpdateAt          string             `json:"update_at" bson:"update_at"`
}

func CreateBankTransaction(client *mongo.Client) BankTransaction {
	databakvirtual := BankTransaction{}
	databakvirtual.Client = client
	return databakvirtual
}
func (databakvirtual BankTransaction) Collection() *mongo.Collection {
	return databakvirtual.Client.Database(config.DATABASEMONGO).Collection("banktransaction")
}

func (databakvirtual BankTransaction) Truncate(ctx context.Context) error {
	return databakvirtual.Collection().Drop(ctx)
}

func (databakvirtual BankTransaction) InsertBankTransaction(ctx context.Context,
	BankAccountNo string,
	BankAccoutName string,
	TransactionAmmount float32,
	Reference string) (*BankTransaction, error) {

	databakvirtual.TransactionID = primitive.NewObjectID()
	databakvirtual.BankAccount_No = BankAccountNo
	databakvirtual.BankAccount_Name = BankAccoutName
	databakvirtual.TransactionAmount = TransactionAmmount
	databakvirtual.Reference = Reference

	_, err := databakvirtual.Collection().InsertOne(ctx, databakvirtual)
	if err != nil {
		return nil, err
	}
	return &databakvirtual, nil

}

func (databakvirtual BankTransaction) FindOne(ctx context.Context, TransactionID string) error {
	id, err := primitive.ObjectIDFromHex(TransactionID)
	if err != nil {
		return err
	}
	filter := bson.D{
		primitive.E{
			Key:   "transaction_id",
			Value: id,
		},
	}
	result := databakvirtual.Collection().FindOne(ctx, filter)
	if result.Err() != nil {
		log.Println(result.Err())
		if result.Err().Error() == config.MONGO_NO_DOCUMENT {
			return errors.New("DataBank_NOT_FOUND")
		}
		return result.Err()
	}
	err = result.Decode(&databakvirtual)
	if err != nil {
		return err
	}
	return nil
}

func (databakvirtual *BankTransaction) UpdateBankTransaction(ctx context.Context,
	TransactionID string,
	TransactionAmount float32) error {
	id, err := primitive.ObjectIDFromHex(TransactionID)
	if err != nil {
		return err
	}
	filter := bson.D{
		primitive.E{
			Key:   "transaction_id",
			Value: id,
		},
	}
	set := bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{
					Key:   "transaction_amount",
					Value: TransactionAmount,
				},
				primitive.E{
					Key:   "update_at",
					Value: time.Now().UTC().Format(time.RFC1123),
				},
			},
		},
	}
	optionsAfter := options.After
	updateOptions := &options.FindOneAndUpdateOptions{
		ReturnDocument: &optionsAfter,
	}
	result := databakvirtual.Collection().FindOneAndUpdate(ctx, filter, set, updateOptions)
	if result.Err() != nil {
		return result.Err()
	}
	return result.Decode(&databakvirtual)
}

func (databakvirtual BankTransaction) ListBankTransaction(ctx context.Context) ([]BankTransaction, error) {
	BankTransactions := []BankTransaction{}
	filter := bson.D{}
	cursor, err := databakvirtual.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &BankTransactions)
	if err != nil {
		return nil, err
	}
	return BankTransactions, nil

}
