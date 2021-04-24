package collectors

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

// AccountInfoPubkey struct which contains a
// list of pubkeys
type AccountInfoPubkey struct {
	Pubkey []string `json:"account_info_pubkey"`
}

type AccountInfoCollector struct {
	RpcClient *rpc.RPCClient

	contextSlot       *prometheus.Desc
	accountData       *prometheus.Desc
	accountExecutable *prometheus.Desc
	accountLamports   *prometheus.Desc
	accountOwner      *prometheus.Desc
	accountRentEpoch  *prometheus.Desc
}

func NewAccountInfoCollector(rpcAddr string) *AccountInfoCollector {
	return &AccountInfoCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"solana_account_info_context_slot",
			"Account Info Context Slot",
			[]string{"pubkey"}, nil),
		accountData: prometheus.NewDesc(
			"solana_account_data",
			"Data associated with the account, either as encoded binary data or JSON format {<program>: <state>}, depending on encoding parameter",
			[]string{"data", "pubkey"}, nil),
		accountExecutable: prometheus.NewDesc(
			"solana_account_executable",
			"Boolean indicating if the account contains a program (and is strictly read-only)",
			[]string{"executable", "pubkey"}, nil),
		accountLamports: prometheus.NewDesc(
			"solana_account_lamports",
			"Number of lamports assigned to this account, as a u64",
			[]string{"pubkey"}, nil),
		accountOwner: prometheus.NewDesc(
			"solana_account_owner",
			"Base-58 encoded Pubkey of the program this account has been assigned to",
			[]string{"owner", "pubkey"}, nil),
		accountRentEpoch: prometheus.NewDesc(
			"solana_account_rent_epoch",
			"The epoch at which this account will next owe rent, as u64",
			[]string{"pubkey"}, nil),
	}
}

func (c *AccountInfoCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *AccountInfoCollector) mustEmitAccountInfoMetrics(ch chan<- prometheus.Metric, response *rpc.AccountInfo, pubkey string) {
	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(response.Context.Slot), pubkey)

	if response.Value.Data != nil {
		ch <- prometheus.MustNewConstMetric(c.accountData, prometheus.GaugeValue, 0, response.Value.Data[0], pubkey)
		ch <- prometheus.MustNewConstMetric(c.accountExecutable, prometheus.GaugeValue, 0, strconv.FormatBool(response.Value.Executable), pubkey)
		ch <- prometheus.MustNewConstMetric(c.accountLamports, prometheus.GaugeValue, float64(response.Value.Lamports), pubkey)
		ch <- prometheus.MustNewConstMetric(c.accountOwner, prometheus.GaugeValue, 0, response.Value.Owner, pubkey)
		ch <- prometheus.MustNewConstMetric(c.accountRentEpoch, prometheus.GaugeValue, float64(response.Value.RentEpoch), pubkey)
	} else {
		ch <- prometheus.MustNewConstMetric(c.accountData, prometheus.GaugeValue, 0, "the requested account doesn't exist", pubkey)
		ch <- prometheus.MustNewConstMetric(c.accountExecutable, prometheus.GaugeValue, 0, "the requested account doesn't exist", pubkey)
		ch <- prometheus.MustNewConstMetric(c.accountLamports, prometheus.GaugeValue, float64(-1), pubkey)
		ch <- prometheus.MustNewConstMetric(c.accountOwner, prometheus.GaugeValue, 0, "the requested account doesn't exist", pubkey)
		ch <- prometheus.MustNewConstMetric(c.accountRentEpoch, prometheus.GaugeValue, float64(-1), pubkey)
	}
}

func (c *AccountInfoCollector) Collect(ch chan<- prometheus.Metric) {

	jsonData, err := GetKeys()
	if err != nil {
		klog.V(2).Infof("accountInfo response: %v", err)
	}

	var keys AccountInfoPubkey
	// we unmarshal our jsonData which contains our
	// jsonFile's content into type which we defined above
	if err = json.Unmarshal(jsonData, &keys); err != nil {
		klog.V(2).Infof("failed to decode response body: %w", err)
	}

	for _, pubkey := range keys.Pubkey {
		ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
		defer cancel()

		info, err := c.RpcClient.GetAccountInfo(ctx, pubkey)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(-1), pubkey)
			ch <- prometheus.MustNewConstMetric(c.accountData, prometheus.GaugeValue, 0, err.Error(), pubkey)
			ch <- prometheus.MustNewConstMetric(c.accountExecutable, prometheus.GaugeValue, 0, err.Error(), pubkey)
			ch <- prometheus.MustNewConstMetric(c.accountLamports, prometheus.GaugeValue, float64(-1), pubkey)
			ch <- prometheus.MustNewConstMetric(c.accountOwner, prometheus.GaugeValue, 0, err.Error(), pubkey)
			ch <- prometheus.MustNewConstMetric(c.accountRentEpoch, prometheus.GaugeValue, float64(-1), pubkey)
		} else {
			c.mustEmitAccountInfoMetrics(ch, info, pubkey)
		}
	}
}
