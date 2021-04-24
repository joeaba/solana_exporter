package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	StackActivationInfo struct {
		Active   int64  `json:"active"`
		Inactive int64  `json:"inactive"`
		State    string `json:"state"`
	}
)
type (
	GetStackActivationResponse struct {
		Result StackActivationInfo `json:"result"`
		Error  rpcError            `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#getconfirmedblocks
func (c *RPCClient) GetStackActivation(ctx context.Context, pubkey string) (*GetStackActivationResponse, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getStakeActivation", []interface{}{pubkey}))

	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("get Stake Activation response: %v", string(body))

	var resp GetStackActivationResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}
