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

// Wallet addresses per network (TRON removed)
var walletAddresses = map[string]string{
	"ethereum": "0x4E94F10F0a34a0DF229e68d5902644046258D678",
	"bnb":      "0x4E94F10F0a34a0DF229e68d5902644046258D678",
	"solana":   "D5WGQzdd6NrAWKfcSCL29DT238SVPkobKrm6VTN3AH9r",
	"bitcoin":  "bc1qqxzu6vvzy55x2n48wpg2drfyu947pfqpcsvjsd",
}

// Pricing plans
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
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&sort=desc&apikey=%s", address, apiKey)

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to create request for Ethereum")
		return
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch Ethereum transactions")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to read Ethereum response")
		return
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		logrus.WithError(err).Error("Failed to parse Ethereum response")
		return
	}
	if status, ok := result["status"].(string); ok && status == "1" {
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
	apiKey := utils.Config.BSCscanAPIKey
	url := fmt.Sprintf("https://api.bscscan.com/api?module=account&action=txlist&address=%s&sort=desc&apikey=%s", address, apiKey)

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to create request for BNB")
		return
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch BNB transactions")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to read BNB response")
		return
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		logrus.WithError(err).Error("Failed to parse BNB response")
		return
	}
	if status, ok := result["status"].(string); ok && status == "1" {
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
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to create request for Solana")
		return
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch Solana transactions")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to read Solana response")
		return
	}
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
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to create request for Bitcoin")
		return
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch Bitcoin transactions")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to read Bitcoin response")
		return
	}
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

	// 3. Extract amount
	var amount float64
	switch network {
	case "ethereum", "bnb":
		valueStr, ok := tx["value"].(string)
		if !ok {
			logrus.Warnf("No value field in tx %s", txHash)
			return
		}
		valueWei, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			logrus.Warnf("Invalid value format in tx %s", txHash)
			return
		}
		amount = valueWei / 1e18
	case "solana":
		if lamports, ok := tx["lamports"].(float64); ok {
			amount = lamports / 1e9
		} else {
			logrus.Warnf("No lamports field in Solana tx %s", txHash)
			return
		}
	case "bitcoin":
		if valueSat, ok := tx["value"].(float64); ok {
			amount = valueSat / 1e8
		} else {
			logrus.Warnf("No value field in Bitcoin tx %s", txHash)
			return
		}
	default:
		logrus.Warnf("Unknown network %s for tx %s", network, txHash)
		return
	}

	// 4. Match with pricing plan
	matchedPlan := matchPlan(amount)
	if matchedPlan == "" {
		logrus.Debugf("Transaction %s amount %.6f does not match any plan", txHash, amount)
		return
	}

	// 5. Extract email from memo / data field
	email := extractEmailFromTx(tx)
	if email == "" || !utils.IsValidEmail(email) {
		logrus.Warnf("No valid email found in transaction %s", txHash)
		return
	}

	// 6. Generate API key using utils function
	apiKey, err := utils.GenerateRandomAPIKey()
	if err != nil {
		logrus.WithError(err).Error("Failed to generate API key")
		return
	}

	// 7. Store API key in DB
	if err := db.SaveAPIKey(apiKey, email, matchedPlan, getPlanLimit(matchedPlan)); err != nil {
		logrus.WithError(err).Error("Failed to save API key")
		return
	}

	// 8. Send email to user
	emailBody := fmt.Sprintf("Your API key: %s\nPlan: %s\nRate limit: %d requests per minute\n\nThank you for supporting Beaconchain!", apiKey, matchedPlan, getPlanLimit(matchedPlan))
	if err := mail.SendTextMail(email, "Your Beaconchain API Key", emailBody, nil); err != nil {
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
		if amount >= price-0.001 && amount <= price+0.001 {
			return plan
		}
	}
	return ""
}

func extractEmailFromTx(tx map[string]interface{}) string {
	// Check memo field
	if memo, ok := tx["memo"].(string); ok && strings.Contains(memo, "@") {
		return strings.TrimSpace(memo)
	}
	// Check input data (Ethereum style)
	if input, ok := tx["input"].(string); ok && strings.Contains(input, "@") {
		// Try to find email pattern in hex data (simplified)
		return extractEmailFromHex(input)
	}
	// Check data field
	if data, ok := tx["data"].(string); ok && strings.Contains(data, "@") {
		return extractEmailFromHex(data)
	}
	// Check contractData for TRON or other
	if contractData, ok := tx["contractData"].(map[string]interface{}); ok {
		if em, ok := contractData["email"].(string); ok && strings.Contains(em, "@") {
			return em
		}
	}
	return ""
}

func extractEmailFromHex(hexStr string) string {
	// Very basic: look for ASCII email pattern in hex string
	// This is a placeholder – in real implementation you'd decode hex and search
	// For now, just return empty (or implement as needed)
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