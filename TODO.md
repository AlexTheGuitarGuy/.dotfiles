# TODO

Personal open tasks. Work tasks live in the work repo.

## Windows

- [ ] Move to LibreWolf: follow `modules/home-manager/librewolf/README.md`, section "Windows setup" (install,
      `windows.ps1`, sync sign-in, Tridactyl without new tab, per-extension settings, set as default browser).
- [ ] Check the clock fix: `RealTimeIsUniversal=1` was set on 2026-11-11 and Windows rebooted. If the clock is off
      by the UTC offset, run `w32tm /resync /force` in an admin shell (if it reports no time data, first run
      `w32tm /config /syncfromflags:manual /manualpeerlist:time.windows.com /update`). Then boot Linux and back:
      Windows time must still be right.
- [ ] Samsung T7: move game saves into Windows and hook them into Ludusavi (`Games\Saves\<game>`, emulators
      in `Games\Emulators`: BizHawk, PCSX2, PSP). List what else of the Windows setup should take part.
- [ ] Samsung T7: move music production data into the DAW locations (`Music - song ideas\Ableton User Library`,
      `FL Studio and plugin data`, `muzica mea\Samples`).
- [ ] Samsung T7: delete `Misc\Old PC\old-windows-configs` or rotate what is in it (`.aws`, `.azure`,
      `.1password`, `.docker` from the old PC).

## Linux

- [ ] Earphones (Nothing Ear (a), `3C:B0:ED:BA:2F:26`): Windows still needed re-pairing after booting Linux
      (checked 2026-11-11). Both systems share one adapter address, so the earphones keep one link key and
      every new pairing breaks the other system. Windows holds the newest pairing. On Linux, do not pair again:
      check with `lsblk` that the Windows partition (870 EVO, partition 2) is `/dev/sda2`, then run
      `sudo bt-sync-windows-key` (`PART=/dev/<partition>` if it differs). Confirm the earphones connect on Linux,
      then boot Windows and confirm they connect there without re-pairing. If the mount fails, disable Fast
      Startup in Windows. If Windows still needs re-pairing, check whether the LE keys need copying too (the
      script copies only the classic link key).
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
