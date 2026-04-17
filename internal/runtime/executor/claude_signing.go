package executor

import (
	"strings"

	"github.com/router-for-me/CLIProxyAPI/v7/internal/config"
	cliproxyauth "github.com/router-for-me/CLIProxyAPI/v7/sdk/cliproxy/auth"
)

func signAnthropicMessagesBody(body []byte) []byte {
	return body
}

func resolveClaudeKeyConfig(cfg *config.Config, auth *cliproxyauth.Auth) *config.ClaudeKey {
	if cfg == nil || auth == nil {
		return nil
	}

	apiKey, baseURL := claudeCreds(auth)
	if apiKey == "" {
		return nil
	}

	for i := range cfg.ClaudeKey {
		entry := &cfg.ClaudeKey[i]
		cfgKey := strings.TrimSpace(entry.APIKey)
		cfgBase := strings.TrimSpace(entry.BaseURL)
		if !strings.EqualFold(cfgKey, apiKey) {
			continue
		}
		if baseURL != "" && cfgBase != "" && !strings.EqualFold(cfgBase, baseURL) {
			continue
		}
		return entry
	}

	return nil
}

// resolveClaudeKeyCloakConfig finds the matching ClaudeKey config and returns its CloakConfig.
func resolveClaudeKeyCloakConfig(cfg *config.Config, auth *cliproxyauth.Auth) *config.CloakConfig {
	entry := resolveClaudeKeyConfig(cfg, auth)
	if entry == nil {
		return nil
	}
	return entry.Cloak
}

func experimentalCCHSigningEnabled(cfg *config.Config, auth *cliproxyauth.Auth) bool {
	entry := resolveClaudeKeyConfig(cfg, auth)
	return entry != nil && entry.ExperimentalCCHSigning
}
