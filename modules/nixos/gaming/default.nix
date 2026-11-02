{
  lib,
  config,
  pkgs,
  user,
  ...
}: {
  options.dotfiles.gaming.enable = lib.mkEnableOption "Steam with gamemode and mangohud";

  config = lib.mkIf config.dotfiles.gaming.enable {
    programs.steam.enable = true;
    programs.gamemode.enable = true;
    hardware.graphics.enable32Bit = true;
    users.users.${user}.packages = [pkgs.mangohud];
  };
}
