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
in {
  imports = [inputs.caelestia-shell.homeManagerModules.default];

  config = lib.mkIf config.dotfiles.hyprland.enable {
    programs.caelestia = {
      enable = true;
      systemd.enable = true;

      settings = {
        appearance.transparency = {
          enabled = true;
          base = 0.9;
          layers = 0.55;
        };

        background = {
          enabled = true;
          wallpaperEnabled = true;
        };

        bar = {
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
            terminal = ["ghostty"];
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

        paths.wallpaperDir = wallpaperDir;

        # Empty location: caelestia finds the city from the IP address.
        services = {
          weatherLocation = "";
          useFahrenheit = false;
          useTwelveHourClock = false;
          maxVolume = 1.0;
          defaultPlayer = "Spotify";
        };
      };

      # Stylix and each tool's own config do the theming, not the caelestia cli.
      cli = {
        enable = true;
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

    xdg.dataFile."caelestia/wallpapers/wallpaper.png".source = ./wallpaper.png;

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
    systemd.user.services.caelestia.Service.ExecStartPre = [
      "${cli}/bin/caelestia scheme set -n catppuccin -f mocha -m dark"
      "${cli}/bin/caelestia wallpaper -f ${wallpaper} -N"
    ];
  };
}
