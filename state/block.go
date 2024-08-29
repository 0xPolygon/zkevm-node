package state

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Block struct
type Block struct {
	BlockNumber uint64
	BlockHash   common.Hash
	ParentHash  common.Hash
	ReceivedAt  time.Time
	Checked     bool
}

func (b *Block) String() string {
	return fmt.Sprintf("BlockNumber: %d, BlockHash: %s, ParentHash: %s, ReceivedAt: %s", b.BlockNumber, b.BlockHash, b.ParentHash, b.ReceivedAt)
}

// NewBlock creates a block with the given data.
func NewBlock(blockNumber uint64) *Block {
	return &Block{BlockNumber: blockNumber}
}
