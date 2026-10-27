{
  lib,
  config,
  ...
}: {
  options.dotfiles.git.enable = lib.mkEnableOption "Git";

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
      };
    };

    programs.delta = {
      enable = true;
      enableGitIntegration = true;
      options.navigate = true;
    };
  };
}
