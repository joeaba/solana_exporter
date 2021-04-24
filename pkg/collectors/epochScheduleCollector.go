package collectors

import (
	"context"
	"strconv"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type EpochScheduleCollector struct {
	RpcClient *rpc.RPCClient

	firstNormalEpoch         *prometheus.Desc
	firstNormalSlot          *prometheus.Desc
	leaderScheduleSlotOffset *prometheus.Desc
	slotsPerEpoch            *prometheus.Desc
	epochsWarmup             *prometheus.Desc
}

func NewEpochScheduleCollector(rpcAddr string) *EpochScheduleCollector {
	return &EpochScheduleCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		firstNormalEpoch: prometheus.NewDesc(
			"solana_first_normal_epoch",
			"First normal-length epoch, log2(slotsPerEpoch) - log2(MINIMUM_SLOTS_PER_EPOCH)",
			nil, nil),
		firstNormalSlot: prometheus.NewDesc(
			"solana_first_normal_slot",
			"MINIMUM_SLOTS_PER_EPOCH * (2.pow(firstNormalEpoch) - 1)",
			nil, nil),
		leaderScheduleSlotOffset: prometheus.NewDesc(
			"solana_leader_schedule_slot_offset",
			"The number of slots before beginning of an epoch to calculate a leader schedule for that epoch",
			nil, nil),
		slotsPerEpoch: prometheus.NewDesc(
			"solana_slots_per_epoch",
			"The maximum number of slots in each epoch",
			nil, nil),
		epochsWarmup: prometheus.NewDesc(
			"solana_epoch_schedule_warmup",
			"Whether epochs start short and grow",
			[]string{"warmup"}, nil),
	}
}

func (c *EpochScheduleCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *EpochScheduleCollector) mustEmitEpochScheduleMetrics(ch chan<- prometheus.Metric, response *rpc.EpochScheduleInfo) {
	ch <- prometheus.MustNewConstMetric(c.firstNormalEpoch, prometheus.GaugeValue, float64(response.FirstNormalEpoch))
	ch <- prometheus.MustNewConstMetric(c.firstNormalSlot, prometheus.GaugeValue, float64(response.FirstNormalSlot))
	ch <- prometheus.MustNewConstMetric(c.leaderScheduleSlotOffset, prometheus.GaugeValue, float64(response.LeaderScheduleSlotOffset))
	ch <- prometheus.MustNewConstMetric(c.slotsPerEpoch, prometheus.GaugeValue, float64(response.SlotsPerEpoch))
	ch <- prometheus.MustNewConstMetric(c.epochsWarmup, prometheus.GaugeValue, 0, strconv.FormatBool(response.Warmup))
}

func (c *EpochScheduleCollector) Collect(ch chan<- prometheus.Metric) {

	ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)

	defer cancel()

	schedule, err := c.RpcClient.GetEpochSchedule(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.firstNormalEpoch, prometheus.GaugeValue, float64(-1))
		ch <- prometheus.MustNewConstMetric(c.firstNormalSlot, prometheus.GaugeValue, float64(-1))
		ch <- prometheus.MustNewConstMetric(c.leaderScheduleSlotOffset, prometheus.GaugeValue, float64(-1))
		ch <- prometheus.MustNewConstMetric(c.slotsPerEpoch, prometheus.GaugeValue, float64(-1))
		ch <- prometheus.MustNewConstMetric(c.epochsWarmup, prometheus.GaugeValue, 0, err.Error())
	} else {
		c.mustEmitEpochScheduleMetrics(ch, schedule)
	}
}
