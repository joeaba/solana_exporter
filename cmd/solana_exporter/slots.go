package main

import (
	"collectors"
	"context"
	"fmt"
	"time"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

func WatchSlots(c *collectors.SolanaCollector) {
	var (
		// Current mapping of relative slot numbers to leader public keys.
		epochSlots map[int64]string
		// Current epoch number corresponding to epochSlots.
		epochNumber int64
		// Last slot number we generated ticks for.
		watermark int64
	)

	ticker := time.NewTicker(slotPacerSchedule)

	for {
		<-ticker.C

		ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)

		// Get current slot height and epoch info
		info, err := c.RpcClient.GetEpochInfo(ctx, rpc.CommitmentMax)
		if err != nil {
			klog.Infof("failed to fetch epoch info, retrying: %v", err)
			cancel()
			continue
		}

		// Calculate first and last slot in epoch.
		firstSlot := info.AbsoluteSlot - info.SlotIndex
		lastSlot := firstSlot + info.SlotsInEpoch

		totalTransactionsTotal.Set(float64(info.TransactionCount))
		confirmedSlotHeight.Set(float64(info.AbsoluteSlot))
		currentEpochNumber.Set(float64(info.Epoch))
		epochFirstSlot.Set(float64(firstSlot))
		epochLastSlot.Set(float64(lastSlot))

		// Get the slot of the lowest confirmed block that has not been purged from the ledger
		block, err := c.RpcClient.GetFirstAvailableBlock(ctx)
		if err != nil {
			klog.Infof("failed to fetch first available block, retrying: %v", err)
		}

		// Get the max slot seen from retransmit stage
		maxSlot, err := c.RpcClient.GetMaxRetransmitSlot(ctx)
		if err != nil {
			klog.Infof("failed to fetch max retransmit slot, retrying: %v", err)
		}

		// Returns the current slot the node is processing
		slot, err := c.RpcClient.GetSlot(ctx)
		if err != nil {
			klog.Infof("failed to fetch slot, retrying: %v", err)
		}

		leader, err := c.RpcClient.GetSlotLeader(ctx)
		if err != nil {
			klog.Infof("failed to fetch leader, retrying: %v", err)
		}

		minLedgerSlot, err := c.RpcClient.GetMinimumLedgerSlot(ctx)
		if err != nil {
			klog.Infof("failed to fetch minimum ledger slot, retrying: %v", err)
		}

		count, err := c.RpcClient.GetTransactionCount(ctx)
		if err != nil {
			klog.Infof("failed to fetch transaction count, retrying: %v", err)
		}

		firstAvailableBlock.Set(float64(block))
		maxRetransmitSlot.Set(float64(maxSlot))
		currentSlot.Set(float64(slot))
		slotLeader.With(prometheus.Labels{"leader": leader}).Add(0)
		minimumLedgerSlot.Set(float64(minLedgerSlot))
		transactionCount.Set(float64(count))

		cancel()

		// Check whether we need to fetch a new leader schedule
		if epochNumber != info.Epoch {
			klog.Infof("new epoch at slot %d: %d (previous: %d)", firstSlot, info.Epoch, epochNumber)

			epochSlots, err = fetchLeaderSlots(c, firstSlot)
			if err != nil {
				klog.Errorf("failed to request leader schedule, retrying: %v", err)
				continue
			}

			klog.V(1).Infof("%d leader slots in epoch %d", len(epochSlots), info.Epoch)

			epochNumber = info.Epoch
			klog.V(1).Infof("we're still in epoch %d, not fetching leader schedule", info.Epoch)

			// Reset watermark to current offset on new epoch (we do not backfill slots we missed at startup)
			watermark = info.SlotIndex
		} else if watermark == info.SlotIndex {
			klog.Infof("slot has not advanced at %d, skipping", info.AbsoluteSlot)
			continue
		}

		klog.Infof("confirmed slot %d (offset %d, +%d), epoch %d (from slot %d to %d, %d remaining)",
			info.AbsoluteSlot, info.SlotIndex, info.SlotIndex-watermark, info.Epoch, firstSlot, lastSlot, lastSlot-info.AbsoluteSlot)

		// Get list of confirmed blocks since the last request. This is totally undocumented, but the result won't
		// contain missed blocks, allowing us to figure out block production success rate.
		rangeStart := firstSlot + watermark
		rangeEnd := firstSlot + info.SlotIndex - 1

		ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		cfm, err := c.RpcClient.GetConfirmedBlocks(ctx, rangeStart, rangeEnd)
		if err != nil {
			klog.Errorf("failed to request confirmed blocks at %d, retrying: %v", watermark, err)
			cancel()
			continue
		}
		cancel()

		klog.V(1).Infof("confirmed blocks: %d -> %d: %v", rangeStart, rangeEnd, cfm)

		// Figure out leaders for each block in range
		for i := watermark; i < info.SlotIndex; i++ {
			leader, ok := epochSlots[i]
			abs := firstSlot + i
			if !ok {
				// This cannot happen with a well-behaved node and is a programming error in either Solana or the exporter.
				klog.Fatalf("slot %d (offset %d) missing from epoch %d leader schedule",
					abs, i, info.Epoch)
			}

			// Check if block was included in getConfirmedBlocks output, otherwise, it was skipped.
			var present bool
			for _, s := range cfm {
				if abs == s {
					present = true
				}
			}

			var skipped string
			var label string
			if present {
				skipped = "(valid)"
				label = "valid"
			} else {
				skipped = "(SKIPPED)"
				label = "skipped"
			}

			leaderSlotsTotal.With(prometheus.Labels{"status": label, "nodekey": leader}).Add(1)
			klog.V(1).Infof("slot %d (offset %d) with leader %s %s", abs, i, leader, skipped)
		}

		watermark = info.SlotIndex
	}
}

func fetchLeaderSlots(c *collectors.SolanaCollector, epochSlot int64) (map[int64]string, error) {
	sch, err := c.RpcClient.GetLeaderSchedule(context.Background(), epochSlot)
	if err != nil {
		return nil, fmt.Errorf("failed to get leader schedule: %w", err)
	}

	slots := make(map[int64]string)

	for pk, sch := range sch {
		for _, i := range sch {
			slots[int64(i)] = pk
		}
	}

	return slots, err
}
