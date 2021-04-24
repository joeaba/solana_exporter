package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	// Context struct {
	// 	Slot int `json:"slot"`
	// }

	SupplyInfo struct {
		CirculatingSupply      int   `json:"circulating"`
		NonCirculatingSupply   int   `json:"nonCirculating"`
		NonCirculatingAccounts []string `json:"nonCirculatingAccounts"`
		TotalSupply            int   `json:"total"`
	}

	// GetSupplyResponse struct {
	// 	Result struct {
	// 		Value       SupplyInfo 		`json:"value"`
	// 	} `json:"result"`
	// 	Error rpcError `json:"error"`
	// }

	GetSupplyResponse struct {
			Result struct {
				ContextSlot int        `json:"context.slot"`
				Value       SupplyInfo `json:"value"`
			} `json:"result"`
			Error rpcError `json:"error"`
	}

	// GetSupplyResponse struct {
	// 	Result SupplyInfo `json:"result"`
	// 	Error  rpcError   `json:"error"`
	// }
)

func (c *RPCClient) GetSupply(ctx context.Context) (*GetSupplyResponse, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getSupply", []interface{}{}))
	
	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("getSupply response: %v", string(body))

	var resp GetSupplyResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}
