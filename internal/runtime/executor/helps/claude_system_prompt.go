package helps

import (
	_ "embed"

	"github.com/tidwall/gjson"
)

// Claude Code interactive-mode system prompt sections (captured from Claude Code v2.1.112).
// These sections are injected into system[] for OAuth cloaking and must keep byte fidelity.

//go:embed claude_interactive_intro.txt
var ClaudeCodeIntro string

//go:embed claude_interactive_system.txt
var ClaudeCodeSystem string

//go:embed claude_interactive_doing_tasks.txt
var ClaudeCodeDoingTasks string

//go:embed claude_interactive_actions.txt
var ClaudeCodeActions string

//go:embed claude_interactive_tone_and_style.txt
var ClaudeCodeToneAndStyle string

//go:embed claude_interactive_text_output.txt
var ClaudeCodeTextOutput string

// BuildUsingToolsSection dynamically builds the interactive "Using your tools" section.
// taskToolName should be the resolved task management tool name (e.g. "TodoWrite" or "TaskCreate"),
// or empty if neither is present.
func BuildUsingToolsSection(taskToolName string) string {
	result := "# Using your tools\n" +
		" - Prefer dedicated tools over Bash when one fits (Read, Edit, Write, Glob, Grep) — reserve Bash for shell-only operations.\n"

	if taskToolName != "" {
		result += " - Use " + taskToolName + " to plan and track work. Mark each task completed as soon as it's done; don't batch.\n"
	}

	result += " - You can call multiple tools in a single response. If you intend to call multiple tools and there are no dependencies between them, make all independent tool calls in parallel. Maximize use of parallel tool calls where possible to increase efficiency. However, if some tool calls depend on previous calls to inform dependent values, do NOT call these tools in parallel and instead call them sequentially. For instance, if one operation must complete before another starts, run these operations sequentially instead.\n\n"

	return result
}

// ResolveTaskToolName checks which task management tool is present in the request.
// Returns the TitleCase tool name to use in the system prompt, or empty string if none.
//
// IMPORTANT: This is called BEFORE tool name remapping (remapOAuthToolNames), so it
// must accept both pre-remap (lowercase) and post-remap (TitleCase) names.
// The returned value is always TitleCase for use in the system prompt text.
func ResolveTaskToolName(payload []byte) string {
	tools := gjson.GetBytes(payload, "tools")
	if !tools.Exists() || !tools.IsArray() {
		return ""
	}
	hasTodoWrite := false
	hasTaskCreate := false
	tools.ForEach(func(_, tool gjson.Result) bool {
		name := tool.Get("name").String()
		switch name {
		case "todowrite", "TodoWrite":
			hasTodoWrite = true
		case "taskcreate", "TaskCreate":
			hasTaskCreate = true
		}
		return true
	})
	if hasTodoWrite {
		return "TodoWrite"
	}
	if hasTaskCreate {
		return "TaskCreate"
	}
	return ""
}
