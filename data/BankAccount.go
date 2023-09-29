package data

import (
	"bootchamp-codeid/helpers"
	"database/sql"
	"errors"
	"log"
	"time"
)

type BankAccount struct {
	BaseData
	BankAccount_No    string  `json:"NoBankAccount"`
	BankAccount_Owner string  `json:"BankAccountOwner"`
	Saldo             float32 `json:"Saldo"`
	CreatedAt         string  `json:"CreateAt"`
}

func CreateBankAccount(db *sql.DB) BankAccount {
	accountbank := BankAccount{}
	accountbank.DB = db
	return accountbank
}
func CreateBankAccountWithTransaction(transaction *sql.Tx) BankAccount {
	accountbank := BankAccount{}
	accountbank.Transaction = transaction
	accountbank.UseTransaction = true
	return accountbank
}
func (bankacc BankAccount) Migrate() error {

	sqlDropTable := "DROP TABLE IF EXISTS public.bankaccount CASCADE"

	sqlCreateTable := `
CREATE TABLE public.bankaccount
(
    bank_account_no character varying(40) NOT NULL,
    bank_account_owner character varying(60) NOT NULL,
	saldo double precision NOT NULL,
    created_at character varying NOT NULL
)

TABLESPACE pg_default;

ALTER TABLE public.bankaccount
    OWNER to postgres;
	`

	_, err := bankacc.DB.Exec(sqlDropTable)
	if err != nil {
		return err
	}

	_, err = bankacc.DB.Exec(sqlCreateTable)
	if err != nil {
		return err
	}

	return nil

}

func (bankacc BankAccount) selectQuery() string {
	sql := `select  
	bank_account_no, 
	bank_account_owner,  
	saldo, 
	created_at from bankaccount`
	return sql
}
func (bankacc BankAccount) fetchRow(cursor *sql.Rows) (BankAccount, error) {
	cek := BankAccount{}
	err := cursor.Scan(
		&cek.BankAccount_No,
		&cek.BankAccount_Owner,
		&cek.Saldo,
		&cek.CreatedAt)
	if err != nil {
		return BankAccount{}, err
	}
	return cek, nil
}
func (bankacc BankAccount) InsertBankAccount(
	BankAccount_No string,
	BankAccount_Owner string,
	saldo float32) error {
	sql := `insert into bankaccount (
	bank_account_no, 
	bank_account_owner,  
	saldo, 
	created_at)
	values ($1,$2,$3,$4)`
	result, err := bankacc.Exec(sql,
		BankAccount_No,
		BankAccount_Owner,
		saldo,
		time.Now().UTC().Format(time.RFC1123))
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("BankAccount.Insert.AffectedRows=", affectedRows)
	return nil

}
func (bankacc BankAccount) ReadBankAccount() ([]BankAccount, error) {
	//                0             1
	sql := bankacc.selectQuery()
	cursor, err := bankacc.Query(sql)
	if err != nil {
		return nil, err
	}

	bankaccs := []BankAccount{}

	for cursor.Next() { //looping rows
		b, err := bankacc.fetchRow(cursor)
		if err != nil {
			return nil, err
		}
		bankaccs = append(bankaccs, b)
	}

	return bankaccs, nil
}

/*
	account_no,
	no_hp, email, account_name, pin,
	password, saldo,
	created_at,
	updated_at*/

func (bankacc BankAccount) UpdateBankAccount(
	BankAccount_No string,
	BankAccount_Owner string,
	Saldo float32) error {

	sql := `update bankaccount
	set bank_account_owner=$2, 
	saldo=$3

	where bank_account_no=$1
	`
	result, err := bankacc.Exec(sql,
		BankAccount_No,
		BankAccount_Owner,
		Saldo)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("BankAccount.Update.AffectedRows=", affectedRows)

	return nil

}

func (bankacc BankAccount) Remove(Account_no string) error {

	sql := "delete from bankaccount where bank_account_no=$1"
	result, err := bankacc.Exec(sql, Account_no)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("BankAccount.Remove.AffectedRows=", affectedRows)
	return nil
}
func (bankacc BankAccount) ResetData() error {

	sql := `delete from bankaccount`
	result, err := bankacc.Exec(sql)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("BankAccount.ResetData.AffectedRows=", affectedRows)
	return nil

}
func (bankacc BankAccount) FindOne(Account_no string) (*BankAccount, error) {

	if helpers.Empty(Account_no) {
		return nil, errors.New("invalid_no_BankAccount")
	}

	sql := bankacc.selectQuery()
	sql += ` where bank_account_no=$1`

	cursor, err := bankacc.Query(sql, Account_no)
	if err != nil {
		return nil, err
	}

	if cursor.Next() {
		b, err := bankacc.fetchRow(cursor)
		if err != nil {
			return nil, err
		}
		return &b, nil
	}

	return nil, errors.New("BankAccount tidak ditemukan")

}
