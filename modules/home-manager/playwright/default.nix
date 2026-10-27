{
  lib,
  config,
  pkgs,
  ...
}: let
  playwright-mcp = pkgs.writeShellApplication {
    name = "playwright-mcp";
    runtimeInputs = [pkgs.nodejs];
    text = ''
      exec npx -y @playwright/mcp@latest --browser chrome --executable-path ${lib.getExe pkgs.google-chrome} "$@"
    '';
  };
in {
  options.dotfiles.playwright.enable = lib.mkEnableOption "Playwright MCP with Chrome and the URL bar init script";

  config = lib.mkIf config.dotfiles.playwright.enable {
    home.packages = [playwright-mcp];

    xdg.configFile."playwright/url-bar.js".source =
      config.lib.file.mkOutOfStoreSymlink "${config.dotfiles.repo}/modules/home-manager/playwright/config/url-bar.js";
  };
}
