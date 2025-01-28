package sequencesender_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	coretype "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"

	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/sequencesender"
)

func TestBuildSequenceBatchesTxData(t *testing.T) {
	mockEtherman := sequencesender.NewEthermanMock(t)

	sender := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	sequences := []types.Sequence{}
	maxSequenceTimestamp := uint64(1234567890)
	initSequenceBatchNumber := uint64(1)
	l2Coinbase := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdef")

	mockEtherman.On("BuildSequenceBatchesTxData", sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase).
		Return(&sender, []byte{0x01, 0x02, 0x03}, nil)

	addr, data, err := mockEtherman.BuildSequenceBatchesTxData(sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)

	assert.NoError(t, err)
	assert.Equal(t, &sender, addr)
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, data)

	mockEtherman.AssertCalled(t, "BuildSequenceBatchesTxData", sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)
}

func TestEstimateGasSequenceBatches(t *testing.T) {
	mockEtherman := sequencesender.NewEthermanMock(t)

	sender := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	sequences := []types.Sequence{}
	maxSequenceTimestamp := uint64(1234567890)
	initSequenceBatchNumber := uint64(1)
	l2Coinbase := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdef")
	tx := &coretype.Transaction{}

	mockEtherman.On("EstimateGasSequenceBatches", sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase).
		Return(tx, nil)

	result, err := mockEtherman.EstimateGasSequenceBatches(sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)

	assert.NoError(t, err)
	assert.Equal(t, tx, result)

	mockEtherman.AssertCalled(t, "EstimateGasSequenceBatches", sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)
}

func TestGetLatestBatchNumber(t *testing.T) {
	mockEtherman := sequencesender.NewEthermanMock(t)

	mockEtherman.On("GetLatestBatchNumber").Return(uint64(42), nil)

	batchNumber, err := mockEtherman.GetLatestBatchNumber()

	assert.NoError(t, err)
	assert.Equal(t, uint64(42), batchNumber)

	mockEtherman.AssertCalled(t, "GetLatestBatchNumber")
}

func TestGetLatestBlockHeader(t *testing.T) {
	mockEtherman := sequencesender.NewEthermanMock(t)

	ctx := context.Background()
	header := &coretype.Header{}

	mockEtherman.On("GetLatestBlockHeader", ctx).Return(header, nil)

	result, err := mockEtherman.GetLatestBlockHeader(ctx)

	assert.NoError(t, err)
	assert.Equal(t, header, result)

	mockEtherman.AssertCalled(t, "GetLatestBlockHeader", ctx)
}
