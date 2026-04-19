package executor

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

func TestGenerateBillingHeader_UsesCCHPlaceholderWhenSigningEnabled(t *testing.T) {
	header := generateBillingHeader([]byte(`{"messages":[{"role":"user","content":"hi"}]}`), true, "2.1.114", "hello", "cli", "")
	if !strings.Contains(header, "cch="+claudeBillingCCHPlaceholder+";") {
		t.Fatalf("billing header should include cch placeholder, got %q", header)
	}
}

func TestSignAnthropicMessagesBody_ReplacesPlaceholderFromFinalSerializedBody(t *testing.T) {
	body := []byte(`{"system":[{"type":"text","text":"x-anthropic-billing-header: cc_version=2.1.114.e13; cc_entrypoint=cli; cch=00000;"}],"messages":[{"role":"user","content":"hi"}],"tools":[{"name":"Read","description":"r","input_schema":{}}]}`)

	h := sha256.Sum256(body)
	wantCCH := hex.EncodeToString(h[:])[:5]
	wantBody := bytes.Replace(body, []byte("cch=00000;"), []byte("cch="+wantCCH+";"), 1)

	signed := signAnthropicMessagesBody(append([]byte(nil), body...))
	if !bytes.Equal(signed, wantBody) {
		t.Fatalf("signed body mismatch\nwant: %s\n got: %s", wantBody, signed)
	}
}

func TestSignAnthropicMessagesBody_NoPlaceholderNoChange(t *testing.T) {
	body := []byte(`{"messages":[{"role":"user","content":"hi"}]}`)
	signed := signAnthropicMessagesBody(append([]byte(nil), body...))
	if !bytes.Equal(signed, body) {
		t.Fatalf("expected body unchanged when no cch placeholder is present")
	}
}
