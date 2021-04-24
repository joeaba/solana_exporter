package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	GetMaxRetransmitSlotResponse struct {
		Result int64    `json:"result"`
		Error  rpcError `json:"error"`
	}
)

func (c *RPCClient) GetMaxRetransmitSlot(ctx context.Context) (int64, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getMaxRetransmitSlot", []interface{}{}))
	
	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return 0, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return 0, fmt.Errorf("RPC call failed: %w", err)
	}
	klog.V(2).Infof("getMaxRetransmitSlot response: %v", string(body))
	
	var resp GetMaxRetransmitSlotResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return 0, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result, nil
}
