{
  lib,
  config,
  options,
  pkgs,
  inputs,
  ...
}: let
  spicePkgs = inputs.spicetify-nix.legacyPackages.${pkgs.stdenv.hostPlatform.system};
in {
  imports = [inputs.spicetify-nix.homeManagerModules.spicetify];

  options.dotfiles.spotify.enable = lib.mkEnableOption "Spotify patched with spicetify";

  config = lib.mkIf config.dotfiles.spotify.enable ({
      programs.spicetify = {
        enable = true;
        theme = spicePkgs.themes.catppuccin;
        colorScheme = "mocha";
        enabledExtensions = [spicePkgs.extensions.adblock];
      };
    }
    // lib.optionalAttrs (options ? stylix) {
      stylix.targets.spicetify.enable = false;
    });
}
