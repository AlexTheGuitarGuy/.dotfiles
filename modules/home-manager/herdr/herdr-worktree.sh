trap 'read -rp "herdr-worktree failed, press enter"' ERR

template="$(herdr plugin config-dir sessionizer)/config.json"
ws=$(herdr pane list | jq -r '[.result.panes[] | select(.focused)][0].workspace_id')

pick=$(herdr worktree list --workspace "$ws" \
  | jq -r '.result.worktrees[] | select(.is_bare | not) | "\(.path)\t\(.open_workspace_id // "")"' \
  | fzf --prompt 'worktree> ' --delimiter '\t' --with-nth 1) || exit 0
dir=${pick%%$'\t'*}
open=${pick#*$'\t'}

if [ -n "$open" ]; then
  exec herdr workspace focus "$open" >/dev/null
fi

opened=$(herdr worktree open --workspace "$ws" --path "$dir" --focus)
new_ws=$(jq -r '.result.workspace.workspace_id' <<<"$opened")
first_tab=$(jq -r '.result.tab.tab_id' <<<"$opened")
pane=$(jq -r '.result.root_pane.pane_id' <<<"$opened")

i=0
while IFS=$'\t' read -r label cmd; do
  if [ "$i" -gt 0 ]; then
    pane=$(herdr tab create --workspace "$new_ws" --cwd "$dir" --label "$label" --no-focus | jq -r '.result.root_pane.pane_id')
  else
    herdr tab rename "$first_tab" "$label" >/dev/null
  fi
  if [ -n "$cmd" ]; then herdr pane run "$pane" "$cmd" >/dev/null; fi
  i=$((i + 1))
done < <(jq -r '.template[] | [.tab, .command // ""] | @tsv' "$template")
herdr tab focus "$first_tab" >/dev/null
