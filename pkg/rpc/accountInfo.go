package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	DataObject struct {
		Nonce struct {
			Initialized struct {
				Authority     string `json:"authority"`
				Blockhash     string `json:"blockhash"`
				FeeCalculator struct {
					LamportsPerSignature int64 `json:"lamportsPerSignature"`
				}
			} `json:"initialized"`
		} `json:"nonce"`
	}

	ValueObject struct {
		// number of lamports assigned to this account, as a u64
		Lamports int64 `json:"lamports"`
		// base-58 encoded Pubkey of the program this account has been assigned to
		Owner string `json:"owner"`
		// data associated with the account, either as encoded binary data or JSON format {<program>: <state>}, depending on encoding parameter
		Data []string `json:"data"`
		// boolean indicating if the account contains a program (and is strictly read-only)
		Executable bool `json:"executable"`
		// the epoch at which this account will next owe rent, as u64
		RentEpoch int64 `json:"rentEpoch"`
	}

	AccountInfo struct {
		Context struct {
			Slot int64 `json:"slot"`
		} `json:"context"`
		Value ValueObject `json:"value"`
	}

	GetAccountInfoResponse struct {
		Result AccountInfo `json:"result"`
		Error  rpcError    `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#getaccountinfo
func (c *RPCClient) GetAccountInfo(ctx context.Context, pubkey string) (*AccountInfo, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getAccountInfo", []interface{}{pubkey, map[string]string{"encoding": "base64"}}), c.rpcAddr)

	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	klog.V(2).Infof("accountInfo response: %v", string(body))

	var resp GetAccountInfoResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}
