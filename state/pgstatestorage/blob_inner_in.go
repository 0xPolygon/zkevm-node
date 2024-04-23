package pgstatestorage

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

const blobInnerFields = "blob_sequence_index, blob_inner_num, blob_type, max_sequence_timestamp, zk_gas_limit, l1_info_tree_leaf_index, l1_info_tree_root, updated_at"
const blobInnerFieldsTypeBlob = "blob_type_index,blob_type_z, blob_type_y,blob_type_commitment,blob_type_proof"

func (p *PostgresStorage) AddBlobInner(ctx context.Context, blobInner *state.BlobInner, dbTx pgx.Tx) error {
	sql := "INSERT INTO state.blob_inner_in (" + blobInnerFields
	if blobInner.Type == state.TypeBlobTransaction {
		sql += "," + blobInnerFieldsTypeBlob
	}
	sql += ") VALUES ($1, $2, $3, $4, $5, $6, $7, $8"
	if blobInner.Type == state.TypeBlobTransaction {
		sql += ", $9, $10,$11,$12,$13"
	}
	sql += ")"
	e := p.getExecQuerier(dbTx)
	arguments := []interface{}{blobInner.BlobSequenceIndex, blobInner.BlobInnerNum, blobInner.Type.String(), blobInner.MaxSequenceTimestamp, blobInner.ZkGasLimit, blobInner.L1InfoLeafIndex, blobInner.L1InfoTreeRoot.String(), time.Now()}
	if blobInner.Type == state.TypeBlobTransaction {
		commitment, err := blobInner.BlobBlobTypeParams.Commitment.MarshalText()
		if err != nil {
			return err
		}
		proof, err := blobInner.BlobBlobTypeParams.Proof.MarshalText()
		if err != nil {
			return err
		}
		arguments = append(arguments, blobInner.BlobBlobTypeParams.BlobIndex, common.Bytes2Hex(blobInner.BlobBlobTypeParams.Z), common.Bytes2Hex(blobInner.BlobBlobTypeParams.Y), commitment, proof)
	}
	_, err := e.Exec(ctx, sql, arguments...)
	return err

}
