# ~/.dotfiles

Personal, generic NixOS and Home Manager config for the `dell-g5` host (Hyprland, caelestia shell,
nushell, herdr, Ghostty). `hosts/host-modules.nix` builds the module list; `modules/nixos` and
`modules/home-manager` hold one directory per feature, each behind a `dotfiles.<name>.enable`
option that `hosts/dell-g5` switches on.

The private work flake (path in `$NH_FLAKE`) takes this flake as its `dotfiles` input. Nothing
work-related goes into this repo: no employer names, hostnames, ticket keys or internal paths. Keep
commit messages generic too.

Before changing anything here, check whether the work flake needs the same change and warn the user when
a change here breaks it (the work flake uses `dotfiles.lib.mkHost`, `dotfiles.ai.mcpServers`, `dotfiles.ai.agents`,
`dotfiles.claude.settingsOverlays`, `dotfiles.librewolf.profilePath` and the nvim, git and nushell
modules).

Nix workflow:

1. Edit, then `git add -A` here. Flakes only see tracked files, so a new file that is not staged is
   invisible to the build.
2. Build through the work flake, which is what gets switched to:
   `cd $NH_FLAKE && nixos-rebuild build --flake .#dell-g5 --override-input dotfiles git+file:///home/alex/.dotfiles`.
3. When Hyprland or caelestia changed, verify the generated config: find `hm_hyprhyprland.lua` in
   the result closure and run `<hyprland package>/bin/Hyprland --verify-config -c <that file>`.
4. `nix flake check --no-build` in both repos and `nix fmt -- --check .` here.
5. Commit here first, then in the work flake run `nix flake update dotfiles`, stage `flake.lock` and the
   changed files by explicit path (never `git add -A` there: it holds local data files that
   must stay out of commits) and commit.
6. `rm -f result`. Never push, and never run `nixos-rebuild switch`, `nh os switch` or `rebuild`:
   the user switches.

Files that their apps or the user edit (nvim config, claude rules and skills, opencode config) are
linked live into this checkout with `mkOutOfStoreSymlink`, so the repo must stay at `~/.dotfiles`.

On Windows the same clone feeds Claude Code: `windows/link-claude.ps1` junctions
`%USERPROFILE%\.claude\rules` and `skills` to the folders here. Keep rules OS-neutral in wording; Linux
commands in them are examples only.

Checks (VM and specs, run one at a time):

- `nix build .#checks.x86_64-linux.url-bar`: the Playwright URL bar script test.
- `.#checks.x86_64-linux.desktop`: the desktop VM test; `.#checks.x86_64-linux.apps` builds every
  desktop app and is heavy (about 20 GB).
