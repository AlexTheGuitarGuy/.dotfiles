# TODO

Personal open tasks. Work tasks live in the work repo.

## Windows

- [ ] Move to LibreWolf: follow `modules/home-manager/librewolf/README.md`, section "Windows setup" (install,
      `windows.ps1`, sync sign-in, Tridactyl without new tab, per-extension settings, set as default browser).
- [ ] Samsung T7: move game saves into Windows and hook them into Ludusavi (`Games\Saves\<game>`, emulators
      in `Games\Emulators`: BizHawk, PCSX2, PSP). List what else of the Windows setup should take part.
- [ ] Samsung T7: move music production data into the DAW locations (`Music - song ideas\Ableton User Library`,
      `FL Studio and plugin data`, `muzica mea\Samples`).
- [ ] Samsung T7: delete `Misc\Old PC\old-windows-configs` or rotate what is in it (`.aws`, `.azure`,
      `.1password`, `.docker` from the old PC).
- [ ] Kensington SlimBlade Pro firmware V4.49 (fixes unstable cursor, Windows only): trackball switch to OFF
      (centre), connect with the USB-C cable, run
      https://www.kensington.com/siteassets/product-support/kensington-slimblad-support-page/kensington-slimblade-pro-v4.49_1678403906.exe,
      click Update, wait for "Success" (up to 5 min, do not unplug).

## Linux

- [ ] Delete the full T7 copy `~/t7-backup-2026-11-09` (368 GB) and `~/t7-cleanup` once nothing is missing
      on the reorganized drive.
- [ ] Backups with restic onto the old Kubuntu 970 EVO (wipes that disk).
- [ ] Tailscale and mosh for phone access, ssh only over tailscale0.
- [ ] Sunshine and Moonlight to stream to the TV.
- [ ] Pick a phone link app: KDE Connect (favoured; `programs.kdeconnect.enable` opens ports 1714-1764,
      plus HM `services.kdeconnect` with `indicator = true`, F-Droid app on the phone) or Valent (GTK,
      same protocol, lighter, less mature). Unknown: whether phone-as-mouse works on Hyprland (needs the
      RemoteDesktop portal). Over Tailscale, add the PC's tailnet IP by hand in the app.
- [ ] Decide whether to switch opencode to Pi.
- [ ] Node version per repo through a devShell.
- [ ] ngrok: run `ngrok config add-authtoken <token>`, then test a tunnel.
- [ ] Lag on battery: unplugged, any load stutters a bit, most visible in Hyprland animations (window and
      workspace switching), although the power profile is performance. Hyprland renders on the Intel iGPU
      (no `AQ_DRM_DEVICES` dGPU without an external screen). Seen 2026-11-10 on battery: one busy core
      only reaches 1.5 GHz (turbo is 4.5), all cores 1.3 GHz, at 58°C, so it is a clamp, not heat. Likely
      the firmware cap on battery, which may also hold back the iGPU. On AC, same settings: one core
      4.2 GHz, all cores 3.9 GHz, so the cap is confirmed. Unknown: whether setting the `alienware-wmi` platform profile
      (`/sys/class/platform-profile/platform-profile-0`, at balanced-performance) to performance like
      `dell-pc` helps. Do not clear BD PROCHOT with `wrmsr`.
- [ ] ssh-agent starts empty since the switch from the gcr agent to plain ssh-agent (commit "replace the gcr
      SSH agent with plain ssh-agent"): git over ssh fails with "Permission denied (publickey)" in shells that
      cannot ask for the passphrase (agents, scripts). Plain `ssh-add` loads nothing, since the keys are not
      named `id_*`; `ssh-add ~/.ssh/github-personal ~/.ssh/gitlab-professional` works by hand. Fix so the keys
      get loaded on their own (for example `AddKeysToAgent yes` in the ssh config, or a login step).

## Phone

- [ ] IronFox uBlock filters: push from desktop LibreWolf to the cloud, pull on IronFox (may already be done).
