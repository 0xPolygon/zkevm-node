package l1_check_block

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
)

type L1BlockCheckerIntegration struct {
	forceCheckOnStart bool
	checker           syncinterfaces.AsyncL1BlockChecker
	sync              SyncCheckReorger
}

type SyncCheckReorger interface {
	ExecuteReorg(blockNumber uint64, reason string) error
	OnDetectedMismatchL1BlockReorg()
}

func NewL1BlockCheckerIntegration(checker syncinterfaces.AsyncL1BlockChecker, sync SyncCheckReorger, forceCheckOnStart bool) *L1BlockCheckerIntegration {
	return &L1BlockCheckerIntegration{
		forceCheckOnStart: forceCheckOnStart,
		checker:           checker,
		sync:              sync,
	}
}

func (v *L1BlockCheckerIntegration) OnStart(ctx context.Context) error {
	if v.forceCheckOnStart {
		log.Infof("%s Forcing L1BlockChecker check before start", logPrefix)
		var result syncinterfaces.IterationResult
		for {
			result = v.checker.RunSynchronous(ctx)
			if result.Err == nil {
				break
			} else {
				time.Sleep(1 * time.Second)
			}
		}
		if result.ReorgDetected {
			v.executeResult(ctx, result)
		}
	}
	v.launch(ctx)
	return nil
}

func (v *L1BlockCheckerIntegration) OnStartL1Sync(ctx context.Context) bool {
	return v.checkBackgroundResult(ctx, "before start L1 sync")

}

func (v *L1BlockCheckerIntegration) OnStartL2Sync(ctx context.Context) bool {
	return v.checkBackgroundResult(ctx, "before start 2 sync")
}

func (v *L1BlockCheckerIntegration) OnCheckReorg(ctx context.Context, latestBlock *state.Block) bool {
	return v.checkBackgroundResult(ctx, "OnCheckReorg")
}

func (v *L1BlockCheckerIntegration) checkBackgroundResult(ctx context.Context, positionMessage string) bool {
	log.Debugf("%s Checking L1BlockChecker %s", logPrefix, positionMessage)
	result := v.checker.GetResponse()
	if result != nil {
		if result.ReorgDetected {
			log.Warnf("%s Checking L1BlockChecker %s: reorg detected %s", logPrefix, positionMessage, result.String())
			v.executeResult(ctx, *result)
		}
		v.launch(ctx)
		return result.ReorgDetected
	}
	return false
}

func (v *L1BlockCheckerIntegration) launch(ctx context.Context) {
	log.Infof("%s L1BlockChecker: starting background process...", logPrefix)
	v.checker.Run(ctx, func() {
		log.Infof("%s L1BlockChecker: finished background process, calling to synchronizer", logPrefix)
		v.sync.OnDetectedMismatchL1BlockReorg()
	})
}

func (v *L1BlockCheckerIntegration) executeResult(ctx context.Context, result syncinterfaces.IterationResult) bool {
	if result.ReorgDetected {
		for {
			err := v.sync.ExecuteReorg(result.BlockNumber, result.ReorgMessage)
			if err == nil {
				return true
			}
			log.Errorf("%s Error executing reorg: %s", logPrefix, err)
			time.Sleep(1 * time.Second)
		}
	}
	return false
}
