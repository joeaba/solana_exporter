package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type GetHealthResponse struct {
	// the current health of the node
	Result string   `json:"result"`
	Error  rpcError `json:"error"`
}

//https://docs.solana.com/developing/clients/jsonrpc-api#gethealth
func (c *RPCClient) GetHealth(ctx context.Context, rpcAddr string) (string, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getHealth", []interface{}{}), rpcAddr)
	if err != nil {
		klog.V(2).Infof("health response: %v", err)
		return "", fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("health response: %v", string(body))

	var resp GetHealthResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return "", fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result, nil
}
