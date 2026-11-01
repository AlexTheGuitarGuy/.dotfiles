{
  lib,
  config,
  pkgs,
  ...
}: let
  script = name: runtimeInputs:
    pkgs.writeShellApplication {
      inherit name runtimeInputs;
      text = builtins.readFile ./${name}.sh;
    };
in {
  options.dotfiles.shots.enable = lib.mkEnableOption "term-shot and k9s-shot proof screenshots";

  config = lib.mkIf config.dotfiles.shots.enable {
    wayland.windowManager.hyprland.settings.window_rule = [
      {
        match.class = "^(shot\\..+)$";
        workspace = "7";
      }
      {
        match.class = "^(google-chrome)$";
        workspace = "7 silent";
      }
    ];

    home.packages =
      (with pkgs; [wtype grim slurp imagemagick ffmpeg])
      ++ [
        (script "term-shot" (with pkgs; [ghostty jq wtype grim]))
        (script "k9s-shot" (with pkgs; [ghostty jq grim tmux]))
      ];
  };
}
