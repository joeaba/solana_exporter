package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	
	"strconv"
	"time"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"k8s.io/klog/v2"
)

const (
	httpTimeout = 5 * time.Second
)

type Pubkey struct {
	value string
}

var (
	rpcAddr = flag.String("rpcURI", "", "Solana RPC URI (including protocol and path)")
	addr    = flag.String("addr", ":8080", "Listen address")
	pubkeys []Pubkey
)

func init() {
	klog.InitFlags(nil)
}

type solanaCollector struct {
	rpcClient *rpc.RPCClient

	totalValidatorsDesc     *prometheus.Desc
	validatorActivatedStake *prometheus.Desc
	validatorLastVote       *prometheus.Desc
	validatorRootSlot       *prometheus.Desc
	validatorDelinquent     *prometheus.Desc
}

type accountCollector struct {
	rpcClient   *rpc.RPCClient
	contextSlot *prometheus.Desc
	value       *prometheus.Desc
	//addressAcc  *prometheus.Desc
}

type supplyCollector struct {
	rpcClient *rpc.RPCClient

	contextSlot            *prometheus.Desc
	totalSupply            *prometheus.Desc
	circulatingSupply      *prometheus.Desc
	nonCirculatingSupply   *prometheus.Desc
	nonCirculatingAccounts *prometheus.Desc
}
type balanceCollector struct {
	rpcClient   *rpc.RPCClient
	contextSlot *prometheus.Desc
	value       *prometheus.Desc
}

type stakeactivationCollector struct {
	rpcClient *rpc.RPCClient
	active    *prometheus.Desc
	inactive  *prometheus.Desc
	state     *prometheus.Desc
}

type tokenaccountbyownerCollector struct {
	rpcClient *rpc.RPCClient

	contextSlot     *prometheus.Desc
	program         *prometheus.Desc
	accountType     *prometheus.Desc
	amount          *prometheus.Desc
	decimals        *prometheus.Desc
	uiAmount        *prometheus.Desc
	uiAmountString  *prometheus.Desc
	delegate        *prometheus.Desc
	delegatedAmount *prometheus.Desc
	isInitialized   *prometheus.Desc
	isNative        *prometheus.Desc
	mint            *prometheus.Desc
	ownerInfo       *prometheus.Desc
	executable      *prometheus.Desc
	lamports        *prometheus.Desc
	owner           *prometheus.Desc
	rentEpoch       *prometheus.Desc
}

type tokensupplyCollector struct {
	rpcClient      *rpc.RPCClient
	contextSlot    *prometheus.Desc
	amount         *prometheus.Desc
	decimals       *prometheus.Desc
	uiAmount       *prometheus.Desc
	uiAmountString *prometheus.Desc
}

type accountinfobase64Collector struct {
	rpcClient   *rpc.RPCClient
	contextSlot *prometheus.Desc
	data        *prometheus.Desc
	executable  *prometheus.Desc
	lamports    *prometheus.Desc
	owner       *prometheus.Desc
	rentEpoch   *prometheus.Desc
}

type getaccountinfojsonparsedCollector struct {
	rpcClient            *rpc.RPCClient
	contextSlot          *prometheus.Desc
	authority            *prometheus.Desc
	blockhash            *prometheus.Desc
	executable           *prometheus.Desc
	lamportsPerSignature *prometheus.Desc
	lamports             *prometheus.Desc
	owner                *prometheus.Desc
	rentEpoch            *prometheus.Desc
	//addressAcc  *prometheus.Desc
}

func NewSolanaCollector(rpcAddr string) *solanaCollector {
	return &solanaCollector{
		rpcClient: rpc.NewRPCClient(rpcAddr),
		totalValidatorsDesc: prometheus.NewDesc(
			"solana_active_validators",
			"Total number of active validators by state",
			[]string{"state"}, nil),
		validatorActivatedStake: prometheus.NewDesc(
			"solana_validator_activated_stake",
			"Activated stake per validator",
			[]string{"pubkey", "nodekey"}, nil),
		validatorLastVote: prometheus.NewDesc(
			"solana_validator_last_vote",
			"Last voted slot per validator",
			[]string{"pubkey", "nodekey"}, nil),
		validatorRootSlot: prometheus.NewDesc(
			"solana_validator_root_slot",
			"Root slot per validator",
			[]string{"pubkey", "nodekey"}, nil),
		validatorDelinquent: prometheus.NewDesc(
			"solana_validator_delinquent",
			"Whether a validator is delinquent",
			[]string{"pubkey", "nodekey"}, nil),
	}
}

func NewAccCollector(rpcAddr string) *accountCollector {
	return &accountCollector{
		rpcClient: rpc.NewRPCClient(rpcAddr),
		contextSlot: prometheus.NewDesc(
			"solana_acc_context_slot",
			"Total number of Solana Context Slot",
			nil, nil),

		value: prometheus.NewDesc(
			"solana_value",
			"Activated Solana Value",
			[]string{"lamports"}, nil),

		// lamportsAcc: prometheus.NewDesc(
		// 	"solana_value",
		// 	"Activated Solana Value",
		// 	[]string{"lamports"}, nil),

		// addressAcc: prometheus.NewDesc(
		// 	"solana_value",
		// 	"Activated Solana Value",
		// 	[]string{"address"}, nil),
	}
}

func NewSupplyCollector(rpcAddr string) *supplyCollector {
	return &supplyCollector{
		rpcClient: rpc.NewRPCClient(rpcAddr),
		contextSlot: prometheus.NewDesc(
			"solana_context_slot",
			"Current slot in context",
			nil, nil),
		totalSupply: prometheus.NewDesc(
			"solana_total_supply",
			"Total supply in lamports",
			nil, nil),
		circulatingSupply: prometheus.NewDesc(
			"solana_circulating_supply",
			"Circulating supply in lamports",
			nil, nil),
		nonCirculatingSupply: prometheus.NewDesc(
			"solana_non_circulating_supply",
			"Non-circulating supply in lamports",
			nil, nil),
		nonCirculatingAccounts: prometheus.NewDesc(
			"solana_non_circulating_accounts",
			"Account addresses of non-circulating accounts",
			[]string{"account_address"}, nil),
	}
}

func NewBalanceCollector(rpcAddr string) *balanceCollector {
	return &balanceCollector{
		rpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"balancecollector_context_slot",
			"balancecollector in context",
			[]string{"pubkey"}, nil),
		value: prometheus.NewDesc(
			"balancecollector_balance",
			"balance",
			[]string{"pubkey"}, nil),
	}
}

func NewTokenAccByOwnerCollector(rpcAddr string) *tokenaccountbyownerCollector {
	return &tokenaccountbyownerCollector{
		rpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"Context_Slot_Acc_Owner",
			"Context Slot",
			[]string{"pubkey"}, nil),

		program: prometheus.NewDesc(
			"Program_Type_Token_Owner",
			"Program Type",
			[]string{"pubkey", "program"}, nil),

		accountType: prometheus.NewDesc(
			"Acc_Type_Token_Owner",
			"Account Type",
			[]string{"pubkey", "accountType"}, nil),

		amount: prometheus.NewDesc(
			"Token_Amount",
			"Token amount",
			[]string{"pubkey"}, nil),

		decimals: prometheus.NewDesc(
			"Decimals_Token_Amount",
			"Decimal for Token Amount",
			[]string{"pubkey"}, nil),

		uiAmount: prometheus.NewDesc(
			"Token_UiAmount",
			"UiAmount for Token Amount",
			[]string{"pubkey"}, nil),

		uiAmountString: prometheus.NewDesc(
			"Token_UiAmountString",
			"UiAmountString for Token Amount",
			[]string{"pubkey"}, nil),

		delegate: prometheus.NewDesc(
			"Delegate_Info_Acc_Owner",
			"Delegate info for acc owner",
			[]string{"pubkey", "delegate"}, nil),

		delegatedAmount: prometheus.NewDesc(
			"Delegated_Amount_Acc_Owner",
			"Delegated Amount for Acc",
			[]string{"pubkey"}, nil),

		isInitialized: prometheus.NewDesc(
			"Acc_Is_Initialized",
			"Is Token Account initialized",
			[]string{"pubkey", "isInitialized"}, nil),

		isNative: prometheus.NewDesc(
			"Acc_Is_Native",
			"Is Token Account native",
			[]string{"pubkey", "isNative"}, nil),

		mint: prometheus.NewDesc(
			"Token_Account_Mint",
			"Mint token for Account",
			[]string{"pubkey", "mint"}, nil),

		ownerInfo: prometheus.NewDesc(
			"Token_Account_Owner_Info",
			"Owner Info for Token Account",
			[]string{"pubkey", "ownerInfo"}, nil),

		executable: prometheus.NewDesc(
			"Executable_Token_Account",
			"Token Account is Executable",
			[]string{"pubkey", "executbale"}, nil),

		lamports: prometheus.NewDesc(
			"Token_Account_Lamports",
			"Lamports for token account",
			[]string{"pubkey"}, nil),

		owner: prometheus.NewDesc(
			"Owner_Token_Account",
			"Owner of Token Account",
			[]string{"pubkey", "owner"}, nil),

		rentEpoch: prometheus.NewDesc(
			"Token_Account_Rent_Epoch",
			"Rent epoch for Token Account",
			[]string{"pubkey"}, nil),
	}
}

func NewTokenSuppyCollector(rpcAddr string) *tokensupplyCollector {
	return &tokensupplyCollector{
		rpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"Token_Supply_context_slot",
			"Total Supply For Context Slot",
			nil, nil),

		amount: prometheus.NewDesc(
			"Token_Supply_Amount",
			"Supply Amount For Token Supply",
			[]string{"pubkey", "amount"}, nil),

		decimals: prometheus.NewDesc(
			"Token_Supply_Decimals",
			"Decimals For Token Supply",
			[]string{"pubkey"}, nil),

		uiAmount: prometheus.NewDesc(
			"Token_Supply_UiAmount",
			"Supply Amount For UiAmount",
			[]string{"pubkey"}, nil),

		uiAmountString: prometheus.NewDesc(
			"Token_Supply_UiAmountString",
			"Supply Amount For UiAmount String",
			[]string{"pubkey", "uiAmounts"}, nil),
	}
}

func NewStakeActivationCollector(rpcAddr string) *stakeactivationCollector {
	return &stakeactivationCollector{

		rpcClient: rpc.NewRPCClient(rpcAddr),

		active: prometheus.NewDesc(
			"Active_Stake",
			"Current Active stake",
			[]string{"pubkey"}, nil),
		inactive: prometheus.NewDesc(
			"Get_In_Active_Stake",
			"Current Inactive Stake",
			[]string{"pubkey"}, nil),
		state: prometheus.NewDesc(
			"Get_State",
			"Current State",
			[]string{"pubkey", "state"}, nil),
	}
}

func NewAccountInfoCollector(rpcAddr string) *accountinfobase64Collector {
	return &accountinfobase64Collector{
		rpcClient: rpc.NewRPCClient(rpcAddr),
		contextSlot: prometheus.NewDesc(
			"solana_account_information_context_slot",
			"account Context Slot",
			[]string{"pubkey"}, nil),

		data: prometheus.NewDesc(
			"Account_Info_Data",
			"Data for Account Info",
			[]string{"pubkey", "accountinfodatastring"}, nil),

		executable: prometheus.NewDesc(
			"Executable_Info_Data",
			"AccInfo Data for Executable",
			[]string{"pubkey", "executable"}, nil),

		lamports: prometheus.NewDesc(
			"Acc_Info_Lamports",
			"Lamports Data For Accounts",
			[]string{"pubkey"}, nil),

		owner: prometheus.NewDesc(
			"Acc_Info_Owner",
			"Owner for account info",
			[]string{"pubkey", "owner"}, nil),

		rentEpoch: prometheus.NewDesc(
			"Acc_Info_RentEpoch",
			"Rent epoch for account info",
			[]string{"pubkey"}, nil),
	}
}

func NewAccountInfoJsonParsedCollector(rpcAddr string) *getaccountinfojsonparsedCollector {
	return &getaccountinfojsonparsedCollector{
		rpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"Account_Info_Json_Parsed_Context",
			"Acc Info Json Parsed Info",
			[]string{"pubkey"}, nil),

		authority: prometheus.NewDesc(
			"Account_Info_Json_Authority",
			"authority for accifojson",
			[]string{"pubkey", "authority"}, nil),

		blockhash: prometheus.NewDesc(
			"Account_Info_Json_Block_Hash",
			"BlockHash for acc info",
			[]string{"pubkey", "blockhash"}, nil),

		executable: prometheus.NewDesc(
			"Account_Info_Json_Executable",
			"Executable for acc info",
			[]string{"pubkey", "executable"}, nil),

		lamportsPerSignature: prometheus.NewDesc(
			"Account_Info_Json_lamportspersignature",
			"lamportsPerSignature for acc info",
			[]string{"pubkey"}, nil),

		lamports: prometheus.NewDesc(
			"Account_Info_Json_lamports",
			"lamports for acc info",
			[]string{"pubkey"}, nil),

		owner: prometheus.NewDesc(
			"Account_Info_Json_owner",
			"owner for acc info",
			[]string{"pubkey", "owner"}, nil),

		rentEpoch: prometheus.NewDesc(
			"Account_Info_Json_rentepoch",
			"rentEpoch for acc info",
			[]string{"pubkey"}, nil),

		// lamportsAcc: prometheus.NewDesc(
		// 	"solana_value",
		// 	"Activated Solana Value",
		// 	[]string{"lamports"}, nil),

		// addressAcc: prometheus.NewDesc(
		// 	"solana_value",
		// 	"Activated Solana Value",
		// 	[]string{"address"}, nil),
	}
}

func (c *solanaCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalValidatorsDesc
}

func (c *supplyCollector) Describe(ch chan<- *prometheus.Desc) {
	//ch <- c.contextSlot
}
func (c *accountCollector) Describe(ch chan<- *prometheus.Desc) {
	//ch <- c.contextSlot
}

func (c *stakeactivationCollector) Describe(ch chan<- *prometheus.Desc) {
	//ch <- c.contextSlot
}

func (c *balanceCollector) Describe(ch chan<- *prometheus.Desc) {
	//ch <- c.contextSlot
}

func (c *tokenaccountbyownerCollector) Describe(ch chan<- *prometheus.Desc) {
	//ch <- c.contextSlot
}

func (c *tokensupplyCollector) Describe(ch chan<- *prometheus.Desc) {
	//ch <- c.contextSlot
}
func (c *accountinfobase64Collector) Describe(ch chan<- *prometheus.Desc) {
	//ch <- c.contextSlot
}
func (c *getaccountinfojsonparsedCollector) Describe(ch chan<- *prometheus.Desc) {
	//ch <- c.contextSlot
}

func (c *solanaCollector) mustEmitMetrics(ch chan<- prometheus.Metric, response *rpc.GetVoteAccountsResponse) {
	ch <- prometheus.MustNewConstMetric(c.totalValidatorsDesc, prometheus.GaugeValue,
		float64(len(response.Result.Delinquent)), "delinquent")
	ch <- prometheus.MustNewConstMetric(c.totalValidatorsDesc, prometheus.GaugeValue,
		float64(len(response.Result.Current)), "current")

	for _, account := range append(response.Result.Current, response.Result.Delinquent...) {
		ch <- prometheus.MustNewConstMetric(c.validatorActivatedStake, prometheus.GaugeValue,
			float64(account.ActivatedStake), account.VotePubkey, account.NodePubkey)
		ch <- prometheus.MustNewConstMetric(c.validatorLastVote, prometheus.GaugeValue,
			float64(account.LastVote), account.VotePubkey, account.NodePubkey)
		ch <- prometheus.MustNewConstMetric(c.validatorRootSlot, prometheus.GaugeValue,
			float64(account.RootSlot), account.VotePubkey, account.NodePubkey)
		pubkeys = append(pubkeys, Pubkey{value: account.VotePubkey})
	}
	fmt.Println(pubkeys)
	for _, account := range response.Result.Current {
		ch <- prometheus.MustNewConstMetric(c.validatorDelinquent, prometheus.GaugeValue,
			0, account.VotePubkey, account.NodePubkey)
	}
	for _, account := range response.Result.Delinquent {
		ch <- prometheus.MustNewConstMetric(c.validatorDelinquent, prometheus.GaugeValue,
			1, account.VotePubkey, account.NodePubkey)
	}
}

func (c *supplyCollector) mustSupplyMetrics(ch chan<- prometheus.Metric, response *rpc.GetSupplyResponse) {

	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue,
		float64(response.Result.ContextSlot))
	ch <- prometheus.MustNewConstMetric(c.totalSupply, prometheus.GaugeValue,
		float64(response.Result.Value.TotalSupply))
	ch <- prometheus.MustNewConstMetric(c.circulatingSupply, prometheus.GaugeValue,
		float64(response.Result.Value.CirculatingSupply))
	ch <- prometheus.MustNewConstMetric(c.nonCirculatingSupply, prometheus.GaugeValue,
		float64(response.Result.Value.NonCirculatingSupply))

	for _, account := range response.Result.Value.NonCirculatingAccounts {
		ch <- prometheus.MustNewConstMetric(c.nonCirculatingAccounts, prometheus.GaugeValue,
			0, account)
	}
	// ch <- prometheus.MustNewConstMetric(c.nonCirculatingAccounts, prometheus.GaugeValue,
	// 	0, response.Result.Value.NonCirculatingAccounts...)
}

func (c *accountCollector) mustAccountMetrics(ch chan<- prometheus.Metric, response *rpc.GetLargestAccountsResponse) {

	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue,
		float64(response.Result.ContextSlot))

	for _, account := range response.Result.Value {
		ch <- prometheus.MustNewConstMetric(c.value, prometheus.GaugeValue,
			float64(account.Lamports), account.Address)
	}

	// for _, account := range response.Result.Value {
	// 	ch <- prometheus.MustNewConstMetric(c.addressAcc, prometheus.GaugeValue,
	// 		0, account.Address)
	// }
}

func (c *balanceCollector) mustBalanceMetrics(ch chan<- prometheus.Metric, response *rpc.GetBalanceResponse, pubkey string) {

	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue,
		float64(response.Result.ContextSlot.Slot), pubkey)
	ch <- prometheus.MustNewConstMetric(c.value, prometheus.GaugeValue,
		float64(response.Result.Value), pubkey)

}

func (c *tokensupplyCollector) mustTokenSupplyMetrics(ch chan<- prometheus.Metric, response *rpc.GetTokenSupplyResponse, pubkey string) {

	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue,
		float64(response.Result.ContextSlot))

	for _, account := range response.Result.Value {
		ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue,
			0, account.Amount)
		ch <- prometheus.MustNewConstMetric(c.amount, prometheus.GaugeValue,
			0, account.Amount)
		ch <- prometheus.MustNewConstMetric(c.decimals, prometheus.GaugeValue,
			float64(account.Decimals))
		ch <- prometheus.MustNewConstMetric(c.uiAmount, prometheus.GaugeValue,
			float64(account.UiAmount))
		ch <- prometheus.MustNewConstMetric(c.uiAmountString, prometheus.GaugeValue,
			0, account.UiAmountString)
	}
}

func (c *tokenaccountbyownerCollector) mustTokenAccByOwnerMetrics(ch chan<- prometheus.Metric, response *rpc.GetTokenAccountsbyownerRes, pubkey string) {
	for _, account := range response.Result.Value {
		amount, _ := strconv.ParseFloat(account.Data.Parsed.Info.TokenAmount.Amount, 64)
		uiAmountString, _ := strconv.ParseFloat(account.Data.Parsed.Info.TokenAmount.UiAmountString, 64)
		isInitialized := strconv.FormatBool(account.Data.Parsed.Info.IsInitialized)
		isNative := strconv.FormatBool(account.Data.Parsed.Info.IsNative)
		executable := strconv.FormatBool(account.Executable)
		ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue,
			float64(response.Result.Context.Slot), pubkey)
		ch <- prometheus.MustNewConstMetric(c.program, prometheus.GaugeValue,
			0, pubkey, account.Data.Program)
		ch <- prometheus.MustNewConstMetric(c.accountType, prometheus.GaugeValue,
			0, pubkey, account.Data.Parsed.AccountType)
		ch <- prometheus.MustNewConstMetric(c.amount, prometheus.GaugeValue,
			amount, pubkey)
		ch <- prometheus.MustNewConstMetric(c.decimals, prometheus.GaugeValue,
			float64(account.Data.Parsed.Info.TokenAmount.Decimals), pubkey)
		ch <- prometheus.MustNewConstMetric(c.uiAmount, prometheus.GaugeValue,
			float64(account.Data.Parsed.Info.TokenAmount.UiAmount), pubkey)
		ch <- prometheus.MustNewConstMetric(c.uiAmountString, prometheus.GaugeValue,
			uiAmountString, pubkey)
		ch <- prometheus.MustNewConstMetric(c.delegate, prometheus.GaugeValue,
			0, pubkey, account.Data.Parsed.Info.Delegate)
		ch <- prometheus.MustNewConstMetric(c.delegatedAmount, prometheus.GaugeValue,
			float64(account.Data.Parsed.Info.DelegatedAmount), pubkey)
		ch <- prometheus.MustNewConstMetric(c.isInitialized, prometheus.GaugeValue,
			0, pubkey, isInitialized)
		ch <- prometheus.MustNewConstMetric(c.isNative, prometheus.GaugeValue,
			0, pubkey, isNative)
		ch <- prometheus.MustNewConstMetric(c.mint, prometheus.GaugeValue,
			0, pubkey, account.Data.Parsed.Info.Mint)
		ch <- prometheus.MustNewConstMetric(c.ownerInfo, prometheus.GaugeValue,
			0, pubkey, account.Data.Parsed.Info.Owner)
		ch <- prometheus.MustNewConstMetric(c.executable, prometheus.GaugeValue,
			0, pubkey, executable)
		ch <- prometheus.MustNewConstMetric(c.lamports, prometheus.GaugeValue,
			float64(account.Lamports), pubkey)
		ch <- prometheus.MustNewConstMetric(c.owner, prometheus.GaugeValue,
			0, pubkey, account.Owner)
		ch <- prometheus.MustNewConstMetric(c.rentEpoch, prometheus.GaugeValue,
			float64(account.RentEpoch), pubkey)
	}
}

func (c *accountinfobase64Collector) mustAccountInfo64Metric(ch chan<- prometheus.Metric, response *rpc.GetAccountInfoBase64Res, pubkey string) {

	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue,
		float64(response.Result.ContextSlot.Slot), pubkey)

	ch <- prometheus.MustNewConstMetric(c.executable, prometheus.GaugeValue,
		0, pubkey, strconv.FormatBool(response.Result.Value.Executable))

	ch <- prometheus.MustNewConstMetric(c.lamports, prometheus.GaugeValue,
		float64(response.Result.Value.Lamports), pubkey)

	ch <- prometheus.MustNewConstMetric(c.owner, prometheus.GaugeValue,
		0, pubkey, response.Result.Value.Owner)

	ch <- prometheus.MustNewConstMetric(c.rentEpoch, prometheus.GaugeValue,
		float64(response.Result.Value.RentEpoch), pubkey)

	for _, account := range response.Result.Value.Data {
		ch <- prometheus.MustNewConstMetric(c.data, prometheus.GaugeValue,
			0, pubkey, account)
	}
	// ch <- prometheus.MustNewConstMetric(c.nonCirculatingAccounts, prometheus.GaugeValue,
	// 	0, response.Result.Value.NonCirculatingAccounts...)
}

func (c *getaccountinfojsonparsedCollector) mustAccountInfoJsonParsedCollector(ch chan<- prometheus.Metric, response *rpc.GetAccountInfoJsonParsedRes, pubkey string) {

	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue,
		float64(response.Result.ContextSlot.Slot), pubkey)
	ch <- prometheus.MustNewConstMetric(c.blockhash, prometheus.GaugeValue,
		0, pubkey, response.Result.Value.Data.Nonce.Initialized.Blockhash)
	ch <- prometheus.MustNewConstMetric(c.authority, prometheus.GaugeValue,
		0, pubkey, response.Result.Value.Data.Nonce.Initialized.Authority)
	ch <- prometheus.MustNewConstMetric(c.lamportsPerSignature, prometheus.GaugeValue,
		float64(response.Result.Value.Data.Nonce.Initialized.FeeCalculator.LamportsPerSignature), pubkey)
	ch <- prometheus.MustNewConstMetric(c.executable, prometheus.GaugeValue,
		0, pubkey, strconv.FormatBool(response.Result.Value.Executable))
	ch <- prometheus.MustNewConstMetric(c.lamports, prometheus.GaugeValue,
		float64(response.Result.Value.Lamports), pubkey)
	ch <- prometheus.MustNewConstMetric(c.owner, prometheus.GaugeValue,
		0, response.Result.Value.Owner, pubkey)
	ch <- prometheus.MustNewConstMetric(c.rentEpoch, prometheus.GaugeValue,
		float64(response.Result.Value.RentEpoch), pubkey)
}

func (c *solanaCollector) Collect(ch chan<- prometheus.Metric) {

	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	accs, err := c.rpcClient.GetVoteAccounts(ctx, rpc.CommitmentRecent)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.totalValidatorsDesc, err)
		ch <- prometheus.NewInvalidMetric(c.validatorActivatedStake, err)
		ch <- prometheus.NewInvalidMetric(c.validatorLastVote, err)
		ch <- prometheus.NewInvalidMetric(c.validatorRootSlot, err)
		ch <- prometheus.NewInvalidMetric(c.validatorDelinquent, err)
	} else {
		c.mustEmitMetrics(ch, accs)
	}
}

func (c *supplyCollector) Collect(ch chan<- prometheus.Metric) {

	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	resp, err := c.rpcClient.GetSupply(ctx)
	klog.Infof("Get Supply value is: %v", resp)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.contextSlot, err)
		ch <- prometheus.NewInvalidMetric(c.totalSupply, err)
		ch <- prometheus.NewInvalidMetric(c.circulatingSupply, err)
		ch <- prometheus.NewInvalidMetric(c.nonCirculatingSupply, err)
		ch <- prometheus.NewInvalidMetric(c.nonCirculatingAccounts, err)
	} else {
		c.mustSupplyMetrics(ch, resp)
	}
}

func (c *accountCollector) Collect(ch chan<- prometheus.Metric) {

	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	accountCollector, err := c.rpcClient.GetLargestAcc(ctx)
	klog.Infof("Get account is: %v", accountCollector)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.contextSlot, err)
		ch <- prometheus.NewInvalidMetric(c.value, err)
		//ch <- prometheus.NewInvalidMetric(c.lamportsAcc, err)
		//ch <- prometheus.NewInvalidMetric(c.addressAcc, err)
	} else {
		c.mustAccountMetrics(ch, accountCollector)
	}
}

func (c *balanceCollector) Collect(ch chan<- prometheus.Metric) {

	// var myarr [4]string
	// myarr[0] = "83astBRguLMdt2h5U1Tpdq5tjFoJ6noeGwaY3mDLVcri"
	// myarr[1] = "vines1vzrYbzLMRdu58ou5XTby4qAqVRLmqo36NKPTg"
	// myarr[2] = "4fYNw3dojWmQ4dXtSGE9epjRGy9pFSx62YypT7avPYvA"
	// myarr[3] = "6H94zdiaYfRfPfKjYLjyr2VFBg6JHXygy84r3qhc3NsC"

	var ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	response, _ := c.rpcClient.GetVoteAccounts(ctx, rpc.CommitmentRecent)

	for index, account := range append(response.Result.Current, response.Result.Delinquent...) {
		// pubkeys = append(pubkeys, Pubkey{value: account.VotePubkey})
		// }
		// fmt.Println(pubkeys)
		// fmt.Println(len(pubkeys))

		//Get Balance Response

		// for i := 0; i < len(pubkeys); i++ {
		if index >= 100 {
			break
		}
		var ctx, cancel = context.WithTimeout(context.Background(), httpTimeout)
		defer cancel()
		balanceCollector, err := c.rpcClient.GetBalance(ctx, account.VotePubkey)

		klog.Infof("Get Balance value is: %v", balanceCollector.Result.ContextSlot)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(c.contextSlot, err)
			ch <- prometheus.NewInvalidMetric(c.value, err)
		} else {
			c.mustBalanceMetrics(ch, balanceCollector, account.VotePubkey)
		}
	}
}

func (c *stakeactivationCollector) mustStakeActivationMetrics(ch chan<- prometheus.Metric, response *rpc.GetStackActivationResponse, pubkey string) {

	ch <- prometheus.MustNewConstMetric(c.active, prometheus.GaugeValue,
		float64(response.Result.Active), pubkey)
	ch <- prometheus.MustNewConstMetric(c.inactive, prometheus.GaugeValue,
		float64(response.Result.Inactive), pubkey)
	ch <- prometheus.MustNewConstMetric(c.state, prometheus.GaugeValue,
		0, pubkey, response.Result.State)

}

func (c *stakeactivationCollector) Collect(ch chan<- prometheus.Metric) {

	var myarr [2]string
	myarr[0] = "AjuuY2XHwQoSufRW9ttiGGhnp6R5CxMKmhvEQpTWYjq3"
	myarr[1] = "5Lpj3wJ34StS2Fd4AeR3fHRaA12JCZcwXSCn4KiKsZxh"

	for i := 0; i < 2; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
		defer cancel()
		stakeactivationCollector, err := c.rpcClient.GetStackActivation(ctx, myarr[i])
		klog.Infof("Get StackActivation Detail is: %v", stakeactivationCollector)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(c.active, err)
			ch <- prometheus.NewInvalidMetric(c.inactive, err)
			ch <- prometheus.NewInvalidMetric(c.state, err)
		} else {
			c.mustStakeActivationMetrics(ch, stakeactivationCollector, myarr[i])
		}
	}
}

func (c *tokenaccountbyownerCollector) Collect(ch chan<- prometheus.Metric) {

	var pubkeys [1]string
	pubkeys[0] = "122FAHxVFQDQjzgSBSNUmLJXJxG4ooUwLdYvgf3ASs2K"

	var mints [1]string
	mints[0] = "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R"

	for i := 0; i < len(pubkeys); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
		defer cancel()
		tokenaccsbyownerResp, err := c.rpcClient.GetTokenAccountOwner(ctx, pubkeys[i], mints[i])
		klog.Infof("Get Account By Owner is: %v", tokenaccsbyownerResp)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(c.contextSlot, err)
			ch <- prometheus.NewInvalidMetric(c.program, err)
			ch <- prometheus.NewInvalidMetric(c.accountType, err)
			ch <- prometheus.NewInvalidMetric(c.amount, err)
			ch <- prometheus.NewInvalidMetric(c.decimals, err)
			ch <- prometheus.NewInvalidMetric(c.uiAmount, err)
			ch <- prometheus.NewInvalidMetric(c.uiAmountString, err)
			ch <- prometheus.NewInvalidMetric(c.delegate, err)
			ch <- prometheus.NewInvalidMetric(c.delegatedAmount, err)
			ch <- prometheus.NewInvalidMetric(c.isInitialized, err)
			ch <- prometheus.NewInvalidMetric(c.isNative, err)
			ch <- prometheus.NewInvalidMetric(c.mint, err)
			ch <- prometheus.NewInvalidMetric(c.ownerInfo, err)
			ch <- prometheus.NewInvalidMetric(c.executable, err)
			ch <- prometheus.NewInvalidMetric(c.lamports, err)
			ch <- prometheus.NewInvalidMetric(c.owner, err)
			ch <- prometheus.NewInvalidMetric(c.rentEpoch, err)

		} else {
			c.mustTokenAccByOwnerMetrics(ch, tokenaccsbyownerResp, pubkeys[i])
		}
	}
}

func (c *tokensupplyCollector) Collect(ch chan<- prometheus.Metric) {

	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	tokensupplyres, err := c.rpcClient.GetTokenSupply(ctx, "JCHsvHwF6TgeM1fapxgAkhVKDU5QtPox3bfCR5sjWirP")
	klog.Infof("Token Supply value is: %v", tokensupplyres)

	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.contextSlot, err)
		ch <- prometheus.NewInvalidMetric(c.amount, err)
		ch <- prometheus.NewInvalidMetric(c.decimals, err)
		ch <- prometheus.NewInvalidMetric(c.uiAmount, err)
		ch <- prometheus.NewInvalidMetric(c.uiAmountString, err)
	} else {
		c.mustTokenSupplyMetrics(ch, tokensupplyres, "JCHsvHwF6TgeM1fapxgAkhVKDU5QtPox3bfCR5sjWirP")
	}
}

func (c *accountinfobase64Collector) Collect(ch chan<- prometheus.Metric) {

	var myarr [2]string
	myarr[0] = "AjuuY2XHwQoSufRW9ttiGGhnp6R5CxMKmhvEQpTWYjq3"
	myarr[1] = "5Lpj3wJ34StS2Fd4AeR3fHRaA12JCZcwXSCn4KiKsZxh"

	for i := 0; i < len(myarr); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
		defer cancel()
		accountinfobase64res, err := c.rpcClient.GetAccountInfoBase64(ctx, myarr[i])
		klog.Infof("Get Account Info Base64 Detail is: %v", accountinfobase64res)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(c.contextSlot, err)
			ch <- prometheus.NewInvalidMetric(c.data, err)
			ch <- prometheus.NewInvalidMetric(c.executable, err)
			ch <- prometheus.NewInvalidMetric(c.lamports, err)
			ch <- prometheus.NewInvalidMetric(c.owner, err)
			ch <- prometheus.NewInvalidMetric(c.rentEpoch, err)
		} else {
			c.mustAccountInfo64Metric(ch, accountinfobase64res, myarr[i])
		}
	}
}

func (c *getaccountinfojsonparsedCollector) Collect(ch chan<- prometheus.Metric) {

	var myarr [2]string
	myarr[0] = "AjuuY2XHwQoSufRW9ttiGGhnp6R5CxMKmhvEQpTWYjq3"
	myarr[1] = "5Lpj3wJ34StS2Fd4AeR3fHRaA12JCZcwXSCn4KiKsZxh"

	for i := 0; i < len(myarr); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
		defer cancel()
		accjsonparseres, err := c.rpcClient.GetAccountInfoJsonParsed(ctx, myarr[i])
		klog.Infof("Get Account Info Json Parsed Detail is: %v", accjsonparseres)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(c.contextSlot, err)
			ch <- prometheus.NewInvalidMetric(c.authority, err)
			ch <- prometheus.NewInvalidMetric(c.blockhash, err)
			ch <- prometheus.NewInvalidMetric(c.executable, err)
			ch <- prometheus.NewInvalidMetric(c.lamportsPerSignature, err)
			ch <- prometheus.NewInvalidMetric(c.lamports, err)
			ch <- prometheus.NewInvalidMetric(c.owner, err)
			ch <- prometheus.NewInvalidMetric(c.rentEpoch, err)
		} else {
			c.mustAccountInfoJsonParsedCollector(ch, accjsonparseres, myarr[i])
		}
	}
}

func main() {
	flag.Parse()

	if *rpcAddr == "" {
		klog.Fatal("Please specify -rpcURI")
	}

	collector := NewSolanaCollector(*rpcAddr)
	sCollector := NewSupplyCollector(*rpcAddr)
	accountCollector := NewAccCollector(*rpcAddr)
	balanceCollector := NewBalanceCollector(*rpcAddr)
	tokenaccbyownerCollector := NewTokenAccByOwnerCollector(*rpcAddr)
	tokensupplyCollector := NewTokenSuppyCollector(*rpcAddr)
	stakeactivationCollector := NewStakeActivationCollector(*rpcAddr)
	accountinfobase64Collector := NewAccountInfoCollector(*rpcAddr)
	accountinfojsonparsedCollector := NewAccountInfoJsonParsedCollector(*rpcAddr)

	go collector.WatchSlots()

	prometheus.MustRegister(collector)
	prometheus.MustRegister(sCollector)
	prometheus.MustRegister(accountCollector)
	prometheus.MustRegister(balanceCollector)
	prometheus.MustRegister(tokenaccbyownerCollector)
	prometheus.MustRegister(tokensupplyCollector)
	prometheus.MustRegister(stakeactivationCollector)
	prometheus.MustRegister(accountinfobase64Collector)
	prometheus.MustRegister(accountinfojsonparsedCollector)

	http.Handle("/metrics", promhttp.Handler())

	klog.Infof("listening on %s", *addr)
	klog.Fatal(http.ListenAndServe(*addr, nil))
}
