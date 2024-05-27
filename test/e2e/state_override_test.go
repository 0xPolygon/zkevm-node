package e2e

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/StateOverride"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
)

func TestStateOverride(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	setup()
	defer teardown()

	ctx := context.Background()

	for _, network := range networks {
		log.Infof("Network %s", network.Name)

		ethereumClient, err := ethclient.Dial(network.URL)
		require.NoError(t, err)

		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		// deploy a smart contract
		scAddr, tx, _, err := StateOverride.DeployStateOverride(auth, ethereumClient)
		require.NoError(t, err)
		err = operations.WaitTxToBeMined(ctx, ethereumClient, tx, operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		gasPrice, err := ethereumClient.SuggestGasPrice(ctx)
		require.NoError(t, err)

		type testCase struct {
			name    string
			execute func(*testing.T, *testCase)
		}

		testCases := []testCase{
			{
				name: "override address balance",
				execute: func(t *testing.T, tc *testCase) {
					address := common.HexToAddress("0x123456789")
					currentBalance, err := ethereumClient.BalanceAt(ctx, address, nil)
					require.NoError(t, err)

					assert.Equal(t, uint64(0), currentBalance.Uint64())

					methodSignature := []byte("addrBalance(address)")
					hash := sha3.NewLegacyKeccak256()
					hash.Write(methodSignature)
					methodID := hash.Sum(nil)[:4]

					paddedAddress := common.LeftPadBytes(address.Bytes(), 32)

					var data []byte
					data = append(data, methodID...)
					data = append(data, paddedAddress...)
					msg := map[string]interface{}{
						"from":     auth.From.String(),
						"to":       scAddr.String(),
						"gasPrice": hex.EncodeBig(gasPrice),
						"data":     hex.EncodeToHex(data),
					}

					res, err := client.JSONRPCCall(network.URL, "eth_estimateGas", msg, "latest")
					require.NoError(t, err)
					require.Nil(t, res.Error)
					require.NotNil(t, res.Result)

					var gasEstimationHex string
					err = json.Unmarshal(res.Result, &gasEstimationHex)
					require.NoError(t, err)

					msg["gas"] = gasEstimationHex

					res, err = client.JSONRPCCall(network.URL, "eth_call", msg, "latest")
					require.NoError(t, err)
					require.Nil(t, res.Error)
					require.NotNil(t, res.Result)

					var notPaddedResultHex string
					err = json.Unmarshal(res.Result, &notPaddedResultHex)
					require.NoError(t, err)
					result := hex.DecodeBig(notPaddedResultHex)

					s1, s2 := currentBalance.String(), result.String()
					log.Debug(s1, " | ", s2)
					assert.Equal(t, s1, s2)

					newBalance := common.Big0.SetUint64(1234567890)
					stateOverride := map[string]interface{}{
						address.String(): map[string]interface{}{
							"balance": hex.EncodeBig(newBalance),
						},
					}

					res, err = client.JSONRPCCall(network.URL, "eth_call", msg, "latest", stateOverride)
					require.NoError(t, err)
					require.Nil(t, res.Error)
					require.NotNil(t, res.Result)

					err = json.Unmarshal(res.Result, &notPaddedResultHex)
					require.NoError(t, err)
					result = hex.DecodeBig(notPaddedResultHex)

					s1, s2 = newBalance.String(), result.String()
					log.Debug(s1, " | ", s2)
					assert.Equal(t, s1, s2)
				},
			},
			// {
			// 	name: "override sender and sc balance in a single call",
			// },
			// {
			// 	name: "override sc code",
			// },
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				tc := testCase
				testCase.execute(t, &tc)
			})
		}
	}
}
