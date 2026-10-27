{
  lib,
  config,
  ...
}: {
  options.dotfiles.nushell.enable = lib.mkEnableOption "Nushell";

  config = lib.mkIf config.dotfiles.nushell.enable {
    programs.nushell = {
      enable = true;
      settings = {
        edit_mode = "vi";
        show_banner = false;
        highlight_resolved_externals = true;
        color_config = {
          shape_external_resolved = "green";
          shape_external = "red";
        };
      };
      environmentVariables.EDITOR = "nvim";
      shellAliases = {
        gst = "git status";
        gsh = "git stash";
        gsha = "git stash apply";
        gl = "git log | less";
        ga = "git add .";
        gp = "git push";
        gpl = "git pull";
        gch = "git checkout";
        gcd = "git checkout dev";
        gcb = "git checkout -b";
        grv = "git checkout .";
        gfc = "git fetch";
        gfd = "git fetch origin dev:dev";
        gmd = "git merge origin/dev";
        grb = "git rebase origin/dev";
      };
      extraConfig = builtins.readFile ./config.nu;
    };

    programs.carapace.enable = true;

    programs.zoxide = {
      enable = true;
      options = ["--cmd cd"];
    };
  };
}
