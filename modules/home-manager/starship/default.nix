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
        palette = "catppuccin_mocha";
        palettes.catppuccin_mocha = {
          green = "#a6e3a1";
          lavender = "#b4befe";
          mauve = "#cba6f7";
          peach = "#fab387";
          red = "#f38ba8";
        };
        character = {
          success_symbol = "[[󰄛](green) ❯](peach)";
          error_symbol = "[[󰄛](red) ❯](peach)";
        };
        git_branch.style = "bold mauve";
        directory = {
          truncation_length = 4;
          style = "bold lavender";
        };
      };
    };
  };
}
