package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type GetSlotLeaderResponse struct {
	// the current slot leader
	Result string   `json:"result"`
	Error  rpcError `json:"error"`
}

// https://docs.solana.com/developing/clients/jsonrpc-api#getslotleader
func (c *RPCClient) GetSlotLeader(ctx context.Context) (string, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getSlotLeader", []interface{}{}), c.rpcAddr)

	if body == nil {
		return "", fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return "", fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("slotLeader response: %v", string(body))

	var resp GetSlotLeaderResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return "", fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result, nil
}
