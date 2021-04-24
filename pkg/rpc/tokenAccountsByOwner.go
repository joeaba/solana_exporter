package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	TokenAmountObject struct {
		Amount         string  `json:"amount"`
		Decimals       int64   `json:"decimals"`
		UiAmount       float64 `json:"uiAmount"`
		UiAmountString string  `json:"uiAmountString"`
	}

	InfoObject struct {
		TokenAmount   TokenAmountObject `json:"tokenAmount"`
		IsInitialized string            `json:"state"`
		IsNative      bool              `json:"isNative"`
		Mint          string            `json:"mint"`
		Owner         string            `json:"owner"`
	}

	ParsedObject struct {
		AccountType string     `json:"type"`
		Info        InfoObject `json:"info"`
	}

	AccountData struct {
		Program string       `json:"program"`
		Parsed  ParsedObject `json:"parsed"`
	}

	TokenAccounts struct {
		Account struct {
			// number of lamports assigned to this account, as a u64
			Lamports int64 `json:"lamports"`
			// base-58 encoded Pubkey of the program this account has been assigned to
			Owner string `json:"owner"`
			// token state data associated with the account, either as encoded binary data or in JSON format
			Data AccountData `json:"data"`
			// boolean indicating if the account contains a program (and is strictly read-only)
			Executable bool `json:"executable"`
			// the epoch at which this account will next owe rent, as u64
			RentEpoch int64 `json:"rentEpoch"`
		} `json:"account"`
	}

	TokenAccountsByOwnerInfo struct {
		Context struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		Value []TokenAccounts `json:"value"`
	}

	GetTokenAccountsByOwnerResponse struct {
		Result TokenAccountsByOwnerInfo `json:"result"`
		Error  rpcError                 `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#gettokenaccountsbyowner
func (c *RPCClient) GetTokenAccountsByOwner(ctx context.Context, pubkey string, mint string) (*TokenAccountsByOwnerInfo, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getTokenAccountsByOwner", []interface{}{pubkey, map[string]string{"mint": mint}, map[string]string{"encoding": "jsonParsed"}}), c.rpcAddr)

	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("tokenAccountsByOwner response: %v", string(body))

	var resp GetTokenAccountsByOwnerResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}
