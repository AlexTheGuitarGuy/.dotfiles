# usage: term-shot <out.png> <wait-seconds-per-command> <cmd>...
out=$(realpath -m "$1"); wait=$2; shift 2
class="shot.term.p$$"

for dir in "${XDG_RUNTIME_DIR:-/run/user/$(id -u)}"/hypr/*/; do
  sig=$(basename "$dir")
  HYPRLAND_INSTANCE_SIGNATURE=$sig hyprctl activeworkspace -j >/dev/null 2>&1 && { export HYPRLAND_INSTANCE_SIGNATURE=$sig; break; }
done
prev=$(hyprctl activeworkspace -j | jq -r .id)
clip=$(mktemp)
wl-paste --no-newline >"$clip" 2>/dev/null || true
restore() {
  wl-copy <"$clip" >/dev/null 2>&1 || true
  rm -f "$clip"
  hyprctl eval "hl.config({ cursor = { invisible = false } })" >/dev/null
  hyprctl dispatch "hl.dsp.focus({ workspace = \"$prev\" })" >/dev/null
}
trap restore EXIT
trap 'exit 130' INT TERM

client() { hyprctl clients -j | jq -c --arg c "$class" '.[] | select(.class == $c)'; }
sleep 1
ghostty --class="$class" --working-directory="$PWD" -e nu >/dev/null 2>&1 &
w=""
for _ in $(seq 1 60); do
  w=$(client)
  [ -n "$w" ] && break; sleep 0.5
done
[ -z "$w" ] && { echo "term-shot: no ghostty window" >&2; exit 1; }
addr=$(jq -r .address <<<"$w")
focus() { hyprctl dispatch "hl.dsp.focus({ window = \"address:$addr\" })" >/dev/null; }
geom() { client | jq -r '"\(.size[0])x\(.size[1])"'; }
screen() {
  hyprctl monitors -j | jq -r --argjson id "$(client | jq .monitor)" \
    '.[] | select(.id == $id) | "\(.width / .scale | floor)x\(.height / .scale | floor)"'
}
descendants() { local c; for c in $(pgrep -P "$1"); do echo "$c"; descendants "$c"; done; }
shell_pid() {
  local p
  for p in $(descendants "$(jq -r .pid <<<"$w")"); do
    [ "$(ps -o comm= -p "$p")" = nu ] && { echo "$p"; return; }
  done
}
busy() { [ -n "$(pgrep -P "$(shell_pid)")" ]; }
shot() {
  focus; sleep 1
  hyprctl eval "hl.config({ cursor = { invisible = true } })" >/dev/null; sleep 1
  grim -g "$(client | jq -r '"\(.at[0]),\(.at[1]) \(.size[0])x\(.size[1])"')" "$1" && echo "$1"
  hyprctl eval "hl.config({ cursor = { invisible = false } })" >/dev/null
}

sleep 1; focus
[ "$(client | jq .fullscreen)" = 0 ] && hyprctl dispatch 'hl.dsp.window.fullscreen()' >/dev/null
sleep 1.5
[ "$(geom)" != "$(screen)" ] && { echo "term-shot: not fullscreen ($(geom) vs $(screen))" >&2; exit 1; }
sleep 2
for cmd in "$@"; do
  focus; printf %s "$cmd" | wl-copy; sleep 0.3
  wtype -M ctrl -M shift -k v -m shift -m ctrl; sleep 0.3; wtype -k Return; sleep 1
  while busy; do sleep 1; done
done
if [ "$wait" = "done" ]; then
  while busy; do sleep 1; done; sleep 2
else
  sleep "$wait"
fi
shot "$out"
if [ "${TERM_SHOT_TOP:-}" = 1 ]; then
  focus; for _ in $(seq 1 60); do wtype -M shift -k Prior -m shift; done; sleep 0.5
  shot "${out%.png}-top.png"
fi
focus; wtype -- exit; wtype -k Return
