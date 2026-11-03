{
  lib,
  config,
  pkgs,
  ...
}: let
  p = config.dotfiles.theme.palette;
in {
  options.dotfiles.posting.enable = lib.mkEnableOption "Posting";

  config = lib.mkIf config.dotfiles.posting.enable {
    home.packages = [pkgs.posting];

    xdg.configFile."posting/config.yaml".source = (pkgs.formats.yaml {}).generate "posting-config" {
      editor = "nvim";
      pager = "nvim -R";
      pager_json = "nvim -R";
      response = {prettify_json = true;};
      theme = "catppuccin";
      url_bar = {hide_secrets_in_value_preview = true;};
      use_host_environment = true;
      watch_env_files = true;
    };

    xdg.dataFile."posting/themes/catppuccin.yaml".source = (pkgs.formats.yaml {}).generate "posting-catppuccin" {
      name = "catppuccin";
      primary = p.accent;
      secondary = p.accent;
      accent = p.lavender;
      background = p.base;
      surface = p.surface0;
      error = p.red;
      success = p.green;
      warning = p.yellow;
      dark = config.dotfiles.theme.polarity == "dark";
      text_area = {
        cursor = "reverse ${p.rosewater}";
        cursor_line = "underline ${p.text}";
        selection = "reverse ${p.overlay2}";
        gutter = "bold ${p.green}";
        matched_bracket = "reverse ${p.subtext1}";
      };
      url = {
        base = "italic ${p.blue}";
        protocol = "bold ${p.teal}";
        separator = p.text;
      };
      syntax = {
        json_key = "italic ${p.blue}";
        json_number = p.peach;
        json_string = p.green;
        json_boolean = p.sky;
        json_null = p.overlay2;
      };
      method = {
        get = "bold ${p.blue}";
        post = "bold ${p.green}";
        put = "bold ${p.yellow}";
        delete = "bold ${p.red}";
        patch = "bold ${p.teal}";
        options = "bold ${p.lavender}";
        head = "bold ${p.mauve}";
      };
    };
  };
}
