package service

import (
	"bootchamp-codeid/data"
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type VATransaction struct {
	BaseService
}

func CreateVirtualAccount(mongoClient *mongo.Client, db *sql.DB) VATransaction {
	virtualAccount := VATransaction{}
	virtualAccount.MongoClient = mongoClient
	virtualAccount.DB = db
	return virtualAccount
}

func (virtualAccount VATransaction) Login(ctx context.Context, NoHp string, Password string) (string, error) {
	// GET VIRTUAL ACCOUNT
	ModelVirtualAccount := data.CreateVirtualAccount(virtualAccount.DB)
	userVirtualAccount, err := ModelVirtualAccount.FindByPhone(NoHp)
	if err != nil {
		return "", errors.New("INVALID_VIRTUAL_ACCOUNT")
	}
	//GET Password
	err = bcrypt.CompareHashAndPassword([]byte(userVirtualAccount.Password), []byte(Password))
	if err != nil {
		return "", errors.New("INVALID_PASSWORD")
	}

	// CREATE TOKEN
	id := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	token := id + userVirtualAccount.Account_no + now
	encryptedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	encodtoken := base64.StdEncoding.EncodeToString(encryptedToken)
	login := data.CreateLoginNoSql(virtualAccount.MongoClient)
	err = login.AddLogin(ctx, encodtoken, userVirtualAccount.AccountName)
	if err != nil {
		return "", err
	}

	encodedToken := base64.StdEncoding.EncodeToString(encryptedToken)

	return encodedToken, nil
}

func (virtualAccount VATransaction) ParseToken(ctx context.Context, token string) (*data.LoginNoSql, error) {

	login := data.CreateLoginNoSql(virtualAccount.MongoClient)
	err := login.FindOneByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if login.Expired {
		return nil, errors.New("login_is_expired")
	}

	return &login, nil

}
