package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	InitializedInfo struct {
		Authority     string `json:"authority"`
		Blockhash     string `json:"blockhash"`
		FeeCalculator struct {
			LamportsPerSignature int64
		} `json:"feeCalculator"`
	}

	AccountDataInfo struct {
		Nonce struct {
			Initialized InitializedInfo `json:"initialized"`
		} `json:"nonce"`
	}

	AccountInfoJsonParsed struct {
		Data       AccountDataInfo `json:"data"`
		Executable bool            `json:"executable"`
		Lamports   int64           `json:"lamports"`
		Owner      string          `json:"owner"`
		RentEpoch  int64           `json:"rentEpoch"`
	}

	GetAccountInfoJsonParsedRes struct {
		Result struct {
			ContextSlot struct {
				Slot int64 `json:"slot"`
			} `json:"context"`
			Value AccountInfoJsonParsed `json:"value"`
		} `json:"result"`
		Error rpcError `json:"error"`
	}

	AccountInfoBase64 struct {
		Data       []string `json:"data"`
		Executable bool     `json:"executable"`
		Lamports   int64    `json:"lamports"`
		Owner      string   `json:"owner"`
		RentEpoch  int64    `json:"rentEpoch"`
	}

	GetAccountInfoBase64Res struct {
		Result struct {
			ContextSlot struct {
				Slot int64 `json:"slot"`
			} `json:"context"`
			Value AccountInfoBase64 `json:"value"`
		} `json:"result"`
		Error rpcError `json:"error"`
	}
)

func (c *RPCClient) GetAccountInfoJsonParsed(ctx context.Context, pubkey string) (*GetAccountInfoJsonParsedRes, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getAccountInfo", []interface{}{pubkey, map[string]string{"encoding": "jsonParsed"}}))
	
	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}
	klog.V(2).Infof("getAccountInfoJsonParesed response: %v", string(body))

	var resp GetAccountInfoJsonParsedRes
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}

func (c *RPCClient) GetAccountInfoBase64(ctx context.Context, pubkey string) (*GetAccountInfoBase64Res, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getAccountInfo", []interface{}{pubkey, map[string]string{"encoding": "base64"}}))
	fmt.Println("~~Body: %w~~", body)
	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("getAccountInfoBase64 response: %v", string(body))

	var resp GetAccountInfoBase64Res
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}
