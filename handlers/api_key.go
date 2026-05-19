package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gobitfly/eth2-beaconchain-explorer/db"
	"github.com/gobitfly/eth2-beaconchain-explorer/utils"
)

type ApiKeyRequest struct {
	Email string `json:"email"`
}

// RequestApiKeyHandler handles API key requests for users who have paid
func RequestApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	var req ApiKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Email == "" || !utils.IsValidEmail(req.Email) {
		http.Error(w, "Valid email required", http.StatusBadRequest)
		return
	}

	// Check if API key already exists for this email
	existingKey, err := db.GetAPIKeyByEmail(req.Email)
	if err == nil && existingKey != "" {
		json.NewEncoder(w).Encode(map[string]string{"apiKey": existingKey})
		return
	}

	// Generate new API key
	apiKey, err := utils.GenerateRandomAPIKey()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	// Store in DB (plan can be updated later after payment)
	err = db.SaveAPIKey(apiKey, req.Email, "pending", 5)
	if err != nil {
		http.Error(w, "Failed to save API key", http.StatusInternalServerError)
		return
	}
	// Send email with API key
	_ = utils.SendAPIKeyEmail(req.Email, apiKey, "pending")
	json.NewEncoder(w).Encode(map[string]string{"message": "API key sent to your email"})
}