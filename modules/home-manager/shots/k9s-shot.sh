# usage: k9s-shot <k9s args...> -- <action>...
#   actions: key:<key> | type:<text> | sleep:<sec> | shot:<out.png>
#   key: slash, Return, Escape or a tmux key name (d, l, 6, C-c)
# k9s runs in a private tmux inside a fullscreen Ghostty; keys go through tmux send-keys,
# so nothing depends on window focus and no keystroke can land in another window.
args=(); while [ "$1" != "--" ]; do args+=("$1"); shift; done; shift
cfg=$(mktemp -d); sock="k9s-shot-$$"; class="shot.k9s.p$$"
client() { hyprctl clients -j | jq -c --arg c "$class" '.[] | select(.class == $c)'; }
trap 'tmux -L "$sock" kill-server 2>/dev/null; rm -rf "$cfg"; for _ in $(seq 1 20); do [ -n "$(client)" ] || break; sleep 0.25; done' EXIT
cp -rL "${K9S_CONFIG_DIR:-$HOME/.config/k9s}/." "$cfg/"
chmod -R u+w "$cfg"
sed -i -E 's/^(\s*tail:\s*)[0-9]+/\150000/; s/^(\s*buffer:\s*)[0-9]+/\150000/; s/^(\s*skipLatestRevCheck:\s*).*/\1true/' "$cfg/config.yaml"
cat > "$cfg/tmux.conf" <<'CONF'
set -g default-shell /bin/sh
set -g status off
set -g escape-time 0
set -g default-terminal "tmux-256color"
set -ga terminal-overrides ",*:RGB"
CONF
t() { tmux -L "$sock" "$@"; }
ghostty --class="$class" -e tmux -L "$sock" -f "$cfg/tmux.conf" new-session -s shot \
  "env K9S_CONFIG_DIR=$(printf %q "$cfg") k9s --readonly $(printf '%q ' "${args[@]}") 2>$(printf %q "$cfg/k9s.err")" >/dev/null 2>&1 &
w=""
for _ in $(seq 1 60); do
  w=$(client)
  [ -n "$w" ] && break; sleep 0.5
done
[ -z "$w" ] && { echo "k9s-shot: no ghostty window" >&2; exit 1; }
addr=$(jq -r .address <<<"$w")
focus() { hyprctl dispatch "hl.dsp.focus({ window = \"address:$addr\" })" >/dev/null; }
active() { [ "$(hyprctl activewindow -j | jq -r .address)" = "$addr" ]; }
sleep 1
for _ in $(seq 1 10); do
  focus
  [ "$(client | jq .fullscreen)" = 0 ] && hyprctl dispatch 'hl.dsp.window.fullscreen()' >/dev/null
  sleep 1
  [ "$(client | jq .fullscreen)" != 0 ] && break
done
[ "$(client | jq .fullscreen)" != 0 ] || { echo "k9s-shot: not fullscreen" >&2; exit 1; }
sleep 6
for a in "$@"; do
  t has-session -t shot 2>/dev/null || { echo "k9s-shot: k9s exited before '$a'" >&2; cp "$cfg/k9s.err" "${XDG_CACHE_HOME:-$HOME/.cache}/k9s-shot.err"; grep -m3 -iE "panic|error|fatal" "$cfg/k9s.err" >&2; exit 1; }
  case "$a" in
    key:*) k=${a#key:}; case "$k" in slash) k=/ ;; Return) k=Enter ;; esac
           t send-keys -t shot "$k"; sleep 0.8 ;;
    type:*) t send-keys -t shot -l -- "${a#type:}"; sleep 0.8 ;;
    sleep:*) sleep "${a#sleep:}" ;;
    shot:*) for _ in $(seq 1 10); do
              focus; sleep 0.5
              active && break
            done
            active || { echo "k9s-shot: window not on top, no shot" >&2; exit 1; }
            grim -g "$(client | jq -r '"\(.at[0]),\(.at[1]) \(.size[0])x\(.size[1])"')" "$(realpath -m "${a#shot:}")"
            active || { rm -f "$(realpath -m "${a#shot:}")"; echo "k9s-shot: window lost focus during shot, discarded" >&2; exit 1; }
            echo "${a#shot:}" ;;
  esac
done
