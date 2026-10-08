# dvag-vpn

Runs the DVAG Fortinet VPN connect flow as a `systemd --user` service
instead of a foreground terminal job, so it survives locking/idle/terminal
close and restarts itself on disconnect.

## Known limitation: every tunnel drop means a new SSO login

The SSO cookie is not single-use, the gateway accepts it again after a drop.
What kills it is the client's own logout:

- openconnect detects a dead peer, reconnects with the same cookie, the
  gateway accepts it but hands out a different IP. openconnect refuses that
  (`Reconnect gave different Legacy IP address`), sends `GET /remote/logout`
  and exits. systemd reruns the script, so a new SSO popup.
- openfortivpn `--persistent` was tried (2026-10-07) and is worse: on a drop
  it also logs out, then loops forever on `Could not get VPN configuration`
  with the dead cookie and never exits, so systemd never restarts it.

The gateway also reports "reconnect-after-drop is allowed within 30 seconds,
but only from the same source IP address".

Most drops were the laptop suspending (anything over 30s ends the session)
or wifi roaming between access points. The unit wraps the script in
`systemd-inhibit --what=idle`, so Plasma's idle suspend is blocked while the
VPN runs. Lid close and manual suspend still work and still need a new
login. Prefer ethernet: a wifi flap or a switch between interfaces changes
the source IP.

The cookie is passed with `--cookie-on-stdin` so sudo does not write it into
the system journal.

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
