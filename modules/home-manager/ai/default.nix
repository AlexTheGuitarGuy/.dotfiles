{
  lib,
  config,
  pkgs,
  ...
}: let
  cfg = config.dotfiles.ai;
  json = pkgs.formats.json {};

  claudeServer = s:
    if s.claude != null
    then s.claude
    else if s.url != null
    then {
      type = "http";
      inherit (s) url;
    }
    else {
      type = "stdio";
      inherit (s) command args;
      env = {};
    };

  opencodeServer = s:
    if s.url != null
    then {
      type = "remote";
      inherit (s) url enabled;
    }
    else {
      type = "local";
      command = [s.command] ++ s.args;
      inherit (s) enabled;
    }
    // lib.optionalAttrs (s.timeout != null) {inherit (s) timeout;};

  claudeAgent = name: a: ''
    ---
    name: ${name}
    description: ${a.description}
    tools: ${lib.concatStringsSep ", " a.claudeTools}
    model: ${a.claudeModel}
    ---
    ${a.instructions}'';

  opencodeAgent = _: a: ''
    ---
    description: ${a.description}
    mode: subagent
    permission:
      edit: deny
    ---
    ${a.instructions}'';

  readOnly = {
    "*" = "deny";
    glob = "allow";
    grep = "allow";
    list = "allow";
    lsp = "allow";
    read = "allow";
    skill = "allow";
    webfetch = "allow";
    websearch = "allow";
  };

  opencodeBase = builtins.fromJSON (builtins.readFile ../opencode/config/opencode.json);
in {
  options.dotfiles.ai = {
    mcpServers = lib.mkOption {
      default = {};
      description = "MCP servers shared by Claude Code and OpenCode. Rendered into ~/.claude.json and the OpenCode mcp block.";
      type = lib.types.attrsOf (lib.types.submodule {
        options = {
          command = lib.mkOption {
            type = lib.types.nullOr lib.types.str;
            default = null;
          };
          args = lib.mkOption {
            type = lib.types.listOf lib.types.str;
            default = [];
          };
          url = lib.mkOption {
            type = lib.types.nullOr lib.types.str;
            default = null;
          };
          enabled = lib.mkOption {
            type = lib.types.bool;
            default = true;
            description = "OpenCode enabled flag.";
          };
          timeout = lib.mkOption {
            type = lib.types.nullOr lib.types.ints.positive;
            default = null;
            description = "OpenCode request timeout in ms.";
          };
          claude = lib.mkOption {
            type = lib.types.nullOr json.type;
            default = null;
            description = "Replaces the derived Claude entry, for settings only Claude knows (headersHelper).";
          };
        };
      });
    };

    agents = lib.mkOption {
      default = {};
      description = "Read-only subagents rendered as Claude agents and OpenCode subagents.";
      type = lib.types.attrsOf (lib.types.submodule {
        options = {
          description = lib.mkOption {type = lib.types.str;};
          instructions = lib.mkOption {type = lib.types.lines;};
          claudeModel = lib.mkOption {
            type = lib.types.str;
            default = "sonnet";
          };
          claudeTools = lib.mkOption {
            type = lib.types.listOf lib.types.str;
            default = ["Read" "Grep" "Glob" "Bash"];
          };
        };
      });
    };
  };

  imports = [./agents.nix];

  config = lib.mkMerge [
    (lib.mkIf config.dotfiles.claude.enable {
      dotfiles.claude.mcpServers = lib.mapAttrs (_: claudeServer) cfg.mcpServers;
      home.file = lib.mapAttrs' (name: a: lib.nameValuePair ".claude/agents/${name}.md" {text = claudeAgent name a;}) cfg.agents;
    })

    (lib.mkIf config.dotfiles.opencode.enable {
      xdg.configFile =
        lib.mapAttrs' (name: a: lib.nameValuePair "opencode/agents/${name}.md" {text = opencodeAgent name a;}) cfg.agents
        // {
          "opencode/opencode.json".source = json.generate "opencode.json" (opencodeBase
            // {
              mcp = lib.mapAttrs (_: opencodeServer) cfg.mcpServers;
              agent = {
                general.permission = readOnly;
                explore.permission = readOnly;
              };
            });
        };
    })

    {
      dotfiles.ai.mcpServers = {
        gitnexus = {
          command = "gitnexus";
          args = ["mcp"];
        };
        playwright = {
          command = "playwright-mcp";
          timeout = 60000;
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
        };
        context7 = {
          command = "npx";
          args = ["-y" "@upstash/context7-mcp"];
          claude = {
            type = "http";
            url = "https://mcp.context7.com/mcp";
            headersHelper = "${lib.getExe pkgs.sops} -d --extract '[\"CONTEXT7_API_KEY\"]' ${config.dotfiles.repo}/secrets/secrets.yaml | ${lib.getExe pkgs.jq} -Rc '{CONTEXT7_API_KEY: .}'";
          };
        };
      };
    }
  ];
}
