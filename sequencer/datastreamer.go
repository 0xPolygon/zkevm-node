package sequencer

import (
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/datastream"
	"github.com/ethereum/go-ethereum/common"
)

func (f *finalizer) DSSendL2Block(batchNumber uint64, blockResponse *state.ProcessBlockResponse, l1InfoTreeIndex uint32, minTimestamp uint64) error {
	forkID := f.stateIntf.GetForkIDByBatchNumber(batchNumber)

	// Send data to streamer
	if f.streamServer != nil {
		l2Block := state.DSL2Block{
			BatchNumber:     batchNumber,
			L2BlockNumber:   blockResponse.BlockNumber,
			Timestamp:       blockResponse.Timestamp,
			Min_timestamp:   minTimestamp,
			L1InfoTreeIndex: l1InfoTreeIndex,
			L1BlockHash:     blockResponse.BlockHashL1,
			GlobalExitRoot:  blockResponse.GlobalExitRoot,
			Coinbase:        f.sequencerAddress,
			ForkID:          forkID,
			BlockHash:       blockResponse.BlockHash,
			StateRoot:       blockResponse.BlockHash, //From etrog, the blockhash is the block root
			BlockInfoRoot:   blockResponse.BlockInfoRoot,
		}

		l2Transactions := []state.DSL2Transaction{}

		for i, txResponse := range blockResponse.TransactionResponses {
			binaryTxData, err := txResponse.Tx.MarshalBinary()
			if err != nil {
				return err
			}

			l2Transaction := state.DSL2Transaction{
				L2BlockNumber:               blockResponse.BlockNumber,
				EffectiveGasPricePercentage: uint8(txResponse.EffectivePercentage),
				Index:                       uint64(i),
				IsValid:                     1,
				EncodedLength:               uint32(len(binaryTxData)),
				Encoded:                     binaryTxData,
				StateRoot:                   txResponse.StateRoot,
			}

			if txResponse.Logs != nil && len(txResponse.Logs) > 0 {
				l2Transaction.Index = uint64(txResponse.Logs[0].TxIndex)
			}

			l2Transactions = append(l2Transactions, l2Transaction)
		}

		log.Infof("[ds-debug] sending l2block %d to datastream channel", blockResponse.BlockNumber)
		f.dataToStream <- state.DSL2FullBlock{
			DSL2Block: l2Block,
			Txs:       l2Transactions,
		}
	}

	return nil
}

func (f *finalizer) DSSendBatchBookmark(batchNumber uint64) {
	// Check if stream server enabled
	if f.streamServer != nil {
		// Send batch bookmark to the streamer
		f.dataToStream <- datastream.BookMark{
			Type:  datastream.BookmarkType_BOOKMARK_TYPE_BATCH,
			Value: batchNumber,
		}
	}
}

func (f *finalizer) DSSendBatchStart(batchNumber uint64, isForced bool) {
	forkID := f.stateIntf.GetForkIDByBatchNumber(batchNumber)

	batchStart := datastream.BatchStart{
		Number: batchNumber,
		ForkId: forkID,
	}

	if isForced {
		batchStart.Type = datastream.BatchType_BATCH_TYPE_FORCED
	} else {
		batchStart.Type = datastream.BatchType_BATCH_TYPE_REGULAR
	}

	if f.streamServer != nil {
		// Send batch start to the streamer
		f.dataToStream <- batchStart
	}
}

func (f *finalizer) DSSendBatchEnd(batchNumber uint64, stateRoot common.Hash, localExitRoot common.Hash) {
	if f.streamServer != nil {
		// Send batch end to the streamer
		f.dataToStream <- datastream.BatchEnd{
			Number:        batchNumber,
			StateRoot:     stateRoot.Bytes(),
			LocalExitRoot: localExitRoot.Bytes(),
		}
	}
}
