{
  lib,
  config,
  pkgs,
  ...
}: {
  options.dotfiles.hyprland.enable = lib.mkEnableOption "Hyprland user config";

  config = lib.mkIf config.dotfiles.hyprland.enable {
    wayland.windowManager.hyprland = {
      enable = true;
      package = null;
      portalPackage = null;
      settings = {
        mod._var = "SUPER";
        bind = map (args: {_args = map lib.generators.mkLuaInline args;}) [
          [''mod .. " + RETURN"'' ''hl.dsp.exec_cmd("${lib.getExe pkgs.ghostty}")'']
          [''mod .. " + Q"'' "hl.dsp.window.close()"]
          [''mod .. " + SHIFT + E"'' "hl.dsp.exit()"]
        ];
      };
    };

    home.packages = [pkgs.ghostty];
  };
}
