package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	GetTransectionCountResponse struct {
		Result int      `json:"result"`
		Error  rpcError `json:"error"`
	}
)

func (c *RPCClient) GetTransectionCount(ctx context.Context) (*GetTransectionCountResponse, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getTransactionCount", []interface{}{}))

	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("getTransectionCount response: %v", string(body))

	var resp GetTransectionCountResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}
