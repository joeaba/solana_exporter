package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type LamportsPerSignature struct {
	LamportsPerSignature int `json:"lamportsPerSignature"`
}
type GetRecentBlockHashRes struct {
	Result struct {
		ContextSlot struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		Value struct {
			Blockhash     string               `json:"blockhash"`
			FeeCalculator LamportsPerSignature `json:"feeCalculator"`
		}
	} `json:"result"`
	Error rpcError `json:"error"`
}

func (c *RPCClient) GetRecentBlockHash(ctx context.Context, commitment Commitment) (*GetRecentBlockHashRes, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getBalance", []interface{}{commitment}))

	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}
	klog.V(2).Infof("getRecentBlockHash response: %v", string(body))

	var resp GetRecentBlockHashRes
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}
