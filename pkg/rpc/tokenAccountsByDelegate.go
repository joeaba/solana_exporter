package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type TokenObj struct {
	Amount         string  `json:"amount"`
	Decimals       int     `json:"decimals"`
	UiAmount       float64 `json:"uiAmount"`
	UiAmountString string  `json:"uiAmountString"`
}

type InfoObject struct {
	TokenAmount     TokenObj `json:"tokenAmount"`
	Delegate        string   `json:"delegate"`
	DelegatedAmount int      `json:"delegatedAmount"`
	IsInitialized   bool     `json:"isInitialized"`
	IsNative        bool     `json:"isNative"`
	Mint            string   `json:"mint"`
	Owner           string   `json:"owner"`
}

type ParsedInfo struct {
	AccountType string     `json:"account"`
	Info        InfoObject `json:"info"`
}

type DataInfo struct {
	Program string     `json:"program"`
	Parsed  ParsedInfo `json:"parsed"`
}

type TokenAccInfo struct {
	Data       DataInfo `json:"data"`
	Executable bool     `json:"executable"`
	Lamports   int64    `json:"lamports"`
	Owner      string   `json:"owner"`
	RentEpoch  int64    `json:"rentEpoch"`
}

type GetTokenAccountsDelegateResponse struct {
	Result struct {
		ContextSlot int            `json:"context.slot"`
		Value       []TokenAccInfo `json:"value"`
	} `json:"result"`
	Error rpcError `json:"error"`
}

func (c *RPCClient) GetTokenAccDelegate(ctx context.Context) (*GetTokenAccountsDelegateResponse, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getTokenAccountsByDelegate", []interface{}{"4Qkev8aNZcqFNSRhQzwyLMFSsi94jHqE8WNVTJzTP99F", map[string]string{"mint": "3wyAj7Rt1TWVPZVteFJPLa26JmLvdb1CAKEFZm3NY75E"}}))
	
	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("GetTokenAccDelegate response: %v", string(body))

	var resp GetTokenAccountsDelegateResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}
