# TODO

Personal open tasks. Work tasks live in the work repo.

## Windows

- [ ] Move to LibreWolf: follow `modules/home-manager/librewolf/README.md`, section "Windows setup" (install,
      `windows.ps1`, sync sign-in, Tridactyl without new tab, per-extension settings, set as default browser).
- [ ] Finish the app reinstalls from the D: to C: migration (list in `C:\my_files\programs\installed-apps.md`,
      settings already restored in AppData). Still missing on 2026-11-12: Native Access (needed if Guitar Rig 7
      opens in demo mode), Cycling '74 Max (only if needed outside Ableton's Max for Live), MAXON Cinema 4D,
      Canon driver for the exact printer or scanner model. FortiClient: decide whether it is still needed.
- [ ] Samsung T7: move game saves into Windows and hook them into Ludusavi (`Games\Saves\<game>`, emulators
      in `Games\Emulators`: BizHawk, PCSX2, PSP). List what else of the Windows setup should take part.
- [ ] Samsung T7: move music production data into the DAW locations (`Music - song ideas\Ableton User Library`,
      `FL Studio and plugin data`, `muzica mea\Samples`).
- [ ] Samsung T7: delete `Misc\Old PC\old-windows-configs` or rotate what is in it (`.aws`, `.azure`,
      `.1password`, `.docker` from the old PC).
- [ ] Kensington SlimBlade Pro firmware V4.49 (fixes unstable cursor, Windows only): trackball switch to OFF
      (centre), connect with the USB-C cable, run
      https://www.kensington.com/siteassets/product-support/kensington-slimblad-support-page/kensington-slimblade-pro-v4.49_1678403906.exe,
      click Update, wait for "Success" (up to 5 min, do not unplug). See the Linux entry for the lag.

## Linux

- [ ] Delete the full T7 copy `~/t7-backup-2026-11-09` (368 GB) and `~/t7-cleanup` once nothing is missing
      on the reorganized drive.
- [ ] Monitor profiles for laptop and docked, plus closing the lid turns off the laptop screen (kanshi);
      the one unticked item from the switch-day list.
- [ ] Old agent conversations: raise the 30-day cleanup limit and, if wanted, copy back the 53 still on the
      old Kubuntu SSD (before restic wipes it).
- [ ] Backups with restic onto the old Kubuntu 970 EVO (wipes that disk).
- [ ] Tailscale and mosh for phone access, ssh only over tailscale0.
- [ ] Sunshine and Moonlight to stream to the TV.
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
- [ ] Kensington SlimBlade Pro lags or drops (2026-11-11). Kernel log stays silent while it lags, so it is
      the 2.4 GHz radio range, not USB or battery. Kensington's dongle guide says to keep the trackball within
      20 cm of the dongle; better in the UGREEN hub near the trackball, worse in the laptop. It lags the same
      on Windows (2026-11-12), so it is not Linux; the user suspects the laptop hardware. Next: firmware
      update (Windows section), then a short USB 2.0 extension cable that puts the dongle within 20 cm, and
      re-pair the dongle if it still lags (switch on 2.4 GHz, when the left LED flashes green hold upper left
      button + Bluetooth pairing button 3 s until both LEDs flash fast).

## Phone

- [ ] IronFox uBlock filters: push from desktop LibreWolf to the cloud, pull on IronFox (may already be done).
