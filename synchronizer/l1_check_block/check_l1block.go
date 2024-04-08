package actions

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jackc/pgx/v4"
)

const (
	logPrefix = "checkL1block:"
)

// This object check old L1block to double-check that the L1block hash is correct
// - Get first not checked block
// - Get last block on L1 (safe/finalized/ or minus -n)

type L1Requester interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

type StateInterfacer interface {
	GetFirstUncheckedBlock(ctx context.Context, dbTx pgx.Tx) (*state.Block, error)
	//GetFirstUncheckedBlockSinceNumber(ctx context.Context, firstBlockNum uint64, dbTx pgx.Tx) (*state.Block, error)
	//GetUncheckedBlocksSinceNumber(ctx context.Context, firstBlockNum uint64, dbTx pgx.Tx) ([]state.Block, error)
	UpdateCheckedBlockByNumber(ctx context.Context, blockNumber uint64, newCheckedStatus bool, dbTx pgx.Tx) error
}

// CheckL1BlockHash is a struct that implements a checker of L1Block hash
type CheckL1BlockHash struct {
	L1Client L1Requester
	State    StateInterfacer
}

func (p *CheckL1BlockHash) Step(ctx context.Context) error {
	stateBlock, err := p.State.GetFirstUncheckedBlock(ctx, nil)
	if errors.Is(err, state.ErrNotFound) {
		log.Debugf("%s: No unchecked blocks to check", logPrefix)
		return nil
	}
	if err != nil {
		return err
	}

	l1SafePointBlock, err := p.L1Client.HeaderByNumber(ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	if err != nil {
		log.Errorf("%s: Error getting L1 block %d. err: %s", logPrefix, rpc.FinalizedBlockNumber, err.Error())
		return err
	}
	return p.DoAllBlocks(ctx, stateBlock, l1SafePointBlock.Number.Uint64())
}

func (p *CheckL1BlockHash) DoAllBlocks(ctx context.Context, stateBlock *state.Block, safeBlockNumber uint64) error {
	var err error
	for {
		if stateBlock.BlockNumber > safeBlockNumber {
			log.Debugf("%s: firtst block %d to check is not still safe enough %d ", stateBlock.BlockNumber, l1SafePointBlock.Number.Uint64(), logPrefix)
			return nil
		}
		err = p.DoBlock(ctx, stateBlock)
		if err != nil {
			return err
		}
		stateBlock, err = p.State.GetFirstUncheckedBlock(ctx, nil)
		if errors.Is(err, state.ErrNotFound) {
			log.Debugf("%s: checked all blocks (safe Block: %d)", logPrefix, safeBlockNumber)
			return nil
		}

	}
}

func (p *CheckL1BlockHash) ReorgDetected(ctx context.Context, blockNumber uint64) error {
	return fmt.Errorf("Reorg detected at block %d", blockNumber)
}

func (p *CheckL1BlockHash) DoBlock(ctx context.Context, stateBlock *state.Block) error {
	if stateBlock == nil {
		log.Warn("%s: function CheckL1Block receive a nil pointer", logPrefix)
		return nil
	}
	l1Block, err := p.L1Client.HeaderByNumber(ctx, big.NewInt(int64(stateBlock.BlockNumber)))
	if err != nil {
		return err
	}
	if l1Block.Hash() != stateBlock.BlockHash {
		log.Errorf("%s: Reorg detected at block %d l1Block.Hash=%s != stateBlock.Hash=%s", logPrefix, stateBlock.BlockNumber,
			l1Block.Hash().String(), stateBlock.BlockHash.String())
		return p.ReorgDetected(ctx, stateBlock.BlockNumber)
	}
	log.Infof("%s: L1Block %d hash %s is correct marking as checked", logPrefix, stateBlock.BlockHash.String(), stateBlock.BlockNumber)
	err = p.State.UpdateCheckedBlockByNumber(ctx, stateBlock.BlockNumber, true, nil)
	if err != nil {
		log.Errorf("%s: Error updating block %d as checked. err: %s", logPrefix, stateBlock.BlockNumber, err.Error())
		return err
	}
}
