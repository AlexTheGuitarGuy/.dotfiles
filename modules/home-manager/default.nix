{
  lib,
  config,
  ...
}: {
  imports = [
    ./btop
    ./fastfetch
    ./gh
    ./gh-dash
    ./ghostty
    ./git
    ./hyprland
    ./k9s
    ./lazygit
    ./nushell
    ./nvim
    ./posting
    ./secrets
    ./starship
  ];

  options.dotfiles.repo = lib.mkOption {
    type = lib.types.str;
    default = "${config.home.homeDirectory}/.dotfiles";
    description = "Checkout of this repo. Configs that their apps rewrite are linked here instead of into the read-only store.";
  };
}
