package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	GetMaxRetransmitSlotResponse struct {
		// the max slot seen from retransmit stage
		Result int64    `json:"result"`
		Error  rpcError `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#getmaxretransmitslot
func (c *RPCClient) GetMaxRetransmitSlot(ctx context.Context) (int64, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getMaxRetransmitSlot", []interface{}{}), c.rpcAddr)

	if body == nil {
		return 0, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return 0, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("maxRetransmitSlot response: %v", string(body))

	var resp GetMaxRetransmitSlotResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return 0, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result, nil
}
