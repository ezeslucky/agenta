import type { AutopilotTrigger } from "../types";


export function buildAutopilotWebhookUrl(params: {
  trigger: Pick<AutopilotTrigger, "kind" | "webhook_token" | "webhook_path" | "webhook_url">;
  apiBaseUrl?: string;
  currentOrigin?: string;
}): string | null {
  const { trigger, apiBaseUrl, currentOrigin } = params;

  if (trigger.kind !== "webhook") return null;

  if (typeof trigger.webhook_url === "string" && trigger.webhook_url) {
    return trigger.webhook_url;
  }

  const path =
    (typeof trigger.webhook_path === "string" && trigger.webhook_path) ||
    (trigger.webhook_token ? `/api/webhooks/autopilots/${trigger.webhook_token}` : null);
  if (!path) return null;

  const base = stripTrailingSlash(apiBaseUrl) || stripTrailingSlash(currentOrigin);
  if (!base) return path; // last resort — relative path will still work in-browser
  return base + path;
}

function stripTrailingSlash(s: string | undefined): string {
  if (!s) return "";
  return s.endsWith("/") ? s.slice(0, -1) : s;
}
