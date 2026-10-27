case "${1:-}" in
  push) ;;
  pull) [[ "${GIT_REFLOG_ACTION:-}" == pull* ]] || exit 0 ;;
  *) exit 0 ;;
esac
command -v gitnexus >/dev/null || exit 0

here="$(git rev-parse --show-toplevel)"
meta="$here/.gitnexus/meta.json"
[ -f "$meta" ] || exit 0

args=(analyze --index-only)
[ "$(jq -r '.stats.embeddings // 0' "$meta")" -gt 0 ] && args+=(--embeddings)

cd "$here"
setsid flock -n .gitnexus/.reindex.lock \
  bash -c 'exec nice -n 19 gitnexus "$@" >.gitnexus/reindex.log 2>&1' _ "${args[@]}" </dev/null >/dev/null 2>&1 &
