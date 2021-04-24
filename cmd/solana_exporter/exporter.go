package main

import (
	"collectors"
	"flag"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"k8s.io/klog/v2"
)

const (
	httpTimeout = 5 * time.Second
)

var (
	rpcAddr = flag.String("rpcURI", "", "Solana RPC URI (including protocol and path)")
	addr    = flag.String("addr", ":8080", "Listen address")
)

func init() {
	klog.InitFlags(nil)
}

func main() {
	flag.Parse()

	if *rpcAddr == "" {
		klog.Fatal("Please specify -rpcURI")
	}

	Collector := collectors.NewSolanaCollector(*rpcAddr)
	CollectHealth := collectors.NewHealthCollector(*rpcAddr)
	CollectInflation := collectors.NewInflationCollector(*rpcAddr)
	CollectLargestAccounts := collectors.NewLargestAccountsCollector(*rpcAddr)
	CollectStakeActivation := collectors.NewStakeActivationCollector(*rpcAddr)
	CollectSupply := collectors.NewSupplyCollector(*rpcAddr)
	CollectTokenAccountBalance := collectors.NewTokenAccountBalanceCollector(*rpcAddr)
	CollectTokenAccountsByOwner := collectors.NewTokenAccountsByOwnerCollector(*rpcAddr)
	CollectVersion := collectors.NewVersionCollector(*rpcAddr)
	CollectEpochSchedule := collectors.NewEpochScheduleCollector(*rpcAddr)
	CollectTokenSupply := collectors.NewTokenSupplyCollector(*rpcAddr)
	CollectAccountInfo := collectors.NewAccountInfoCollector(*rpcAddr)
	CollectBalance := collectors.NewBalanceCollector(*rpcAddr)
	CollectRecentBlockhash := collectors.NewRecentBlockhashCollector(*rpcAddr)

	go WatchSlots(Collector)

	prometheus.MustRegister(Collector)
	prometheus.MustRegister(CollectHealth)
	prometheus.MustRegister(CollectInflation)
	prometheus.MustRegister(CollectLargestAccounts)
	prometheus.MustRegister(CollectStakeActivation)
	prometheus.MustRegister(CollectSupply)
	prometheus.MustRegister(CollectTokenAccountBalance)
	prometheus.MustRegister(CollectTokenAccountsByOwner)
	prometheus.MustRegister(CollectVersion)
	prometheus.MustRegister(CollectEpochSchedule)
	prometheus.MustRegister(CollectTokenSupply)
	prometheus.MustRegister(CollectAccountInfo)
	prometheus.MustRegister(CollectBalance)
	prometheus.MustRegister(CollectRecentBlockhash)

	http.Handle("/metrics", promhttp.Handler())

	klog.Infof("listening on %s", *addr)
	klog.Fatal(http.ListenAndServe(*addr, nil))
}
