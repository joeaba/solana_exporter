package collectors

import (
	"context"
	"encoding/json"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

// TokenAccountPubkey struct which contains a
// list of pubkeys
type TokenAccountPubkey struct {
	Pubkey []string `json:"token_account_pubkey"`
}

type TokenAccountBalanceCollector struct {
	RpcClient *rpc.RPCClient

	contextSlot           *prometheus.Desc
	balanceAmount         *prometheus.Desc
	balanceDecimals       *prometheus.Desc
	balanceUiAmountString *prometheus.Desc
}

func NewTokenAccountBalanceCollector(rpcAddr string) *TokenAccountBalanceCollector {
	return &TokenAccountBalanceCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"solana_token_account_balance_context_slot",
			"Token Account Balance Context Slot",
			[]string{"pubkey"}, nil),
		balanceAmount: prometheus.NewDesc(
			"solana_token_account_balance_amount",
			"The raw balance without decimals, a string representation of u64",
			[]string{"amount", "pubkey"}, nil),
		balanceDecimals: prometheus.NewDesc(
			"solana_token_account_balance_decimals",
			"Number of base 10 digits to the right of the decimal place",
			[]string{"pubkey"}, nil),
		balanceUiAmountString: prometheus.NewDesc(
			"solana_token_account_balance_amount_string",
			"The balance as a string, using mint-prescribed decimals",
			[]string{"amountString", "pubkey"}, nil),
	}
}

func (c *TokenAccountBalanceCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *TokenAccountBalanceCollector) mustEmitTokenAccountBalanceMetrics(ch chan<- prometheus.Metric, response *rpc.TokenAccountBalanceInfo, pubkey string) {
	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(response.Context.Slot), pubkey)
	ch <- prometheus.MustNewConstMetric(c.balanceAmount, prometheus.GaugeValue, 0, response.Value.Amount, pubkey)
	ch <- prometheus.MustNewConstMetric(c.balanceDecimals, prometheus.GaugeValue, float64(response.Value.Decimals), pubkey)
	ch <- prometheus.MustNewConstMetric(c.balanceUiAmountString, prometheus.GaugeValue, 0, response.Value.UiAmountString, pubkey)
}

func (c *TokenAccountBalanceCollector) Collect(ch chan<- prometheus.Metric) {
	jsonData, err := GetKeys()
	if err != nil {
		klog.V(2).Infof("tokenAccountBalance response: %v", err)
	}

	var keys TokenAccountPubkey
	// we unmarshal our jsonData which contains our
	// jsonFile's content into type which we defined above
	if err = json.Unmarshal(jsonData, &keys); err != nil {
		klog.V(2).Infof("failed to decode response body: %w", err)
	}

	for _, pubkey := range keys.Pubkey {
		ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
		defer cancel()

		balance, err := c.RpcClient.GetTokenAccountBalance(ctx, pubkey)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(-1), pubkey)
			ch <- prometheus.MustNewConstMetric(c.balanceAmount, prometheus.GaugeValue, 0, err.Error(), pubkey)
			ch <- prometheus.MustNewConstMetric(c.balanceDecimals, prometheus.GaugeValue, float64(-1), pubkey)
			ch <- prometheus.MustNewConstMetric(c.balanceUiAmountString, prometheus.GaugeValue, 0, err.Error(), pubkey)
		} else {
			c.mustEmitTokenAccountBalanceMetrics(ch, balance, pubkey)
		}
	}
}
