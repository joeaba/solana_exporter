module github.com/certusone/solana_exporter

go 1.13

require (
	github.com/prometheus/client_golang v1.4.0
	k8s.io/klog/v2 v2.4.0
)

replace github.com/certusone/solana_exporter/cmd/solana_exporter/collector => ./cmd/solana_exporter/collector/
