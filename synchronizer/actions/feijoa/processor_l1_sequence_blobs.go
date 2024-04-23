package feijoa

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

// stateProcessorSequenceBlobsInterface interface required from state
type stateProcessorSequenceBlobsInterface interface {
	AddBlobSequence(ctx context.Context, blobSequence *state.BlobSequence, dbTx pgx.Tx) error
	GetLastBlobSequence(ctx context.Context, dbTx pgx.Tx) (*state.BlobSequence, error)
}

// ProcessorSequenceBlobs processor for SequenceBlobs
type ProcessorSequenceBlobs struct {
	actions.ProcessorBase[ProcessorL1InfoTreeUpdate]
	state stateProcessorSequenceBlobsInterface
}

// NewProcessorSequenceBlobs new processor for SequenceBlobs
func NewProcessorSequenceBlobs(state stateProcessorSequenceBlobsInterface) *ProcessorSequenceBlobs {
	return &ProcessorSequenceBlobs{
		ProcessorBase: *actions.NewProcessorBase[ProcessorL1InfoTreeUpdate](
			[]etherman.EventOrder{etherman.SequenceBlobsOrder},
			actions.ForksIdOnlyFeijoa),
		state: state}
}

// Process process event
// - Store BlobSequence
// - Split Data into BlobInner (executor)
// - Store BlobInner
func (p *ProcessorSequenceBlobs) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	seqBlobs := l1Block.SequenceBlobs[order.Pos]

	err := p.storeBlobSequence(ctx, dbTx, &seqBlobs, l1Block.ReceivedAt)
	if err != nil {
		return err
	}

	for idx := range seqBlobs.Blobs {
		log.Infof("Blob %d: %s", idx, seqBlobs.Blobs[idx].String())
	}
	return nil
}

func (p *ProcessorSequenceBlobs) storeBlobSequence(ctx context.Context, dbTx pgx.Tx, seqBlobs *etherman.SequenceBlobs, createAt time.Time) error {
	if seqBlobs == nil || seqBlobs.EventData == nil {
		return fmt.Errorf("sequence blobs is nil or EventData is nil")
	}

	nextIndex := uint64(1)
	previousBlobSequenceOnState, err := p.state.GetLastBlobSequence(ctx, dbTx)
	if err != nil {
		return err
	}
	if previousBlobSequenceOnState != nil {
		nextIndex = previousBlobSequenceOnState.Index + 1
	}
	stateBlobSequence := state.BlobSequence{
		Index:             nextIndex,
		L2Coinbase:        seqBlobs.L2Coinbase,
		FinalAccInputHash: seqBlobs.FinalAccInputHash,
		LastBlobSequenced: seqBlobs.EventData.LastBlobSequenced,
		CreateAt:          createAt,
	}
	return p.state.AddBlobSequence(ctx, &stateBlobSequence, dbTx)
}
