package collectors

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/certusone/solana_exporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

// TokenAccountPubkey struct which contains a
// list of pubkeys
type AccountOwnerPubkey struct {
	Pubkey [][]string `json:"account_owner_pubkey_mint"`
}

type TokenAccountsByOwnerCollector struct {
	RpcClient *rpc.RPCClient

	contextSlot         *prometheus.Desc
	dataProgram         *prometheus.Desc
	accountType         *prometheus.Desc
	tokenAmount         *prometheus.Desc
	tokenAmountDecimals *prometheus.Desc
	tokenUiAmount       *prometheus.Desc
	tokenAmountString   *prometheus.Desc
	isInitialized       *prometheus.Desc
	isNative            *prometheus.Desc
	mintInfo            *prometheus.Desc
	ownerInfo           *prometheus.Desc
	executableProgram   *prometheus.Desc
	accountLamports     *prometheus.Desc
	accountOwner        *prometheus.Desc
	accountRentEpoch    *prometheus.Desc
}

func NewTokenAccountsByOwnerCollector(rpcAddr string) *TokenAccountsByOwnerCollector {
	return &TokenAccountsByOwnerCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"solana_token_accounts_by_owner_context_slot",
			"Token Accounts By Owner - Context Slot",
			[]string{"owner_pubkey", "mint_pubkey", "index"}, nil),
		dataProgram: prometheus.NewDesc(
			"solana_token_accounts_by_owner_program",
			"Token Accounts By Owner - Program",
			[]string{"program", "owner_pubkey", "mint_pubkey", "index"}, nil),
		accountType: prometheus.NewDesc(
			"solana_token_accounts_by_owner_account_type",
			"Token Accounts By Owner - Account Type",
			[]string{"type", "owner_pubkey", "mint_pubkey", "index"}, nil),
		tokenAmount: prometheus.NewDesc(
			"solana_token_accounts_by_owner_token_amount",
			"Token Accounts By Owner - Token Amount",
			[]string{"amount", "owner_pubkey", "mint_pubkey", "index"}, nil),
		tokenAmountDecimals: prometheus.NewDesc(
			"solana_token_accounts_by_owner_token_amount_decimals",
			"Token Accounts By Owner - Token Amount Decimals",
			[]string{"owner_pubkey", "mint_pubkey", "index"}, nil),
		tokenUiAmount: prometheus.NewDesc(
			"solana_token_accounts_by_owner_token_ui_amount",
			"Token Accounts By Owner - Token UI Amount",
			[]string{"owner_pubkey", "mint_pubkey", "index"}, nil),
		tokenAmountString: prometheus.NewDesc(
			"solana_token_accounts_by_owner_token_amount_string",
			"Token Accounts By Owner - Token Amount String",
			[]string{"amount", "owner_pubkey", "mint_pubkey", "index"}, nil),
		isInitialized: prometheus.NewDesc(
			"solana_token_accounts_by_owner_is_initialized",
			"Token Accounts By Owner - Is Initialized",
			[]string{"isInitialized", "owner_pubkey", "mint_pubkey", "index"}, nil),
		isNative: prometheus.NewDesc(
			"solana_token_accounts_by_owner_is_native",
			"Token Accounts By Owner - Is Native",
			[]string{"isNative", "owner_pubkey", "mint_pubkey", "index"}, nil),
		mintInfo: prometheus.NewDesc(
			"solana_token_accounts_by_owner_mint",
			"Token Accounts By Owner - Mint",
			[]string{"mint", "owner_pubkey", "mint_pubkey", "index"}, nil),
		ownerInfo: prometheus.NewDesc(
			"solana_token_accounts_by_owner_info_owner",
			"Token Accounts By Owner - Owner Info",
			[]string{"owner", "owner_pubkey", "mint_pubkey", "index"}, nil),
		executableProgram: prometheus.NewDesc(
			"solana_token_accounts_by_owner_executable",
			"Boolean indicating if the account contains a program (and is strictly read-only)",
			[]string{"executable", "owner_pubkey", "mint_pubkey", "index"}, nil),
		accountLamports: prometheus.NewDesc(
			"solana_token_accounts_by_owner_lamports",
			"Number of lamports assigned to this account, as a u64",
			[]string{"owner_pubkey", "mint_pubkey", "index"}, nil),
		accountOwner: prometheus.NewDesc(
			"solana_token_accounts_by_owner_account_owner",
			"Base-58 encoded Pubkey of the program this account has been assigned to",
			[]string{"owner", "owner_pubkey", "mint_pubkey", "index"}, nil),
		accountRentEpoch: prometheus.NewDesc(
			"solana_token_accounts_by_owner_rent_epoch",
			"The epoch at which this account will next owe rent, as u64",
			[]string{"owner_pubkey", "mint_pubkey", "index"}, nil),
	}
}

func (c *TokenAccountsByOwnerCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *TokenAccountsByOwnerCollector) mustEmitTokenAccountsByOwnerMetrics(ch chan<- prometheus.Metric, response *rpc.TokenAccountsByOwnerInfo, pubkey string, mint string) {
	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(response.Context.Slot), pubkey, mint, "")

	for accountIndex, account := range response.Value {
		accountIndex := strconv.Itoa(accountIndex)
		ch <- prometheus.MustNewConstMetric(c.dataProgram, prometheus.GaugeValue, 0, account.Account.Data.Program, pubkey, mint, accountIndex)
		ch <- prometheus.MustNewConstMetric(c.accountType, prometheus.GaugeValue, 0, account.Account.Data.Parsed.AccountType, pubkey, mint, accountIndex)
		ch <- prometheus.MustNewConstMetric(c.tokenAmount, prometheus.GaugeValue, 0, account.Account.Data.Parsed.Info.TokenAmount.Amount, pubkey, mint, accountIndex)
		ch <- prometheus.MustNewConstMetric(c.tokenAmountDecimals, prometheus.GaugeValue, float64(account.Account.Data.Parsed.Info.TokenAmount.Decimals), pubkey, mint, accountIndex)
		ch <- prometheus.MustNewConstMetric(c.tokenUiAmount, prometheus.GaugeValue, float64(account.Account.Data.Parsed.Info.TokenAmount.UiAmount), pubkey, mint, accountIndex)
		ch <- prometheus.MustNewConstMetric(c.tokenAmountString, prometheus.GaugeValue, 0, account.Account.Data.Parsed.Info.TokenAmount.UiAmountString, pubkey, mint, accountIndex)
		ch <- prometheus.MustNewConstMetric(c.isInitialized, prometheus.GaugeValue, 0, account.Account.Data.Parsed.Info.IsInitialized, pubkey, mint, accountIndex)
		ch <- prometheus.MustNewConstMetric(c.isNative, prometheus.GaugeValue, 0, strconv.FormatBool(account.Account.Data.Parsed.Info.IsNative), pubkey, mint, accountIndex)
		ch <- prometheus.MustNewConstMetric(c.mintInfo, prometheus.GaugeValue, 0, account.Account.Data.Parsed.Info.Mint, pubkey, mint, accountIndex)
		ch <- prometheus.MustNewConstMetric(c.ownerInfo, prometheus.GaugeValue, 0, account.Account.Data.Parsed.Info.Owner, pubkey, mint, accountIndex)
		ch <- prometheus.MustNewConstMetric(c.executableProgram, prometheus.GaugeValue, 0, strconv.FormatBool(account.Account.Executable), pubkey, mint, accountIndex)
		ch <- prometheus.MustNewConstMetric(c.accountLamports, prometheus.GaugeValue, float64(account.Account.Lamports), pubkey, mint, accountIndex)
		ch <- prometheus.MustNewConstMetric(c.accountOwner, prometheus.GaugeValue, 0, account.Account.Owner, pubkey, mint, accountIndex)
		ch <- prometheus.MustNewConstMetric(c.accountRentEpoch, prometheus.GaugeValue, float64(account.Account.RentEpoch), pubkey, mint, accountIndex)
	}
}

func (c *TokenAccountsByOwnerCollector) Collect(ch chan<- prometheus.Metric) {
	jsonData, err := GetKeys()
	if err != nil {
		klog.V(2).Infof("tokenAccountsByOwner response: %v", err)
	}

	var keys AccountOwnerPubkey
	// we unmarshal our jsonData which contains our
	// jsonFile's content into type which we defined above
	if err = json.Unmarshal(jsonData, &keys); err != nil {
		klog.V(2).Infof("failed to decode response body: %w", err)
	}

	for _, pubkeys := range keys.Pubkey {
		ownerPubkey := string(pubkeys[0])
		mintPubkey := string(pubkeys[1])
		ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
		defer cancel()

		accounts, err := c.RpcClient.GetTokenAccountsByOwner(ctx, ownerPubkey, mintPubkey)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(-1), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.dataProgram, prometheus.GaugeValue, 0, err.Error(), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.accountType, prometheus.GaugeValue, 0, err.Error(), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.tokenAmount, prometheus.GaugeValue, 0, err.Error(), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.tokenAmountDecimals, prometheus.GaugeValue, float64(-1), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.tokenUiAmount, prometheus.GaugeValue, float64(-1), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.tokenAmountString, prometheus.GaugeValue, 0, err.Error(), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.isInitialized, prometheus.GaugeValue, 0, err.Error(), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.isNative, prometheus.GaugeValue, 0, err.Error(), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.mintInfo, prometheus.GaugeValue, 0, err.Error(), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.ownerInfo, prometheus.GaugeValue, 0, err.Error(), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.executableProgram, prometheus.GaugeValue, 0, err.Error(), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.accountLamports, prometheus.GaugeValue, float64(-1), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.accountOwner, prometheus.GaugeValue, 0, err.Error(), ownerPubkey, mintPubkey, "")
			ch <- prometheus.MustNewConstMetric(c.accountRentEpoch, prometheus.GaugeValue, float64(-1), ownerPubkey, mintPubkey, "")
		} else {
			c.mustEmitTokenAccountsByOwnerMetrics(ch, accounts, ownerPubkey, mintPubkey)
		}
	}
}
