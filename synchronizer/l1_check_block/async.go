package l1_check_block

import (
	"context"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
)

type L1BlockChecker interface {
	Step(ctx context.Context) error
}

type AsyncCheck struct {
	checker     L1BlockChecker
	mutex       sync.Mutex
	lastResult  *syncinterfaces.IterationResult
	onFnishCall func()
}

func NewAsyncCheck(checker L1BlockChecker) *AsyncCheck {
	return &AsyncCheck{
		checker: checker,
	}
}

func (a *AsyncCheck) Run(ctx context.Context, onFinish func()) {
	a.lastResult = nil
	a.onFnishCall = onFinish
	a.launchChecker(ctx)
}

// RunSynchronous is a method that forces the check to be synchronous before starting the async check
func (a *AsyncCheck) RunSynchronous(ctx context.Context) syncinterfaces.IterationResult {
	return a.executeIteration(ctx)
}

func (a *AsyncCheck) GetResponse() *syncinterfaces.IterationResult {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.lastResult
}

func (a *AsyncCheck) setResult(result syncinterfaces.IterationResult) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.lastResult = &result
}

func (a *AsyncCheck) launchChecker(ctx context.Context) {
	go func() {
		log.Infof("%s L1BlockChecker: starting background process", logPrefix)
		for {
			result := a.step(ctx)
			if result != nil {
				a.setResult(*result)
				break
			}
		}
		log.Infof("%s L1BlockChecker: finished background process", logPrefix)
		if a.onFnishCall != nil {
			a.onFnishCall()
		}
	}()
}

func (a *AsyncCheck) step(ctx context.Context) *syncinterfaces.IterationResult {
	select {
	case <-ctx.Done():
		return &syncinterfaces.IterationResult{Err: ctx.Err()}
	case <-time.After(1 * time.Second):
		result := a.executeIteration(ctx)
		if result.ReorgDetected {
			return &result
		}
	}
	return nil
}

// executeIteration executes a single iteration of the checker
func (a *AsyncCheck) executeIteration(ctx context.Context) syncinterfaces.IterationResult {
	res := syncinterfaces.IterationResult{}
	res.Err = a.checker.Step(ctx)
	if res.Err != nil {
		log.Errorf("%s Fail check L1 Blocks: %s", logPrefix, res.Err)
		if common.IsReorgError(res.Err) {
			// log error
			blockNumber := common.GetReorgErrorBlockNumber(res.Err)
			log.Infof("%s Reorg detected at block %d", logPrefix, blockNumber)
			// It keeps blocked until the channel is read
			res.BlockNumber = blockNumber
			res.ReorgDetected = true
			res.ReorgMessage = res.Err.Error()
			res.Err = nil
		}
	}
	return res
}
