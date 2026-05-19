package bootstrap

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gobitfly/eth2-beaconchain-explorer/db"
	"github.com/gobitfly/eth2-beaconchain-explorer/handlers"
	"github.com/gobitfly/eth2-beaconchain-explorer/services"
)

// AutoInit initialises all Code Intersection components automatically
func AutoInit(ctx context.Context, router *gin.Engine) error {
	// 1. Create database tables if missing
	if err := ensureDatabaseTables(); err != nil {
		log.Printf("⚠️ DB init warning: %v", err)
	}

	// 2. Start the payment scanner in background
	go services.StartPaymentScanner(ctx)

	// 3. Register API routes – no need to touch the router manually anymore
	registerAPIRoutes(router)

	log.Println("✅ Code Intersection (BeaconchainHorizon) auto‑initialised")
	return nil
}

func ensureDatabaseTables() error {
	_, err := db.WriterDb.Exec(`
		CREATE TABLE IF NOT EXISTS api_keys (
			key TEXT PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			plan_name TEXT NOT NULL,
			limit_ip INT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS processed_txs (
			tx_hash TEXT PRIMARY KEY,
			network TEXT NOT NULL,
			processed_at TIMESTAMP DEFAULT NOW()
		);
	`)
	return err
}

func registerAPIRoutes(router *gin.Engine) {
	router.GET("/api/wallet-addresses", handlers.ApiWalletAddresses)
	router.POST("/api/request-api-key", handlers.RequestApiKeyHandler)
}