package executor

import (
	"bytes"
	"strings"
	"testing"
)

func TestGenerateBillingHeader_UsesCCHPlaceholderWhenSigningEnabled(t *testing.T) {
	header := generateBillingHeader([]byte(`{"messages":[{"role":"user","content":"hi"}]}`), true, "2.1.114", "hello", "cli", "")
	if !strings.Contains(header, "cch="+claudeBillingCCHPlaceholder+";") {
		t.Fatalf("billing header should include cch placeholder, got %q", header)
	}
}

func TestSignAnthropicMessagesBody_DoesNotReplacePlaceholder(t *testing.T) {
	body := []byte(`{"system":[{"type":"text","text":"x-anthropic-billing-header: cc_version=2.1.114.e13; cc_entrypoint=cli; cch=00000;"}],"messages":[{"role":"user","content":"hi"}],"tools":[{"name":"Read","description":"r","input_schema":{}}]}`)

	signed := signAnthropicMessagesBody(append([]byte(nil), body...))
	if !bytes.Equal(signed, body) {
		t.Fatalf("expected signing to leave body unchanged\nwant: %s\n got: %s", body, signed)
	}
}

func TestSignAnthropicMessagesBody_NoPlaceholderNoChange(t *testing.T) {
	body := []byte(`{"messages":[{"role":"user","content":"hi"}]}`)
	signed := signAnthropicMessagesBody(append([]byte(nil), body...))
	if !bytes.Equal(signed, body) {
		t.Fatalf("expected body unchanged when no cch placeholder is present")
	}
}
