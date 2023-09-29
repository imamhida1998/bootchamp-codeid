package data

import (
	"bootchamp-codeid/helpers"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type VirtualAccount struct {
	BaseData
	Account_no  string  `json:"NoAccount" bson:"account___no"`
	No_hp       string  `json:"NoHandphone" bson:"NoHandphone"`
	Email       string  `json:"Email" bson:"email"`
	Seq_no      int     `json:"seq_no" bson:"seq_no"`
	AccountName string  `json:"AccountName" bson:"account_name"`
	Pin         string  `json:"PIN" bson:"pin"`
	Password    string  `json:"password" bson:"password"`
	Saldo       float32 `json:"Saldo" bson:"saldo"`
	CreatedAt   string  `json:"CreatedAt" bson:"created_at"`
	//	CreatedAt   string  `json:"CreatedAt" bson:"created_at"`
	UpdatedAt *string `json:"UpdatedAt" bson:"updated_at"`
}

func CreateVirtualAccount(db *sql.DB) VirtualAccount {
	virtualaccount := VirtualAccount{}
	virtualaccount.DB = db
	return virtualaccount
}
func CreateVirtualAccountWithTransaction(transaction *sql.Tx) VirtualAccount {
	virtualaccount := VirtualAccount{}
	virtualaccount.Transaction = transaction
	virtualaccount.UseTransaction = true
	return virtualaccount
}

// func (virtualaccount VirtualAccount)MaskPadding()
func (virtualaccount VirtualAccount) Migrate() error {

	sqlDropTable := "DROP TABLE IF EXISTS public.virtual_accounts CASCADE"

	sqlCreateTable := `
CREATE TABLE public.virtual_accounts
(
    account_no character varying(40) NOT NULL,
    no_hp character varying NOT NULL,
	email character varying(40) NOT NULL,
	account_name character varying(40) NOT NULL,
	pin character varying(25) NOT NULL,
	password character varying(60) NOT NULL,
	saldo double precision NOT NULL,
	created_at character varying NOT NULL,
    updated_at character varying,
    CONSTRAINT virtual_accounts_pkey PRIMARY KEY (account_no)
)

TABLESPACE pg_default;

ALTER TABLE public.virtual_accounts
    OWNER to postgres;
	`

	_, err := virtualaccount.DB.Exec(sqlDropTable)
	if err != nil {
		return err
	}

	_, err = virtualaccount.DB.Exec(sqlCreateTable)
	if err != nil {
		return err
	}

	return nil

}

func (virtualaccount VirtualAccount) selectQuery() string {
	sql := `select account_no, 
	no_hp, email, account_name, pin, 
	password, saldo, 
	created_at, 
	updated_at from virtual_accounts`
	return sql
}
func (virtualaccount VirtualAccount) fetchRow(cursor *sql.Rows) (VirtualAccount, error) {
	cek := VirtualAccount{}
	err := cursor.Scan(
		&cek.Account_no,
		&cek.No_hp,
		&cek.Email,
		&cek.AccountName,
		&cek.Pin,
		&cek.Password,
		&cek.Saldo,
		&cek.CreatedAt,
		&cek.UpdatedAt)
	if err != nil {
		return VirtualAccount{}, err
	}
	return cek, nil
}
func (virtualaccount VirtualAccount) InsertVirtualAccount(Account_no string, No_hp string, Email string, AccountName string, Pin string, Password string, saldo float32) (error, error) {
	sql := `insert into virtual_accounts (account_no, 
		no_hp, email, account_name, pin, 
		password, saldo, 
		created_at) values ($1,$2,$3,$4,$5,$6,$7,$8)`

	//https://pkg.go.dev/github.com/shopspring/decimal#example-NewFromFloat
	// Ancount_no

	//n := number.Decimal(20,number.Pad('0'), number.FormatWidth(6))
	sds, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	epin := base64.StdEncoding.EncodeToString([]byte(Pin))
	result, err := virtualaccount.Exec(sql,
		Account_no,
		No_hp,
		Email,
		AccountName,
		epin,
		sds,
		saldo,
		time.Now().UTC().Format(time.RFC1123))
	if err != nil {
		return err, nil
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err, nil
	}

	log.Println("VirtualAccount.Add.AffectedRows=", affectedRows)
	return nil, nil

}
func (virtualaccount VirtualAccount) ReadVirtualAccount() ([]VirtualAccount, error) {
	//                0             1
	sql := virtualaccount.selectQuery()
	cursor, err := virtualaccount.Query(sql)
	if err != nil {
		return nil, err
	}

	virtualaccounts := []VirtualAccount{}

	for cursor.Next() { //looping rows
		b, err := virtualaccount.fetchRow(cursor)
		if err != nil {
			return nil, err
		}
		virtualaccounts = append(virtualaccounts, b)
	}

	return virtualaccounts, nil
}

/*
	account_no,
	no_hp, email, account_name, pin,
	password, saldo,
	created_at,
	updated_at*/

func (virtualaccount VirtualAccount) UpdateVirtualAccount(Account_no string, No_hp string, Email string, AccountName string, Pin string, Password string, Saldo float32) error {

	sql := `update virtual_accounts
	set No_hp=$2,
	email=$3,
	account_name=$4,
	pin=$5,
	password=$6,
	saldo=$7,
	updated_at=$8

	where account_no=$1
	`
	result, err := virtualaccount.Exec(sql,
		Account_no,
		No_hp,
		Email,
		AccountName,
		Pin,
		Password,
		Saldo,
		time.Now().UTC().Format(time.RFC1123))
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("VirtualAccount.Update.AffectedRows=", affectedRows)

	return nil

}

func (virtualaccount VirtualAccount) UpdateAccountPin(Account_no string, AccountName string, Pin string) error {
	sql := `update virtual_accounts 
			set account_name =$2,pin=$3,updated_at=$4
			where account_no=$1`
	result, err := virtualaccount.Exec(sql, Account_no, AccountName, Pin, time.Now().Format(time.RFC1123))
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("VirtualAccount.Update.AffectedRows=", affectedRows)
	return nil

}

func (virtualaccount VirtualAccount) Remove(Account_no string) (error, error) {

	sql := "delete from virtual_accounts where account_no=$1"
	result, err := virtualaccount.Exec(sql, Account_no)
	if err != nil {
		return err, nil
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err, nil
	}

	log.Println("VirtualAccount.Remove.AffectedRows=", affectedRows)
	return nil, nil
}
func (virtualaccount VirtualAccount) ResetData() error {

	sql := `delete from virtual_accounts`
	result, err := virtualaccount.Exec(sql)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("VirtualAccount.ResetData.AffectedRows=", affectedRows)
	return nil

}
func (virtualaccount VirtualAccount) FindByPhone(NoHp string) (*VirtualAccount, error) {

	if helpers.Empty(NoHp) {
		return nil, errors.New("invalid_no_VirtualAccount")
	}

	sql := `select * from virtual_accounts where no_hp=$1`

	cursor, err := virtualaccount.Query(sql, NoHp)
	if err != nil {
		return nil, err
	}

	if cursor.Next() {
		b, err := virtualaccount.fetchRow(cursor)
		if err != nil {
			return nil, err
		}
		return &b, nil
	}

	return nil, errors.New("NoHp tidak ditemukan")

}

func (virtualaccount VirtualAccount) FindOne(Account_no string) (*VirtualAccount, error) {

	if helpers.Empty(Account_no) {
		return nil, errors.New("invalid_no_VirtualAccount")
	}

	sql := `select * from virtual_accounts where account_no=$1`

	cursor, err := virtualaccount.Query(sql, Account_no)
	if err != nil {
		return nil, err
	}

	if cursor.Next() {
		b, err := virtualaccount.fetchRow(cursor)
		if err != nil {
			return nil, err
		}
		return &b, nil
	}

	return nil, errors.New("VirtualAccount tidak ditemukan")

}
