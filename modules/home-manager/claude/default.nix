{
  lib,
  config,
  pkgs,
  ...
}: let
  cfg = config.dotfiles.claude;
  live = config.lib.file.mkOutOfStoreSymlink;
  repo = "${config.dotfiles.repo}/modules/home-manager/claude";
  node = lib.getExe pkgs.nodejs;

  script = name: runtimeInputs:
    pkgs.writeShellApplication {
      inherit name runtimeInputs;
      text = builtins.readFile ./${name}.sh;
    };
  claude-settings-build = script "claude-settings-build" [pkgs.jq pkgs.diffutils];
  claude-worktree = script "claude-worktree" [pkgs.git pkgs.jq];
  ai-commit-msgs = script "ai-commit-msgs" [pkgs.git];

  settings = {
    attribution = {
      commit = "";
      pr = "";
    };
    hooks = {
      SessionStart = [
        {
          hooks = [
            {
              type = "command";
              command = "\"${node}\" \"${config.home.homeDirectory}/.claude/hooks/caveman-activate.js\"";
              timeout = 5;
              statusMessage = "Loading caveman mode...";
            }
          ];
        }
        {
          matcher = "*";
          hooks = [
            {
              type = "command";
              command = "bash '${config.home.homeDirectory}/.claude/hooks/herdr-agent-state.sh' session";
              timeout = 10;
            }
          ];
        }
      ];
      UserPromptSubmit = [
        {
          hooks = [
            {
              type = "command";
              command = "\"${node}\" \"${config.home.homeDirectory}/.claude/hooks/caveman-mode-tracker.js\"";
              timeout = 5;
              statusMessage = "Tracking caveman mode...";
            }
          ];
        }
      ];
      PreToolUse = [
        {
          matcher = "Grep|Glob|Bash";
          hooks = [
            {
              type = "command";
              command = "${node} '${config.home.homeDirectory}/.claude/hooks/gitnexus/gitnexus-hook.cjs'";
              timeout = 10;
              statusMessage = "Enriching with GitNexus graph context...";
            }
          ];
        }
      ];
      WorktreeCreate = [
        {
          hooks = [
            {
              type = "command";
              command = "${lib.getExe claude-worktree} create";
              timeout = 120;
            }
          ];
        }
      ];
      WorktreeRemove = [
        {
          hooks = [
            {
              type = "command";
              command = "${lib.getExe claude-worktree} remove";
              timeout = 60;
            }
          ];
        }
      ];
    };
    statusLine = {
      type = "command";
      command = "bash \"${config.home.homeDirectory}/.claude/hooks/caveman-statusline.sh\"";
    };
    enabledPlugins = {
      "typescript-lsp@claude-plugins-official" = true;
      "superpowers@claude-plugins-official" = true;
      "code-review@claude-plugins-official" = true;
      "code-simplifier@claude-plugins-official" = true;
      "pr-review-toolkit@claude-plugins-official" = true;
      "context7@claude-plugins-official" = true;
      "explanatory-output-style@claude-plugins-official" = false;
      "lua-lsp@claude-plugins-official" = true;
      "caveman@caveman" = true;
      "ponytail@ponytail" = true;
    };
    extraKnownMarketplaces = {
      caveman.source = {
        source = "github";
        repo = "JuliusBrussee/caveman";
      };
      ponytail.source = {
        source = "github";
        repo = "DietrichGebert/ponytail";
      };
    };
    tui = "fullscreen";
    voice = {
      enabled = true;
      mode = "hold";
    };
    preferredNotifChannel = "ghostty";
    inputNeededNotifEnabled = true;
    agentPushNotifEnabled = true;
    autoMode = {
      allow = ["$defaults"];
      soft_deny = ["$defaults"];
    };
    voiceEnabled = true;
    model = "opus";
  };
in {
  options.dotfiles.claude = {
    enable = lib.mkEnableOption "Claude Code";
    settingsOverlays = lib.mkOption {
      type = lib.types.listOf lib.types.str;
      default = [];
      description = "JSON files merged over the base settings into ~/.claude/settings.json, in order.";
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [pkgs.unstable.claude-code ai-commit-msgs claude-settings-build];

    home.file =
      {
        ".claude/hooks".source = live "${repo}/hooks";
      }
      // lib.mapAttrs' (name: _: lib.nameValuePair ".claude/rules/${name}" {source = live "${repo}/rules/${name}";}) (builtins.readDir ./rules)
      // lib.mapAttrs' (name: _: lib.nameValuePair ".claude/skills/${name}" {source = live "${repo}/skills/${name}";}) (builtins.readDir ./skills);

    home.activation.claudeSettings = lib.hm.dag.entryAfter ["writeBoundary"] ''
      run mkdir -p ${config.home.homeDirectory}/.claude
      if ! run ${lib.getExe claude-settings-build} ${(pkgs.formats.json {}).generate "claude-settings.base.json" settings} ${lib.escapeShellArgs cfg.settingsOverlays}; then
        warnEcho "claude settings.json left unchanged, see the message above"
      fi
    '';
  };
}
