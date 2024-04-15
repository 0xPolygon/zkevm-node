//go:generate oapi-codegen -package=examplepkg -generate=types,client,spec -o=examplepkg/example-client.go beacon-node-oapi.json
package beaconclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const jsonRPCVersion = "2.0"

// BeaconAPIClient client of Beacon API
// https://ethereum.github.io/beacon-APIs/
type BeaconAPIClient struct {
	urlBase string
}

// NewClient creates an instance of client
func NewBeaconAPIClient(url string) *BeaconAPIClient {
	return &BeaconAPIClient{
		urlBase: url,
	}
}

type BeaconAPIResponse struct {
	Result json.RawMessage
}

// JSONRPCBeaconCall executes restapi call to beacon-api node
func JSONRPCBeaconCall(ctx context.Context, urlBase, methodPath string) (BeaconAPIResponse, error) {
	//url := path.Join(urlBase, methodPath)
	url := fmt.Sprintf("%s%s", urlBase, methodPath)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return BeaconAPIResponse{}, err
	}
	httpReq.Header.Add("Content-type", "application/json")

	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return BeaconAPIResponse{}, err
	}

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return BeaconAPIResponse{}, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode != http.StatusOK {
		return BeaconAPIResponse{}, fmt.Errorf("BeaconClient fails url:%s status_code:%v response:%v", url, httpRes.StatusCode, string(resBody))
	}

	return BeaconAPIResponse{
		Result: resBody,
	}, nil
}

func unserializeGenericResponse[T any](response BeaconAPIResponse) (T, error) {
	result := struct {
		Data T `json:"data"`
	}{}
	err := json.Unmarshal(response.Result, &result)
	if err != nil {
		var zero T
		return zero, err
	}
	return result.Data, nil
}

func unserializeGenericResponse2[T any](response BeaconAPIResponse) (T, error) {
	var result T
	err := json.Unmarshal(response.Result, &result)
	if err != nil {
		var zero T
		return zero, err
	}
	return result, nil
}
