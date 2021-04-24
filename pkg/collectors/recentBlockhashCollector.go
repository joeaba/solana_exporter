package collectors

import (
	"context"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type RecentBlockhashCollector struct {
	RpcClient *rpc.RPCClient

	contextSlot          *prometheus.Desc
	blockhash            *prometheus.Desc
	lamportsPerSignature *prometheus.Desc
}

func NewRecentBlockhashCollector(rpcAddr string) *RecentBlockhashCollector {
	return &RecentBlockhashCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"solana_recent_blockhash_context_slot",
			"Recent Blockhash Context Slot",
			nil, nil),
		blockhash: prometheus.NewDesc(
			"solana_recent_blockhash",
			"A Hash as base-58 encoded string",
			[]string{"hash"}, nil),
		lamportsPerSignature: prometheus.NewDesc(
			"solana_recent_blockhash_lamports_per_signature",
			"FeeCalculator object, the fee schedule for this block hash",
			nil, nil),
	}
}

func (c *RecentBlockhashCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *RecentBlockhashCollector) mustEmitRecentBlockhashMetrics(ch chan<- prometheus.Metric, response *rpc.RecentBlockhashInfo) {
	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(response.Context.Slot))
	ch <- prometheus.MustNewConstMetric(c.blockhash, prometheus.GaugeValue, 0, response.Value.Blockhash)
	ch <- prometheus.MustNewConstMetric(c.lamportsPerSignature, prometheus.GaugeValue, float64(response.Value.FeeCalculator.LamportsPerSignature))
}

func (c *RecentBlockhashCollector) Collect(ch chan<- prometheus.Metric) {

	ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
	defer cancel()

	blockhash, err := c.RpcClient.GetRecentBlockhash(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(-1))
		ch <- prometheus.MustNewConstMetric(c.blockhash, prometheus.GaugeValue, 0, err.Error())
		ch <- prometheus.MustNewConstMetric(c.lamportsPerSignature, prometheus.GaugeValue, float64(-1))
	} else {
		c.mustEmitRecentBlockhashMetrics(ch, blockhash)
	}
}
