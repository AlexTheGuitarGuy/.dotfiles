{
  lib,
  config,
  pkgs,
  ...
}: let
  script = name: runtimeInputs:
    pkgs.writeShellApplication {
      inherit name runtimeInputs;
      text = builtins.readFile ./${name}.sh;
    };
  git-todos = script "git-todos" [pkgs.git];
  worktree-gitnexus = script "worktree-gitnexus" [pkgs.git pkgs.jq];
  gitnexus-reindex = script "gitnexus-reindex" [pkgs.git pkgs.jq pkgs.util-linux];
  git-wt-clean = script "git-wt-clean" [pkgs.git pkgs.jq pkgs.gawk];
  worktree-pnpm = pkgs.writeShellApplication {
    name = "worktree-pnpm";
    runtimeInputs = [pkgs.git pkgs.gnugrep];
    text = "PNPM=${lib.escapeShellArg config.dotfiles.git.worktreePnpm}\n" + builtins.readFile ./worktree-pnpm.sh;
  };
in {
  options.dotfiles.git = {
    enable = lib.mkEnableOption "Git";
    worktreePnpm = lib.mkOption {
      type = lib.types.str;
      default = lib.getExe pkgs.pnpm;
      description = "pnpm command run in new worktrees of pnpm repos; override to add registry credentials.";
    };
  };

  config = lib.mkIf config.dotfiles.git.enable {
    programs.git = {
      enable = true;
      ignores = [
        ".excluded/"
        ".gitnexus/"
        ".99-tmp/"
        "**/.claude/settings.local.json"
      ];
      settings = {
        user = {
          name = "AlexSeremet";
          email = "seremet.alexandru.work@gmail.com";
        };
        pull.rebase = true;
        merge.conflictstyle = "diff3";
        hook = {
          todos = {
            event = "post-checkout";
            command = lib.getExe git-todos;
          };
          gitnexus-worktree = {
            event = "post-checkout";
            command = lib.getExe worktree-gitnexus;
          };
          pnpm-worktree = {
            event = "post-checkout";
            command = lib.getExe worktree-pnpm;
          };
          gitnexus-push = {
            event = "pre-push";
            command = "${lib.getExe gitnexus-reindex} push";
          };
          gitnexus-pull = {
            event = ["post-merge" "post-rewrite"];
            command = "${lib.getExe gitnexus-reindex} pull";
          };
        };
      };
    };

    home.packages = [git-wt-clean pkgs.gitnexus];

    programs.delta = {
      enable = true;
      enableGitIntegration = true;
      options.navigate = true;
    };
  };
}
