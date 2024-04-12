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
	mockChecker    *mock_syncinterfaces.AsyncL1BlockChecker
	mockPreChecker *mock_syncinterfaces.AsyncL1BlockChecker
	mockSync       *mock_l1_check_block.SyncCheckReorger
	sut            *l1_check_block.L1BlockCheckerIntegration
	ctx            context.Context
	resultOk       syncinterfaces.IterationResult
	resultError    syncinterfaces.IterationResult
	resultReorg    syncinterfaces.IterationResult
}

func newDataIntegration(t *testing.T, forceCheckOnStart bool) *testDataIntegration {
	return newDataIntegrationPreChecker(t, forceCheckOnStart, nil)
}

func newDataIntegrationWithPreChecker(t *testing.T, forceCheckOnStart bool) *testDataIntegration {
	return newDataIntegrationPreChecker(t, forceCheckOnStart, mock_syncinterfaces.NewAsyncL1BlockChecker(t))
}

func newDataIntegrationPreChecker(t *testing.T, forceCheckOnStart bool, mockPreChecker *mock_syncinterfaces.AsyncL1BlockChecker) *testDataIntegration {
	mockChecker := mock_syncinterfaces.NewAsyncL1BlockChecker(t)
	mockSync := mock_l1_check_block.NewSyncCheckReorger(t)
	sut := l1_check_block.NewL1BlockCheckerIntegration(mockChecker, mockPreChecker, mockSync, forceCheckOnStart, time.Millisecond)
	return &testDataIntegration{
		mockChecker:    mockChecker,
		mockPreChecker: mockPreChecker,
		mockSync:       mockSync,
		sut:            sut,
		ctx:            context.Background(),
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

// OnStart if check and preCheck execute both, and launch both in background
func TestIntegrationCheckAndPreCheckOnStartForceCheck(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultOk)
	data.mockPreChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultOk)
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	data.mockPreChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

// OnStart if mainChecker returns reorg doesnt need to run preCheck
func TestIntegrationCheckAndPreCheckOnStartMainCheckerReturnReorg(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultReorg)
	data.mockSync.EXPECT().ExecuteReorg(uint64(1234), mock.Anything).Return(nil).Once()
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	data.mockPreChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

// If mainCheck is OK, but preCheck returns reorg, it should execute reorg
func TestIntegrationCheckAndPreCheckOnStartPreCheckerReturnReorg(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultOk)
	data.mockPreChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultReorg)
	data.mockSync.EXPECT().ExecuteReorg(uint64(1234), mock.Anything).Return(nil).Once()
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	data.mockPreChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

// The process is running on background, no results yet
func TestIntegrationCheckAndPreCheckOnOnCheckReorgRunningOnBackground(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().GetResult().Return(nil)
	data.mockPreChecker.EXPECT().GetResult().Return(nil)
	reorgExecuted := data.sut.OnCheckReorg(data.ctx, nil)
	require.False(t, reorgExecuted)
}

func TestIntegrationCheckAndPreCheckOnOnCheckReorgOneProcessHaveResultOK(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().GetResult().Return(&data.resultOk)
	data.mockPreChecker.EXPECT().GetResult().Return(nil)
	// One have been stopped, so must relaunch both
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	data.mockPreChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	reorgExecuted := data.sut.OnCheckReorg(data.ctx, nil)
	require.False(t, reorgExecuted)
}

func TestIntegrationCheckAndPreCheckOnOnCheckReorgMainCheckerReorg(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().GetResult().Return(&data.resultReorg)
	data.mockPreChecker.EXPECT().GetResult().Return(nil)
	data.mockSync.EXPECT().ExecuteReorg(uint64(1234), mock.Anything).Return(nil).Once()
	// One have been stopped, so must relaunch both
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	data.mockPreChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	reorgExecuted := data.sut.OnCheckReorg(data.ctx, nil)
	require.True(t, reorgExecuted)
}

func TestIntegrationCheckAndPreCheckOnOnCheckReorgPreCheckerReorg(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().GetResult().Return(nil)
	data.mockPreChecker.EXPECT().GetResult().Return(&data.resultReorg)
	data.mockSync.EXPECT().ExecuteReorg(uint64(1234), mock.Anything).Return(nil).Once()
	// One have been stopped, so must relaunch both
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	data.mockPreChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	reorgExecuted := data.sut.OnCheckReorg(data.ctx, nil)
	require.True(t, reorgExecuted)
}

func TestIntegrationCheckAndPreCheckOnOnCheckReorgBothReorgWinOldest1(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	reorgMain := data.resultReorg
	reorgMain.BlockNumber = 1235
	data.mockChecker.EXPECT().GetResult().Return(&reorgMain)
	reorgPre := data.resultReorg
	reorgPre.BlockNumber = 1236
	data.mockPreChecker.EXPECT().GetResult().Return(&reorgPre)
	data.mockSync.EXPECT().ExecuteReorg(uint64(1235), mock.Anything).Return(nil).Once()
	// One have been stopped, so must relaunch both
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	data.mockPreChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	reorgExecuted := data.sut.OnCheckReorg(data.ctx, nil)
	require.True(t, reorgExecuted)
}

func TestIntegrationCheckAndPreCheckOnOnCheckReorgBothReorgWinOldest2(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	reorgMain := data.resultReorg
	reorgMain.BlockNumber = 1236
	data.mockChecker.EXPECT().GetResult().Return(&reorgMain)
	reorgPre := data.resultReorg
	reorgPre.BlockNumber = 1235
	data.mockPreChecker.EXPECT().GetResult().Return(&reorgPre)
	data.mockSync.EXPECT().ExecuteReorg(uint64(1235), mock.Anything).Return(nil).Once()
	// One have been stopped, so must relaunch both
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	data.mockPreChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	reorgExecuted := data.sut.OnCheckReorg(data.ctx, nil)
	require.True(t, reorgExecuted)
}
