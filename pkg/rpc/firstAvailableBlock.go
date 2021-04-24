package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type GetFirstAvailableBlockResponse struct {
	// the slot of the lowest confirmed block that has not been purged from the ledger
	Result int64    `json:"result"`
	Error  rpcError `json:"error"`
}

// https://docs.solana.com/developing/clients/jsonrpc-api#getfirstavailableblock
func (c *RPCClient) GetFirstAvailableBlock(ctx context.Context) (int64, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getFirstAvailableBlock", []interface{}{}), c.rpcAddr)

	if body == nil {
		return -1, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return -1, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("firstAvailableBlock response: %v", string(body))

	var resp GetFirstAvailableBlockResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return -1, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return -1, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result, nil
}
