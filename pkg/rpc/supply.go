package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	CurrentSupply struct {
		// total supply in lamports
		TotalSupply int64 `json:"total"`
		// circulating supply in lamports
		CirculatingSupply int64 `json:"circulating"`
		// non-circulating supply in lamports
		NonCirculatingSupply int64 `json:"nonCirculating"`
		// an array of account addresses of non-circulating accounts, as strings
		NonCirculatingAccounts []string `json:"nonCirculatingAccounts"`
	}

	SupplyInfo struct {
		Context struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		Value CurrentSupply `json:"value"`
	}

	GetSupplyResponse struct {
		Result SupplyInfo `json:"result"`
		Error  rpcError   `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#getsupply
func (c *RPCClient) GetSupply(ctx context.Context) (*SupplyInfo, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getSupply", []interface{}{}), c.rpcAddr)

	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("supply response: %v", string(body))

	var resp GetSupplyResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}
