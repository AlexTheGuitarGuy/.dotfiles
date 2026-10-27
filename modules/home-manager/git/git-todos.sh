[ "${1:-}" = 0000000000000000000000000000000000000000 ] || exit 0
top="$(git rev-parse --show-toplevel 2>/dev/null)" || exit 0
[ -e "$top/.excluded/todos.md" ] && exit 0
mkdir -p "$top/.excluded"
printf '# TODOs\n' > "$top/.excluded/todos.md"
