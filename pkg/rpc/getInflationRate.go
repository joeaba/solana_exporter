package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	Inflationinfo struct {
		// Current absolute slot in epoch
		Epoch int64 `json:"epoch"`
		// Current block height
		Foundation float64 `json:"foundation"`
		// Current epoch number
		Total float64 `json:"total"`
		// Current slot relative to the start of the current epoch
		Validator float64 `json:"validator"`
	}

	GetInflationRateRes struct {
		Result Inflationinfo `json:"result"`
		Error  rpcError      `json:"error"`
	}
)

//https://docs.solana.com/developing/clients/jsonrpc-api#getinflationrate
func (c *RPCClient) GetInflationRate(ctx context.Context, commitment Commitment) (*Inflationinfo, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getInflationRate", []interface{}{}))
	
	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}
	klog.V(2).Infof("getInfaltionRate response: %v", string(body))

	var resp GetInflationRateRes
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}
