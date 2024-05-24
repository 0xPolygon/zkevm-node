package sequencer

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/state"
)

func (f *finalizer) DSSendL2Block(ctx context.Context, batchNumber uint64, blockResponse *state.ProcessBlockResponse, l1InfoTreeIndex uint32) error {
	forkID := f.stateIntf.GetForkIDByBatchNumber(batchNumber)

	// Send data to streamer
	if f.streamServer != nil {
		l2Block := state.DSL2Block{
			BatchNumber:     batchNumber,
			L2BlockNumber:   blockResponse.BlockNumber,
			Timestamp:       int64(blockResponse.Timestamp),
			L1InfoTreeIndex: l1InfoTreeIndex,
			L1BlockHash:     blockResponse.BlockHashL1,
			GlobalExitRoot:  blockResponse.GlobalExitRoot,
			Coinbase:        f.sequencerAddress,
			ForkID:          uint16(forkID),
			BlockHash:       blockResponse.BlockHash,
			StateRoot:       blockResponse.BlockHash, //From etrog, the blockhash is the block root
		}

		l2Transactions := []state.DSL2Transaction{}

		for _, txResponse := range blockResponse.TransactionResponses {
			binaryTxData, err := txResponse.Tx.MarshalBinary()
			if err != nil {
				return err
			}

			l2Transaction := state.DSL2Transaction{
				L2BlockNumber:               blockResponse.BlockNumber,
				EffectiveGasPricePercentage: uint8(txResponse.EffectivePercentage),
				IsValid:                     1,
				EncodedLength:               uint32(len(binaryTxData)),
				Encoded:                     binaryTxData,
				StateRoot:                   txResponse.StateRoot,
			}

			l2Transactions = append(l2Transactions, l2Transaction)
		}

		f.checkDSBufferIsFull(ctx)

		f.dataToStream <- state.DSL2FullBlock{
			DSL2Block: l2Block,
			Txs:       l2Transactions,
		}

		f.dataToStreamCount.Add(1)
	}

	return nil
}

func (f *finalizer) DSSendBatchBookmark(ctx context.Context, batchNumber uint64) {
	// Check if stream server enabled
	if f.streamServer != nil {
		f.checkDSBufferIsFull(ctx)

		// Send batch bookmark to the streamer
		f.dataToStream <- state.DSBookMark{
			Type:  state.BookMarkTypeBatch,
			Value: batchNumber,
		}

		f.dataToStreamCount.Add(1)
	}
}

func (f *finalizer) checkDSBufferIsFull(ctx context.Context) {
	if f.dataToStreamCount.Load() == datastreamChannelBufferSize {
		f.Halt(ctx, fmt.Errorf("datastream channel buffer full"), true)
	}
}
