package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	TokenBalance struct {
		// the raw balance without decimals, a string representation of u64
		Amount string `json:"amount"`
		// number of base 10 digits to the right of the decimal place
		Decimals int64 `json:"decimals"`
		// the balance as a string, using mint-prescribed decimals
		UiAmountString string `json:"uiAmountString"`
	}

	TokenAccountBalanceInfo struct {
		Context struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		Value TokenBalance `json:"value"`
	}

	GetTokenAccountBalanceResponse struct {
		Result TokenAccountBalanceInfo `json:"result"`
		Error  rpcError                `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#gettokenaccountbalance
func (c *RPCClient) GetTokenAccountBalance(ctx context.Context, pubkey string) (*TokenAccountBalanceInfo, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getTokenAccountBalance", []interface{}{pubkey}), c.rpcAddr)

	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("tokenAccountBalance: %v", string(body))

	var resp GetTokenAccountBalanceResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}
