# herdr setup

Replaces tmux (`tmux/.tmux.conf` + `scripts/tmux-session-generator.sh`). tmux is
kept as a fallback and is not removed.

Tested against herdr 0.8.2, stable channel (`herdr channel show`).

## Binary

`herdr` is vendored at `i3/.local/bin/herdr` and gitignored (not tracked). On a
new machine, install it out of band:

    herdr update            # self-update once a build is on PATH
    # or download from https://herdr.dev to ~/.local/bin/herdr

## Files in this package (stowed)

- `.config/herdr/config.toml`             keybindings, theme, status bar, flags
- `.config/herdr/status.sh`               weather widget for the tab bar (needs curl)
- `.config/herdr/sessionizer.config.json` reference config for the sessionizer plugin
- `.config/systemd/user/herdr.service`    headless server unit (replaces tmux-continuum)

The package is stowed (`stow herdr/` from the repo root), so `~/.config/herdr/*`
and `~/.config/systemd/user/herdr.service` are symlinks into this repo. Edit here,
then `herdr server reload-config`. No copying.

The one thing stow cannot place: the sessionizer plugin reads its config from its
own (gitignored) plugin dir, so after install/changes:

    cp ~/.config/herdr/sessionizer.config.json "$(herdr plugin config-dir sessionizer)/config.json"
    herdr server reload-config

## Plugins

`herdr plugin ...` works on stable (just hidden from `herdr --help`). Plugin code
is cloned under `~/.config/herdr/plugins/` and gitignored. Reinstall on a new
machine:

    herdr plugin install salkhalil/herdr-sessionizer -y
    herdr plugin install furuhashin/herdr-synchronize-panes -y
    cp ~/.config/herdr/sessionizer.config.json "$(herdr plugin config-dir sessionizer)/config.json"
    herdr server reload-config

- **sessionizer** (`salkhalil/herdr-sessionizer`): tgs-equivalent. `prefix+f`
  opens an fzf picker over open workspaces + zoxide dirs; picking a dir creates a
  workspace with the template tabs from `sessionizer.config.json` (4 plain tabs).
  Needs `fzf >= 0.45`, `jq`, `bash >= 4`, `zoxide`.
- **synchronize-panes** (`furuhashin/herdr-synchronize-panes`): `prefix+m`
  broadcast overlay, one command into every pane of the current tab (old tmux
  `bind m`). One-shot, not live keystroke mirroring. Needs `node` on PATH.

`scripts/tzt` (alias `tzt`) wraps the sessionizer's standalone `bin/sessionizer`
then attaches the herdr TUI, so `tzt` from a bare shell behaves like `tgs`: pick a
project, land inside herdr. `tzt <dir>` skips the picker.

## Status bar

`[ui] tab_bar_right` in `config.toml` shows Chisinau weather (via `status.sh`,
refreshed every 15 min) and a clock. `status.sh` runs on the herdr server, so the
command path is absolute. Replaces the old dracula/tmux weather widget.

## Autostart and persistence

    herdr server stop                    # if an ad-hoc server is already running
    systemctl --user daemon-reload
    systemctl --user enable --now herdr
    systemctl --user status herdr        # confirm it stays "active (running)"
    loginctl enable-linger "$USER"       # start at boot, not just first login

On a server restart herdr restores workspace/tab/pane layout, cwd, and (with
`experimental.pane_history = true`) recent scrollback. Running processes are NOT
restored, unlike tmux-resurrect. `session.resume_agents_on_restore` (default true)
re-resumes recognized agent conversations only.

If `status` shows the service exiting immediately the server is daemonizing:
change the unit `Type=simple` to `Type=exec` (or `forking`) and retry.

## Agent integration

`claude/.claude/hooks/herdr-agent-state.sh` plus a `SessionStart` entry in the
matching `settings.json` are herdr's agent integration for the CLI under
`claude/` (added by `herdr integration install`). The hook reports the running
session to the herdr sidebar for live state and exact-session resume. herdr
rewrites the hook on integration version bumps, so expect an occasional diff.
`herdr integration status` to inspect, `herdr integration uninstall` to remove.

## Prefix key

Prefix is `ctrl+a` (tmux parity). herdr 0.8.x has **no send-prefix**, so a literal
`ctrl+a` (readline beginning-of-line) cannot be passed through cleanly. Use `Home`,
or add:

```toml
[[keys.command]]
key = "prefix+ctrl+a"
type = "shell"
command = "herdr pane send-keys $(herdr pane current --current | jq -r .result.pane.pane_id) ctrl+a"
```

(`herdr pane send-keys <pane> ctrl+a` is verified; whether a `type = "shell"`
binding resolves `--current` to the triggering pane is not.)

`stty -ixon` in `zsh/.zshrc` and `nushell/.config/nushell/config.nu` disables
XON/XOFF flow control so `ctrl+s` / `ctrl+q` are usable keys. It was required when
the prefix was `ctrl+s`; now it is just a sane default.

## tmux bindings not carried over

- `bind r` reload            -> `herdr server reload-config`
- `bind c` name-on-create    -> `prefix+c` then `prefix+shift+t`
- `M-u/i/o/p` window jump     -> `alt+1..9` (herdr has no letter-key tab bind)
- `C-S-Left/Right` tab swap   -> dropped (use `herdr tab` CLI)
- `escape-time`, `focus-events`, `terminal-features RGB` -> native, not needed
