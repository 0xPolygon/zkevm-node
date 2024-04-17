package pgstatestorage

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// AddBlobSequence adds a new blob sequence to the state.
func (p *PostgresStorage) AddBlobSequence(ctx context.Context, blobSequence *state.BlobSequence, dbTx pgx.Tx) error {
	const addBlobSequenceSQL = "INSERT INTO state.blob_sequence (index, coinbase, final_acc_input_hash, last_blob_sequenced, create_at) VALUES ($1, $2, $3, $4, $5)"

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addBlobSequenceSQL, blobSequence.Index, blobSequence.L2Coinbase, blobSequence.FinalAccInputHash, blobSequence.LastBlobSequenced, blobSequence.CreateAt)
	return err
}

// GetLastBlobSequence returns the last blob sequence stored in the state.
func (p *PostgresStorage) GetLastBlobSequence(ctx context.Context, dbTx pgx.Tx) (*state.BlobSequence, error) {
	var (
		coinbase          string
		finalAccInputHash string
		lastBlobSequenced uint64
		createAt          time.Time
		blobSequence      state.BlobSequence
	)
	const getLastBlobSequenceSQL = "SELECT index, coinbase, final_acc_input_hash, last_blob_sequenced, create_at FROM state.blob_sequence ORDER BY index DESC LIMIT 1"

	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getLastBlobSequenceSQL).Scan(&blobSequence.Index, &coinbase, &finalAccInputHash, &lastBlobSequenced, &createAt)
	if errors.Is(err, pgx.ErrNoRows) {
		// If none on database return a nil object
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	blobSequence.L2Coinbase = common.HexToAddress(coinbase)
	blobSequence.FinalAccInputHash = common.HexToHash(finalAccInputHash)
	blobSequence.LastBlobSequenced = lastBlobSequenced
	blobSequence.CreateAt = createAt
	return &blobSequence, nil
}
