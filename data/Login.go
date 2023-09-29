package data

import (
	"bootchamp-codeid/db/mangodb/config"
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

type LoginNoSql struct {
	BaseNoSql
	LoginID     string `json:"login" bson:"login_id"`
	Token       string `json:"token" bson:"token"`
	AccountNo   string `json:"NoAccount" bson:"account_no"`
	CreatedAt   string `json:"CreateAt" bson:"created_at"`
	Expired     bool   `json:"Expired" bson:"expired"`
	ExpiredTime string `json:"TimeExpired" bson:"expired_time"`
}

// constructor
func CreateLoginNoSql(client *mongo.Client) LoginNoSql {
	login := LoginNoSql{}
	login.Client = client
	return login
}

func (login LoginNoSql) Collection() *mongo.Collection {
	return login.Client.Database(config.DATABASEMONGO).Collection("login")
}

func (login LoginNoSql) Truncate(ctx context.Context) error {
	return login.Collection().Drop(ctx)
}

func (login LoginNoSql) AddLogin(ctx context.Context,
	Token string,
	//Login_id string,
	Account_no string,
) error {

	login.Token = Token
	login.LoginID = uuid.New().String()
	login.AccountNo = Account_no
	login.CreatedAt = time.Now().Format(time.RFC1123)
	login.ExpiredTime = time.Now().Format(time.Hour.String())

	_, err := login.Collection().InsertOne(ctx, login)
	if err != nil {
		return err
	}
	return nil

}

func (login *LoginNoSql) FindOneByToken(ctx context.Context, token string) error {

	filter := bson.D{
		primitive.E{Key: "token", Value: token},
	}

	result := login.Collection().FindOne(ctx, filter)
	if result.Err() != nil {

		log.Println(result.Err())
		if result.Err().Error() == config.MONGO_NO_DOCUMENT {
			return errors.New("login_not_found")
		}
		return result.Err()
	}

	err := result.Decode(&login)
	if err != nil {
		return err
	}

	return nil
}

func (login LoginNoSql) Delete(ctx context.Context, id string) error {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		primitive.E{Key: "_id", Value: oid},
	}

	result, err := login.Collection().DeleteOne(ctx, filter, nil)
	if err != nil {
		return err
	}
	log.Println(result)
	return nil

}

func (login LoginNoSql) ListLogins(ctx context.Context) ([]LoginNoSql, error) {

	logins := []LoginNoSql{}
	filter := bson.D{}
	cursor, err := login.Collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &logins)
	if err != nil {
		return nil, err
	}

	return logins, nil

}
