package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

const (
	slotPacerSchedule = 1 * time.Second
)

var (
	totalTransactionsTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "solana_confirmed_transactions_total",
		Help: "Total number of transactions processed since genesis (max confirmation)",
	})

	confirmedSlotHeight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "solana_confirmed_slot_height",
		Help: "Last confirmed slot height processed by watcher routine (max confirmation)",
	})

	currentEpochNumber = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "solana_confirmed_epoch_number",
		Help: "Current epoch (max confirmation)",
	})

	epochFirstSlot = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "solana_confirmed_epoch_first_slot",
		Help: "Current epoch's first slot (max confirmation)",
	})

	epochLastSlot = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "solana_confirmed_epoch_last_slot",
		Help: "Current epoch's last slot (max confirmation)",
	})

	leaderSlotsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "solana_leader_slots_total",
			Help: "Number of leader slots per leader, grouped by skip status (max confirmation)",
		},
		[]string{"status", "nodekey"})

	getHealth = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "solana_health",
		Help: "Current Health",
	})

	getFirstAvailableBlock = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "first_available_block",
		Help: "Current First_Block",
	})

	getInflationEpoch = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Infaltion_Epoch",
		Help: "Current Infaltion_Epoch",
	})

	getInfaltionFoundation = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Infaltion_Foundation",
		Help: "Current Inflation foundation",
	})

	getInfaltionTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Infaltion_Total",
		Help: "Current Infaltion total",
	})

	getInfaltionValidator = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Infaltion_Validator",
		Help: "Current Infaltion validator",
	})

	getMaxRetransmitSlot = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Max_Retransmit_slot",
		Help: "Current Retransmit Slot",
	})

	getVersion = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "Get_Version",
			Help: "Current version",
		},
		[]string{"version"})
	// getLargestAcc = prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Name: "getSupply",
	// 	Help: "Current getSupply",
	// })

	getTokenAccountBalance = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Get_Token_Acc_Balance",
		Help: "Current Token Account Balance",
	})

	getEpochSceduleInfoBool = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "Get_Epoch_Schedule_Info_Bool",
		Help: "Current Epoch Schedule info bool",
	},
		[]string{"warmup"})

	getFirstNormalEpoch = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Get_Get_First_Normal_Epoch",
		Help: "Current GetFirstNormalEpoch",
	})

	getFirstNormalSlot = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Get_Get_First_Normal_Slot",
		Help: "Current GetFirstNormalSlot",
	})

	getLeaderScheduleSlotOffset = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Get_Get_LeaderSchedule_SlotOffset",
		Help: "Current GetLeaderScheduleSlotOffset",
	})

	getSlotsPerEpoch = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Get_GetSlotsPerEpoch",
		Help: "Current GetSlotsPerEpoch",
	})

	getSlot = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Get_Slot_Info",
		Help: "Current Slot",
	})

	getBalance = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Get_Balance_Info",
		Help: "Current Balance",
	})

	getTransactionCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Get_Transection_Info",
		Help: "Transection Info",
	})

	getSlotleader = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "Get_Slot_Leader",
		Help: "Slot Leader",
	},
		[]string{"slotleader"})

	getMinimumLeadger = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Get_Min_Leadger",
		Help: "Min Leadger",
	})
	getRecentContext = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Get_Recent_Blockhash_Context",
		Help: "Recent Context",
	})
	lamportsPerSignature = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Get_Lamport_Per_Signature",
		Help: "Lamport Per Signature",
	})
	getBlockHash = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "Get_Blockhash_ContextBlock",
		Help: "Hash Context",
	},
		[]string{"blockhash"})
)

func init() {

	prometheus.MustRegister(totalTransactionsTotal)
	prometheus.MustRegister(confirmedSlotHeight)
	prometheus.MustRegister(currentEpochNumber)
	prometheus.MustRegister(epochFirstSlot)
	prometheus.MustRegister(epochLastSlot)
	prometheus.MustRegister(leaderSlotsTotal)
	prometheus.MustRegister(getHealth)
	prometheus.MustRegister(getFirstAvailableBlock)
	prometheus.MustRegister(getInflationEpoch)
	prometheus.MustRegister(getInfaltionFoundation)
	prometheus.MustRegister(getInfaltionTotal)
	prometheus.MustRegister(getInfaltionValidator)
	prometheus.MustRegister(getEpochSceduleInfoBool)
	//prometheus.MustRegister(getLargestAcc)
	prometheus.MustRegister(getTokenAccountBalance)
	prometheus.MustRegister(getVersion)
	prometheus.MustRegister(getSlot)
	prometheus.MustRegister(getBalance)
	prometheus.MustRegister(getFirstNormalEpoch)
	prometheus.MustRegister(getFirstNormalSlot)
	prometheus.MustRegister(getSlotsPerEpoch)
	prometheus.MustRegister(getLeaderScheduleSlotOffset)
	prometheus.MustRegister(getTransactionCount)
	prometheus.MustRegister(getSlotleader)
	prometheus.MustRegister(getMinimumLeadger)
	prometheus.MustRegister(getRecentContext)
	prometheus.MustRegister(getBlockHash)
	prometheus.MustRegister(lamportsPerSignature)

}

func (c *solanaCollector) WatchSlots() {
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

		// Get current slot height and epoch info
		ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
		info, err := c.rpcClient.GetEpochInfo(ctx, rpc.CommitmentMax)
		if err != nil {
			klog.Infof("failed to fetch info info, retrying: %v", err)
			cancel()
			continue
		}
		cancel()

		// Calculate first and last slot in epoch.
		firstSlot := info.AbsoluteSlot - info.SlotIndex
		lastSlot := firstSlot + info.SlotsInEpoch

		totalTransactionsTotal.Set(float64(info.TransactionCount))
		confirmedSlotHeight.Set(float64(info.AbsoluteSlot))
		currentEpochNumber.Set(float64(info.Epoch))
		epochFirstSlot.Set(float64(firstSlot))
		epochLastSlot.Set(float64(lastSlot))

		// Check whether we need to fetch a new leader schedule
		if epochNumber != info.Epoch {
			klog.Infof("new epoch at slot %d: %d (previous: %d)", firstSlot, info.Epoch, epochNumber)

			epochSlots, err = c.fetchLeaderSlots(firstSlot)
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

		// Get Health

		ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		health, err := c.rpcClient.GetHealth(ctx)
		klog.Infof("Health is: %v", health)
		fmt.Println(health + "HEALth full information")
		if err != nil {
			klog.Infof("failed to fetch info info, retrying: %v", err)
			cancel()
			continue
		}

		cancel()

		if health == "ok" {
			getHealth.Set(1)
		} else {
			getHealth.Set(0)
		}

		// Get First Available Block

		ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		firstavailableblock, err := c.rpcClient.GetFirstAvailableBlock(ctx)
		klog.Infof("firstavailableblock is: %v", firstavailableblock)

		if err != nil {
			klog.Infof("failed to fetch info info, retrying: %v", err)
			cancel()
			continue
		}

		cancel()

		getFirstAvailableBlock.Set(float64(firstavailableblock))

		// Get Transection Count

		ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		gettransactioncount, err := c.rpcClient.GetTransectionCount(ctx)
		klog.Infof("Transection Count is: %v", gettransactioncount)
		fmt.Println(health + "Transection Count information")
		if err != nil {
			klog.Infof("failed to fetch info info, retrying: %v", err)
			cancel()
			continue
		}

		cancel()

		getTransactionCount.Set(float64(gettransactioncount.Result))

		// Get Inflation Rate

		ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		inflationrate, err := c.rpcClient.GetInflationRate(ctx, rpc.CommitmentMax)

		if err != nil {
			klog.Infof("failed to fetch info info, retrying: %v", err)
			cancel()
			continue
		}
		klog.Infof("Infaltion Rate is: %v", inflationrate)

		cancel()

		// Calculate first and last slot in epoch.
		Epoch := inflationrate.Epoch
		Foundation := inflationrate.Foundation
		Total := inflationrate.Total
		Validator := inflationrate.Validator

		getInflationEpoch.Set(float64(Epoch))
		getInfaltionFoundation.Set(Foundation)
		getInfaltionTotal.Set(Total)
		getInfaltionValidator.Set(Validator)

		// Get Max Retransmit Slot

		ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		retransmitslot, err := c.rpcClient.GetMaxRetransmitSlot(ctx)

		klog.Infof("Retransmit Slot is: %v", retransmitslot)

		if err != nil {
			klog.Infof("failed to fetch info info, retrying: %v", err)
			cancel()
			continue
		}

		getMaxRetransmitSlot.Set(float64(retransmitslot))

		//Get Version

		ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		getversion, err := c.rpcClient.GetVersion(ctx)
		klog.Infof("Get Token Account value is: %v", getversion)

		if err != nil {
			klog.Infof("failed to fetch info info, retrying: %v", err)
			cancel()
			continue
		}

		getVersion.With(prometheus.Labels{"version": getversion.Result.SolonaCore}).Add(0)

		//Get token acc by delegate

		ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		gettokenacc, err := c.rpcClient.GetTokenAccDelegate(ctx)
		klog.Infof("Get Token Account Delegate is: %v", gettokenacc)

		if err != nil {
			klog.Infof("failed to fetch info info, retrying: %v", err)
			cancel()
			continue
		}

		//Get Eopch InfoSchedule

		ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		epochschedule, err := c.rpcClient.GetEpochSchedule(ctx)
		klog.Infof("Get Epoch Schedule is: %v", epochschedule)

		if err != nil {
			klog.Infof("failed to fetch info info, retrying: %v", err)
			cancel()
			continue
		}

		cancel()

		getFirstNormalEpoch.Set(float64(epochschedule.Result.FirstNormalEpoch))
		getFirstNormalSlot.Set(float64(epochschedule.Result.FirstNormalSlot))
		getLeaderScheduleSlotOffset.Set(float64(epochschedule.Result.LeaderScheduleSlotOffset))
		getSlotsPerEpoch.Set(float64(epochschedule.Result.SlotsPerEpoch))
		getEpochSceduleInfoBool.With(prometheus.Labels{"warmup": strconv.FormatBool(epochschedule.Result.Warmup)}).Add(0)

		//Get Slot Response

		ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		getslot, err := c.rpcClient.GetSlot(ctx)
		klog.Infof("Get Slot: %v", getslot)

		if err != nil {
			klog.Infof("failed to fetch info info, retrying: %v", err)
			cancel()
			continue
		}

		cancel()

		getSlot.Set(float64(getslot.Result))

		//Get Slot Leader

		ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		getslotleader, err := c.rpcClient.GetSlotleader(ctx)
		klog.Infof("Get Slot Leader: %v", getslotleader)

		if err != nil {
			klog.Infof("failed to fetch info info, retrying: %v", err)
			cancel()
			continue
		}

		cancel()

		getSlotleader.With(prometheus.Labels{"slotleader": getslotleader.Result}).Add(0)

		// Get Recent Block Hash

		// ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		// recentblockhash, err := c.rpcClient.GetRecentBlockHash(ctx, rpc.CommitmentMax)
		// klog.Infof("Recent Block Hash Is: %v", recentblockhash)

		// if err != nil {
		// 	klog.Infof("failed to fetch info info, retrying: %v", err)
		// 	cancel()
		// 	continue
		// }

		// cancel()

		// getRecentContext.Set(float64(recentblockhash.Result.ContextSlot.Slot))
		// getBlockHash.With(prometheus.Labels{"blockhash": recentblockhash.Result.Value.Blockhash}).Add(0)
		// lamportsPerSignature.Set(float64(recentblockhash.Result.Value.FeeCalculator.LamportsPerSignature))
		// //Ger Minimum Leadger Slot

		ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		minimumleadgerslot, err := c.rpcClient.GetMinimunLeadegerSlot(ctx)
		klog.Infof("Get Minimum Leadger Slot: %v", minimumleadgerslot)

		if err != nil {
			klog.Infof("failed to fetch info, retrying: %v", err)
			cancel()
			continue
		}

		cancel()

		getMinimumLeadger.Set(float64(minimumleadgerslot.Result))

		// var myarr [4]string
		// myarr[0] = "83astBRguLMdt2h5U1Tpdq5tjFoJ6noeGwaY3mDLVcri"
		// myarr[1] = "vines1vzrYbzLMRdu58ou5XTby4qAqVRLmqo36NKPTg"
		// myarr[2] = "4fYNw3dojWmQ4dXtSGE9epjRGy9pFSx62YypT7avPYvA"
		// myarr[3] = "6H94zdiaYfRfPfKjYLjyr2VFBg6JHXygy84r3qhc3NsC"

		// //Get Balance Response

		// for i := 0; i <= len(myarr)-1; i++ {
		// 	ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		// 	getbalance, err := c.rpcClient.GetBalance(ctx, myarr[i])
		// 	klog.Infof("Get Balance: %v", getbalance)

		// 	if err != nil {
		// 		klog.Infof("failed to fetch info info, retrying: %v", err)
		// 		cancel()
		// 		continue
		// 	}
		// 	cancel()

		// 	getBalance.Set(float64(getbalance.Result.ContextSlot))
		// 	getBalance.Set(float64(getbalance.Result.Value))

		// }
		//ver, _ := strconv.ParseInt(getversion.Result, 10, 32)

		// getMaxRetransmitSlot.Add(float64(retransmitslot))

		// Get list of confirmed blocks since the last request. This is totally undocumented, but the result won't
		// contain missed blocks, allowing us to figure out block production success rate.

		rangeStart := firstSlot + watermark
		rangeEnd := firstSlot + info.SlotIndex - 1

		ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		cfm, err := c.rpcClient.GetConfirmedBlocks(ctx, rangeStart, rangeEnd)
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

func (c *solanaCollector) fetchLeaderSlots(epochSlot int64) (map[int64]string, error) {
	sch, err := c.rpcClient.GetLeaderSchedule(context.Background(), epochSlot)
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
