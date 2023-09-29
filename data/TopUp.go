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

type TopUp struct {
	BaseNoSql `bson:"-"`

	TopUpID            primitive.ObjectID `json:"TopID" bson:"TopID"`
	BankCode           string             `json:"BankCode" bson:"bankCode"`
	BankName           string             `json:"BankName" bson:"BankName"`
	NoVirtualAccount   string             `json:"no_virtual_account" bson:"no_virtual_account"`
	NameVirtualAccount string             `json:"name_virtual_account" bson:"name_virtual_account"`
	TopUpAmount        float32            `json:"top_up_amount" bson:"top_up_amount"`
	CreateAt           string             `json:"create_at" bson:"create_at"`
}

func CreateTopUp(client *mongo.Client) TopUp {
	topup := TopUp{}
	topup.Client = client
	return topup
}
func (topup TopUp) Collection() *mongo.Collection {
	return topup.Client.Database(config.DATABASEMONGO).Collection("topups")
}

func (topup TopUp) Truncate(ctx context.Context) error {
	return topup.Collection().Drop(ctx)
}

func (topup TopUp) AddTopUp(ctx context.Context,
	backcode string,
	backname string,
	NoVirtualAccount string,
	NameVirtualAccount string,
	TopUpAmount float32,
) (*TopUp, error) {

	topup.TopUpID = primitive.NewObjectID()
	topup.BankCode = backcode
	topup.BankName = backname
	topup.NoVirtualAccount = NoVirtualAccount
	topup.NameVirtualAccount = NameVirtualAccount
	topup.TopUpAmount = TopUpAmount
	topup.CreateAt = time.Now().Format(time.RFC1123)

	//topup.CreatedAt = time.Now().UTC()

	_, err := topup.Collection().InsertOne(ctx, topup)
	if err != nil {
		return nil, err
	}
	return &topup, nil

}

func (topup *TopUp) FindOne(ctx context.Context, Top_id string) (error, error) {

	oid, err := primitive.ObjectIDFromHex(Top_id)
	if err != nil {
		return err, nil
	}
	filter := bson.D{
		primitive.E{Key: "TopID", Value: oid},
	}

	result := topup.Collection().FindOne(ctx, filter)

	if result.Err() != nil {

		log.Println(result.Err())
		if result.Err().Error() == config.MONGO_NO_DOCUMENT {
			return errors.New("topup_not_found"), nil
		}
		return result.Err(), nil
	}

	err = result.Decode(&topup)
	if err != nil {
		return err, nil
	}

	return nil, nil

}

func (topup *TopUp) Update(ctx context.Context,
	id string,
	bankcode string,
	bankname string) error {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		primitive.E{Key: "TopID", Value: oid},
	}

	set := bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "BankCode", Value: bankcode},
				primitive.E{Key: "BankName", Value: bankname},
				//	primitive.E{Key: "updatedAt", Value: time.Now().UTC()},
			},
		},
	}

	optionsAfter := options.After
	updateOptions := &options.FindOneAndUpdateOptions{
		ReturnDocument: &optionsAfter,
	}

	result := topup.Collection().FindOneAndUpdate(ctx, filter, set, updateOptions)
	if result.Err() != nil {
		return result.Err()
	}

	return result.Decode(&topup)

}

func (topup TopUp) Delete(ctx context.Context, id string) (error, error) {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err, nil
	}

	filter := bson.D{
		primitive.E{Key: "TopID", Value: oid},
	}

	result, err := topup.Collection().DeleteOne(ctx, filter, nil)
	if err != nil {
		return err, nil
	}
	log.Println(result)
	return nil, nil

}

func (topup TopUp) ListTopUps(ctx context.Context) ([]TopUp, error) {

	topups := []TopUp{}
	filter := bson.D{}
	cursor, err := topup.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &topups)
	if err != nil {
		return nil, err
	}

	return topups, nil

}
