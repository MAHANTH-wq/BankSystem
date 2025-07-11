package db

import (
	"context"
	"os"
	"simplebank/util"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testDB *pgxpool.Pool
var testQueries *Queries

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	ctx := context.Background()
	testDB, err = pgxpool.New(ctx, config.DBSource)
	if err != nil {
		panic(err)
	}
	defer testDB.Close()
	testQueries = New(testDB)

	// Run the tests
	os.Exit(m.Run())
}
