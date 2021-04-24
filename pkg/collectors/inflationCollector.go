package collectors

import (
	"context"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type InflationCollector struct {
	RpcClient *rpc.RPCClient

	totalInflation      *prometheus.Desc
	validatorInflation  *prometheus.Desc
	foundationInflation *prometheus.Desc
	epochInflation      *prometheus.Desc
}

func NewInflationCollector(rpcAddr string) *InflationCollector {
	return &InflationCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		totalInflation: prometheus.NewDesc(
			"solana_total_inflation",
			"Total inflation",
			nil, nil),
		validatorInflation: prometheus.NewDesc(
			"solana_validator_inflation",
			"Inflation allocated to validators",
			nil, nil),
		foundationInflation: prometheus.NewDesc(
			"solana_foundation_inflation",
			"Inflation allocated to the foundation",
			nil, nil),
		epochInflation: prometheus.NewDesc(
			"solana_epoch_inflation",
			"Epoch for which inflation values are valid",
			nil, nil),
	}
}

func (c *InflationCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *InflationCollector) mustEmitInflationMetrics(ch chan<- prometheus.Metric, response *rpc.InflationInfo) {
	ch <- prometheus.MustNewConstMetric(c.totalInflation, prometheus.GaugeValue, response.Total)
	ch <- prometheus.MustNewConstMetric(c.validatorInflation, prometheus.GaugeValue, response.Validator)
	ch <- prometheus.MustNewConstMetric(c.foundationInflation, prometheus.GaugeValue, response.Foundation)
	ch <- prometheus.MustNewConstMetric(c.epochInflation, prometheus.GaugeValue, response.Epoch)
}

func (c *InflationCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
	defer cancel()

	info, err := c.RpcClient.GetInflationRate(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.totalInflation, prometheus.GaugeValue, float64(-1))
		ch <- prometheus.MustNewConstMetric(c.validatorInflation, prometheus.GaugeValue, float64(-1))
		ch <- prometheus.MustNewConstMetric(c.foundationInflation, prometheus.GaugeValue, float64(-1))
		ch <- prometheus.MustNewConstMetric(c.epochInflation, prometheus.GaugeValue, float64(-1))
	} else {
		c.mustEmitInflationMetrics(ch, info)
	}
}
