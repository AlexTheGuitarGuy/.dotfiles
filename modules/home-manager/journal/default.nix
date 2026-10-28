{
  lib,
  config,
  pkgs,
  ...
}: {
  options.dotfiles.journal.enable = lib.mkEnableOption "journal, a work-session/todo/blocked log";

  config = lib.mkIf config.dotfiles.journal.enable {
    home.packages = [
      (pkgs.writeShellApplication {
        name = "journal";
        runtimeInputs = [pkgs.jq];
        text = builtins.readFile ./journal.sh;
      })
    ];
  };
}
