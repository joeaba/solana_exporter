package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type GetBalanceResponse struct {
	Result struct {
		ContextSlot struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		Value       int64 `json:"value"`
	} `json:"result"`
	Error rpcError `json:"error"`
}

func (c *RPCClient) GetBalance(ctx context.Context, pubkey string) (*GetBalanceResponse, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getBalance", []interface{}{pubkey}))
	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}
	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("getBalance response: %v", string(body))

	var resp GetBalanceResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}