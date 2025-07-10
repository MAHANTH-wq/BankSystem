package db

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testDB *pgxpool.Pool
var testQueries *Queries

func TestMain(m *testing.M) {
	var err error
	ctx := context.Background()
	testDB, err = pgxpool.New(ctx, dbSource)
	if err != nil {
		panic(err)
	}
	defer testDB.Close()
	testQueries = New(testDB)

	// Run the tests
	os.Exit(m.Run())
}
