package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gobitfly/eth2-beaconchain-explorer/db"
	"github.com/gobitfly/eth2-beaconchain-explorer/mail"
	"github.com/sirupsen/logrus"
)

// ---------- Configuration ----------
var (
	// Scan interval in minutes
	scanInterval = 10 * time.Minute

	// Wallet addresses per network
	walletAddresses = map[string]string{
		"ethereum": "0x4E94F10F0a34a0DF229e68d5902644046258D678",
		"bnb":      "0x4E94F10F0a34a0DF229e68d5902644046258D678",
		"solana":   "D5WGQzdd6NrAWKfcSCL29DT238SVPkobKrm6VTN3AH9r",
		"bitcoin":  "bc1qqxzu6vvzy55x2n48wpg2drfyu947pfqpcsvjsd",
		"tron":     "TPYGjDE7qfRkuU4h4DGincc8GsCUPK5Cw9",
	}

	// API keys (should be loaded from env/config)
	etherscanAPIKey = "YOUR_ETHERSCAN_API_KEY"
	bscscanAPIKey   = "YOUR_BSCSCAN_API_KEY"
)

// Pricing plans in USD or token units (simplified – you can adjust)
var plans = map[string]float64{
	"basic":      0.05,
	"pro":        0.20,
	"enterprise": 0.50,
}

// ---------- Main scanner loop ----------
func StartPaymentScanner(ctx context.Context) {
	ticker := time.NewTicker(scanInterval)
	logrus.Info("Payment scanner started, scanning every ", scanInterval)
	for {
		select {
		case <-ticker.C:
			scanAllWallets()
		case <-ctx.Done():
			logrus.Info("Payment scanner stopped")
			return
		}
	}
}

func scanAllWallets() {
	logrus.Info("Scanning wallets for new payments...")
	scanEthereum(walletAddresses["ethereum"])
	scanBNB(walletAddresses["bnb"])
	scanSolana(walletAddresses["solana"])
	scanBitcoin(walletAddresses["bitcoin"])
	scanTron(walletAddresses["tron"])
}

// ---------- Ethereum (Etherscan) ----------
func scanEthereum(address string) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&sort=desc&apikey=%s", address, etherscanAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch Ethereum transactions")
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		logrus.WithError(err).Error("Failed to parse Ethereum response")
		return
	}
	if result["status"] == "1" {
		if txs, ok := result["result"].([]interface{}); ok {
			for _, tx := range txs {
				txMap := tx.(map[string]interface{})
				processTransaction(txMap, "ethereum")
			}
		}
	}
}

// ---------- BNB Chain (BSCscan) ----------
func scanBNB(address string) {
	url := fmt.Sprintf("https://api.bscscan.com/api?module=account&action=txlist&address=%s&sort=desc&apikey=%s", address, bscscanAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch BNB transactions")
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		logrus.WithError(err).Error("Failed to parse BNB response")
		return
	}
	if result["status"] == "1" {
		if txs, ok := result["result"].([]interface{}); ok {
			for _, tx := range txs {
				txMap := tx.(map[string]interface{})
				processTransaction(txMap, "bnb")
			}
		}
	}
}

// ---------- Solana (Solscan) ----------
func scanSolana(address string) {
	url := fmt.Sprintf("https://public-api.solscan.io/account/transactions?account=%s&limit=50", address)
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch Solana transactions")
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var txs []map[string]interface{}
	if err := json.Unmarshal(body, &txs); err != nil {
		logrus.WithError(err).Error("Failed to parse Solana response")
		return
	}
	for _, tx := range txs {
		processTransaction(tx, "solana")
	}
}

// ---------- Bitcoin (Blockchain.com) ----------
func scanBitcoin(address string) {
	url := fmt.Sprintf("https://blockchain.info/rawaddr/%s", address)
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch Bitcoin transactions")
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		logrus.WithError(err).Error("Failed to parse Bitcoin response")
		return
	}
	if txs, ok := result["txs"].([]interface{}); ok {
		for _, tx := range txs {
			txMap := tx.(map[string]interface{})
			processTransaction(txMap, "bitcoin")
		}
	}
}

// ---------- Tron (TronGrid) ----------
func scanTron(address string) {
	url := fmt.Sprintf("https://api.trongrid.io/v1/accounts/%s/transactions?limit=50", address)
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch Tron transactions")
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		logrus.WithError(err).Error("Failed to parse Tron response")
		return
	}
	if data, ok := result["data"].([]interface{}); ok {
		for _, tx := range data {
			txMap := tx.(map[string]interface{})
			processTransaction(txMap, "tron")
		}
	}
}

// ---------- Transaction processing ----------
func processTransaction(tx map[string]interface{}, network string) {
	// 1. Extract transaction hash
	txHash, ok := tx["hash"].(string)
	if !ok {
		logrus.Warn("Transaction hash missing, skipping")
		return
	}

	// 2. Check if already processed
	processed, err := db.IsTxProcessed(txHash, network)
	if err != nil {
		logrus.WithError(err).Error("Failed to check processed transaction")
		return
	}
	if processed {
		logrus.Debugf("Transaction %s already processed", txHash)
		return
	}

	// 3. Extract amount based on network
	var amount float64
	switch network {
	case "ethereum", "bnb":
		valueStr, ok := tx["value"].(string)
		if !ok {
			return
		}
		valueWei, _ := strconv.ParseFloat(valueStr, 64)
		amount = valueWei / 1e18 // convert wei to ETH/BNB
	case "solana":
		// Lamports to SOL (1 SOL = 1e9 lamports)
		if lamports, ok := tx["lamports"].(float64); ok {
			amount = lamports / 1e9
		}
	case "bitcoin":
		// BTC value is in satoshis in "value" field
		if valueSat, ok := tx["value"].(float64); ok {
			amount = valueSat / 1e8
		}
	case "tron":
		// TRX value is in SUN (1 TRX = 1e6 SUN)
		if valueSun, ok := tx["value"].(float64); ok {
			amount = valueSun / 1e6
		}
	}

	// 4. Match with pricing plan
	matchedPlan := matchPlan(amount)
	if matchedPlan == "" {
		logrus.Debugf("Transaction %s amount %.6f does not match any plan", txHash, amount)
		return
	}

	// 5. Extract email from memo / data field
	email := extractEmailFromTx(tx)
	if email == "" {
		logrus.Warnf("No email found in transaction %s", txHash)
		return
	}

	// 6. Generate API key
	apiKey, err := generateAPIKey()
	if err != nil {
		logrus.WithError(err).Error("Failed to generate API key")
		return
	}

	// 7. Store API key in DB (associate with email and plan)
	if err := db.SaveAPIKey(apiKey, email, matchedPlan, getPlanLimit(matchedPlan)); err != nil {
		logrus.WithError(err).Error("Failed to save API key")
		return
	}

	// 8. Send email with API key
	if err := mail.SendTextMail(email, "Your Beaconchain API Key", fmt.Sprintf("Your API key: %s\nPlan: %s\nRate limit: %d requests per minute", apiKey, matchedPlan, getPlanLimit(matchedPlan)), nil); err != nil {
		logrus.WithError(err).Error("Failed to send email")
	} else {
		logrus.Infof("API key sent to %s for plan %s (tx %s)", email, matchedPlan, txHash)
	}

	// 9. Mark transaction as processed
	if err := db.MarkTxProcessed(txHash, network); err != nil {
		logrus.WithError(err).Error("Failed to mark transaction processed")
	}
}

// ---------- Helper functions ----------
func matchPlan(amount float64) string {
	for plan, price := range plans {
		// Allow small floating point differences (0.001)
		if amount >= price-0.001 && amount <= price+0.001 {
			return plan
		}
	}
	return ""
}

func extractEmailFromTx(tx map[string]interface{}) string {
	// Look for "email" field in memo, input, data, or contractData
	if memo, ok := tx["memo"].(string); ok && strings.Contains(memo, "@") {
		return memo
	}
	if input, ok := tx["input"].(string); ok && strings.Contains(input, "@") {
		return input
	}
	if data, ok := tx["data"].(string); ok && strings.Contains(data, "@") {
		return data
	}
	if contractData, ok := tx["contractData"].(map[string]interface{}); ok {
		if em, ok := contractData["email"].(string); ok {
			return em
		}
	}
	return ""
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getPlanLimit(plan string) int {
	switch plan {
	case "basic":
		return 5
	case "pro":
		return 20
	case "enterprise":
		return 100
	default:
		return 5
	}
}