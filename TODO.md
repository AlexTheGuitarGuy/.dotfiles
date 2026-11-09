# TODO

Personal open tasks. Work tasks live in the work repo.

## Windows

- [ ] After pulling this repo, run `powershell -ExecutionPolicy Bypass -File windows\link-claude.ps1` once,
      so Claude Code on Windows gets the rules and skills from this clone.
- [ ] Samsung T7: move game saves into Windows and hook them into Ludusavi (`Games\Saves\<game>`, emulators
      in `Games\Emulators`: BizHawk, PCSX2, PSP). List what else of the Windows setup should take part.
- [ ] Samsung T7: move music production data into the DAW locations (`Music - song ideas\Ableton User Library`,
      `FL Studio and plugin data`, `muzica mea\Samples`).
- [ ] Samsung T7: delete `Misc\Old PC\old-windows-configs` or rotate what is in it (`.aws`, `.azure`,
      `.1password`, `.docker` from the old PC).
- [ ] Optional: set Claude Code plugins and settings by hand (superpowers, caveman, ponytail, model) if missed.

## Linux

- [ ] Delete the full T7 copy `~/t7-backup-2026-11-09` (368 GB) and `~/t7-cleanup` once nothing is missing
      on the reorganized drive.
- [ ] Backups with restic onto the old Kubuntu 970 EVO (wipes that disk).
- [ ] Tailscale and mosh for phone access, ssh only over tailscale0.
- [ ] Sunshine and Moonlight to stream to the TV.
- [ ] Decide whether to switch opencode to Pi.
- [ ] Node version per repo through a devShell.
- [ ] ngrok: run `ngrok config add-authtoken <token>`, then test a tunnel.

## Phone

- [ ] IronFox uBlock filters: push from desktop LibreWolf to the cloud, pull on IronFox (may already be done).
