package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	LargestAccountInfo struct {
		Lamports int    `json:"lamports"`
		Address  string `json:"address"`
	}

	GetLargestAccountsResponse struct {
		Result struct {
			ContextSlot int                  `json:"context.slot"`
			Value       []LargestAccountInfo `json:"value"`
		} `json:"result"`
		Error rpcError `json:"error"`
	}
)

// type GetLargestAccRes struct {
// 	Result struct {
// 		ContextSlot int            `json:"context.slot"`
// 		Value       map[string]int `json:"value"`
// 	} `json:"result"`
// 	Error rpcError `json:"error"`
// }

// type GetLargestAccRes struct {
// 	Result struct {
// 		ContextSlot int            `json:"context.slot"`
// 		Value       map[int]string `json:"value"`
// 	} `json:"result"`
// 	Error rpcError `json:"error"`
// }

func (c *RPCClient) GetLargestAcc(ctx context.Context) (*GetLargestAccountsResponse, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getLargestAccounts", []interface{}{}))
	
	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}
	klog.V(2).Infof("getLargestAcc response: %v", string(body))

	var resp GetLargestAccountsResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}
