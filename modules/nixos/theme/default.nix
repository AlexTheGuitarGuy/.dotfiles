{
  lib,
  config,
  user,
  ...
}: let
  cfg = config.dotfiles.theme;
in {
  options.dotfiles.theme = {
    accent = lib.mkOption {
      type = lib.types.enum (import ../../catppuccin.nix).accents;
      default = "blue";
      description = "Catppuccin accent colour used across all themed apps.";
    };
    polarity = lib.mkOption {
      type = lib.types.enum ["dark" "light"];
      default = "dark";
      description = "Dark uses Catppuccin Mocha, light uses Catppuccin Latte.";
    };
  };

  config.home-manager.users.${user}.dotfiles.theme = {inherit (cfg) accent polarity;};
}
