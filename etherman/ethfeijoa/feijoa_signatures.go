package ethfeijoa

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// Events Feijoa Signatures
	// Events new ZkEvm/RollupBase
	// lastBlobSequenced is the count of blob sequenced after process this event
	//  if the first event have 1 blob -> lastBlobSequenced=1
	eventSequenceBlobsSignatureHash = crypto.Keccak256Hash([]byte("SequenceBlobs(uint64 indexed lastBlobSequenced)"))
)

func processEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	switch vLog.Topics[0] {
	case eventSequenceBlobsSignatureHash:
		var event SequenceBlobs
		err := contract.UnpackLog(&event, "SequenceBlobs", vLog.Data)
		if err != nil {
			return err
		}
		// Process event
		// ...
	}
	return nil
}
