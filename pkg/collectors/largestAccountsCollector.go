package collectors

import (
	"context"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type LargestAccountsCollector struct {
	RpcClient *rpc.RPCClient

	contextSlot     *prometheus.Desc
	accountLamports *prometheus.Desc
}

func NewLargestAccountsCollector(rpcAddr string) *LargestAccountsCollector {
	return &LargestAccountsCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"solana_largest_accounts_context_slot",
			"Context Slot for Largest Accounts",
			nil, nil),
		accountLamports: prometheus.NewDesc(
			"solana_largest_accounts",
			"The 20 largest accounts, by lamport balance",
			[]string{"address"}, nil),
	}
}

func (c *LargestAccountsCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *LargestAccountsCollector) mustEmitLargestAccountsMetrics(ch chan<- prometheus.Metric, response *rpc.LargestAccountsInfo) {
	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(response.Context.Slot))

	for _, account := range response.Value {
		ch <- prometheus.MustNewConstMetric(c.accountLamports, prometheus.GaugeValue, float64(account.Lamports), account.Address)
	}
}

func (c *LargestAccountsCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
	defer cancel()

	info, err := c.RpcClient.GetLargestAccounts(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(-1))
		ch <- prometheus.MustNewConstMetric(c.accountLamports, prometheus.GaugeValue, float64(-1), err.Error())
	} else {
		c.mustEmitLargestAccountsMetrics(ch, info)
	}
}
