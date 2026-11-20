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
- [ ] Sunshine and Moonlight to mirror the PC screen on the Samsung TV (PC on Wi-Fi, TV on Ethernet). At home:
      find the TV model and Tizen version (Settings, Support, About This TV). Moonlight for Tizen needs Tizen 5.5+
      (about 2020+ models) and a sideload in Developer Mode (Samsung-Jellyfin-Installer, Moonlight.wgt). Tizen
      has no Tailscale, so Sunshine's ports are open in the firewall (gaming module) and Sunshine starts only by
      hand. Steps: `systemctl --user start sunshine`, set the web UI login at https://localhost:47990, install
      Moonlight with Apps2Samsung, add the PC by its home IP, enter the PIN in the web UI, stream "Desktop".
      Older TV: a Google TV stick or HDMI. Test 2026-11-20 with the external monitor (Hyprland renders on the
      NVIDIA dGPU then): black screen, "[wayland] Frame capture failed", since frames from the dGPU cannot be
      encoded on the iGPU (VA-API). Untested on the laptop screen alone (iGPU renders and encodes), which is the
      home setup: test that, phone over Tailscale is enough. The external monitor is not needed (it would need
      NVENC: `pkgs.sunshine.override {cudaSupport = true;}`, a local build, plus `adapter_name` on the NVIDIA node).
- [ ] Decide whether to switch opencode to Pi.
- [ ] Try `pass` (the standard Unix password manager, gpg-encrypted files in a git repo) instead of Proton Pass
      (from the old Trello board).
- [ ] Try Emacs (from the old Trello board).
- [ ] Blocked (firmware): lag on battery: unplugged, any load stutters a bit, most visible in Hyprland animations (window and
      workspace switching), although the power profile is performance. Hyprland renders on the Intel iGPU
      (no `AQ_DRM_DEVICES` dGPU without an external screen). Seen 2026-11-10 on battery: one busy core
      only reaches 1.5 GHz (turbo is 4.5), all cores 1.3 GHz, at 58°C, so it is a clamp, not heat. Likely
      the firmware cap on battery, which may also hold back the iGPU. On AC, same settings: one core
      4.2 GHz, all cores 3.9 GHz, so the cap is confirmed. Tested 2026-11-20: setting the `alienware-wmi`
      platform profile to performance changes nothing (still 1.6 GHz flat on every core at 61°C), and every
      software limit is open (scaling_max 4.5 GHz, turbo on, max_perf_pct 100, PL1 90 W, no thermal
      throttles). Battery is at 82% of design capacity. So the firmware or EC clamps on battery; on Dell this
      is by design and not reachable from Linux. Same stutter on the laptop screen and on the external one;
      blur and shadows off made no difference. Only unblocks if the BIOS setup (F2) has a battery or turbo
      option, or Dell ships a BIOS newer than 1.27.0 (check on dell.com with the service tag, from Windows).
      Do not clear BD PROCHOT with `wrmsr`.

## Phone

- [ ] IronFox uBlock filters: push from desktop LibreWolf to the cloud, pull on IronFox (may already be done).
