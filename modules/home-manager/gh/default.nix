{
  lib,
  config,
  pkgs,
  ...
}: {
  options.dotfiles.gh.enable = lib.mkEnableOption "GitHub CLI";

  config = lib.mkIf config.dotfiles.gh.enable {
    home.packages = [pkgs.gh];

    xdg.configFile."gh/config.yml".source = config.lib.file.mkOutOfStoreSymlink "${config.dotfiles.repo}/modules/home-manager/gh/config.yml";
  };
}
