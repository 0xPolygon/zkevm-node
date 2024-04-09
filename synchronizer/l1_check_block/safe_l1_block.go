package l1_check_block

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/rpc"
)

type L1BlockPoint int

const (
	FinalizedBlockNumber L1BlockPoint = 2
	SafeBlockNumber      L1BlockPoint = 1
	LastBlockNumber      L1BlockPoint = 0
)

func (v L1BlockPoint) ToString() string {
	switch v {
	case FinalizedBlockNumber:
		return "finalized"
	case SafeBlockNumber:
		return "safe"
	case LastBlockNumber:
		return "latest"
	}
	return "Unknown"
}

// StringToL1BlockPoint converts a string to a L1BlockPoint
func StringToL1BlockPoint(s string) L1BlockPoint {
	switch s {
	case "finalized":
		return FinalizedBlockNumber
	case "safe":
		return SafeBlockNumber
	case "latest":
		return LastBlockNumber
	default:
		return FinalizedBlockNumber
	}
}

func (v L1BlockPoint) ToGethRequest() *big.Int {
	switch v {
	case FinalizedBlockNumber:
		return big.NewInt(int64(rpc.FinalizedBlockNumber))
	case SafeBlockNumber:
		return big.NewInt(int64(rpc.SafeBlockNumber))
	case LastBlockNumber:
		return nil
	}
	return big.NewInt(int64(v))
}

type SafeL1BlockNumberFetch struct {
	// SafeBlockPoint is the block number that is reference to l1 Block
	SafeBlockPoint L1BlockPoint
	// Offset is a vaule add to the L1 block
	Offset int
}

// NewSafeL1BlockNumberFetch creates a new SafeL1BlockNumberFetch
func NewSafeL1BlockNumberFetch(safeBlockPoint L1BlockPoint, offset int) *SafeL1BlockNumberFetch {
	return &SafeL1BlockNumberFetch{
		SafeBlockPoint: safeBlockPoint,
		Offset:         offset,
	}
}

func (p *SafeL1BlockNumberFetch) GetSafeBlockNumber(ctx context.Context, requester L1Requester) (uint64, error) {
	l1SafePointBlock, err := requester.HeaderByNumber(ctx, p.SafeBlockPoint.ToGethRequest())
	if err != nil {
		log.Errorf("%s: Error getting L1 block %d. err: %s", logPrefix, p.String(), err.Error())
		return uint64(0), err
	}
	return l1SafePointBlock.Number.Uint64() + uint64(p.Offset), nil
}

func (p *SafeL1BlockNumberFetch) String() string {
	return p.SafeBlockPoint.ToString() + " offset:" + string(p.Offset)
}
