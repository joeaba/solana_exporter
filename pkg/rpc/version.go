package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	GetVersionResponse struct {
		Result struct {
			// software version of solana-core
			Version string `json:"solana-core"`
		} `json:"result"`
		Error rpcError `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#getversion
func (c *RPCClient) GetVersion(ctx context.Context, rpcAddr string) (*string, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getVersion", []interface{}{}), rpcAddr)

	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("version response: %v", string(body))

	var resp GetVersionResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result.Version, nil
}
