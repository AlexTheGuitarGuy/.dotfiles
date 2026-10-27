# usage: term-shot <out.png> <wait-seconds-per-command> <cmd>...
out=$(realpath -m "$1"); wait=$2; shift 2
class="shot.term.p$$"
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
sleep 1; focus
[ "$(client | jq .fullscreen)" = 0 ] && hyprctl dispatch 'hl.dsp.window.fullscreen()' >/dev/null
sleep 1.5
[ "$(geom)" != "$(screen)" ] && { echo "term-shot: not fullscreen ($(geom) vs $(screen))" >&2; exit 1; }
sleep 2
for cmd in "$@"; do
  focus; wtype -- "$cmd"; wtype -k Return; sleep "$wait"
done
grim -g "$(client | jq -r '"\(.at[0]),\(.at[1]) \(.size[0])x\(.size[1])"')" "$out" && echo "$out"
focus; wtype -- exit; wtype -k Return
