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
  shellScale = 1.0 / config.dotfiles.hyprland.scale;
  shellPx = px: builtins.floor (px * shellScale + 0.5);

  shellTokens = pkgs.runCommand "caelestia-shell-tokens.json" {nativeBuildInputs = [pkgs.gawk pkgs.jq];} ''
    awk -v scale=${toString shellScale} '
      /^class [A-Za-z]+ : / { cls = $2 }
      cls == "SizeTokens" && match($0, /CONFIG_SUBOBJECT\(([A-Za-z]+), ([A-Za-z]+)\)/, m) { key[m[1]] = m[2] }
      match($0, /CONFIG_(GLOBAL_)?PROPERTY\((int|qreal), ([A-Za-z]+), ([0-9.]+)\)/, m) && m[3] !~ /Sweep$|^heightMult$|^ratio$/ {
        n++; pcls[n] = cls; pname[n] = m[3]; pval[n] = (m[2] == "int") ? int(m[4] * scale + 0.5) : m[4] * scale
      }
      END { for (i = 1; i <= n; i++) if (pcls[i] in key) print key[pcls[i]], pname[i], pval[i] }
    ' ${inputs.caelestia-shell}/plugin/src/Caelestia/Config/tokens.hpp \
      | jq -Rn '{sizes: (reduce (inputs | split(" ")) as [$g, $k, $v] ({}; .[$g][$k] = ($v | tonumber)))}' > $out
  '';
in {
  imports = [inputs.caelestia-shell.homeManagerModules.default];

  config = lib.mkIf config.dotfiles.hyprland.enable {
    programs.caelestia = {
      enable = true;
      systemd.enable = true;
      package = inputs.caelestia-shell.packages.${pkgs.stdenv.hostPlatform.system}.with-cli.overrideAttrs (old: {
        patches = (old.patches or []) ++ [./caelestia-calc.patch ./caelestia-search.patch ./caelestia-tabs-files.patch ./caelestia-gpu-card.patch];
        postInstall =
          old.postInstall
          + ''
            wrapProgram $out/bin/caelestia-shell --prefix PATH : ${lib.makeBinPath [pkgs.brotab pkgs.fd pkgs.fzf pkgs.xdg-utils]}
          '';
      });

      settings = {
        appearance = {
          transparency = {
            enabled = true;
            base = 0.9;
            layers = 0.55;
          };
          font =
            lib.genAttrs ["headline" "title" "body" "label"] (_: {family = "JetBrainsMono Nerd Font";})
            // {
              mono.family = "JetBrainsMono Nerd Font Mono";
              clock = "JetBrainsMono Nerd Font";
              workspaces = "JetBrainsMono Nerd Font";
              scale = shellScale;
            };
          padding.scale = shellScale;
          spacing.scale = shellScale;
          rounding.scale = shellScale;
        };

        border = {
          rounding = shellPx 25;
          thickness = shellPx 10;
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
          weatherUnits = "Celsius";
          clockFormat = "TwentyFourHour";
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

    xdg.dataFile."caelestia/wallpapers".source = config.lib.file.mkOutOfStoreSymlink "${config.dotfiles.repo}/modules/home-manager/hyprland/wallpapers";

    # Caelestia rewrites shell.json at runtime, so replace the store link with a writable copy.
    xdg.configFile."caelestia/shell.json".force = true;
    xdg.configFile."caelestia/shell-tokens.json" = lib.mkIf (config.dotfiles.hyprland.scale != 1.0) {source = shellTokens;};
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
      "${pkgs.bash}/bin/bash -c '[ -s ${config.xdg.stateHome}/caelestia/wallpaper/path.txt ] || ${cli}/bin/caelestia wallpaper -f ${wallpaper} -N'"
    ];
  };
}
