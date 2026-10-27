{
  lib,
  config,
  pkgs,
  inputs,
  ...
}: let
  herdr = inputs.herdr.packages.${pkgs.stdenv.hostPlatform.system}.default;
  pluginDeps = with pkgs; [bash jq fzf zoxide nodejs];

  sessionizer = pkgs.fetchFromGitHub {
    owner = "salkhalil";
    repo = "herdr-sessionizer";
    rev = "0ba66fafbee69628b7215566a94e7cdd02de4e0a";
    hash = "sha256-+iyMABY2RlyF6kC9XkIzZ1Hb7J5M32gw3SNFSdCvz3Y=";
  };
  synchronize-panes = pkgs.fetchFromGitHub {
    owner = "furuhashin";
    repo = "herdr-synchronize-panes";
    rev = "b827d3f3a68ceb01cf72a51546c68f9f33188adb";
    hash = "sha256-ynEMd66VL0Lv+IyvzBA9L86Rbr/sdzB1kqYGJur4PHo=";
  };

  status = pkgs.writeShellApplication {
    name = "herdr-status";
    runtimeInputs = [pkgs.curl];
    bashOptions = [];
    text = builtins.readFile ./status.sh;
  };

  herdr-worktree = pkgs.writeShellApplication {
    name = "herdr-worktree";
    runtimeInputs = [herdr pkgs.jq pkgs.fzf];
    text = builtins.readFile ./herdr-worktree.sh;
  };

  # tgs-equivalent for herdr: fuzzy-pick a project into a herdr workspace
  # (via the sessionizer plugin), then attach the herdr TUI when run from a plain shell.
  #
  #   tzt            interactive picker
  #   tzt <dir>      skip the picker, switch/create directly
  tzt = pkgs.writeShellApplication {
    name = "tzt";
    runtimeInputs = [herdr] ++ pluginDeps;
    text = ''
      bash ${sessionizer}/bin/sessionizer "$@"

      [ -n "''${HERDR_ENV:-}" ] && exit 0

      exec herdr
    '';
  };

  toml = pkgs.formats.toml {};
in {
  options.dotfiles.herdr.enable = lib.mkEnableOption "herdr terminal workspace manager";

  config = lib.mkIf config.dotfiles.herdr.enable {
    home.packages = [herdr herdr-worktree tzt];

    xdg.configFile = {
      "herdr/config.toml".source = toml.generate "herdr-config.toml" {
        onboarding = false;
        keys = {
          prefix = "ctrl+a";
          focus_pane_left = "prefix+h";
          focus_pane_down = "prefix+j";
          focus_pane_up = "prefix+k";
          focus_pane_right = "prefix+l";
          switch_tab = "prefix+1..9";
          switch_workspace = "prefix+shift+1..9";
          detach = "prefix+q";
          workspace_picker = "prefix+w";
          new_tab = "prefix+c";
          close_pane = "prefix+x";
          split_vertical = "prefix+v";
          split_horizontal = "prefix+-";
          zoom = "prefix+z";
          toggle_sidebar = "prefix+b";
          previous_agent = "prefix+ctrl+p";
          next_agent = "prefix+ctrl+n";
          focus_agent = "prefix+alt+1..9";
          last_pane = "prefix+tab";
          copy_mode = "prefix+{";
          previous_workspace = "prefix+[";
          next_workspace = "prefix+]";
          indexed.tabs = "alt";
          command = [
            {
              key = "prefix+f";
              type = "plugin_action";
              command = "sessionizer.pick";
              description = "fuzzy-switch workspace";
            }
            {
              key = "prefix+g";
              type = "popup";
              command = "herdr-worktree";
              width = "60%";
              height = "40%";
            }
            {
              key = "prefix+m";
              type = "plugin_action";
              command = "synchronize-panes.broadcast";
              description = "broadcast a command to all panes in this tab";
            }
            {
              key = "prefix+a";
              type = "plugin_action";
              command = "sessionizer.agents";
              description = "jump to an agent (picker)";
            }
          ];
        };
        ui = {
          mouse_capture = false;
          sidebar_width = 32;
          tab_bar_right = [
            {
              type = "command";
              command = lib.getExe status;
              interval_seconds = 900;
              timeout_seconds = 5;
            }
            {
              type = "datetime";
              format = "%a %H:%M";
            }
          ];
          toast = {
            delivery = "system";
            delay_seconds = 2;
          };
          sound.enabled = true;
        };
        theme = {
          name = "catppuccin";
          auto_switch = true;
          light_name = "catppuccin-latte";
          dark_name = "catppuccin";
        };
        experimental.pane_history = false;
        advanced.scrollback_limit_bytes = 2000000;
      };

      "herdr/plugins/config/sessionizer/config.json".source =
        config.lib.file.mkOutOfStoreSymlink "${config.dotfiles.repo}/modules/home-manager/herdr/sessionizer.config.json";
    };

    home.activation.herdrPlugins = lib.hm.dag.entryAfter ["writeBoundary"] ''
      run ${lib.getExe herdr} plugin link ${sessionizer}
      run ${lib.getExe herdr} plugin link ${synchronize-panes}
    '';

    systemd.user.services.herdr = {
      Unit = {
        Description = "herdr headless terminal multiplexer server";
        Documentation = "https://herdr.dev";
        After = ["graphical-session.target"];
        PartOf = ["graphical-session.target"];
      };
      Service = {
        Type = "simple";
        Environment = [
          "HERDR_LOG=herdr=warn"
          "PATH=${lib.makeBinPath pluginDeps}:${config.home.profileDirectory}/bin:/run/current-system/sw/bin"
        ];
        IOSchedulingClass = "best-effort";
        IOSchedulingPriority = 0;
        ExecStart = "${lib.getExe herdr} server";
        ExecStop = "${lib.getExe herdr} server stop";
        Restart = "on-failure";
        RestartSec = 2;
      };
      Install.WantedBy = ["graphical-session.target"];
    };
  };
}
