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

/*
pay_id,

	merchant_va_account_no,
	merchant_va_account_name,
	src_va_account_no,
	src_va_account_name,
	pay_amount,
	created_at,
*/
type Pay struct {
	BaseNoSql                  `bson:"-"`
	PayID                      primitive.ObjectID `json:"PayID" bson:"PayID"`
	MerchantNoVirtualAccount   string             `json:"merchant_no_virtual_account" bson:"merchant_no_virtual_account"`
	MerchantNameVirtualAccount string             `json:"merchant_name_virtual_account" bson:"merchant_name_virtual_account"`
	SourceNoVirtualAccount     string             `json:"source_no_virtual_account" bson:"source_no_virtual_account"`
	SourceNameVirtualAccount   string             `json:"source_name_virtual_account" bson:"source_name_virtual_account"`
	PayAmount                  float32            `json:"pay_amount" bson:"pay_amount"`
	CreatedAt                  string             `json:"created_at" bson:"created_at"`
}

func CreatePay(client *mongo.Client) Pay {
	payment := Pay{}
	payment.Client = client
	return payment
}
func (payment Pay) Collection() *mongo.Collection {
	return payment.Client.Database(config.DATABASEMONGO).Collection("Pay")
}
func (payment Pay) Truncate(ctx context.Context) error {
	return payment.Collection().Drop(ctx)
}
func (payment Pay) AddPayment(ctx context.Context,
	merchant_va_account_no string,
	merchant_va_account_name string,
	src_va_account_no string,
	src_va_account_name string,
	pay_amount float32) (*Pay, error) {

	payment.PayID = primitive.NewObjectID()
	payment.MerchantNoVirtualAccount = merchant_va_account_no
	payment.MerchantNameVirtualAccount = merchant_va_account_name
	payment.SourceNoVirtualAccount = src_va_account_no
	payment.SourceNameVirtualAccount = src_va_account_name
	payment.PayAmount = pay_amount
	payment.CreatedAt = time.Now().Format(time.RFC1123)
	_, err := payment.Collection().InsertOne(ctx, payment)
	if err != nil {
		return nil, err
	}
	return &payment, nil

}

func (payment Pay) FindOne(ctx context.Context, PayID string) error {
	oid, err := primitive.ObjectIDFromHex(PayID)
	if err != nil {
		return nil
	}
	filter := bson.D{
		primitive.E{Key: "PayID", Value: oid},
	}

	result := payment.Collection().FindOne(ctx, filter)
	if result.Err() != nil {

		log.Println(result.Err())
		if result.Err().Error() == config.MONGO_NO_DOCUMENT {
			return errors.New("topup_not_found")
		}
		return result.Err()
	}
	err = result.Decode(&payment)
	if err != nil {
		return err
	}

	return nil

}
func (payment Pay) Update(ctx context.Context,
	PayID string,
	merchant_va_account_no string,
	merchant_va_account_name string,
	src_va_account_no string,
	src_va_account_name string,
	pay_amount float32) error {

	id, err := primitive.ObjectIDFromHex(PayID)
	if err != nil {
		return err
	}
	filter := bson.D{
		primitive.E{
			Key:   "PayID",
			Value: id,
		},
	}
	set := bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{
					Key:   "merchant_no_virtual_account",
					Value: merchant_va_account_no},

				primitive.E{
					Key:   "merchant_name_virtual_account",
					Value: merchant_va_account_name},
				primitive.E{
					Key:   "source_no_virtual_account",
					Value: src_va_account_no},
				primitive.E{
					Key:   "source_name_virtual_account",
					Value: src_va_account_name},
				primitive.E{
					Key:   "pay_amount",
					Value: pay_amount},
			},
		},
	}

	optionsAfter := options.After

	updateOptions := &options.FindOneAndUpdateOptions{
		ReturnDocument: &optionsAfter,
	}
	result := payment.Collection().FindOneAndUpdate(ctx, filter, set, updateOptions)
	if result.Err() != nil {
		return result.Err()
	}
	return result.Decode(&payment)
}

func (payment Pay) Delete(ctx context.Context, PayID string) (error, error) {

	id, err := primitive.ObjectIDFromHex(PayID)
	if err != nil {
		return err, nil
	}
	filter := bson.D{
		primitive.E{
			Key:   "PayID",
			Value: id},
	}
	result, err := payment.Collection().DeleteOne(ctx, filter, nil)
	if err != nil {
		return err, nil
	}
	log.Println(result)
	return nil, nil
}

func (payment Pay) ListPayment(ctx context.Context) ([]Pay, error) {

	payments := []Pay{}
	filter := bson.D{}
	cursor, err := payment.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &payments)
	if err != nil {
		return nil, err
	}
	return payments, nil

}
