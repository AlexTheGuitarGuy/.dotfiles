{
  lib,
  config,
  pkgs,
  ...
}: let
  mcpVersion = "0.0.84";
  playwright-mcp = pkgs.writeShellApplication {
    name = "playwright-mcp";
    runtimeInputs = [pkgs.nodejs];
    text = ''
      if ! ls -d "$HOME"/.cache/ms-playwright/ffmpeg-* >/dev/null 2>&1; then
        npx -y --prefer-offline -p @playwright/mcp@${mcpVersion} playwright install ffmpeg >&2 || true
      fi
      exec npx -y --prefer-offline @playwright/mcp@${mcpVersion} --browser chrome --executable-path ${lib.getExe pkgs.google-chrome} "$@"
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
