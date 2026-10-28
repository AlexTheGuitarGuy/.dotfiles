{
  lib,
  config,
  pkgs,
  ...
}: {
  options.dotfiles.fastfetch.enable = lib.mkEnableOption "fastfetch";

  config = lib.mkIf config.dotfiles.fastfetch.enable {
    programs.fastfetch = {
      enable = true;
      settings = {
        "$schema" = "https://github.com/fastfetch-cli/fastfetch/raw/dev/doc/json_schema.json";
        display = {separator = " 󰑃  ";};
        logo = {padding = {top = 1;};};
        modules = [
          "break"
          {
            key = " DISTRO";
            keyColor = "yellow";
            type = "os";
          }
          {
            key = "│ ├";
            keyColor = "yellow";
            type = "kernel";
          }
          {
            key = "│ ├󰏖";
            keyColor = "yellow";
            type = "packages";
          }
          {
            key = "│ └";
            keyColor = "yellow";
            type = "shell";
          }
          {
            key = " DE/WM";
            keyColor = "blue";
            type = "wm";
          }
          {
            key = "│ ├󰉼";
            keyColor = "blue";
            type = "wmtheme";
          }
          {
            key = "│ ├󰀻";
            keyColor = "blue";
            type = "icons";
          }
          {
            key = "│ ├";
            keyColor = "blue";
            type = "cursor";
          }
          {
            key = "│ ├";
            keyColor = "blue";
            type = "terminalfont";
          }
          {
            key = "│ └";
            keyColor = "blue";
            type = "terminal";
          }
          {
            key = "󰌢 SYSTEM";
            keyColor = "green";
            type = "host";
          }
          {
            key = "│ ├󰻠";
            keyColor = "green";
            type = "cpu";
          }
          {
            format = "{2}";
            key = "│ ├󰻑";
            keyColor = "green";
            type = "gpu";
          }
          {
            compactType = "original-with-refresh-rate";
            key = "│ ├󰍹";
            keyColor = "green";
            type = "display";
          }
          {
            key = "│ ├󰾆";
            keyColor = "green";
            type = "memory";
          }
          {
            key = "│ ├󰓡";
            keyColor = "green";
            type = "swap";
          }
          {
            key = "│ ├󰅐";
            keyColor = "green";
            type = "uptime";
          }
          {
            key = "│ └󰍹";
            keyColor = "green";
            type = "display";
          }
          {
            format = "{2}";
            key = " AUDIO";
            keyColor = "magenta";
            type = "sound";
          }
          {
            key = "│ ├󰥠";
            keyColor = "magenta";
            type = "player";
          }
          {
            key = "│ └󰝚";
            keyColor = "magenta";
            type = "media";
          }
          {
            format = "[90m  [31m  [32m  [33m  [34m  [35m  [36m  [37m  [38m  [39m  [39m    [38m  [37m  [36m  [35m  [34m  [33m  [32m  [31m  [90m ";
            type = "custom";
          }
          "break"
        ];
      };
    };

    xdg.configFile = {
      "fastfetch/config-compact.jsonc".source = ./config-compact.jsonc;
      "fastfetch/config-pokemon.jsonc".source = ./config-pokemon.jsonc;
      "fastfetch/config-v2.jsonc".source = ./config-v2.jsonc;
      "fastfetch/nixos.png".source = "${pkgs.nixos-icons}/share/icons/hicolor/256x256/apps/nix-snowflake.png";
    };
  };
}
