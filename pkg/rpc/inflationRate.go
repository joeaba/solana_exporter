package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	InflationInfo struct {
		// Total inflation
		Total float64 `json:"total"`
		// Inflation allocated to validators
		Validator float64 `json:"validator"`
		// Inflation allocated to the foundation
		Foundation float64 `json:"foundation"`
		// Epoch for which these values are valid
		Epoch float64 `json:"epoch"`
	}

	GetInflationRateResponse struct {
		Result InflationInfo `json:"result"`
		Error  rpcError      `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#getinflationrate
func (c *RPCClient) GetInflationRate(ctx context.Context) (*InflationInfo, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getInflationRate", []interface{}{}), c.rpcAddr)

	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("infaltionRate response: %v", string(body))

	var resp GetInflationRateResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}
