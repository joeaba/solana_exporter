package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
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
		[]string{"status", "nodekey"},
	)
)

var firstAvailableBlock = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "solana_first_available_block",
	Help: "The slot of the lowest confirmed block that has not been purged from the ledger",
})

var maxRetransmitSlot = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "solana_max_retransmit_slot",
	Help: "The max slot seen from retransmit stage",
})

var currentSlot = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "solana_current_slot",
	Help: "The current slot the node is processing",
})

var slotLeader = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "solana_slot_leader",
	Help: "The current slot leader",
}, []string{"leader"})

var minimumLedgerSlot = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "solana_minimum_ledger_slot",
	Help: "The lowest slot that the node has information about in its ledger, this value may increase over time if the node is configured to purge older ledger data",
})

var transactionCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "solana_transaction_count",
	Help: "The current Transaction count from the ledger",
})

func init() {
	prometheus.MustRegister(totalTransactionsTotal)
	prometheus.MustRegister(confirmedSlotHeight)
	prometheus.MustRegister(currentEpochNumber)
	prometheus.MustRegister(epochFirstSlot)
	prometheus.MustRegister(epochLastSlot)
	prometheus.MustRegister(leaderSlotsTotal)

	prometheus.MustRegister(firstAvailableBlock)
	prometheus.MustRegister(maxRetransmitSlot)
	prometheus.MustRegister(currentSlot)
	prometheus.MustRegister(slotLeader)
	prometheus.MustRegister(minimumLedgerSlot)
	prometheus.MustRegister(transactionCount)
}
