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
        character = {
          success_symbol = "[[󰊠](green) ❯](orange)";
          error_symbol = "[[󰊠](red) ❯](orange)";
        };
        git_branch.style = "bold purple";
        directory = {
          truncation_length = 4;
          style = "bold base07";
        };
      };
    };
  };
}
