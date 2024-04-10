package l1_check_block_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_check_block"
	mock_l1_check_block "github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_check_block/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	genericErrorToTest = fmt.Errorf("error")
)

type testDataIntegration struct {
	mockChecker *mock_syncinterfaces.AsyncL1BlockChecker
	mockSync    *mock_l1_check_block.SyncCheckReorger
	sut         *l1_check_block.L1BlockCheckerIntegration
	ctx         context.Context
	resultOk    syncinterfaces.IterationResult
	resultError syncinterfaces.IterationResult
	resultReorg syncinterfaces.IterationResult
}

func newDataIntegration(t *testing.T, forceCheckOnStart bool) *testDataIntegration {
	mockChecker := mock_syncinterfaces.NewAsyncL1BlockChecker(t)
	mockSync := mock_l1_check_block.NewSyncCheckReorger(t)
	sut := l1_check_block.NewL1BlockCheckerIntegration(mockChecker, mockSync, forceCheckOnStart, time.Millisecond)
	return &testDataIntegration{
		mockChecker: mockChecker,
		mockSync:    mockSync,
		sut:         sut,
		ctx:         context.Background(),
		resultReorg: syncinterfaces.IterationResult{
			ReorgDetected: true,
			BlockNumber:   1234,
		},
		resultOk: syncinterfaces.IterationResult{
			ReorgDetected: false,
		},
		resultError: syncinterfaces.IterationResult{
			Err:           genericErrorToTest,
			ReorgDetected: false,
		},
	}
}

func TestIntegrationIfNoForceCheckOnlyLaunchBackgroudChecker(t *testing.T) {
	data := newDataIntegration(t, false)
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

func TestIntegrationIfForceCheckRunsSynchronousOneTimeAndAfterLaunchBackgroudChecker(t *testing.T) {
	data := newDataIntegration(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultOk)
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

func TestIntegrationIfSyncCheckReturnsReorgExecuteIt(t *testing.T) {
	data := newDataIntegration(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultReorg)
	data.mockSync.EXPECT().ExecuteReorg(uint64(1234), "").Return(nil)
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

func TestIntegrationIfSyncCheckReturnErrorRetry(t *testing.T) {
	data := newDataIntegration(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultError).Once()
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultOk).Once()
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

func TestIntegrationIfSyncCheckReturnsReorgExecuteItAndFailsRetry(t *testing.T) {
	data := newDataIntegration(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultReorg)
	data.mockSync.EXPECT().ExecuteReorg(uint64(1234), mock.Anything).Return(genericErrorToTest).Once()
	data.mockSync.EXPECT().ExecuteReorg(uint64(1234), mock.Anything).Return(nil).Once()
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}
