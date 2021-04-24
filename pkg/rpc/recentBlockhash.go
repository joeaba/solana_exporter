package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	RecentBlockhashInfo struct {
		Context struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		// the balance of the account of provided Pubkey
		Value struct {
			Blockhash     string `json:"blockhash"`
			FeeCalculator struct {
				LamportsPerSignature int64 `json:"lamportsPerSignature"`
			} `json:"feeCalculator"`
		} `json:"value"`
	}

	GetRecentBlockhashResponse struct {
		Result RecentBlockhashInfo `json:"result"`
		Error  rpcError            `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#getrecentblockhash
func (c *RPCClient) GetRecentBlockhash(ctx context.Context) (*RecentBlockhashInfo, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getRecentBlockhash", []interface{}{}), c.rpcAddr)

	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("recentBlockhash response: %v", string(body))

	var resp GetRecentBlockhashResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}
