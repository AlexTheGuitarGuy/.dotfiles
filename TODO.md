# TODO

Personal open tasks. Work tasks live in the work repo.

## Windows

- [ ] After pulling this repo, run `powershell -ExecutionPolicy Bypass -File windows\link-claude.ps1`, so
      Claude Code on Windows gets the rules, skills, plugins and model from this clone. Rerun it after a pull
      that changed `settings.portable.json`.
- [ ] Move to LibreWolf: follow `modules/home-manager/librewolf/README.md`, section "Windows setup" (install,
      `windows.ps1`, sync sign-in, Tridactyl without new tab, per-extension settings, set as default browser).
- [ ] Check the earphones connect in Windows without removing them first (dual-boot pairing fix from
      2026-11-05, never checked on the Windows side).
- [ ] Check the clock fix: Windows should treat the hardware clock as UTC (registry change given 2026-11-01;
      unknown whether it was applied). Symptom if not: wrong time in Windows after booting Linux.
- [ ] Samsung T7: move game saves into Windows and hook them into Ludusavi (`Games\Saves\<game>`, emulators
      in `Games\Emulators`: BizHawk, PCSX2, PSP). List what else of the Windows setup should take part.
- [ ] Samsung T7: move music production data into the DAW locations (`Music - song ideas\Ableton User Library`,
      `FL Studio and plugin data`, `muzica mea\Samples`).
- [ ] Samsung T7: delete `Misc\Old PC\old-windows-configs` or rotate what is in it (`.aws`, `.azure`,
      `.1password`, `.docker` from the old PC).

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

## Phone

- [ ] IronFox uBlock filters: push from desktop LibreWolf to the cloud, pull on IronFox (may already be done).
