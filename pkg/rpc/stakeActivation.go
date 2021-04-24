package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	StakeActivationInfo struct {
		// the stake account's activation state
		State string `json:"state"`
		// stake active during the epoch
		Active int64 `json:"active"`
		// stake inactive during the epoch
		Inactive int64 `json:"inactive"`
	}
)
type (
	GetStackActivationResponse struct {
		Result StakeActivationInfo `json:"result"`
		Error  rpcError            `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#getstakeactivation
func (c *RPCClient) GetStakeActivation(ctx context.Context, pubkey string) (*StakeActivationInfo, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getStakeActivation", []interface{}{pubkey}), c.rpcAddr)

	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("stakeActivation response: %v", string(body))

	var resp GetStackActivationResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}
