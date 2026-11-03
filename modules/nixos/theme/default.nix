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
    tertiary = lib.mkOption {
      type = lib.types.enum (import ../../catppuccin.nix).accents;
      default = "lavender";
      description = "Second Catppuccin accent, for things that should stand apart from the main accent.";
    };
    polarity = lib.mkOption {
      type = lib.types.enum ["dark" "light"];
      default = "dark";
      description = "Dark uses Catppuccin Mocha, light uses Catppuccin Latte.";
    };
  };

  config.home-manager.users.${user}.dotfiles.theme = {inherit (cfg) accent tertiary polarity;};
}
