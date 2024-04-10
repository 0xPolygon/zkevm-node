package l1_check_block

import (
	"context"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
)

// L1BlockChecker is an interface that defines the method to check L1 blocks
type L1BlockChecker interface {
	Step(ctx context.Context) error
}

const (
	defaultPeriodTime = time.Second
)

// AsyncCheck is a wrapper for L1BlockChecker to become asynchronous
type AsyncCheck struct {
	checker     L1BlockChecker
	mutex       sync.Mutex
	lastResult  *syncinterfaces.IterationResult
	onFnishCall func()
	periodTime  time.Duration
}

// NewAsyncCheck creates a new AsyncCheck
func NewAsyncCheck(checker L1BlockChecker) *AsyncCheck {
	return &AsyncCheck{
		checker:    checker,
		periodTime: defaultPeriodTime,
	}
}

// NewAsyncCheckWithPeriodTime creates a new AsyncCheck with a period time between relaunch checker.Step
func NewAsyncCheckWithPeriodTime(checker L1BlockChecker, periodTime time.Duration) *AsyncCheck {
	return &AsyncCheck{
		checker:    checker,
		periodTime: periodTime,
	}
}

// Run is a method that starts the async check
func (a *AsyncCheck) Run(ctx context.Context, onFinish func()) {
	a.lastResult = nil
	a.onFnishCall = onFinish
	a.launchChecker(ctx)
}

// RunSynchronous is a method that forces the check to be synchronous before starting the async check
func (a *AsyncCheck) RunSynchronous(ctx context.Context) syncinterfaces.IterationResult {
	return a.executeIteration(ctx)
}

// GetResponse returns the last result of the check:
// - Nil -> still running
// - Not nil -> finished, and this is the result. You must call again Run to start a new check
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

// step is a method that executes a until executeItertion
// returns an error or a reorg
func (a *AsyncCheck) step(ctx context.Context) *syncinterfaces.IterationResult {
	select {
	case <-ctx.Done():
		log.Debugf("%s L1BlockChecker: context done", logPrefix)
		return &syncinterfaces.IterationResult{Err: ctx.Err()}
	default:
		result := a.executeIteration(ctx)
		if result.ReorgDetected {
			return &result
		}
		log.Debugf("%s L1BlockChecker:returned %s waiting %s to relaunch", logPrefix, result.String(), a.periodTime)
		time.Sleep(a.periodTime)
	}
	return nil
}

// executeIteration executes a single iteration of the checker
func (a *AsyncCheck) executeIteration(ctx context.Context) syncinterfaces.IterationResult {
	res := syncinterfaces.IterationResult{}
	log.Debugf("%s calling checker.Step(...)", logPrefix)
	res.Err = a.checker.Step(ctx)
	log.Debugf("%s returned checker.Step(...) %w", logPrefix, res.Err)
	if res.Err != nil {
		log.Errorf("%s Fail check L1 Blocks: %w", logPrefix, res.Err)
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
