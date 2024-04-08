package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type migrationTest0018 struct{}

func (m migrationTest0018) InsertData(db *sql.DB) error {
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := db.Exec(addBlock, 1, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"); err != nil {
		return err
	}
	return nil
}

func (m migrationTest0018) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash, checked) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(addBlock, 2, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1", true)
	assert.NoError(t, err)
	_, err = db.Exec(addBlock, 3, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1", false)
	assert.NoError(t, err)
	const sql = `SELECT count(*) FROM state.block WHERE checked = true`
	row := db.QueryRow(sql)
	var result int
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 2, result)

	const sqlCheckedFalse = `SELECT count(*) FROM state.block WHERE checked = false`
	row = db.QueryRow(sqlCheckedFalse)

	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 1, result)
}

func (m migrationTest0018) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	var result int

	// Check column wip doesn't exists in state.batch table
	const sql = `SELECT count(*) FROM state.block`
	row := db.QueryRow(sql)
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 3, result)
}

func TestMigration0018(t *testing.T) {
	runMigrationTest(t, 18, migrationTest0018{})
}
