# LibreWolf

Config for LibreWolf with Sidebery (vertical tabs), Tridactyl (keyboard control) and Catppuccin Mocha styling.

On Linux this package is stowed like the others. Windows has no stow, so follow the steps below.

## Windows setup

### 1. Get the dotfiles

```powershell
winget install Git.Git
git clone https://github.com/AlexTheGuitarGuy/.dotfiles.git $HOME\.dotfiles
```

Already cloned? Run `git -C $HOME\.dotfiles pull` instead.

### 2. Install LibreWolf

```powershell
winget install LibreWolf.LibreWolf
```

Start LibreWolf once and close it, so it creates its profile. The script in the next step needs that profile.

### 3. Apply the config

```powershell
powershell -ExecutionPolicy Bypass -File $HOME\.dotfiles\scripts\librewolf-windows.ps1
```

It copies:

- `librewolf.overrides.cfg` to `%USERPROFILE%\.librewolf\` (stay logged in, sync on, session restore, title bar, fonts, cookie banners, faster loading)
- `tridactylrc` and the Catppuccin theme to `%USERPROFILE%\.config\tridactyl\`
- the CSS (square corners, hidden tab bar, `about:` pages, Sidebery styling) into every profile's `chrome` folder

It should print "copied chrome to ...". The script copies instead of linking, so after a dotfiles change run `git pull` and the script again.

### 4. Sign in to sync

Start LibreWolf, open the menu, choose Sign in and use the Mozilla account (password plus the 2FA code from Proton Pass). Within a few minutes these arrive:

- bookmarks, including the "Search shortcuts" folder
- history
- add-ons from addons.mozilla.org and the Catppuccin theme
- Dark Reader settings (through its own sync)

### 5. Install what sync can't

1. No new tab Tridactyl: open https://tridactyl.cmcaine.co.uk/betas/nonewtab/tridactyl_no_new_tab_beta-latest.xpi and click Add. If sync also installed the regular Tridactyl, remove that one in `about:addons`. Then open Tridactyl in `about:addons` and set Allow automatic updates to Off: the beta updates from the author's server instead of addons.mozilla.org, and with the helper below it can run commands on the machine. To update, read the changes on the Tridactyl GitHub first, then install the new `.xpi` from the link above.
2. Tridactyl helper: type `:nativeinstall`, run the PowerShell command it shows, then check with `:native` and load the config with `:source`. If `:native` reports it missing, LibreWolf on Windows may look in its own registry location instead of Firefox's.

### 6. Per-extension settings

| Extension | What to do |
|---|---|
| Sidebery | Settings, Sync: turn on and pick the profile synced from Linux (settings and keybindings). |
| uBlock Origin | Filter lists, Annoyances: tick the Cookie Notices lists, then Apply changes. Or Settings, Restore from cloud storage if a backup was made on Linux. |
| Violentmonkey | Settings, Sync with the same provider as on Linux, or export the scripts on Linux and import them here. |
| Proton Pass | Sign in, then Settings, theme: Dark (otherwise its in-page dialogs are light). |
| 1Password | Sign in. If the QR code shows stripes, click the canvas icon in the address bar, allow it and reload. |
| SponsorBlock, UnTrap | Import their settings if they were customized. |

### 7. Windows specifics

- Default browser: Settings, Apps, Default apps, LibreWolf, Set default.
- Updates: LibreWolf does not update itself on Windows. Run `winget upgrade LibreWolf.LibreWolf` now and then, or install LibreWolf WinUpdater.
- Tab shortcuts: Ctrl+1 to 9 already switch tabs on Windows.
- Ctrl+Alt+D (duplicate tab): Windows treats Ctrl+Alt as AltGr. If the keyboard layout uses AltGr+D for a character, use Tridactyl's `yd` instead.
- Fonts: the Liberation font settings do nothing on Windows; it keeps its normal fonts. If words look widely spaced, the fingerprint protection font bug also hits Windows and needs a Windows font in the overrides.

### 8. Check

Restart LibreWolf once (Ctrl+Q, then start it) and check:

- square corners, Sidebery styling with the purple outline
- new tabs show the normal LibreWolf page
- Tridactyl's `:` command line uses the Catppuccin colors
- tabs come back after a restart
- search shortcuts work, for example `:yt test`
- `about:support` shows the Mozilla account as signed in

## Things that break under fingerprint protection

| Problem | Fix |
|---|---|
| QR codes or maps look striped | Click the canvas icon in the address bar, allow, reload |
| Wrong timezone or odd layout on a site | Add the site's base domain to `privacy.fingerprintingProtection.granularOverrides` in the profile's `user.js` (kept outside this repo), for example `[{"firstPartyDomain":"example.com","overrides":"-AllTargets"}]`. `localhost` and IPs can't be listed |
| Sites ignore dark mode | Sites see the real dark preference (`-CSSPrefersColorScheme`); Dark Reader for sites without a dark theme |
| A site breaks because something is blocked | Shield icon next to the address bar, turn off tracking protection for that site |
