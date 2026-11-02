{
  lib,
  config,
  ...
}: {
  imports = [
    ./ai
    ./apps
    ./btop
    ./bruno
    ./claude
    ./cli
    ./fastfetch
    ./gh
    ./gh-dash
    ./herdr
    ./journal
    ./ghostty
    ./git
    ./gpg
    ./hyprland
    ./k9s
    ./lazygit
    ./librewolf
    ./nushell
    ./nvim
    ./opencode
    ./playwright
    ./posting
    ./secrets
    ./shots
    ./spotify
    ./ssh
    ./starship
  ];

  options.dotfiles.repo = lib.mkOption {
    type = lib.types.str;
    default = "${config.home.homeDirectory}/.dotfiles";
    description = "Checkout of this repo. Configs that their apps rewrite are linked here instead of into the read-only store.";
  };
}
