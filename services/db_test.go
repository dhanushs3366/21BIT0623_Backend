package dbtest_test

import (
	"database/sql"
	"testing"

	"github.com/dhanushs3366/21BIT0623_Backend.git/services/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *db.Store {
	sqlDB, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	store, err := db.GetNewStore(sqlDB)
	err = store.CreateUserTable()
}
