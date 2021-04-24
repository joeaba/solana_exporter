package collectors

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

// NodeIP struct which contains a
// list of Node IPs
type NodeIP struct {
	IP []string `json:"node_ip"`
}

type HealthCollector struct {
	RpcClient *rpc.RPCClient

	nodeHealth *prometheus.Desc
}

func NewHealthCollector(rpcAddr string) *HealthCollector {
	return &HealthCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		nodeHealth: prometheus.NewDesc(
			"solana_node_health",
			"The current health of the node",
			[]string{"status", "ip"}, nil),
	}
}

func (c *HealthCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *HealthCollector) mustEmitHealthMetrics(ch chan<- prometheus.Metric, status string, IP string) {
	ch <- prometheus.MustNewConstMetric(c.nodeHealth, prometheus.GaugeValue, 0, status, IP)
}

func (c *HealthCollector) Collect(ch chan<- prometheus.Metric) {

	jsonData, err := GetKeys()
	if err != nil {
		klog.V(2).Infof("health response: %v", err)
	}

	var IPs NodeIP
	// we unmarshal our jsonData which contains our
	// jsonFile's content into type which we defined above
	if err = json.Unmarshal(jsonData, &IPs); err != nil {
		klog.V(2).Infof("failed to decode response body: %w", err)
	}

	for _, IP := range IPs.IP {

		match, err := regexp.MatchString(`^[^a-z]`, IP)

		if err != nil {
			c.mustEmitHealthMetrics(ch, err.Error(), IP)
		}

		IP = "http://" + IP
		if match {
			IP = IP + ":8899"
		}

		ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)

		defer cancel()

		status, err := c.RpcClient.GetHealth(ctx, IP)
		if err != nil {
			if strings.Contains(err.Error(), "deadline exceeded") {
				c.mustEmitHealthMetrics(ch, "Node is unhealthy", IP)
			} else {
				c.mustEmitHealthMetrics(ch, err.Error(), IP)
			}
		} else {
			c.mustEmitHealthMetrics(ch, status, IP)
		}
	}
}
