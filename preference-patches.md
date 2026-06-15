# Preference Patches

Custom commits on the `preference` branch, rebased onto official `v7.0.2` (`1fca942b`).

## Status

- **Base**: `v7.0.2` (2026-05-10)
- **Branch**: `origin/preference` (HEAD: `01997bf0`)
- **Deployment**: Currently running official `v7.0.2` (no preference patches)
- **Compact defer_loading**: Official v7.0.2 已自行修复，commit `79c80411` 已从 preference 删除

## Patches by Category

### Claude OAuth Fingerprint (Header Signing)

Align outgoing request headers and signing to match real Claude Code CLI fingerprints, so Anthropic backend treats proxied OAuth requests as genuine Claude Code sessions.

| Commit | Message |
|--------|---------|
| `3b517b73` | fix(claude): align OAuth system prompt and tool mapping with Claude Code v2.1.88 |
| `67166ed3` | fix(claude): update version fingerprint from 2.1.88 to 2.1.108 |
| `b661c90f` | feat(claude): align request fingerprint to Claude Code 2.1.110 |
| `fea6f3b0` | feat(claude): align request fingerprint to Claude Code 2.1.112 and broaden Opus max_tokens rule |
| `8c528e5c` | refactor(claude): emit 3-block interactive system[] matching Claude Code 2.1.112 |
| `fde76c5d` | fix(claude): always emit cch=00000 in billing header matching real Claude Code |
| `8d99a408` | refactor(claude): align Accept-Encoding with interactive capture |
| `a27d0e71` | refactor(claude): align X-Stainless-Timeout default with interactive capture |
| `95620c8f` | refactor(claude): default Claude Code entrypoint to cli for interactive parity |
| `21f85c2b` | fix(claude): align 2.1.114 request fingerprint signing |
| `0b3e2fc7` | fix(claude): default X-Stainless-Timeout to 600 matching Claude Code 2.1.114 |
| `cb090457` | fix(claude): default X-Stainless-Runtime-Version to v24.3.0 matching Claude Code 2.1.114 |

### Claude OAuth Thinking / Reasoning

Ensure thinking/reasoning is visible for non-Claude-Code OAuth clients (OpenCode, etc.) while preserving passthrough for genuine Claude Code.

| Commit | Message |
|--------|---------|
| `dee82966` | refactor(claude): swap advanced-tool-use for redact-thinking beta token |
| `22015952` | fix(claude): force adaptive thinking for Opus on OAuth path |
| `38718176` | fix(claude): force adaptive thinking on all OAuth Claude models |
| `2d5f0a11` | fix(claude): preserve client thinking.display on OAuth path so OpenCode can render summarized reasoning |
| `0101bff4` | fix(claude): drop redact-thinking beta for non-Claude-Code OAuth UAs so OpenCode sees visible reasoning |
| `cdef707d` | fix(claude): preserve Claude Code OAuth passthrough |

### Claude OAuth Tool / Effort / Multi-source Strategy

Three-source OAuth strategy: original Claude Code (passthrough), OpenCode (visible thinking + default effort), OpenClaw (independent defaults).

| Commit | Message |
|--------|---------|
| `aa65a5cf` | feat(claude): add client source detection and OpenClaw-specific OAuth transformations |
| `22dc4391` | feat(claude): default effort to medium for non-OpenCode OAuth clients |
| `9a2c1745` | fix(claude): move ensureDefaultEffort after applyClaudeToolPrefix and add to CountTokens |
| `6530c1ea` | fix(claude): force baseline betas and device profile for ALL OAuth clients |
| `ff8f7880` | fix(claude): strip unmapped tools for ALL OAuth clients including OpenCode |
| `c85e178c` | fix: suppress unused variable isNonOpenCodeOAuth |
| `41a4cb4b` | fix(claude): drop tool_choice when it references an unmapped/stripped tool |
| `b56f15b5` | fix(claude): preserve built-in tool_choice when dropping unmapped tools |
| `ef267835` | fix(claude): remove duplicate tool_choice block and update stale comment |
| `7b1ed703` | feat(claude): force Opus 4.6 max_tokens to 64k (128k at max effort) |

### Codex Compact (REMOVED - v7.0.2 fixed upstream)

Patch `79c80411` 已从 preference 分支删除。如果未来 compact 再次出现 `Deferred tools require tools.tool_search` 错误，按以下方法修复：

**问题**：Codex compact 请求中 tools 数组包含 `defer_loading: true` 的嵌套工具定义（如 `codex_app` / `automation_update` namespace），上游 Codex API 在 compact 路径不支持 tool_search，导致请求被拒。

**修复方法**：在 `internal/runtime/executor/codex_executor.go` 的 `executeCompact` 函数中，`normalizeCodexInstructions(body)` 之后添加：

```go
body = stripCodexToolDeferLoading(body)
```

并添加以下两个辅助函数：

```go
func stripCodexToolDeferLoading(body []byte) []byte {
	return stripCodexToolDeferLoadingAtPath(body, "tools")
}

func stripCodexToolDeferLoadingAtPath(body []byte, path string) []byte {
	tools := gjson.GetBytes(body, path)
	if !tools.IsArray() {
		return body
	}
	for i := range tools.Array() {
		toolPath := fmt.Sprintf("%s.%d", path, i)
		body, _ = sjson.DeleteBytes(body, toolPath+".defer_loading")
		body = stripCodexToolDeferLoadingAtPath(body, toolPath+".tools")
	}
	return body
}
```

**原始 commit**：`79c80411 fix(codex): strip nested deferred tool loading for compact`（已删除）
**PR**：#3167 (https://github.com/router-for-me/CLIProxyAPI/pull/3167)

### Server

| Commit | Message |
|--------|---------|
| `a47d2f1b` | fix(server): exit after printing version |

## How to Rebase onto Newer Official Version

```bash
# 1. Fetch latest upstream
git fetch upstream --tags

# 2. Rebase preference onto new tag (e.g. v7.0.3)
GIT_MASTER=1 git rebase --onto v7.0.3 v7.0.2 HEAD

# 3. Optionally drop obsolete patch
GIT_MASTER=1 git rebase --onto v7.0.3 v7.0.2 HEAD
# then interactive drop of 79c80411 if desired

# 4. Build, verify, deploy
bash ~/.config/opencode/skills/cliproxyapi-claude-oauth/scripts/deploy.sh --auto-version

# 5. Force push
GIT_MASTER=1 git push --force-with-lease origin preference
```

## Relevant Files

- `internal/runtime/executor/claude_executor.go` - Claude OAuth 三来源策略、fingerprint、thinking
- `internal/runtime/executor/claude_executor_test.go` - Claude executor tests
- `internal/runtime/executor/claude_signing.go` - Request fingerprint signing (cch=00000)
- `internal/runtime/executor/helps/claude_device_profile.go` - Device profile defaults
- `internal/runtime/executor/helps/cloak_utils.go` - Cloak/fingerprint utilities
- `internal/runtime/executor/codex_executor.go` - Codex compact defer_loading (obsolete)
- `cmd/server/main.go` - Version exit fix

## PR References

- PR #3167: `fix(codex): strip deferred tool loading from compact requests` (OPEN, obsolete)
  - URL: https://github.com/router-for-me/CLIProxyAPI/pull/3167
  - Branch: `fix/codex-compact-defer-loading`
  - Worktree: `/Users/wykk/Projects/CLIProxyAPI-codex-defer-pr`
