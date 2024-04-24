package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

const (
	// EarliestBlockNumber represents the earliest block number, always 0
	EarliestBlockNumber = BlockNumber(-1)
	// LatestBlockNumber represents the latest block number
	LatestBlockNumber = BlockNumber(-2)
	// PendingBlockNumber represents the pending block number
	PendingBlockNumber = BlockNumber(-3)
	// SafeBlockNumber represents the last verified block number that is safe on Ethereum
	SafeBlockNumber = BlockNumber(-4)
	// FinalizedBlockNumber represents the last verified block number that is finalized on Ethereum
	FinalizedBlockNumber = BlockNumber(-5)

	// EarliestBatchNumber represents the earliest batch number, always 0
	EarliestBatchNumber = BatchNumber(-1)
	// LatestBatchNumber represents the last closed batch number
	LatestBatchNumber = BatchNumber(-2)
	// PendingBatchNumber represents the last batch in the trusted state
	PendingBatchNumber = BatchNumber(-3)
	// SafeBatchNumber represents the last batch verified in a block that is safe on Ethereum
	SafeBatchNumber = BatchNumber(-4)
	// FinalizedBatchNumber represents the last batch verified in a block that has been finalized on Ethereum
	FinalizedBatchNumber = BatchNumber(-5)

	// Earliest contains the string to represent the earliest block known.
	Earliest = "earliest"
	// Latest contains the string to represent the latest block known.
	Latest = "latest"
	// Pending contains the string to represent the pending block known.
	Pending = "pending"
	// Safe contains the string to represent the last virtualized block known.
	Safe = "safe"
	// Finalized contains the string to represent the last verified block known.
	Finalized = "finalized"

	// EIP-1898: https://eips.ethereum.org/EIPS/eip-1898 //

	// BlockNumberKey is the key for the block number for EIP-1898
	BlockNumberKey = "blockNumber"
	// BlockHashKey is the key for the block hash for EIP-1898
	BlockHashKey = "blockHash"
	// RequireCanonicalKey is the key for the require canonical for EIP-1898
	RequireCanonicalKey = "requireCanonical"
)

// Request is a jsonrpc request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response is a jsonrpc  success response
type Response struct {
	JSONRPC string
	ID      interface{}
	Result  json.RawMessage
	Error   *ErrorObject
}

// ErrorObject is a jsonrpc error
type ErrorObject struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    *ArgBytes `json:"data,omitempty"`
}

// RPCError returns an instance of RPCError from the
// data available in the ErrorObject instance
func (e *ErrorObject) RPCError() RPCError {
	var data []byte
	if e.Data != nil {
		data = *e.Data
	}
	rpcError := NewRPCErrorWithData(e.Code, e.Message, data)
	return *rpcError
}

// NewResponse returns Success/Error response object
func NewResponse(req Request, reply []byte, err Error) Response {
	var result json.RawMessage
	if reply != nil {
		result = reply
	}

	var errorObj *ErrorObject
	if err != nil {
		errorObj = &ErrorObject{
			Code:    err.ErrorCode(),
			Message: err.Error(),
		}
		if err.ErrorData() != nil {
			errorObj.Data = ArgBytesPtr(err.ErrorData())
		}
	}

	return Response{
		JSONRPC: req.JSONRPC,
		ID:      req.ID,
		Result:  result,
		Error:   errorObj,
	}
}

// MarshalJSON customizes the JSON representation of the response.
func (r Response) MarshalJSON() ([]byte, error) {
	if r.Error != nil {
		return json.Marshal(struct {
			JSONRPC string       `json:"jsonrpc"`
			ID      interface{}  `json:"id"`
			Error   *ErrorObject `json:"error"`
		}{
			JSONRPC: r.JSONRPC,
			ID:      r.ID,
			Error:   r.Error,
		})
	}

	return json.Marshal(struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      interface{}     `json:"id"`
		Result  json.RawMessage `json:"result"`
	}{
		JSONRPC: r.JSONRPC,
		ID:      r.ID,
		Result:  r.Result,
	})
}

// Bytes return the serialized response
func (s Response) Bytes() ([]byte, error) {
	return json.Marshal(s)
}

// SubscriptionResponse used to push response for filters
// that have an active web socket connection
type SubscriptionResponse struct {
	JSONRPC string                     `json:"jsonrpc"`
	Method  string                     `json:"method"`
	Params  SubscriptionResponseParams `json:"params"`
}

// SubscriptionResponseParams parameters for subscription responses
type SubscriptionResponseParams struct {
	Subscription string          `json:"subscription"`
	Result       json.RawMessage `json:"result"`
}

// Bytes return the serialized response
func (s SubscriptionResponse) Bytes() ([]byte, error) {
	return json.Marshal(s)
}

// BlockNumber is the number of a ethereum block
type BlockNumber int64

// UnmarshalJSON automatically decodes the user input for the block number, when a JSON RPC method is called
func (b *BlockNumber) UnmarshalJSON(buffer []byte) error {
	num, err := StringToBlockNumber(string(buffer))
	if err != nil {
		return err
	}
	*b = num
	return nil
}

// GetNumericBlockNumber returns a numeric block number based on the BlockNumber instance
func (b *BlockNumber) GetNumericBlockNumber(ctx context.Context, s StateInterface, e EthermanInterface, dbTx pgx.Tx) (uint64, Error) {
	bValue := LatestBlockNumber
	if b != nil {
		bValue = *b
	}

	switch bValue {
	case EarliestBlockNumber:
		return 0, nil

	case LatestBlockNumber, PendingBlockNumber:
		lastBlockNumber, err := s.GetLastL2BlockNumber(ctx, dbTx)
		if err != nil {
			return 0, NewRPCError(DefaultErrorCode, "failed to get the last block number from state")
		}

		return lastBlockNumber, nil

	case SafeBlockNumber:
		l1SafeBlockNumber, err := e.GetSafeBlockNumber(ctx)
		if err != nil {
			return 0, NewRPCError(DefaultErrorCode, "failed to get the safe block number from ethereum")
		}

		lastBlockNumber, err := s.GetLastVerifiedL2BlockNumberUntilL1Block(ctx, l1SafeBlockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return 0, nil
		} else if err != nil {
			return 0, NewRPCError(DefaultErrorCode, "failed to get the safe block number from state")
		}

		return lastBlockNumber, nil

	case FinalizedBlockNumber:
		l1FinalizedBlockNumber, err := e.GetFinalizedBlockNumber(ctx)
		if err != nil {
			return 0, NewRPCError(DefaultErrorCode, "failed to get the finalized block number from ethereum")
		}

		lastBlockNumber, err := s.GetLastVerifiedL2BlockNumberUntilL1Block(ctx, l1FinalizedBlockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return 0, nil
		} else if err != nil {
			return 0, NewRPCError(DefaultErrorCode, "failed to get the finalized block number from state")
		}

		return lastBlockNumber, nil

	default:
		if bValue < 0 {
			return 0, NewRPCError(InvalidParamsErrorCode, "invalid block number: %v", bValue)
		}
		return uint64(bValue), nil
	}
}

// StringOrHex returns the block number as a string or hex
// n == -5 = finalized
// n == -4 = safe
// n == -3 = pending
// n == -2 = latest
// n == -1 = earliest
// n >=  0 = hex(n)
func (b *BlockNumber) StringOrHex() string {
	if b == nil {
		return Latest
	}

	switch *b {
	case EarliestBlockNumber:
		return Earliest
	case PendingBlockNumber:
		return Pending
	case LatestBlockNumber:
		return Latest
	case SafeBlockNumber:
		return Safe
	case FinalizedBlockNumber:
		return Finalized
	default:
		return hex.EncodeUint64(uint64(*b))
	}
}

// StringToBlockNumber converts a string like "latest" or "0x1" to a BlockNumber instance
func StringToBlockNumber(str string) (BlockNumber, error) {
	str = strings.Trim(str, "\"")
	switch str {
	case Earliest:
		return EarliestBlockNumber, nil
	case Pending:
		return PendingBlockNumber, nil
	case Latest, "":
		return LatestBlockNumber, nil
	case Safe:
		return SafeBlockNumber, nil
	case Finalized:
		return FinalizedBlockNumber, nil
	}

	n, err := encoding.DecodeUint64orHex(&str)
	if err != nil {
		return 0, err
	}
	return BlockNumber(n), nil
}

// BlockNumberOrHash allows a string value to be parsed
// into a block number or a hash, it's used by methods
// like eth_call that allows the block to be specified
// either by the block number or the block hash
type BlockNumberOrHash struct {
	number           *BlockNumber
	hash             *ArgHash
	requireCanonical bool
}

// IsHash checks if the hash has value
func (b *BlockNumberOrHash) IsHash() bool {
	return b.hash != nil
}

// IsNumber checks if the number has value
func (b *BlockNumberOrHash) IsNumber() bool {
	return b.number != nil
}

// SetHash sets the hash and nullify the number
func (b *BlockNumberOrHash) SetHash(hash ArgHash, requireCanonical bool) {
	t := hash
	b.number = nil
	b.hash = &t
	b.requireCanonical = requireCanonical
}

// SetNumber sets the number and nullify the hash
func (b *BlockNumberOrHash) SetNumber(number BlockNumber) {
	t := number
	b.number = &t
	b.hash = nil
	b.requireCanonical = false
}

// Hash returns the hash
func (b *BlockNumberOrHash) Hash() *ArgHash {
	return b.hash
}

// Number returns the number
func (b *BlockNumberOrHash) Number() *BlockNumber {
	return b.number
}

// UnmarshalJSON automatically decodes the user input for the block number, when a JSON RPC method is called
func (b *BlockNumberOrHash) UnmarshalJSON(buffer []byte) error {
	var number BlockNumber
	err := json.Unmarshal(buffer, &number)
	if err == nil {
		b.SetNumber(number)
		return nil
	}

	var hash ArgHash
	err = json.Unmarshal(buffer, &hash)
	if err == nil {
		b.SetHash(hash, false)
		return nil
	}

	var m map[string]interface{}
	err = json.Unmarshal(buffer, &m)
	if err == nil {
		if v, ok := m[BlockNumberKey]; ok {
			vStr, ok := v.(string)
			if !ok {
				return fmt.Errorf("invalid %v", BlockNumberKey)
			}
			input, err := json.Marshal(vStr)
			if err != nil {
				return err
			}
			err = json.Unmarshal(input, &number)
			if err != nil {
				return fmt.Errorf("invalid %v", BlockNumberKey)
			}
			b.SetNumber(number)
			return nil
		} else if v, ok := m[BlockHashKey]; ok {
			vStr, ok := v.(string)
			if !ok {
				return fmt.Errorf("invalid %v", BlockHashKey)
			}
			input, err := json.Marshal(vStr)
			if err != nil {
				return err
			}
			err = json.Unmarshal(input, &hash)
			if err != nil {
				return fmt.Errorf("invalid %v", BlockHashKey)
			}
			requireCanonical, ok := m[RequireCanonicalKey]
			if ok {
				switch v := requireCanonical.(type) {
				case bool:
					b.SetHash(hash, v)
				default:
					return fmt.Errorf("invalid %v", RequireCanonicalKey)
				}
			} else {
				b.SetHash(hash, false)
			}
			return nil
		} else {
			return errors.New("invalid block or hash")
		}
	}

	return err
}

// Index of a item
type Index int64

// UnmarshalJSON automatically decodes the user input for the block number, when a JSON RPC method is called
func (i *Index) UnmarshalJSON(buffer []byte) error {
	str := strings.Trim(string(buffer), "\"")
	n, err := encoding.DecodeUint64orHex(&str)
	if err != nil {
		return err
	}
	*i = Index(n)
	return nil
}

// BatchNumber is the number of a ethereum block
type BatchNumber int64

// UnmarshalJSON automatically decodes the user input for the block number, when a JSON RPC method is called
func (b *BatchNumber) UnmarshalJSON(buffer []byte) error {
	num, err := stringToBatchNumber(string(buffer))
	if err != nil {
		return err
	}
	*b = num
	return nil
}

// GetNumericBatchNumber returns a numeric batch number based on the BatchNumber instance
func (b *BatchNumber) GetNumericBatchNumber(ctx context.Context, s StateInterface, e EthermanInterface, dbTx pgx.Tx) (uint64, Error) {
	bValue := LatestBatchNumber
	if b != nil {
		bValue = *b
	}

	switch bValue {
	case EarliestBatchNumber:
		return 0, nil

	case LatestBatchNumber:
		batchNumber, err := s.GetLastClosedBatchNumber(ctx, dbTx)
		if err != nil {
			return 0, NewRPCError(DefaultErrorCode, "failed to get the last batch number from state")
		}

		return batchNumber, nil

	case PendingBatchNumber:
		batchNumber, err := s.GetLastBatchNumber(ctx, dbTx)
		if err != nil {
			return 0, NewRPCError(DefaultErrorCode, "failed to get the pending batch number from state")
		}

		return batchNumber, nil

	case SafeBatchNumber:
		l1SafeBlockNumber, err := e.GetSafeBlockNumber(ctx)
		if err != nil {
			return 0, NewRPCError(DefaultErrorCode, "failed to get the safe batch number from ethereum")
		}

		batchNumber, err := s.GetLastVerifiedBatchNumberUntilL1Block(ctx, l1SafeBlockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return 0, nil
		} else if err != nil {
			return 0, NewRPCError(DefaultErrorCode, "failed to get the safe batch number from state")
		}

		return batchNumber, nil

	case FinalizedBatchNumber:
		l1FinalizedBlockNumber, err := e.GetFinalizedBlockNumber(ctx)
		if err != nil {
			return 0, NewRPCError(DefaultErrorCode, "failed to get the finalized batch number from ethereum")
		}

		batchNumber, err := s.GetLastVerifiedBatchNumberUntilL1Block(ctx, l1FinalizedBlockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return 0, nil
		} else if err != nil {
			return 0, NewRPCError(DefaultErrorCode, "failed to get the finalized batch number from state")
		}

		return batchNumber, nil

	default:
		if bValue < 0 {
			return 0, NewRPCError(InvalidParamsErrorCode, "invalid batch number: %v", bValue)
		}
		return uint64(bValue), nil
	}
}

// StringOrHex returns the batch number as a string or hex
// n == -5 = finalized
// n == -4 = safe
// n == -3 = pending
// n == -2 = latest
// n == -1 = earliest
// n >=  0 = hex(n)
func (b *BatchNumber) StringOrHex() string {
	if b == nil {
		return Latest
	}

	switch *b {
	case EarliestBatchNumber:
		return Earliest
	case PendingBatchNumber:
		return Pending
	case LatestBatchNumber:
		return Latest
	case SafeBatchNumber:
		return Safe
	case FinalizedBatchNumber:
		return Finalized
	default:
		return hex.EncodeUint64(uint64(*b))
	}
}

func stringToBatchNumber(str string) (BatchNumber, error) {
	str = strings.Trim(str, "\"")
	switch str {
	case Earliest:
		return EarliestBatchNumber, nil
	case Latest, "":
		return LatestBatchNumber, nil
	}

	n, err := encoding.DecodeUint64orHex(&str)
	if err != nil {
		return 0, err
	}
	return BatchNumber(n), nil
}
