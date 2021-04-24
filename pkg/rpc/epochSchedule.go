package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/klog/v2"
)

type (
	EpochScheduleInfo struct {
		// the maximum number of slots in each epoch
		SlotsPerEpoch int64 `json:"slotsPerEpoch"`
		// the number of slots before beginning of an epoch to calculate a leader schedule for that epoch
		LeaderScheduleSlotOffset int64 `json:"leaderScheduleSlotOffset"`
		// whether epochs start short and grow
		Warmup bool `json:"warmup"`
		// first normal-length epoch, log2(slotsPerEpoch) - log2(MINIMUM_SLOTS_PER_EPOCH)
		FirstNormalEpoch int64 `json:"firstNormalEpoch"`
		// MINIMUM_SLOTS_PER_EPOCH * (2.pow(firstNormalEpoch) - 1)
		FirstNormalSlot int64 `json:"firstNormalSlot"`
	}

	GetEpochScheduleResponse struct {
		Result EpochScheduleInfo `json:"result"`
		Error  rpcError          `json:"error"`
	}
)

// https://docs.solana.com/developing/clients/jsonrpc-api#getepochschedule
func (c *RPCClient) GetEpochSchedule(ctx context.Context) (*EpochScheduleInfo, error) {
	body, err := c.rpcRequest(ctx, formatRPCRequest("getEpochSchedule", []interface{}{}), c.rpcAddr)

	if body == nil {
		return nil, fmt.Errorf("RPC call failed: Body empty")
	}

	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}
	klog.V(2).Infof("epochScedule response: %v", string(body))

	var resp GetEpochScheduleResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("RPC error: %d %v", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}
