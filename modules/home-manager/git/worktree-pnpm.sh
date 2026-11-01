[ "${1:-}" = 0000000000000000000000000000000000000000 ] || exit 0

here="$(git rev-parse --show-toplevel)"
[ -f "$here/pnpm-lock.yaml" ] || grep -q '"packageManager": *"pnpm@' "$here/package.json" 2>/dev/null || exit 0
[ -d "$here/node_modules" ] && exit 0

echo "worktree-pnpm: installing dependencies in $here"
(cd "$here" && "$PNPM" install) || echo "worktree-pnpm: pnpm install failed in $here" >&2
