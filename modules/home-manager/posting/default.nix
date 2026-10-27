{
  lib,
  config,
  pkgs,
  ...
}: {
  options.dotfiles.posting.enable = lib.mkEnableOption "Posting";

  config = lib.mkIf config.dotfiles.posting.enable {
    home.packages = [pkgs.posting];

    xdg.configFile."posting/config.yaml".source = (pkgs.formats.yaml {}).generate "posting-config" {
      editor = "nvim";
      pager = "nvim -R";
      pager_json = "nvim -R";
      response = {prettify_json = true;};
      theme = "catppuccin-mocha-blue";
      url_bar = {hide_secrets_in_value_preview = true;};
      use_host_environment = true;
      watch_env_files = true;
    };

    xdg.dataFile."posting/themes/catppuccin-mocha-blue.yaml".source = ./catppuccin-mocha-blue.yaml;
  };
}
