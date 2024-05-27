package state

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
)

// OverrideAccount indicates the overriding fields of account during the execution
// of a message call.
// Note, state and stateDiff can't be specified at the same time. If state is
// set, message execution will only use the data in the given state. Otherwise
// if statDiff is set, all diff will be applied first and then execute the call
// message.
type OverrideAccount struct {
	Nonce     *uint64                      `json:"nonce"`
	Code      *[]byte                      `json:"code"`
	Balance   *big.Int                     `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

// StateOverride is the collection of overridden accounts.
type StateOverride map[common.Address]OverrideAccount

// toExecutorStateOverride
func (so *StateOverride) toExecutorStateOverride() map[string]*executor.OverrideAccount {
	overrides := map[string]*executor.OverrideAccount{}
	if so == nil {
		return overrides
	}

	for addr, accOverride := range *so {
		var nonce uint64 = 0
		if accOverride.Nonce != nil {
			nonce = *accOverride.Nonce
		}

		var code []byte
		if accOverride.Code != nil {
			code = *accOverride.Code
		}

		var balance []byte
		if accOverride.Balance != nil {
			balance = (*accOverride.Balance).Bytes()
		}

		st := map[string]string{}
		if accOverride.State != nil {
			for k, v := range *accOverride.State {
				st[k.String()] = v.String()
			}
		}

		stDiff := map[string]string{}
		if accOverride.StateDiff != nil {
			for k, v := range *accOverride.StateDiff {
				stDiff[k.String()] = v.String()
			}
		}
		overrides[addr.String()] = &executor.OverrideAccount{
			Balance:   balance,
			Nonce:     nonce,
			Code:      code,
			State:     st,
			StateDiff: stDiff,
		}
	}
	return overrides
}

// toExecutorStateOverride
func (so *StateOverride) toExecutorStateOverrideV2() map[string]*executor.OverrideAccountV2 {
	overrides := map[string]*executor.OverrideAccountV2{}
	if so == nil {
		return overrides
	}

	for addr, accOverride := range *so {
		var nonce uint64 = 0
		if accOverride.Nonce != nil {
			nonce = *accOverride.Nonce
		}

		var code []byte
		if accOverride.Code != nil {
			code = *accOverride.Code
		}

		var balance []byte
		if accOverride.Balance != nil {
			balance = (*accOverride.Balance).Bytes()
		}

		st := map[string]string{}
		if accOverride.State != nil {
			for k, v := range *accOverride.State {
				st[k.String()] = v.String()
			}
		}

		stDiff := map[string]string{}
		if accOverride.StateDiff != nil {
			for k, v := range *accOverride.StateDiff {
				stDiff[k.String()] = v.String()
			}
		}
		overrides[addr.String()] = &executor.OverrideAccountV2{
			Balance:   balance,
			Nonce:     nonce,
			Code:      code,
			State:     st,
			StateDiff: stDiff,
		}
	}
	return overrides
}
