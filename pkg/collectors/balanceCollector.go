package collectors

import (
	"context"
	"encoding/json"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

// AccountBalancePubkey struct which contains a
// list of pubkeys
type AccountBalancePubkey struct {
	Pubkey []string `json:"account_balance_pubkey"`
}

type BalanceCollector struct {
	RpcClient *rpc.RPCClient

	contextSlot    *prometheus.Desc
	accountBalance *prometheus.Desc
}

func NewBalanceCollector(rpcAddr string) *BalanceCollector {
	return &BalanceCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"solana_account_balance_context_slot",
			"Account Balance Context Slot",
			[]string{"pubkey"}, nil),
		accountBalance: prometheus.NewDesc(
			"solana_account_balance",
			"The balance of the account of provided Pubkey",
			[]string{"pubkey"}, nil),
	}
}

func (c *BalanceCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *BalanceCollector) mustEmitBalanceMetrics(ch chan<- prometheus.Metric, response *rpc.BalanceInfo, pubkey string) {
	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(response.Context.Slot), pubkey)
	ch <- prometheus.MustNewConstMetric(c.accountBalance, prometheus.GaugeValue, float64(response.Value), pubkey)
}

func (c *BalanceCollector) Collect(ch chan<- prometheus.Metric) {

	jsonData, err := GetKeys()
	if err != nil {
		klog.V(2).Infof("balance response: %v", err)
	}

	var keys AccountBalancePubkey
	// we unmarshal our jsonData which contains our
	// jsonFile's content into type which we defined above
	if err = json.Unmarshal(jsonData, &keys); err != nil {
		klog.V(2).Infof("failed to decode response body: %w", err)
	}

	for _, pubkey := range keys.Pubkey {
		ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
		defer cancel()

		balance, err := c.RpcClient.GetBalance(ctx, pubkey)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(-1), pubkey)
			ch <- prometheus.MustNewConstMetric(c.accountBalance, prometheus.GaugeValue, float64(-1), pubkey)
		} else {
			c.mustEmitBalanceMetrics(ch, balance, pubkey)
		}
	}
}
