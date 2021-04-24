package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type epochinfo struct {
	FirstNormalEpoch         int64 `json:"firstNormalEpoch"`
	FirstNormalSlot          int64 `json:"firstNormalSlot"`
	LeaderScheduleSlotOffset int64 `json:"leaderScheduleSlotOffset"`
	SlotsPerEpoch            int64 `json:"slotsPerEpoch"`
	Warmup                   bool  `json:"warmup"`
}

type (
	GetEpochScheduleResponse struct {
		Result epochinfo `json:"result"`
		Error  rpcError  `json:"error"`
	}
)

//https://docs.solana.com/developing/clients/jsonrpc-api#getEpochScedule
func (c *RPCClient) GetEpochSchedule(ctx context.Context) (*GetEpochScheduleResponse, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getEpochSchedule", []interface{}{}))
	
	fmt.Println("~~Body: %w~~", body)
	fmt.Println(body == nil)
	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}
	klog.V(2).Infof("getEpochScedule response: %v", string(body))

	var resp GetEpochScheduleResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp, nil
}
