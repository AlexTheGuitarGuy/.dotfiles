{
  lib,
  config,
  ...
}: {
  options.dotfiles.starship.enable = lib.mkEnableOption "Starship prompt";

  config = lib.mkIf config.dotfiles.starship.enable {
    programs.starship = {
      enable = true;
      settings = {
        palette = "catppuccin";
        palettes.catppuccin = {inherit (config.dotfiles.theme.palette) green lavender accent peach red;};
        character = {
          success_symbol = "[[󰊠](green) ❯](peach)";
          error_symbol = "[[󰊠](red) ❯](peach)";
        };
        git_branch.style = "bold accent";
        directory = {
          truncation_length = 4;
          style = "bold lavender";
        };
      };
    };
  };
}
