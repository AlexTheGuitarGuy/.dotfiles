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
    env.GITNEXUS_HOOK_CLI_PATH = "${pkgs.gitnexus}/lib/node_modules/gitnexus/dist/cli/index.js";
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
    mcpServers = lib.mkOption {
      type = lib.types.attrsOf (pkgs.formats.json {}).type;
      default = {};
      description = "User-scope MCP servers, written into ~/.claude.json on activation. Servers added by hand are kept.";
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [pkgs.unstable.claude-code ai-commit-msgs claude-settings-build pkgs.evidence];

    dotfiles.claude.mcpServers = {
      gitnexus = {
        type = "stdio";
        command = "gitnexus";
        args = ["mcp"];
        env = {};
      };
      playwright = {
        type = "stdio";
        command = "playwright-mcp";
        args = [
          "--user-data-dir"
          "${config.home.homeDirectory}/.cache/playwright-mcp"
          "--output-dir"
          "${config.home.homeDirectory}/.cache/playwright-mcp-out"
          "--save-session"
          "--caps"
          "devtools"
          "--init-script"
          "${config.home.homeDirectory}/.config/playwright/url-bar.js"
        ];
        env = {};
      };
      context7 = {
        type = "http";
        url = "https://mcp.context7.com/mcp";
        headersHelper = "${lib.getExe pkgs.sops} -d --extract '[\"CONTEXT7_API_KEY\"]' ${config.dotfiles.repo}/secrets/secrets.yaml | ${lib.getExe pkgs.jq} -Rc '{CONTEXT7_API_KEY: .}'";
      };
    };

    home.file =
      {
        ".claude/hooks".source = live "${repo}/hooks";
      }
      // lib.mapAttrs' (name: _: lib.nameValuePair ".claude/rules/${name}" {source = live "${repo}/rules/${name}";}) (builtins.readDir ./rules)
      // lib.mapAttrs' (name: _: lib.nameValuePair ".claude/skills/${name}" {source = live "${repo}/skills/${name}";}) (builtins.readDir ./skills);

    home.activation.claudeMcpServers = lib.hm.dag.entryAfter ["writeBoundary"] ''
      claude_json=${config.home.homeDirectory}/.claude.json
      [ -e "$claude_json" ] || run sh -c 'umask 077; echo "{}" > "$1"' _ "$claude_json"
      tmp=$(mktemp "$claude_json.XXXXXX")
      run ${lib.getExe pkgs.jq} --slurpfile m ${(pkgs.formats.json {}).generate "claude-mcp-servers.json" cfg.mcpServers} \
        '.mcpServers = ((.mcpServers // {}) + $m[0])' "$claude_json" > "$tmp"
      run mv "$tmp" "$claude_json"
    '';

    # settings.json is a real file, not a store link: Claude rewrites it (/config). Hand edits are kept, with a warning.
    home.activation.claudeSettings =lib.hm.dag.entryAfter ["writeBoundary"] ''
      run mkdir -p ${config.home.homeDirectory}/.claude
      if ! run ${lib.getExe claude-settings-build} ${(pkgs.formats.json {}).generate "claude-settings.base.json" settings} ${lib.escapeShellArgs cfg.settingsOverlays}; then
        warnEcho "claude settings.json left unchanged, see the message above"
      fi
    '';
  };
}
