package collectors

import (
	"context"
	"encoding/json"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

// TokenMintPubkey struct which contains a
// list of pubkeys
type TokenMintPubkey struct {
	Pubkey []string `json:"token_mint_pubkey"`
}

type TokenSupplyCollector struct {
	RpcClient *rpc.RPCClient

	contextSlot          *prometheus.Desc
	supplyAmount         *prometheus.Desc
	supplyAmountDecimals *prometheus.Desc
	supplyAmountString   *prometheus.Desc
}

func NewTokenSupplyCollector(rpcAddr string) *TokenSupplyCollector {
	return &TokenSupplyCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"solana_token_supply_context_slot",
			"Token Supply Context Slot",
			[]string{"pubkey"}, nil),
		supplyAmount: prometheus.NewDesc(
			"solana_token_supply_amount",
			"The raw total token supply without decimals, a string representation of u64",
			[]string{"amount", "pubkey"}, nil),
		supplyAmountDecimals: prometheus.NewDesc(
			"solana_token_supply_amount_decimals",
			"Number of base 10 digits to the right of the decimal place",
			[]string{"pubkey"}, nil),
		supplyAmountString: prometheus.NewDesc(
			"solana_token_suuply_amount_string",
			"The total token supply as a string, using mint-prescribed decimals",
			[]string{"amount", "pubkey"}, nil),
	}
}

func (c *TokenSupplyCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *TokenSupplyCollector) mustEmitTokenSupplyMetrics(ch chan<- prometheus.Metric, response *rpc.TokenSupplyInfo, pubkey string) {
	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(response.Context.Slot), pubkey)
	ch <- prometheus.MustNewConstMetric(c.supplyAmount, prometheus.GaugeValue, 0, response.Value.Amount, pubkey)
	ch <- prometheus.MustNewConstMetric(c.supplyAmountDecimals, prometheus.GaugeValue, float64(response.Value.Decimals), pubkey)
	ch <- prometheus.MustNewConstMetric(c.supplyAmountString, prometheus.GaugeValue, 0, response.Value.UiAmountString, pubkey)
}

func (c *TokenSupplyCollector) Collect(ch chan<- prometheus.Metric) {

	jsonData, err := GetKeys()
	if err != nil {
		klog.V(2).Infof("tokenSupply response: %v", err)
	}

	var keys TokenMintPubkey
	// we unmarshal our jsonData which contains our
	// jsonFile's content into type which we defined above
	if err = json.Unmarshal(jsonData, &keys); err != nil {
		klog.V(2).Infof("failed to decode response body: %w", err)
	}

	for _, pubkey := range keys.Pubkey {
		ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
		defer cancel()

		supply, err := c.RpcClient.GetTokenSupply(ctx, pubkey)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(-1), pubkey)
			ch <- prometheus.MustNewConstMetric(c.supplyAmount, prometheus.GaugeValue, 0, err.Error(), pubkey)
			ch <- prometheus.MustNewConstMetric(c.supplyAmountDecimals, prometheus.GaugeValue, float64(-1), pubkey)
			ch <- prometheus.MustNewConstMetric(c.supplyAmountString, prometheus.GaugeValue, 0, err.Error(), pubkey)
		} else {
			c.mustEmitTokenSupplyMetrics(ch, supply, pubkey)
		}
	}
}
