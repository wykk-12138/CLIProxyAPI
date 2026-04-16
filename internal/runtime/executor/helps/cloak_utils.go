package helps

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

// claudeUserIDPayload is the JSON structure for Claude Code 2.1.110+ metadata.user_id.
type claudeUserIDPayload struct {
	DeviceID    string `json:"device_id"`
	AccountUUID string `json:"account_uuid"`
	SessionID   string `json:"session_id"`
}

// generateFakeUserID generates a fake user ID in Claude Code 2.1.110+ JSON format.
// Format: {"device_id":"<64-hex>","account_uuid":"<uuid>","session_id":"<uuid>"}
// If deviceID or accountUUID are provided (from config), they are used instead of random values.
func generateFakeUserID(deviceID, accountUUID string) string {
	if deviceID == "" {
		hexBytes := make([]byte, 32)
		_, _ = rand.Read(hexBytes)
		deviceID = hex.EncodeToString(hexBytes)
	}
	if accountUUID == "" {
		accountUUID = uuid.New().String()
	}
	sessionID := uuid.New().String()

	p := claudeUserIDPayload{
		DeviceID:    deviceID,
		AccountUUID: accountUUID,
		SessionID:   sessionID,
	}
	b, _ := json.Marshal(p)
	return string(b)
}

// isValidUserID checks if a user ID matches Claude Code 2.1.110+ JSON format
// or the legacy string format.
func isValidUserID(userID string) bool {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return false
	}
	// New JSON format: {"device_id":"...","account_uuid":"...","session_id":"..."}
	if strings.HasPrefix(userID, "{") {
		var p claudeUserIDPayload
		if err := json.Unmarshal([]byte(userID), &p); err != nil {
			return false
		}
		return p.DeviceID != "" && p.AccountUUID != "" && p.SessionID != ""
	}
	// Legacy format: user_[64-hex]_account_[uuid]_session_[uuid]
	return strings.HasPrefix(userID, "user_") && strings.Contains(userID, "_account_") && strings.Contains(userID, "_session_")
}

// GenerateFakeUserID generates a fake user ID with random device_id and account_uuid.
func GenerateFakeUserID() string {
	return generateFakeUserID("", "")
}

// GenerateFakeUserIDWithConfig generates a fake user ID, using config values when available.
func GenerateFakeUserIDWithConfig(deviceID, accountUUID string) string {
	return generateFakeUserID(deviceID, accountUUID)
}

func IsValidUserID(userID string) bool {
	return isValidUserID(userID)
}

// ShouldCloak determines if request should be cloaked based on config and client User-Agent.
// Returns true if cloaking should be applied.
func ShouldCloak(cloakMode string, userAgent string) bool {
	switch strings.ToLower(cloakMode) {
	case "always":
		return true
	case "never":
		return false
	default: // "auto" or empty
		// If client is Claude Code, don't cloak
		return !strings.HasPrefix(userAgent, "claude-cli")
	}
}

// isClaudeCodeClient checks if the User-Agent indicates a Claude Code client.
func isClaudeCodeClient(userAgent string) bool {
	return strings.HasPrefix(userAgent, "claude-cli")
}
