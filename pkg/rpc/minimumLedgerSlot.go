package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type GetMinimumLedgerSlotResponse struct {
	// the lowest slot that the node has information about in its ledger, this value may increase over time if the node is configured to purge older ledger data
	Result int64    `json:"result"`
	Error  rpcError `json:"error"`
}

// https://docs.solana.com/developing/clients/jsonrpc-api#minimumledgerslot
func (c *RPCClient) GetMinimumLedgerSlot(ctx context.Context) (int64, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("minimumLedgerSlot", []interface{}{}), c.rpcAddr)

	if body == nil {
		return -1, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return -1, fmt.Errorf("RPC call failed: %w", err)
	}
	klog.V(2).Infof("minimumLedgerSlot response: %v", string(body))

	var resp GetMinimumLedgerSlotResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return -1, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return -1, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result, nil
}
