{
  lib,
  config,
  ...
}: let
  cfg = config.dotfiles.theme;
  catppuccin = import ../../catppuccin.nix;
  withHash = flavor: lib.mapAttrs (_: hex: "#${hex}") (catppuccin.${flavor} // {accent = catppuccin.${flavor}.${cfg.accent};});
  target = catppuccin.${cfg.flavor} // {mauve = catppuccin.${cfg.flavor}.${cfg.accent};};
  rgb = hex: lib.concatMapStringsSep ", " (i: toString (lib.fromHexString (builtins.substring i 2 hex))) [0 2 4];
  variants = hex: [hex (lib.toUpper hex) (rgb hex)];
  names = builtins.attrNames catppuccin.mocha;
  recolor = builtins.replaceStrings (lib.concatMap (n: variants catppuccin.mocha.${n}) names) (lib.concatMap (n: variants target.${n}) names);
  readOnly = description: default:
    lib.mkOption {
      inherit description default;
      readOnly = true;
    };
in {
  options.dotfiles.theme = {
    accent = lib.mkOption {
      type = lib.types.enum catppuccin.accents;
      default = "blue";
      description = "Catppuccin accent colour used across all themed apps.";
    };
    polarity = lib.mkOption {
      type = lib.types.enum ["dark" "light"];
      default = "dark";
      description = "Dark uses Catppuccin Mocha, light uses Catppuccin Latte.";
    };
    flavor = readOnly "Catppuccin flavor for the current polarity." (
      if cfg.polarity == "light"
      then "latte"
      else "mocha"
    );
    palette = readOnly "Colours of the current flavor as #rrggbb, plus `accent`." (withHash cfg.flavor);
    mocha = readOnly "Mocha colours as #rrggbb, plus `accent`." (withHash "mocha");
    latte = readOnly "Latte colours as #rrggbb, plus `accent`." (withHash "latte");
    recolor = readOnly "Turns text written in Mocha with a mauve accent into the current flavor and accent." recolor;
  };
}
