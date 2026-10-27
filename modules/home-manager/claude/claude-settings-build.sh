base="$1"; shift
out="$HOME/.claude/settings.json"

if [ -e "$out" ] && [ ! -L "$out" ] && [ -z "${FORCE:-}" ]; then
  if ! diff -q <(jq -S . "$out") <(jq -S . "$out.last-build" 2>/dev/null || echo '{}') >/dev/null; then
    echo "claude-settings-build: $out was changed since the last build (e.g. by /config)." >&2
    echo "Move those changes into $base or an overlay, or rerun with FORCE=1 to overwrite." >&2
    exit 1
  fi
fi

tmp="$(mktemp "$out.XXXXXX")"
jq -s '
  def merge($a; $b):
    if ($a | type) == "object" and ($b | type) == "object" then
      reduce ($b | keys_unsorted[]) as $k ($a; .[$k] = merge(.[$k]; $b[$k]))
    elif ($a | type) == "array" and ($b | type) == "array" then $a + ($b - $a)
    elif $b == null then $a
    else $b end;
  reduce .[1:][] as $overlay (.[0]; merge(.; $overlay))
' "$base" "$@" > "$tmp"
cp "$tmp" "$out.last-build"
mv -f "$tmp" "$out"
echo "wrote $out from $(basename "$base")${*:+ + $*}"
