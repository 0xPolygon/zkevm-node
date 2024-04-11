package l1_check_block_test

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_check_block"
)

func TestPreCheckL1BlockStart(t *testing.T) {
	l1_check_block.NewPreCheckL1BlockHash(nil, nil, nil, nil)
}
