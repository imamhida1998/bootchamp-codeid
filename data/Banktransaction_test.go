package data

import (
	"bootchamp-codeid/db/mangodb"
	"context"
	"testing"
)

func TestBankTransaction_InsertBankTransaction(t *testing.T) {
	ctx := context.TODO()
	conn, err := mangodb.OpenMongoDb(ctx)
	if err != nil {
		return
	}
	defer mangodb.CloseMongoDb(ctx, conn)

}
