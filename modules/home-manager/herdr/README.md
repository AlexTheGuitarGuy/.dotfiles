# herdr setup

Replaces tmux.

Everything comes from `default.nix` in this directory:

- `herdr` from the upstream flake (`inputs.herdr`, pinned to a release tag)
- `~/.config/herdr/config.toml`, generated from Nix (keybindings, theme, status bar, flags).
  Edit `default.nix`, rebuild, then `herdr server reload-config`.
- `status.sh`, packaged as `herdr-status` for the tab bar (needs curl)
- `sessionizer.config.json`, linked live into the sessionizer plugin's config dir
  (`herdr plugin config-dir sessionizer`), so edits apply on the next picker run
- `herdr-worktree` and `tzt` as packages on PATH
- the `herdr` systemd user service

## Plugins

Both plugins are pinned with `fetchFromGitHub` and registered with `herdr plugin link`
on every home-manager activation. Bump `rev` and `hash` in `default.nix` to update.
The herdr service's PATH carries their dependencies (bash, jq, fzf, zoxide, node).

- **sessionizer** (`salkhalil/herdr-sessionizer`): tgs-equivalent. `prefix+f`
  opens an fzf picker over open workspaces + zoxide dirs; picking a dir creates a
  workspace with the template tabs from `sessionizer.config.json` (nvim, lazygit,
  agent, shell). `prefix+g` (`herdr-worktree`) reuses the same tabs for worktrees.
- **synchronize-panes** (`furuhashin/herdr-synchronize-panes`): `prefix+m`
  broadcast overlay, one command into every pane of the current tab (old tmux
  `bind m`). One-shot, not live keystroke mirroring.

`tzt` runs the sessionizer's standalone `bin/sessionizer`, then attaches the herdr
TUI, so `tzt` from a bare shell behaves like `tgs`: pick a project, land inside
herdr. `tzt <dir>` skips the picker.

## Status bar

`[ui] tab_bar_right` in `config.toml` shows Chisinau weather (via `status.sh`,
refreshed every 15 min) and a clock. `status.sh` runs on the herdr server, so the
command path is absolute. Replaces the old dracula/tmux weather widget.

## Autostart and persistence

The `herdr` user service starts with the graphical session. Check it with
`systemctl --user status herdr`.

On a server restart herdr restores workspace/tab/pane layout and cwd. Running
processes are NOT restored, unlike tmux-resurrect. `session.resume_agents_on_restore`
(default true) re-resumes recognized agent conversations only.

`experimental.pane_history` is off: with it on, herdr rewrites a ~800 KB
`session-history.json` on every activity tick, which on a busy machine shows up as
input lag while typing. The trade-off is no scrollback restore across a full server
restart (rare, and the processes are gone anyway). `HERDR_LOG=herdr=warn` in the
unit silences the per-event INFO log for the same reason.

## Agent integration

`claude/.claude/hooks/herdr-agent-state.sh` plus a `SessionStart` entry in the
matching `settings.json` are herdr's agent integration for the CLI under
`claude/` (added by `herdr integration install`). The hook reports the running
session to the herdr sidebar for live state and exact-session resume. herdr
rewrites the hook on integration version bumps, so expect an occasional diff.
`herdr integration status` to inspect, `herdr integration uninstall` to remove.

## Prefix key

Prefix is `ctrl+a` (tmux parity). herdr has no send-prefix, so a bare `ctrl+a` at
the shell prompt is caught by herdr, not readline; use `Home` for
beginning-of-line.

`stty -ixon` in the nushell config disables
XON/XOFF flow control so `ctrl+s` / `ctrl+q` are usable keys. Was required for the
old `ctrl+s` prefix, now just a sane default.

## tmux bindings not carried over

- `bind r` reload            -> `herdr server reload-config`
- `bind c` name-on-create    -> `prefix+c` then `prefix+shift+t`
- `M-u/i/o/p` window jump     -> `alt+1..9` (herdr has no letter-key tab bind)
- `C-S-Left/Right` tab swap   -> unbound; native `move_tab_previous` / `move_tab_next` in `[keys]` if wanted
- `escape-time`, `focus-events`, `terminal-features RGB` -> native, not needed
