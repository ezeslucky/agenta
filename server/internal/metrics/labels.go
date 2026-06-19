package metrics

import (
	"regexp"
	"strings"

	"github.com/ezeslucky/agenta/server/pkg/taskfailure"
)

const (
	labelSource         = "source"
	labelRuntimeMode    = "runtime_mode"
	labelProvider       = "provider"
	labelTerminalStatus = "terminal_status"
	labelFailureReason  = "failure_reason"
	labelTokenType      = "token_type"
	labelModel          = "model"
	labelModelAlias     = "model_alias"

	// PR3 labels (funnel / community / commercial).
	labelSignupSource = "signup_source"
	labelPlatform     = "platform"
	labelPath         = "path"
	labelCadence      = "cadence"
	labelTriggerKind  = "trigger_kind"
	labelReason       = "reason"
	labelRecoverable  = "recoverable"
	labelKind         = "kind"
	labelStatus       = "status"
	labelEventKind    = "event_kind"
	labelAction       = "action"
	labelResult       = "result"
	labelOp           = "op"
)

var businessMetricLabels = map[string][]string{
	"agenta_agent_task_enqueued_total":     {labelSource, labelRuntimeMode},
	"agenta_agent_task_dispatched_total":   {labelSource, labelRuntimeMode},
	"agenta_agent_task_started_total":      {labelSource, labelRuntimeMode, labelProvider},
	"agenta_agent_task_terminal_total":     {labelSource, labelRuntimeMode, labelTerminalStatus},
	"agenta_agent_task_failed_total":       {labelSource, labelRuntimeMode, labelFailureReason},
	"agenta_agent_task_queue_wait_seconds": {labelSource, labelRuntimeMode},
	"agenta_agent_task_run_seconds":        {labelSource, labelRuntimeMode, labelTerminalStatus},
	"agenta_agent_task_total_seconds":      {labelSource, labelRuntimeMode, labelTerminalStatus},
	"agenta_agent_task_in_progress":        {labelSource, labelRuntimeMode},
	"agenta_agent_task_iteration_count":    {labelSource, labelTerminalStatus},
	"agenta_llm_tokens_total":              {labelProvider, labelModel, labelTokenType, labelRuntimeMode, labelSource},
	"agenta_llm_cost_usd_total":            {labelProvider, labelModel, labelTokenType, labelRuntimeMode, labelSource},
	"agenta_llm_unpriced_tokens_total":     {labelProvider, labelModelAlias, labelTokenType},
	"agenta_llm_request_total":             {labelProvider, labelModel, labelRuntimeMode},
	"agenta_task_queued_expired_total":     {labelSource, labelRuntimeMode},
	"agenta_task_lease_expired_total":      {labelSource},

	// PR3 funnel / community / commercial.
	"agenta_signup_total":                             {labelSignupSource},
	"agenta_workspace_created_total":                  {labelSource},
	"agenta_team_invite_sent_total":                   {},
	"agenta_team_invite_accepted_total":               {},
	"agenta_onboarding_started_total":                 {labelPlatform},
	"agenta_onboarding_questionnaire_submitted_total": {},
	"agenta_onboarding_completed_total":               {labelPath},
	"agenta_cloud_waitlist_joined_total":              {},
	"agenta_issue_created_total":                      {labelSource, labelPlatform},
	"agenta_chat_message_sent_total":                  {labelPlatform},
	"agenta_agent_created_total":                      {labelRuntimeMode, labelSource},
	"agenta_squad_created_total":                      {},
	"agenta_autopilot_created_total":                  {labelCadence},
	"agenta_issue_executed_total":                     {labelSource},
	"agenta_runtime_registered_total":                 {labelRuntimeMode, labelProvider},
	"agenta_runtime_ready_total":                      {labelRuntimeMode, labelProvider},
	"agenta_runtime_ready_seconds":                    {labelRuntimeMode, labelProvider},
	"agenta_runtime_failed_total":                     {labelRuntimeMode, labelProvider, labelFailureReason, labelRecoverable},
	"agenta_runtime_offline_total":                    {labelRuntimeMode, labelProvider},
	"agenta_daemon_ws_message_received_total":         {labelKind},
	"agenta_autopilot_run_started_total":              {labelCadence, labelTriggerKind},
	"agenta_autopilot_run_terminal_total":             {labelCadence, labelTriggerKind, labelTerminalStatus},
	"agenta_autopilot_run_skipped_total":              {labelCadence, labelReason},
	"agenta_webhook_delivery_total":                   {labelProvider, labelStatus},
	"agenta_github_event_received_total":              {labelEventKind, labelAction},
	"agenta_github_pr_review_total":                   {labelResult},
	"agenta_cloudruntime_request_total":               {labelOp, labelStatus},
	"agenta_cloudruntime_request_duration_seconds":    {labelOp},
	"agenta_feedback_submitted_total":                 {labelKind, labelPlatform},
	"agenta_contact_sales_submitted_total":            {labelSource},
}

var forbiddenMetricLabels = map[string]struct{}{
	"workspace_id": {},
	"user_id":      {},
	"agent_id":     {},
	"task_id":      {},
	"issue_id":     {},
	"runtime_id":   {},
	"session_id":   {},
	"ip":           {},
}

var (
	knownSources = map[string]string{
		"issue":           "issue",
		"chat":            "chat",
		"autopilot":       "autopilot",
		"autopilot_issue": "autopilot_issue",
		"quick_create":    "quick_create",
		"manual":          "manual",
		"api":             "api",
		"other":           "other",
	}
	knownRuntimeModes = map[string]string{
		"local":   "local",
		"cloud":   "cloud",
		"unknown": "unknown",
	}
	knownRuntimeProviders = map[string]string{
		"antigravity":   "antigravity",
		"claude":        "claude",
		"codebuddy":     "codebuddy",
		"codex":         "codex",
		"copilot":       "copilot",
		"cursor":        "cursor",
		"gemini":        "gemini",
		"hermes":        "hermes",
		"kiro":          "kiro",
		"kimi":          "kimi",
		"agenta_agent": "agenta_agent",
		"openclaw":      "openclaw",
		"opencode":      "opencode",
		"pi":            "pi",
		"other":         "other",
	}
	knownTerminalStatuses = map[string]string{
		"completed": "completed",
		"failed":    "failed",
		"cancelled": "cancelled",
		"blocked":   "blocked",
		"other":     "other",
	}
	knownTokenTypes = map[string]string{
		"input":       "input",
		"output":      "output",
		"cache_read":  "cache_read",
		"cache_write": "cache_write",
	}
	knownFailureReasons = map[string]string{}
	modelAliasUnsafeRe  = regexp.MustCompile(`[^a-z0-9._:/+-]+`)
)

func init() {
	for _, reason := range taskfailure.AllReasons() {
		knownFailureReasons[reason.String()] = reason.String()
	}
}

func validateBusinessMetricLabels() {
	for metric, labels := range businessMetricLabels {
		for _, label := range labels {
			if _, forbidden := forbiddenMetricLabels[label]; forbidden {
				panic("forbidden high-cardinality label " + label + " on " + metric)
			}
		}
	}
}

func metricLabels(metric string) []string {
	labels, ok := businessMetricLabels[metric]
	if !ok {
		panic("missing business metric label definition for " + metric)
	}
	return labels
}

func NormalizeTaskSource(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if normalized, ok := knownSources[value]; ok {
		return normalized
	}
	return "other"
}

func NormalizeRuntimeMode(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if normalized, ok := knownRuntimeModes[value]; ok {
		return normalized
	}
	return "unknown"
}

func NormalizeRuntimeProvider(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if normalized, ok := knownRuntimeProviders[value]; ok {
		return normalized
	}
	return "other"
}

func NormalizeTerminalStatus(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if normalized, ok := knownTerminalStatuses[value]; ok {
		return normalized
	}
	return "other"
}

func NormalizeFailureReason(value string) string {
	value = strings.TrimSpace(value)
	if normalized, ok := knownFailureReasons[value]; ok {
		return normalized
	}
	return taskfailure.Classify(value).String()
}

func NormalizeTokenType(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if normalized, ok := knownTokenTypes[value]; ok {
		return normalized
	}
	return "input"
}

func NormalizeModelAlias(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return "unknown"
	}
	value = modelAliasUnsafeRe.ReplaceAllString(value, "_")
	if len(value) > 128 {
		return value[:128]
	}
	return value
}
