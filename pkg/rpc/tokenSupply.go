package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	SupplyObject struct {
		// The raw total token supply without decimals, a string representation of u64
		Amount string `json:"amount"`
		// Number of base 10 digits to the right of the decimal place
		Decimals int64 `json:"decimals"`
		// The total token supply as a string, using mint-prescribed decimals
		UiAmountString string `json:"uiAmountString"`
	}

	TokenSupplyInfo struct {
		Context struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		Value SupplyObject `json:"value"`
	}

	GetTokenSupplyResponse struct {
		Result TokenSupplyInfo `json:"result"`
		Error  rpcError        `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#gettokensupply
func (c *RPCClient) GetTokenSupply(ctx context.Context, pubkey string) (*TokenSupplyInfo, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getTokenSupply", []interface{}{pubkey}), c.rpcAddr)

	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}
	klog.V(2).Infof("tokenSupply response: %v", string(body))

	var resp GetTokenSupplyResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}
