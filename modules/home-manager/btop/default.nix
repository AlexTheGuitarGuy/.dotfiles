{
  lib,
  config,
  pkgs,
  ...
}: {
  options.dotfiles.btop.enable = lib.mkEnableOption "btop";

  config = lib.mkIf config.dotfiles.btop.enable {
    home.packages = [pkgs.btop];

    xdg.dataFile."btop/themes/catppuccin.theme".text = config.dotfiles.theme.recolor (builtins.readFile ./config/themes/catppuccin_mocha.theme);
    xdg.configFile.btop.source = config.lib.file.mkOutOfStoreSymlink "${config.dotfiles.repo}/modules/home-manager/btop/config";
  };
}
