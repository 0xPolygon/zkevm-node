package beaconclient

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// /eth/v1/beacon/genesis
const beaconGenesisPath = "/eth/v1/beacon/genesis"

type BeaconGenesisResponse struct {
	GenesisTime           uint64
	GenesisValidatorsRoot common.Address
	GenesisForkVersion    string
}

type beaconGenesisResponseInternal struct {
	GenesisTime           string `json:"GENESIS_TIME"`
	GenesisValidatorsRoot string `json:"GENESIS_VALIDATORS_ROOT"`
	GenesisForkVersion    string `json:"GENESIS_FORK_VERSION"`
}

func convertBeaconGenesisResponseInternal(data beaconGenesisResponseInternal) (BeaconGenesisResponse, error) {
	res := BeaconGenesisResponse{
		GenesisTime:           0,
		GenesisValidatorsRoot: common.HexToAddress(data.GenesisValidatorsRoot),
		GenesisForkVersion:    data.GenesisForkVersion,
	}
	return res, nil
}

func (c *BeaconAPIClient) BeaconGenesis(ctx context.Context) (*BeaconGenesisResponse, error) {
	response, err := JSONRPCBeaconCall(ctx, c.urlBase, beaconGenesisPath)
	if err != nil {
		return nil, err
	}

	internalStruct, err := unserializeGenericResponse[beaconGenesisResponseInternal](response)
	if err != nil {
		return nil, err
	}

	responseData, err := convertBeaconGenesisResponseInternal(internalStruct)
	if err != nil {
		return nil, err
	}
	return &responseData, nil
}
