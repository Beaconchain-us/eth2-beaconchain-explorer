package handlers

import (
	"encoding/json"
	"net/http"
)

type WalletAddress struct {
	Network string `json:"network"`
	Address string `json:"address"`
}

// ApiWalletAddresses returns the list of wallet addresses for API monetization
// These addresses are used for users to send payments (ETH/BNB, SOL, BTC)
func ApiWalletAddresses(w http.ResponseWriter, r *http.Request) {
	wallets := []WalletAddress{
		{"Ethereum / BNB (BEP20)", "0x4E94F10F0a34a0DF229e68d5902644046258D678"},
		{"Solana (SOL/USDC)", "D5WGQzdd6NrAWKfcSCL29DT238SVPkobKrm6VTN3AH9r"},
		{"Bitcoin (BTC)", "bc1qqxzu6vvzy55x2n48wpg2drfyu947pfqpcsvjsd"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallets)
}