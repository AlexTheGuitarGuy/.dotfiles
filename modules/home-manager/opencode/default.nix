{
  lib,
  config,
  pkgs,
  ...
}: let
  live = config.lib.file.mkOutOfStoreSymlink;
  repo = "${config.dotfiles.repo}/modules/home-manager/opencode/config";
  link = dir:
    lib.mapAttrs' (name: _: lib.nameValuePair "opencode/${dir}${name}" {source = live "${repo}/${dir}${name}";})
    (lib.filterAttrs (_: type: type == "regular") (builtins.readDir ./config/${dir}));
in {
  options.dotfiles.opencode.enable = lib.mkEnableOption "OpenCode";

  config = lib.mkIf config.dotfiles.opencode.enable {
    # nodejs for the npx-launched MCP servers in opencode.json
    home.packages = [pkgs.unstable.opencode pkgs.nodejs];

    # Each file is linked on its own: OpenCode installs its plugins into ~/.config/opencode
    # (node_modules, package.json), so the directory itself stays writable.
    xdg.configFile = link "" // link "agents/" // link "plugins/";
  };
}
