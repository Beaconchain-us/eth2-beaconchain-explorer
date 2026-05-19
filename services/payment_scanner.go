package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gobitfly/eth2-beaconchain-explorer/db"
	"github.com/gobitfly/eth2-beaconchain-explorer/mail"
	"github.com/gobitfly/eth2-beaconchain-explorer/utils"
	"github.com/sirupsen/logrus"
)

// Wallet addresses for payment scanning (same as in api_wallet.go)
var walletAddresses = map[string]string{
	"ethereum": "0x4E94F10F0a34a0DF229e68d5902644046258D678",
	"bnb":      "0x4E94F10F0a34a0DF229e68d5902644046258D678",
	"solana":   "D5WGQzdd6NrAWKfcSCL29DT238SVPkobKrm6VTN3AH9r",
	"bitcoin":  "bc1qqxzu6vvzy55x2n48wpg2drfyu947pfqpcsvjsd",
}

// Pricing plans (in ETH/BNB/SOL/BTC equivalent)
var plans = map[string]float64{
	"basic":      0.05,
	"pro":        0.20,
	"enterprise": 0.50,
}

// StartPaymentScanner starts the periodic payment scanner
func StartPaymentScanner(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	logrus.Info("Payment scanner started, scanning every 10 minutes")
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
}

// ---------- Ethereum (Etherscan) ----------
func scanEthereum(address string) {
	apiKey := utils.Config.EtherscanAPIKey
	if apiKey == "" {
		logrus.Warn("Etherscan API key not set")
		return
	}
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&sort=desc&apikey=%s", address, apiKey)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch Ethereum transactions")
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if status, ok := result["status"].(string); ok && status == "1" {
		if txs, ok := result["result"].([]interface{}); ok {
			for _, tx := range txs {
				processTransaction(tx.(map[string]interface{}), "ethereum")
			}
		}
	}
}

// ---------- BNB Chain (BSCscan) ----------
func scanBNB(address string) {
	apiKey := utils.Config.BSCscanAPIKey
	if apiKey == "" {
		logrus.Warn("BSCscan API key not set")
		return
	}
	url := fmt.Sprintf("https://api.bscscan.com/api?module=account&action=txlist&address=%s&sort=desc&apikey=%s", address, apiKey)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch BNB transactions")
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if status, ok := result["status"].(string); ok && status == "1" {
		if txs, ok := result["result"].([]interface{}); ok {
			for _, tx := range txs {
				processTransaction(tx.(map[string]interface{}), "bnb")
			}
		}
	}
}

// ---------- Solana (Solscan) ----------
func scanSolana(address string) {
	url := fmt.Sprintf("https://public-api.solscan.io/account/transactions?account=%s&limit=50", address)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch Solana transactions")
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var txs []map[string]interface{}
	json.Unmarshal(body, &txs)
	for _, tx := range txs {
		processTransaction(tx, "solana")
	}
}

// ---------- Bitcoin (Blockchain.com) ----------
func scanBitcoin(address string) {
	url := fmt.Sprintf("https://blockchain.info/rawaddr/%s", address)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch Bitcoin transactions")
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if txs, ok := result["txs"].([]interface{}); ok {
		for _, tx := range txs {
			processTransaction(tx.(map[string]interface{}), "bitcoin")
		}
	}
}

// ---------- Transaction processing ----------
func processTransaction(tx map[string]interface{}, network string) {
	txHash, _ := tx["hash"].(string)
	if txHash == "" {
		return
	}
	processed, _ := db.IsTxProcessed(txHash, network)
	if processed {
		return
	}

	var amount float64
	switch network {
	case "ethereum", "bnb":
		valStr, _ := tx["value"].(string)
		if valStr == "" {
			return
		}
		valWei, _ := strconv.ParseFloat(valStr, 64)
		amount = valWei / 1e18
	case "solana":
		if lamports, ok := tx["lamports"].(float64); ok {
			amount = lamports / 1e9
		} else {
			return
		}
	case "bitcoin":
		if valSat, ok := tx["value"].(float64); ok {
			amount = valSat / 1e8
		} else {
			return
		}
	}

	matchedPlan := ""
	for plan, price := range plans {
		if amount >= price-0.001 && amount <= price+0.001 {
			matchedPlan = plan
			break
		}
	}
	if matchedPlan == "" {
		return
	}

	email := extractEmailFromTx(tx)
	if email == "" || !utils.IsValidEmail(email) {
		logrus.Warnf("No valid email in tx %s", txHash)
		return
	}

	apiKey, err := utils.GenerateRandomAPIKey()
	if err != nil {
		logrus.WithError(err).Error("Failed to generate API key")
		return
	}
	if err := db.SaveAPIKey(apiKey, email, matchedPlan, getPlanLimit(matchedPlan)); err != nil {
		logrus.WithError(err).Error("Failed to save API key")
		return
	}
	emailBody := fmt.Sprintf("Your API key: %s\nPlan: %s\nRate limit: %d requests per minute", apiKey, matchedPlan, getPlanLimit(matchedPlan))
	if err := mail.SendTextMail(email, "Your Beaconchain API Key", emailBody, nil); err != nil {
		logrus.WithError(err).Error("Failed to send email")
	} else {
		logrus.Infof("API key sent to %s for plan %s (tx %s)", email, matchedPlan, txHash)
	}
	db.MarkTxProcessed(txHash, network)
}

func extractEmailFromTx(tx map[string]interface{}) string {
	if memo, ok := tx["memo"].(string); ok && strings.Contains(memo, "@") {
		return strings.TrimSpace(memo)
	}
	return ""
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