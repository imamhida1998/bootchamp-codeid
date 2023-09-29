package postgresql

import (
	"bootchamp-codeid/db/postgresql/config"
	"testing"
)

func TestConnection_OpenConnection(t *testing.T) {

	db := OpenConnection(config.IMAM)
	defer db.Close()

	err := db.Ping()
	if err != nil {
		t.Fatal("ping error", err.Error())
	}

}
