package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	GetHealthResponse struct {
		Result string   `json:"result"`
		Error  rpcError `json:"error"`
	}
)

//https://docs.solana.com/developing/clients/jsonrpc-api#gethealth
func (c *RPCClient) GetHealth(ctx context.Context) (string, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getHealth", []interface{}{}))
	
	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return "nil", fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return "null", fmt.Errorf("RPC call failed: %w", err)
	}
	klog.V(2).Infof("getHealth response: %v", string(body))

	var resp GetHealthResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return "null", fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return "null", fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result, nil
}
