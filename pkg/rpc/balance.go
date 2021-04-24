package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	BalanceInfo struct {
		Context struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		// the balance of the account of provided Pubkey
		Value int64 `json:"value"`
	}

	GetBalanceResponse struct {
		Result BalanceInfo `json:"result"`
		Error  rpcError    `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#getbalance
func (c *RPCClient) GetBalance(ctx context.Context, pubkey string) (*BalanceInfo, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getBalance", []interface{}{pubkey}), c.rpcAddr)

	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}
	klog.V(2).Infof("balance response: %v", string(body))

	var resp GetBalanceResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}
