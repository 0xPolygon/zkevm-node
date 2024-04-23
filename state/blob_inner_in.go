package state

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/jackc/pgx/v4"
)

// BlobType is the type of the blob type
type BlobType uint8

const (
	// TypeCallData The data is stored on call data directly
	TypeCallData BlobType = 0
	// TypeBlobTransaction The data is stored on a blob
	TypeBlobTransaction BlobType = 1
	// TypeForcedBlob The data is a forced Blob
	TypeForcedBlob BlobType = 2
)

func (b BlobType) String() string {
	switch b {
	case TypeCallData:
		return "call_data"
	case TypeBlobTransaction:
		return "blob"
	case TypeForcedBlob:
		return "forced"
	default:
		return "Unknown"
	}
}

type BlobBlobTypeParams struct {
	BlobIndex  uint64
	Z          []byte
	Y          []byte
	Commitment kzg4844.Commitment
	Proof      kzg4844.Proof
}

// BlobInner struct
type BlobInner struct {
	BlobSequenceIndex    uint64      // Index of the blobSequence in DB (is a internal number)
	BlobInnerNum         uint64      // Incremental value, starts from 1
	Type                 BlobType    // Type of the blob
	MaxSequenceTimestamp time.Time   // it comes from SequenceBlobs call to contract
	ZkGasLimit           uint64      // it comes from SequenceBlobs call to contract
	L1InfoLeafIndex      uint32      // it comes from SequenceBlobs call to contract
	L1InfoTreeRoot       common.Hash // obtained from the L1InfoTree
	//Data                 []byte              // it comes from SequenceBlobs call to contract or from sidecar Blob
	BlobBlobTypeParams *BlobBlobTypeParams // Field only valid if BlobType == BlobTransaction

	// We don't need blockNumber because is in BlobSequence
	//BlockNumber             uint64
	//PreviousL1InfoTreeIndex uint32      // ?? we need that?
	//PreviousL1InfoTreeRoot  common.Hash // ?? we need that?
}

func (s *State) AddBlobInner(ctx context.Context, blobInner *BlobInner, dbTx pgx.Tx) error {
	return s.storage.AddBlobInner(ctx, blobInner, dbTx)
}
