package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4"
)

const (
	// FilterTypeLog represents a filter of type log.
	FilterTypeLog = "log"
	// FilterTypeBlock represents a filter of type block.
	FilterTypeBlock = "block"
	// FilterTypePendingTx represent a filter of type pending Tx.
	FilterTypePendingTx = "pendingTx"
)

// Filter represents a filter.
type Filter struct {
	ID         string
	Type       FilterType
	Parameters interface{}
	LastPoll   time.Time
	WsConn     *concurrentWsConn

	wsQueue       *state.Queue[[]byte]
	wsQueueSignal *sync.Cond
}

// EnqueueSubscriptionDataToBeSent enqueues subscription data to be sent
// via web sockets connection
func (f *Filter) EnqueueSubscriptionDataToBeSent(data []byte) {
	f.wsQueue.Push(data)
	f.wsQueueSignal.Broadcast()
}

// SendEnqueuedSubscriptionData consumes all the enqueued subscription data
// and sends it via web sockets connection.
func (f *Filter) SendEnqueuedSubscriptionData() {
	for {
		// wait for a signal that a new item was
		// added to the queue
		log.Debugf("waiting subscription data signal")
		f.wsQueueSignal.L.Lock()
		f.wsQueueSignal.Wait()
		f.wsQueueSignal.L.Unlock()
		log.Debugf("subscription data signal received, sending enqueued data")
		for {
			d, err := f.wsQueue.Pop()
			if err == state.ErrQueueEmpty {
				break
			} else if err != nil {
				log.Errorf("failed to pop subscription data from queue to be sent via web sockets to filter %v, %s", f.ID, err.Error())
				break
			}
			f.sendSubscriptionResponse(d)
		}
	}
}

// sendSubscriptionResponse send data as subscription response via
// web sockets connection controlled by a mutex
func (f *Filter) sendSubscriptionResponse(data []byte) {
	const errMessage = "Unable to write WS message to filter %v, %s"

	start := time.Now()
	res := types.SubscriptionResponse{
		JSONRPC: "2.0",
		Method:  "eth_subscription",
		Params: types.SubscriptionResponseParams{
			Subscription: f.ID,
			Result:       data,
		},
	}
	message, err := json.Marshal(res)
	if err != nil {
		log.Errorf(fmt.Sprintf(errMessage, f.ID, err.Error()))
		return
	}

	err = f.WsConn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Errorf(fmt.Sprintf(errMessage, f.ID, err.Error()))
		return
	}
	log.Debugf("WS message sent: %v", string(message))
	log.Debugf("[SendSubscriptionResponse] took %v", time.Since(start))
}

// FilterType express the type of the filter, block, logs, pending transactions
type FilterType string

// LogFilter is a filter for logs
type LogFilter struct {
	BlockHash *common.Hash
	FromBlock *types.BlockNumber
	ToBlock   *types.BlockNumber
	Addresses []common.Address
	Topics    [][]common.Hash
	Since     *time.Time
}

// addTopic adds specific topics to the log filter topics
func (f *LogFilter) addTopic(topics ...string) error {
	if f.Topics == nil {
		f.Topics = [][]common.Hash{}
	}

	topicsHashes := []common.Hash{}

	for _, topic := range topics {
		topicHash := common.Hash{}
		if err := topicHash.UnmarshalText([]byte(topic)); err != nil {
			return err
		}

		topicsHashes = append(topicsHashes, topicHash)
	}

	f.Topics = append(f.Topics, topicsHashes)

	return nil
}

// addAddress Adds the address to the log filter
func (f *LogFilter) addAddress(raw string) error {
	if f.Addresses == nil {
		f.Addresses = []common.Address{}
	}

	addr := common.Address{}

	if err := addr.UnmarshalText([]byte(raw)); err != nil {
		return err
	}

	f.Addresses = append(f.Addresses, addr)

	return nil
}

// MarshalJSON allows to customize the JSON representation.
func (f *LogFilter) MarshalJSON() ([]byte, error) {
	var obj types.LogFilterRequest

	obj.BlockHash = f.BlockHash

	if f.FromBlock != nil && (*f.FromBlock == types.LatestBlockNumber) {
		fromBlock := ""
		obj.FromBlock = &fromBlock
	} else if f.FromBlock != nil {
		fromBlock := hex.EncodeUint64(uint64(*f.FromBlock))
		obj.FromBlock = &fromBlock
	}

	if f.ToBlock != nil && (*f.ToBlock == types.LatestBlockNumber) {
		toBlock := ""
		obj.ToBlock = &toBlock
	} else if f.ToBlock != nil {
		toBlock := hex.EncodeUint64(uint64(*f.ToBlock))
		obj.ToBlock = &toBlock
	}

	if f.Addresses != nil {
		if len(f.Addresses) == 1 {
			obj.Address = f.Addresses[0].Hex()
		} else {
			obj.Address = f.Addresses
		}
	}

	obj.Topics = make([]interface{}, 0, len(f.Topics))
	for _, topic := range f.Topics {
		if len(topic) == 0 {
			obj.Topics = append(obj.Topics, nil)
		} else if len(topic) == 1 {
			obj.Topics = append(obj.Topics, topic[0])
		} else {
			obj.Topics = append(obj.Topics, topic)
		}
	}

	return json.Marshal(obj)
}

// UnmarshalJSON decodes a json object
func (f *LogFilter) UnmarshalJSON(data []byte) error {
	var obj types.LogFilterRequest

	err := json.Unmarshal(data, &obj)

	if err != nil {
		return err
	}

	f.BlockHash = obj.BlockHash
	lbb := types.LatestBlockNumber

	if obj.FromBlock != nil && *obj.FromBlock == "" {
		f.FromBlock = &lbb
	} else if obj.FromBlock != nil {
		bn, err := types.StringToBlockNumber(*obj.FromBlock)
		if err != nil {
			return err
		}
		f.FromBlock = &bn
	}

	if obj.ToBlock != nil && *obj.ToBlock == "" {
		f.ToBlock = &lbb
	} else if obj.ToBlock != nil {
		bn, err := types.StringToBlockNumber(*obj.ToBlock)
		if err != nil {
			return err
		}
		f.ToBlock = &bn
	}

	if obj.Address != nil {
		// decode address, either "" or [""]
		switch raw := obj.Address.(type) {
		case string:
			// ""
			if err := f.addAddress(raw); err != nil {
				return err
			}

		case []interface{}:
			// ["", ""]
			for _, addr := range raw {
				if item, ok := addr.(string); ok {
					if err := f.addAddress(item); err != nil {
						return err
					}
				} else {
					return errors.New("address expected")
				}
			}

		default:
			return errors.New("failed to decode address. Expected either '' or ['', '']")
		}
	}

	if obj.Topics != nil {
		// decode topics, either "" or ["", ""] or null
		for _, item := range obj.Topics {
			switch raw := item.(type) {
			case string:
				// ""
				if err := f.addTopic(raw); err != nil {
					return err
				}

			case []interface{}:
				// ["", ""]
				res := []string{}

				for _, i := range raw {
					if item, ok := i.(string); ok {
						res = append(res, item)
					} else {
						return errors.New("hash expected")
					}
				}

				if err := f.addTopic(res...); err != nil {
					return err
				}

			case nil:
				// null
				if err := f.addTopic(); err != nil {
					return err
				}

			default:
				return errors.New("failed to decode topics. Expected '' or [''] or null")
			}
		}
	}

	// decode topics
	return nil
}

// Match returns whether the receipt includes topics for this filter
func (f *LogFilter) Match(log *types.Log) bool {
	// check addresses
	if len(f.Addresses) > 0 {
		match := false

		for _, addr := range f.Addresses {
			if addr == log.Address {
				match = true
			}
		}

		if !match {
			return false
		}
	}
	// check topics
	if len(f.Topics) > len(log.Topics) {
		return false
	}

	for i, sub := range f.Topics {
		match := len(sub) == 0

		for _, topic := range sub {
			if log.Topics[i] == topic {
				match = true

				break
			}
		}

		if !match {
			return false
		}
	}

	return true
}

// GetNumericBlockNumbers load the numeric block numbers from state accordingly
// to the provided from and to block number
func (f *LogFilter) GetNumericBlockNumbers(ctx context.Context, cfg Config, s types.StateInterface, e types.EthermanInterface, dbTx pgx.Tx) (uint64, uint64, types.Error) {
	return getNumericBlockNumbers(ctx, s, e, f.FromBlock, f.ToBlock, cfg.MaxLogsBlockRange, state.ErrMaxLogsBlockRangeLimitExceeded, dbTx)
}

// ShouldFilterByBlockHash if the filter should consider the block hash value
func (f *LogFilter) ShouldFilterByBlockHash() bool {
	return f.BlockHash != nil
}

// ShouldFilterByBlockRange if the filter should consider the block range values
func (f *LogFilter) ShouldFilterByBlockRange() bool {
	return f.FromBlock != nil || f.ToBlock != nil
}

// Validate check if the filter instance is valid
func (f *LogFilter) Validate() error {
	if f.ShouldFilterByBlockHash() && f.ShouldFilterByBlockRange() {
		return ErrFilterInvalidPayload
	}
	return nil
}

// NativeBlockHashBlockRangeFilter is a filter to filter native block hash by block by number
type NativeBlockHashBlockRangeFilter struct {
	FromBlock types.BlockNumber `json:"fromBlock"`
	ToBlock   types.BlockNumber `json:"toBlock"`
}

// GetNumericBlockNumbers load the numeric block numbers from state accordingly
// to the provided from and to block number
func (f *NativeBlockHashBlockRangeFilter) GetNumericBlockNumbers(ctx context.Context, cfg Config, s types.StateInterface, e types.EthermanInterface, dbTx pgx.Tx) (uint64, uint64, types.Error) {
	return getNumericBlockNumbers(ctx, s, e, &f.FromBlock, &f.ToBlock, cfg.MaxNativeBlockHashBlockRange, state.ErrMaxNativeBlockHashBlockRangeLimitExceeded, dbTx)
}

// getNumericBlockNumbers load the numeric block numbers from state accordingly
// to the provided from and to block number
func getNumericBlockNumbers(ctx context.Context, s types.StateInterface, e types.EthermanInterface, fromBlock, toBlock *types.BlockNumber, maxBlockRange uint64, maxBlockRangeErr error, dbTx pgx.Tx) (uint64, uint64, types.Error) {
	var fromBlockNumber uint64 = 0
	if fromBlock != nil {
		fbn, rpcErr := fromBlock.GetNumericBlockNumber(ctx, s, e, dbTx)
		if rpcErr != nil {
			return 0, 0, rpcErr
		}
		fromBlockNumber = fbn
	}

	toBlockNumber, rpcErr := toBlock.GetNumericBlockNumber(ctx, s, e, dbTx)
	if rpcErr != nil {
		return 0, 0, rpcErr
	}

	if toBlockNumber < fromBlockNumber {
		_, rpcErr := RPCErrorResponse(types.InvalidParamsErrorCode, state.ErrInvalidBlockRange.Error(), nil, false)
		return 0, 0, rpcErr
	}

	blockRange := toBlockNumber - fromBlockNumber
	if maxBlockRange > 0 && blockRange > maxBlockRange {
		errMsg := fmt.Sprintf(maxBlockRangeErr.Error(), maxBlockRange)
		_, rpcErr := RPCErrorResponse(types.InvalidParamsErrorCode, errMsg, nil, false)
		return 0, 0, rpcErr
	}

	return fromBlockNumber, toBlockNumber, nil
}
