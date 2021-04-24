package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	ToekenAccBal struct {
		Amount         string  `json:"amount"`
		Decimals       int     `json:"decimals"`
		UiAmount       float64 `json:"uiAmount"`
		UiAmountString string  `json:"uiAmountString"`
	}

	GetTokenAccBalRes struct {
		Result struct {
			Context int          `json:"context.slot"`
			Value   ToekenAccBal `json:"value"`
		} `json:"result"`
		Error rpcError `json:"error"`
	}
)

func (c *RPCClient) GetTokenAccount(ctx context.Context) (*GetTokenAccBalRes, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getTokenAccountBalance", []interface{}{"JCHsvHwF6TgeM1fapxgAkhVKDU5QtPox3bfCR5sjWirP"}))

	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("get Token Account Balance: %v", string(body))

	var resp GetTokenAccBalRes
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}
