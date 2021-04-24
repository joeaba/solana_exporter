package collectors

import (
	"context"
	"encoding/json"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

// StakeAccountPubkey struct which contains a
// list of pubkeys
type StakeAccountPubkey struct {
	Pubkey []string `json:"stake_account_pubkey"`
}

type StakeActivationCollector struct {
	RpcClient *rpc.RPCClient

	accountState  *prometheus.Desc
	stakeActive   *prometheus.Desc
	stakeInactive *prometheus.Desc
}

func NewStakeActivationCollector(rpcAddr string) *StakeActivationCollector {
	return &StakeActivationCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		accountState: prometheus.NewDesc(
			"solana_stake_account_activation_stake",
			"The stake account's activation state",
			[]string{"state", "pubkey"}, nil),
		stakeActive: prometheus.NewDesc(
			"solana_stake_active",
			"Stake active during the epoch",
			[]string{"pubkey"}, nil),
		stakeInactive: prometheus.NewDesc(
			"solana_stake_inactive",
			"Stake inactive during the epoch",
			[]string{"pubkey"}, nil),
	}
}

func (c *StakeActivationCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *StakeActivationCollector) mustEmitStakeActivationMetrics(ch chan<- prometheus.Metric, response *rpc.StakeActivationInfo, pubkey string) {
	ch <- prometheus.MustNewConstMetric(c.accountState, prometheus.GaugeValue, 0, response.State, pubkey)
	ch <- prometheus.MustNewConstMetric(c.stakeActive, prometheus.GaugeValue, float64(response.Active), pubkey)
	ch <- prometheus.MustNewConstMetric(c.stakeInactive, prometheus.GaugeValue, float64(response.Inactive), pubkey)
}

func (c *StakeActivationCollector) Collect(ch chan<- prometheus.Metric) {
	jsonData, err := GetKeys()
	if err != nil {
		klog.V(2).Infof("stakeActivation response: %v", err)
	}

	var keys StakeAccountPubkey
	// we unmarshal our jsonData which contains our
	// jsonFile's content into type which we defined above
	if err = json.Unmarshal(jsonData, &keys); err != nil {
		klog.V(2).Infof("failed to decode response body: %w", err)
	}

	for _, pubkey := range keys.Pubkey {
		ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
		defer cancel()

		info, err := c.RpcClient.GetStakeActivation(ctx, pubkey)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(c.accountState, prometheus.GaugeValue, 0, err.Error(), pubkey)
			ch <- prometheus.MustNewConstMetric(c.stakeActive, prometheus.GaugeValue, float64(-1), pubkey)
			ch <- prometheus.MustNewConstMetric(c.stakeInactive, prometheus.GaugeValue, float64(-1), pubkey)
		} else {
			c.mustEmitStakeActivationMetrics(ch, info, pubkey)
		}
	}
}
