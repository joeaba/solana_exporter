package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	AccountValue struct {
		// Number of lamports in the account
		Lamports int64 `json:"lamports"`
		// Base-58 encoded address of the account
		Address string `json:"address"`
	}

	LargestAccountsInfo struct {
		Context struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		Value []AccountValue `json:"value"`
	}

	GetLargestAccountsResponse struct {
		Result LargestAccountsInfo `json:"result"`
		Error  rpcError            `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#getlargestaccounts
func (c *RPCClient) GetLargestAccounts(ctx context.Context) (*LargestAccountsInfo, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getLargestAccounts", []interface{}{}), c.rpcAddr)

	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("largestAccounts response: %v", string(body))

	var resp GetLargestAccountsResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}
