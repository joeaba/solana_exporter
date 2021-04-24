package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type GetSlotResponse struct {
	// the current slot the node is processing
	Result int64    `json:"result"`
	Error  rpcError `json:"error"`
}

// https://docs.solana.com/developing/clients/jsonrpc-api#getslot
func (c *RPCClient) GetSlot(ctx context.Context) (int64, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getSlot", []interface{}{}), c.rpcAddr)

	if body == nil {
		return -1, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return -1, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("slot response: %v", string(body))

	var resp GetSlotResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return -1, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return -1, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result, nil
}
