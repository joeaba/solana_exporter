package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	TokenSupplyInfo struct {
		Amount         string `json:"amount"`
		Decimals       int    `json:"decimals"`
		UiAmount       int    `json:"uiAmount"`
		UiAmountString string `json:"uiAmountString"`
	}

	GetTokenSupplyResponse struct {
		Result struct {
			ContextSlot int               `json:"context.slot"`
			Value       []TokenSupplyInfo `json:"value"`
		} `json:"result"`
		Error rpcError `json:"error"`
	}
)

func (c *RPCClient) GetTokenSupply(ctx context.Context, pubkey string) (*GetTokenSupplyResponse, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getTokenSupply", []interface{}{pubkey}))

	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}
	klog.V(2).Infof("GetTokenSupply response: %v", string(body))

	var resp GetTokenSupplyResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}
