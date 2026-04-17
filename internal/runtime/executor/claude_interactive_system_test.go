package executor

import (
	"os"
	"strings"
	"testing"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/runtime/executor/helps"
	"github.com/tidwall/gjson"
)

func TestCheckSystemInstructions_InteractiveShape(t *testing.T) {
	payload := []byte(`{"model":"claude-opus-4-7","system":[{"type":"text","text":"You are Claude Code"}],"messages":[{"role":"user","content":"hi"}],"tools":[{"name":"TaskCreate","description":"x","input_schema":{}},{"name":"Read","description":"x","input_schema":{}}]}`)
	out := checkSystemInstructions(payload)

	if got := gjson.GetBytes(out, "system.#").Int(); got != 3 {
		t.Fatalf("expected 3 system blocks, got %d", got)
	}

	if !strings.HasPrefix(gjson.GetBytes(out, "system.0.text").String(), "x-anthropic-billing-header: cc_version=") {
		t.Error("system[0] must start with billing header")
	}
	if gjson.GetBytes(out, "system.0.cache_control").Exists() {
		t.Error("system[0] must not have cache_control")
	}
	if !strings.Contains(gjson.GetBytes(out, "system.0.text").String(), "cc_entrypoint=cli") {
		t.Error("system[0] must contain cc_entrypoint=cli")
	}

	if got := gjson.GetBytes(out, "system.1.text").String(); got != "You are Claude Code, Anthropic's official CLI for Claude." {
		t.Errorf("system[1].text mismatch: %q", got)
	}
	if got := gjson.GetBytes(out, "system.1.cache_control.type").String(); got != "ephemeral" {
		t.Errorf("system[1].cache_control.type = %q, want ephemeral", got)
	}
	if got := gjson.GetBytes(out, "system.1.cache_control.ttl").String(); got != "1h" {
		t.Errorf("system[1].cache_control.ttl = %q, want 1h", got)
	}
	if gjson.GetBytes(out, "system.1.cache_control.scope").Exists() {
		t.Error("system[1].cache_control must not have scope field")
	}

	golden, err := os.ReadFile("testdata/claude_interactive_system2_golden.txt")
	if err != nil {
		t.Fatal(err)
	}
	if got := gjson.GetBytes(out, "system.2.text").String(); got != string(golden) {
		t.Errorf("system[2].text does not match golden fixture\ngolden length=%d, got length=%d", len(golden), len(got))
	}
	if gjson.GetBytes(out, "system.2.cache_control.scope").Exists() {
		t.Error("system[2].cache_control must not have scope field")
	}
	if got := gjson.GetBytes(out, "system.2.cache_control.ttl").String(); got != "1h" {
		t.Errorf("system[2].cache_control.ttl = %q, want 1h", got)
	}
}

func TestCheckSystemInstructions_InteractiveShape_TodoWriteSubstitution(t *testing.T) {
	payload := []byte(`{"model":"claude-opus-4-7","messages":[{"role":"user","content":"hi"}],"tools":[{"name":"TodoWrite","description":"x","input_schema":{}},{"name":"Read","description":"x","input_schema":{}}]}`)
	out := checkSystemInstructions(payload)

	if got := gjson.GetBytes(out, "system.#").Int(); got != 3 {
		t.Fatalf("expected 3 system blocks, got %d", got)
	}

	expected := helps.ClaudeCodeIntro +
		helps.ClaudeCodeSystem +
		helps.ClaudeCodeDoingTasks +
		helps.ClaudeCodeActions +
		helps.BuildUsingToolsSection("TodoWrite") +
		helps.ClaudeCodeToneAndStyle +
		helps.ClaudeCodeTextOutput

	if got := gjson.GetBytes(out, "system.2.text").String(); got != expected {
		t.Fatalf("system[2].text mismatch for TodoWrite substitution\nwant length=%d got length=%d", len(expected), len(got))
	}
}
