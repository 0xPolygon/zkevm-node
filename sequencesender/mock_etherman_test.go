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
	// 创建 Mock 实例
	mockEtherman := sequencesender.NewEthermanMock(t)

	// 测试数据
	sender := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	sequences := []types.Sequence{}
	maxSequenceTimestamp := uint64(1234567890)
	initSequenceBatchNumber := uint64(1)
	l2Coinbase := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdef")

	// 设置期望值
	mockEtherman.On("BuildSequenceBatchesTxData", sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase).
		Return(&sender, []byte{0x01, 0x02, 0x03}, nil)

	// 调用被测函数
	addr, data, err := mockEtherman.BuildSequenceBatchesTxData(sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)

	// 验证返回值
	assert.NoError(t, err)
	assert.Equal(t, &sender, addr)
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, data)

	// 验证 Mock 调用
	mockEtherman.AssertCalled(t, "BuildSequenceBatchesTxData", sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)
}

func TestEstimateGasSequenceBatches(t *testing.T) {
	// 创建 Mock 实例
	mockEtherman := sequencesender.NewEthermanMock(t)

	// 测试数据
	sender := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	sequences := []types.Sequence{}
	maxSequenceTimestamp := uint64(1234567890)
	initSequenceBatchNumber := uint64(1)
	l2Coinbase := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdef")
	tx := &coretype.Transaction{}

	// 设置期望值
	mockEtherman.On("EstimateGasSequenceBatches", sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase).
		Return(tx, nil)

	// 调用被测函数
	result, err := mockEtherman.EstimateGasSequenceBatches(sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)

	// 验证返回值
	assert.NoError(t, err)
	assert.Equal(t, tx, result)

	// 验证 Mock 调用
	mockEtherman.AssertCalled(t, "EstimateGasSequenceBatches", sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)
}

func TestGetLatestBatchNumber(t *testing.T) {
	// 创建 Mock 实例
	mockEtherman := sequencesender.NewEthermanMock(t)

	// 设置期望值
	mockEtherman.On("GetLatestBatchNumber").Return(uint64(42), nil)

	// 调用被测函数
	batchNumber, err := mockEtherman.GetLatestBatchNumber()

	// 验证返回值
	assert.NoError(t, err)
	assert.Equal(t, uint64(42), batchNumber)

	// 验证 Mock 调用
	mockEtherman.AssertCalled(t, "GetLatestBatchNumber")
}

func TestGetLatestBlockHeader(t *testing.T) {
	// 创建 Mock 实例
	mockEtherman := sequencesender.NewEthermanMock(t)

	// 测试数据
	ctx := context.Background()
	header := &coretype.Header{}

	// 设置期望值
	mockEtherman.On("GetLatestBlockHeader", ctx).Return(header, nil)

	// 调用被测函数
	result, err := mockEtherman.GetLatestBlockHeader(ctx)

	// 验证返回值
	assert.NoError(t, err)
	assert.Equal(t, header, result)

	// 验证 Mock 调用
	mockEtherman.AssertCalled(t, "GetLatestBlockHeader", ctx)
}
