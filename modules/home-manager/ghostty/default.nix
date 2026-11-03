{
  lib,
  config,
  ...
}: let
  ansi = {
    mocha = ["45475a" "f38ba8" "a6e3a1" "f9e2af" "89b4fa" "f5c2e7" "94e2d5" "bac2de" "585b70" "f7aec2" "c2ecbf" "fcd682" "aeccfc" "f398da" "b1eae1" "a6adc8"];
    latte = ["bcc0cc" "d20f39" "40a02b" "df8e1d" "1e66f5" "ea76cb" "179299" "5c5f77" "acb0be" "e7103f" "46b02f" "e49931" "3878f6" "ef95d7" "19a1a8" "6c6f85"];
  };
  theme = flavor: let
    p = config.dotfiles.theme.${flavor};
  in {
    palette = lib.imap0 (i: hex: "${toString i}=#${hex}") ansi.${flavor};
    background = p.base;
    foreground = p.text;
    cursor-color = p.accent;
    cursor-text = p.base;
    selection-background = p.accent;
    selection-foreground = p.base;
  };
in {
  options.dotfiles.ghostty.enable = lib.mkEnableOption "Ghostty";

  config = lib.mkIf config.dotfiles.ghostty.enable {
    programs.ghostty = {
      enable = true;
      settings = {
        font-family = "JetBrainsMono Nerd Font";
        font-size = 10;
        app-notifications = "no-clipboard-copy";
        theme = "light:catppuccin-latte,dark:catppuccin-mocha";
        window-decoration = false;
        keybind = ["ctrl+shift+e=unbind" "ctrl+shift+p=unbind"];
      };
      themes = lib.genAttrs ["catppuccin-mocha" "catppuccin-latte"] (name: theme (lib.removePrefix "catppuccin-" name));
    };
  };
}
