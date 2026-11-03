{
  lib,
  config,
  pkgs,
  inputs,
  ...
}: let
  cli = config.programs.caelestia.cli.package;
  wallpaperDir = "${config.xdg.dataHome}/caelestia/wallpapers";
  wallpaper = "${wallpaperDir}/wallpaper.png";

  theme = config.dotfiles.theme;
  accentColours = pkgs.runCommand "caelestia-accent.json" {nativeBuildInputs = [(pkgs.python3.withPackages (p: [p.materialyoucolor]))];} ''
    python3 -I ${./caelestia-accent.py} ${lib.removePrefix "#" theme.palette.accent} ${theme.polarity} > $out
  '';
  applyAccent = pkgs.writeShellScript "caelestia-accent" ''
    scheme="${config.xdg.stateHome}/caelestia/scheme.json"
    tmp="$(${pkgs.coreutils}/bin/mktemp "$scheme.XXXXXX")"
    ${lib.getExe pkgs.jq} --slurpfile accent ${accentColours} '.colours += $accent[0]' "$scheme" > "$tmp"
    ${pkgs.coreutils}/bin/mv "$tmp" "$scheme"
  '';
in {
  imports = [inputs.caelestia-shell.homeManagerModules.default];

  config = lib.mkIf config.dotfiles.hyprland.enable {
    programs.caelestia = {
      enable = true;
      systemd.enable = true;
      package = inputs.caelestia-shell.packages.${pkgs.stdenv.hostPlatform.system}.with-cli.overrideAttrs (old: {
        patches = (old.patches or []) ++ [./caelestia-calc.patch ./caelestia-search.patch ./caelestia-tabs-files.patch ./caelestia-gpu-card.patch ./caelestia-lock-corners.patch ./caelestia-night-light.patch ./caelestia-osk.patch ./caelestia-kb-layout.patch ./caelestia-richtext-font.patch];
        postInstall =
          old.postInstall
          + ''
            wrapProgram $out/bin/caelestia-shell --prefix PATH : ${lib.makeBinPath [pkgs.brotab pkgs.fd pkgs.fzf]}
          '';
      });

      settings = {
        appearance = {
          transparency = {
            enabled = true;
            base = 0.9;
            layers = 0.55;
          };
        };

        border.rounding = 25;

        background = {
          enabled = true;
          wallpaperEnabled = true;
        };

        bar = {
          entries =
            map (id: {
              inherit id;
              enabled = true;
            })
            ["logo" "workspaces" "spacer" "activeWindow" "spacer" "tray" "clock" "statusIcons" "keyboard" "power"];
          workspaces = {
            shown = 4;
            activeIndicator = true;
            occupiedBg = false;
            showWindows = true;
            showWindowsOnSpecialWorkspaces = true;
            maxWindowIcons = 4;
          };
          statusIcons =
            map (id: {
              inherit id;
              enabled = !(builtins.elem id ["audio" "microphone"]);
            })
            ["lockStatus" "audio" "microphone" "kbLayout" "network" "bluetooth" "battery"];
        };

        general = {
          apps = {
            terminal = ["ghostty" "-e"];
            audio = ["pavucontrol"];
            playback = ["mpv"];
            explorer = ["ghostty" "-e" "yazi"];
          };
          idle = {
            lockBeforeSleep = true;
            inhibitWhenAudio = true;
            timeouts = [
              {
                timeout = 900;
                idleAction = "lock";
              }
            ];
          };
        };

        launcher = {
          vimKeybinds = true;
          useFuzzy = {
            apps = true;
            actions = true;
            schemes = true;
            variants = true;
            wallpapers = true;
          };
        };

        lock = {
          enabled = true;
          useWallpaper = true;
          hideNotifs = true;
          enableFprint = false;
          enableHowdy = false;
        };

        notifs = {
          defaultExpireTimeout = 2000;
          fullscreenExpireTimeout = 2000;
        };

        utilities.quickToggles =
          map (id: {
            inherit id;
            enabled = id != "vpn";
          })
          ["wifi" "bluetooth" "mic" "settings" "gameMode" "wallpaperPlaying" "dnd" "nightLight" "vpn"];

        paths.wallpaperDir = wallpaperDir;

        # Empty location: caelestia finds the city from the IP address.
        services = {
          weatherLocation = "";
          weatherUnits = "Celsius";
          clockFormat = "TwentyFourHour";
          maxVolume = 1.0;
          defaultPlayer = "Spotify";
        };
      };

      # Stylix and each tool's own config do the theming, not the caelestia cli.
      cli = {
        enable = true;
        package = inputs.caelestia-shell.inputs.caelestia-cli.packages.${pkgs.stdenv.hostPlatform.system}.default.overrideAttrs (old: {
          postInstall =
            (old.postInstall or "")
            + ''

              printf '\n\xef\xa3\xbf apple logo as apple users see it (iphone ipad mac only, others see an empty box)\n' >> $out/lib/python*/site-packages/caelestia/data/emojis.txt
            '';
        });
        settings.theme = lib.genAttrs [
          "enableTerm"
          "enableHypr"
          "enableDiscord"
          "enableSpicetify"
          "enablePandora"
          "enableFuzzel"
          "enableBtop"
          "enableNvtop"
          "enableHtop"
          "enableGtk"
          "enableQt"
          "enableWarp"
          "enableChromium"
          "enableZed"
          "enableCava"
        ] (_: false);
      };
    };

    xdg.dataFile."caelestia/wallpapers".source = config.lib.file.mkOutOfStoreSymlink "${config.dotfiles.repo}/modules/home-manager/hyprland/wallpapers";

    # Caelestia rewrites shell.json at runtime, so replace the store link with a writable copy.
    xdg.configFile."caelestia/shell.json".force = true;
    home.activation.caelestiaWritableSettings = lib.hm.dag.entryBetween ["reloadSystemd"] ["linkGeneration"] ''
      config_file="${config.xdg.configHome}/caelestia/shell.json"
      if [ -L "$config_file" ]; then
        tmp_file="$(${pkgs.coreutils}/bin/mktemp "''${config_file}.XXXXXX")"
        ${pkgs.coreutils}/bin/cp --dereference "$config_file" "$tmp_file"
        ${pkgs.coreutils}/bin/chmod 0644 "$tmp_file"
        ${pkgs.coreutils}/bin/mv --force --no-target-directory "$tmp_file" "$config_file"
      fi
    '';

    # Seed the color scheme and wallpaper before the shell starts.
    systemd.user.services.caelestia.Unit = {
      After = ["hyprsunset.service"];
      Wants = ["hyprsunset.service"];
    };
    systemd.user.services.caelestia.Service.ExecStartPre = [
      "${cli}/bin/caelestia scheme set -n catppuccin -f ${theme.flavor} -m ${theme.polarity}"
      "${applyAccent}"
      "${pkgs.bash}/bin/bash -c '[ -s ${config.xdg.stateHome}/caelestia/wallpaper/path.txt ] || ${cli}/bin/caelestia wallpaper -f ${wallpaper} -N'"
    ];

    services.blueman-applet.enable = true;
    dconf.settings."org/blueman/general".plugin-list = ["!StatusIcon" "!ShowConnected"];
  };
}
