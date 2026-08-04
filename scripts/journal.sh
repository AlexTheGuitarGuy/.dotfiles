#!/usr/bin/env bash
# Lightweight work-session/todo/blocked log. No sqlite dependency: one active
# session enforced via a state file, entries appended to a JSON-lines log.
set -euo pipefail

STATE_DIR="$HOME/.local/state/journal"
ACTIVE_FILE="$STATE_DIR/active.json"
LOG_FILE="$STATE_DIR/log.jsonl"

mkdir -p "$STATE_DIR"
touch "$LOG_FILE"

now() { date -u +%Y-%m-%dT%H:%M:%SZ; }

append_entry() {
  local type="$1" session="$2" project="$3" text="$4"
  jq -nc --arg ts "$(now)" --arg type "$type" --arg session "$session" \
    --arg project "$project" --arg text "$text" \
    '{timestamp: $ts, type: $type, session: (if $session == "" then null else $session end),
      project: (if $project == "" then null else $project end), text: $text}' \
    >> "$LOG_FILE"
}

active_project() {
  [ -f "$ACTIVE_FILE" ] && jq -r '.project' "$ACTIVE_FILE" || true
}

cmd_start() {
  local project="${1:?usage: journal start <project>}"
  if [ -f "$ACTIVE_FILE" ]; then
    echo "Session already active for '$(active_project)'. Run 'journal stop' first." >&2
    exit 1
  fi
  local session_id
  session_id="$(now)-$project"
  jq -nc --arg project "$project" --arg session "$session_id" --arg started_at "$(now)" \
    '{project: $project, session: $session, started_at: $started_at}' > "$ACTIVE_FILE"
  append_entry "session_start" "$session_id" "$project" "started"
  echo "Started session for '$project'."
}

cmd_stop() {
  if [ ! -f "$ACTIVE_FILE" ]; then
    echo "No active session." >&2
    exit 1
  fi
  local project session
  project="$(jq -r '.project' "$ACTIVE_FILE")"
  session="$(jq -r '.session' "$ACTIVE_FILE")"
  append_entry "session_stop" "$session" "$project" "stopped"
  rm "$ACTIVE_FILE"
  echo "Stopped session for '$project'."
}

cmd_entry() {
  local type="$1" text="${2:?usage: journal $1 <text>}"
  local session="" project=""
  if [ -f "$ACTIVE_FILE" ]; then
    session="$(jq -r '.session' "$ACTIVE_FILE")"
    project="$(jq -r '.project' "$ACTIVE_FILE")"
  fi
  append_entry "$type" "$session" "$project" "$text"
  if [ -n "$project" ]; then
    echo "[$project] $type: $text"
  else
    echo "[inbox] $type: $text"
  fi
}

cmd_log() {
  local filter="${1:-}"
  case "$filter" in
    --session) jq -c 'select(.session != null)' "$LOG_FILE" ;;
    --inbox) jq -c 'select(.session == null)' "$LOG_FILE" ;;
    --today) jq -c --arg today "$(date -u +%Y-%m-%d)" 'select(.timestamp | startswith($today))' "$LOG_FILE" ;;
    "") cat "$LOG_FILE" ;;
    *) echo "Unknown filter: $filter (use --session, --inbox, or --today)" >&2; exit 1 ;;
  esac | jq -r '"\(.timestamp)  \(.type)\t\(.project // "inbox")\t\(.text)"'
}

case "${1:-}" in
  start) shift; cmd_start "$@" ;;
  stop) cmd_stop ;;
  note) shift; cmd_entry "note" "$*" ;;
  todo) shift; cmd_entry "todo" "$*" ;;
  blocked) shift; cmd_entry "blocked" "$*" ;;
  log) shift; cmd_log "${1:-}" ;;
  *)
    echo "Usage: journal <start <project>|stop|note <text>|todo <text>|blocked <text>|log [--session|--inbox|--today]>" >&2
    exit 1
    ;;
esac
