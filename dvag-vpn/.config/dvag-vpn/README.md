# dvag-vpn

Runs the DVAG Fortinet VPN connect flow as a `systemd --user` service
instead of a foreground terminal job, so it survives locking/idle/terminal
close and restarts itself on disconnect.

## Known limitation (deferred, not fixed)

`Restart=on-failure` reruns the **entire** `dvag-vpn-connect` script on any
disconnect, including the interactive SSO login window (Electron webview),
not just the `openconnect` reconnect. This was a deliberate choice, not an
oversight: each `SVPNCOOKIE` from the SSO step is treated as single-use, so
retrying `openconnect` with the same stale cookie after a drop would just
fail forever. Re-running the whole flow means you get a fresh cookie (and a
fresh popup to click through) on every reconnect.

**This single-use assumption was never actually verified against a real
disconnect/reconnect cycle** (inferred from the Fortinet SSO cookie
mechanics, not observed). If it turns out the same cookie *can* be reused
across multiple `openconnect` invocations, this should change to: cache the
cookie after the first successful auth, have the systemd unit try
reconnecting with the cached cookie first, and only fall back to the
interactive SSO popup once that cached cookie actually fails. Left
unfixed until confirmed with a real test.

## Portability: what's in dotfiles vs. machine-local

**Tracked in this repo (stow-managed, reproducible on a new machine):**
- `.local/bin/dvag-vpn-connect` - the wrapper script
- `.config/systemd/user/dvag-vpn.service` - the systemd unit
  (`WantedBy=graphical-session.target`, same pattern as `herdr.service`)

This machine runs **KDE Plasma (Kubuntu)**, not i3, see
[[kde-plasma-not-i3]] in memory - an earlier version of this feature
wrongly assumed i3 and added a manual `systemctl --user
import-environment DISPLAY XAUTHORITY` to i3's config, which did nothing
since i3 isn't the active session. Verified KDE Plasma's own systemd
integration already imports `DISPLAY`/`XAUTHORITY` into `systemd --user`
automatically (`systemctl --user show-environment` shows both correctly
set with no extra glue) - so no DISPLAY-propagation workaround is needed
here at all. A different WM/DE that lacks that integration (i3, a
from-scratch Hyprland setup) would need its own equivalent of that
import line, added to *that* WM's config, not assumed.

**NOT tracked, must already exist on the machine:**
- `~/Apps/openfortivpn-webview/openfortivpn-webview-electron` - the actual
  Electron app + its `node_modules` (needs its own `npm install`), lives
  outside `~/.dotfiles` entirely, nothing here installs or clones it
- The `openconnect` binary itself (`/usr/sbin/openconnect`) - an apt
  package, not managed here
- Passwordless sudo (`(ALL) NOPASSWD: ALL` in this machine's sudoers) -
  the script runs `sudo openconnect` with no password handling of its own;
  on a machine without passwordless sudo already configured, the systemd
  unit will just hang/fail on that line (no TTY for a password prompt)
- Working `node`/`npm`/`npx` on PATH (this machine gets that from
  nvm, itself set up in `zsh/.config/zsh/exports.zsh`, so partially
  portable, but nothing here runs `npm install` for the Electron app)
- A desktop environment/WM that actually fires `graphical-session.target`
  on login with DISPLAY/XAUTHORITY already in the systemd --user
  environment (KDE Plasma does this natively here; a lighter WM like i3
  or Hyprland needs a manual `systemctl --user import-environment
  DISPLAY XAUTHORITY` added to its own startup config instead)

So: the automation/service-management layer is fully portable, the actual
VPN client app and its system-level prerequisites are not, and would need
to be set up by hand on any new machine before this becomes useful there.
