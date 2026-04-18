package executor

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestEnsureOpusAdaptiveThinking_SetsAdaptiveForOpus(t *testing.T) {
	in := []byte(`{"messages":[{"role":"user","content":"hi"}]}`)
	out := ensureOpusAdaptiveThinking(in, "claude-opus-4-7")
	if got := gjson.GetBytes(out, "thinking.type").String(); got != "adaptive" {
		t.Errorf("thinking.type = %q, want %q", got, "adaptive")
	}
}

func TestEnsureOpusAdaptiveThinking_OverridesDisabled(t *testing.T) {
	in := []byte(`{"thinking":{"type":"disabled"},"messages":[{"role":"user","content":"hi"}]}`)
	out := ensureOpusAdaptiveThinking(in, "claude-opus-4-7")
	if got := gjson.GetBytes(out, "thinking.type").String(); got != "adaptive" {
		t.Errorf("thinking.type = %q, want %q (disabled must be overridden)", got, "adaptive")
	}
}

func TestEnsureOpusAdaptiveThinking_DropsBudgetTokensAndDisplay(t *testing.T) {
	in := []byte(`{"thinking":{"type":"enabled","budget_tokens":16000,"display":"summarized"},"messages":[{"role":"user","content":"hi"}]}`)
	out := ensureOpusAdaptiveThinking(in, "claude-opus-4-7")
	if got := gjson.GetBytes(out, "thinking.type").String(); got != "adaptive" {
		t.Errorf("thinking.type = %q, want adaptive", got)
	}
	if gjson.GetBytes(out, "thinking.budget_tokens").Exists() {
		t.Error("thinking.budget_tokens must be deleted (incompatible with adaptive)")
	}
	if gjson.GetBytes(out, "thinking.display").Exists() {
		t.Error("thinking.display must be deleted (non-official fingerprint)")
	}
}

func TestEnsureOpusAdaptiveThinking_PreservesEffort(t *testing.T) {
	in := []byte(`{"thinking":{"type":"disabled"},"output_config":{"effort":"max"},"messages":[{"role":"user","content":"hi"}]}`)
	out := ensureOpusAdaptiveThinking(in, "claude-opus-4-7")
	if got := gjson.GetBytes(out, "output_config.effort").String(); got != "max" {
		t.Errorf("output_config.effort = %q, want %q (must be preserved)", got, "max")
	}
}

func TestEnsureOpusAdaptiveThinking_NoopForNonOpus(t *testing.T) {
	in := []byte(`{"thinking":{"type":"disabled"},"messages":[{"role":"user","content":"hi"}]}`)
	out := ensureOpusAdaptiveThinking(in, "claude-sonnet-4-5")
	if got := gjson.GetBytes(out, "thinking.type").String(); got != "disabled" {
		t.Errorf("thinking.type = %q, want %q (non-opus must be untouched)", got, "disabled")
	}
}

func TestEnsureOpusAdaptiveThinking_MatchesAnyOpusVersion(t *testing.T) {
	for _, model := range []string{"claude-opus-4-5", "claude-opus-4-6", "claude-opus-4-7", "Claude-Opus-5", "claude-opus-4-7-max"} {
		in := []byte(`{"messages":[{"role":"user","content":"hi"}]}`)
		out := ensureOpusAdaptiveThinking(in, model)
		if got := gjson.GetBytes(out, "thinking.type").String(); got != "adaptive" {
			t.Errorf("model %q: thinking.type = %q, want adaptive", model, got)
		}
	}
}
