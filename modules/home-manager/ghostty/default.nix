{
  lib,
  config,
  ...
}: {
  options.dotfiles.ghostty.enable = lib.mkEnableOption "Ghostty";

  config = lib.mkIf config.dotfiles.ghostty.enable {
    programs.ghostty = {
      enable = true;
      settings = {
        font-family = "JetBrainsMono Nerd Font";
        font-size = 10;
        app-notifications = "no-clipboard-copy";
        theme = "light:Catppuccin Latte,dark:Catppuccin Mocha";
      };
    };
  };
}
