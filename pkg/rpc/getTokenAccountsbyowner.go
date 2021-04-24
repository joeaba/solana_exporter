package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type TokenOwnerObj struct {
	Amount         string  `json:"amount"`
	Decimals       int     `json:"decimals"`
	UiAmount       float64 `json:"uiAmount"`
	UiAmountString string  `json:"uiAmountString"`
}

type InfoOwnerObject struct {
	TokenAmount     TokenOwnerObj `json:"tokenAmount"`
	Delegate        string        `json:"delegate"`
	DelegatedAmount int           `json:"delegatedAmount"`
	IsInitialized   bool          `json:"isInitialized"`
	IsNative        bool          `json:"isNative"`
	Mint            string        `json:"mint"`
	Owner           string        `json:"owner"`
}

type ParsedOwnerInfo struct {
	AccountType string          `json:"account"`
	Info        InfoOwnerObject `json:"info"`
}

type DataOwnerInfo struct {
	Program string          `json:"program"`
	Parsed  ParsedOwnerInfo `json:"parsed"`
}

type TokenAccOwnerInfo struct {
	Data       DataOwnerInfo `json:"data"`
	Executable bool          `json:"executable"`
	Lamports   int64         `json:"lamports"`
	Owner      string        `json:"owner"`
	RentEpoch  int64         `json:"rentEpoch"`
}

type GetTokenAccountsbyownerRes struct {
	Result struct {
		Context struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		Value []TokenAccOwnerInfo `json:"value"`
	} `json:"result"`
	Error rpcError `json:"error"`
}

func (c *RPCClient) GetTokenAccountOwner(ctx context.Context, pubkey string, mint string) (*GetTokenAccountsbyownerRes, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getTokenAccountsByOwner", []interface{}{pubkey, map[string]string{"mint": mint}, map[string]string{"encoding": "jsonParsed"}}))

	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}
	klog.V(2).Infof("GetTokenAccountsByOwner response: %v", string(body))

	var resp GetTokenAccountsbyownerRes
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}
