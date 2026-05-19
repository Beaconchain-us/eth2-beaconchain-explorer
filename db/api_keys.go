package db

import (
	"database/sql"
)

// GetAPIKeyByEmail retrieves an API key for a given email
func GetAPIKeyByEmail(email string) (string, error) {
	var key string
	err := FrontendReaderDB.Get(&key, "SELECT key FROM api_keys WHERE email=$1", email)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return key, err
}

// SaveAPIKey stores or updates an API key for a user
func SaveAPIKey(key, email, plan string, limitIP int) error {
	_, err := FrontendWriterDB.Exec(`
		INSERT INTO api_keys (key, email, plan_name, limit_ip)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO UPDATE
		SET key = EXCLUDED.key, plan_name = EXCLUDED.plan_name, limit_ip = EXCLUDED.limit_ip
	`, key, email, plan, limitIP)
	return err
}

// IsTxProcessed checks if a transaction has already been processed
func IsTxProcessed(txHash, network string) (bool, error) {
	var count int
	err := FrontendReaderDB.Get(&count, "SELECT COUNT(*) FROM processed_txs WHERE tx_hash=$1 AND network=$2", txHash, network)
	return count > 0, err
}

// MarkTxProcessed records a transaction as processed
func MarkTxProcessed(txHash, network string) error {
	_, err := FrontendWriterDB.Exec("INSERT INTO processed_txs (tx_hash, network) VALUES ($1, $2)", txHash, network)
	return err
}