{
  lib,
  config,
  pkgs,
  user,
  ...
}: {
  options.dotfiles.hyprland.enable = lib.mkEnableOption "Hyprland desktop";

  config = lib.mkIf config.dotfiles.hyprland.enable {
    home-manager.users.${user}.dotfiles.hyprland.enable = true;

    programs.hyprland = {
      enable = true;
      xwayland.enable = true;
    };

    services.greetd = {
      enable = true;
      settings.default_session = {
        command = "${lib.getExe pkgs.tuigreet} --time --cmd ${lib.getExe' config.programs.hyprland.package "start-hyprland"}";
        user = "greeter";
      };
    };

    environment.sessionVariables.NIXOS_OZONE_WL = "1";
  };
}
