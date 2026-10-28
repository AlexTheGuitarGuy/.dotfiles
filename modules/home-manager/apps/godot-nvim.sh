pipe=/tmp/godot-nvim.pipe

setup() {
  local settings
  settings="$(printf '%s\n' "$HOME"/.config/godot/editor_settings-4.*.tres | sort -V | tail -1)"
  if [ ! -f "$settings" ]; then
    echo "godot-nvim: no editor settings yet, start Godot once first" >&2
    exit 1
  fi
  if pgrep -x godot >/dev/null; then
    echo "godot-nvim: close Godot first, it rewrites $settings on exit" >&2
    exit 1
  fi
  local flags="open {project} {file} {line} {col}"
  python3 - "$settings" "$GODOT_NVIM_EXEC" "$flags" <<'EOF'
import re, sys
path, exec_path, flags = sys.argv[1], sys.argv[2], sys.argv[3]
want = {
    'text_editor/external/use_external_editor': 'true',
    'text_editor/external/exec_path': '"' + exec_path + '"',
    'text_editor/external/exec_flags': '"' + flags.replace('"', '\\"') + '"',
}
text = open(path).read()
for key, value in want.items():
    line = f'{key} = {value}'
    pattern = re.compile(rf'^{re.escape(key)} = .*$', re.M)
    text = pattern.sub(lambda _: line, text) if pattern.search(text) else text.rstrip('\n') + f'\n{line}\n'
open(path, 'w').write(text)
EOF
  echo "godot-nvim: external editor set in $settings"
}

herdr_workspace() {
  herdr pane list 2>/dev/null | jq -r --arg p "$1" \
    '[.result.panes[] | select(.cwd == $p or (.cwd | startswith($p + "/")))][0].workspace_id // empty'
}

start_in_herdr() {
  local project="$1" ws out
  ws="$(herdr_workspace "$project")"
  if [ -n "$ws" ]; then
    out="$(herdr tab create --workspace "$ws" --cwd "$project" --label nvim --focus)"
  else
    out="$(herdr workspace create --cwd "$project" --label "$(basename "$project")" --focus)"
  fi
  herdr pane run "$(jq -r '.result.root_pane.pane_id' <<<"$out")" "nvim --listen $pipe" >/dev/null
  pgrep -f '^herdr( |$)' >/dev/null || setsid ghostty -e herdr >/dev/null 2>&1 &
}

ghostty_window() {
  hyprctl clients -j | jq -r --arg title "$1" \
    '[.[] | select(.class == "com.mitchellh.ghostty" and (.title | startswith($title)))][-1].address // empty'
}

focus_terminal() {
  local addr=""
  for _ in $(seq 20); do
    addr="$(ghostty_window "$(hostname): ")"
    [ -n "$addr" ] && break
    [ -n "$1" ] || break
    sleep 0.2
  done
  [ -n "$addr" ] || addr="$(ghostty_window "")"
  [ -n "$addr" ] && hyprctl dispatch "hl.dsp.focus({ window = \"address:$addr\" })" >/dev/null 2>&1 || true
}

open_file() {
  local project="$1" file="$2" line="${3:-1}" col="${4:-1}" herdr_up="" ws
  command -v herdr >/dev/null && herdr status server 2>/dev/null | grep -q 'status: running' && herdr_up=1
  if nvim --server "$pipe" --remote-expr 1 >/dev/null 2>&1; then
    if [ -n "$herdr_up" ]; then
      ws="$(herdr_workspace "$project")"
      [ -n "$ws" ] && herdr workspace focus "$ws" >/dev/null
    fi
  else
    rm -f "$pipe"
    cd "$project"
    if [ -n "$herdr_up" ]; then
      start_in_herdr "$project"
    else
      setsid ghostty -e nvim --listen "$pipe" >/dev/null 2>&1 &
    fi
    for _ in $(seq 50); do
      nvim --server "$pipe" --remote-expr 1 >/dev/null 2>&1 && break
      sleep 0.2
    done
  fi
  nvim --server "$pipe" --remote "$file"
  nvim --server "$pipe" --remote-send "<C-\\><C-N>:call cursor($line, $col)<CR>"
  focus_terminal "$herdr_up"
}

case "${1:-}" in
  setup) setup ;;
  open) shift; open_file "$@" ;;
  *) echo "usage: godot-nvim setup | open <project> <file> [line] [col]" >&2; exit 2 ;;
esac
