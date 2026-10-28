{
  lib,
  config,
  pkgs,
  ...
}: let
  bruno = pkgs.unstable.bruno;

  # Bruno reads OAuth2 secrets straight from process.env, and GUI launches (app menu,
  # launcher) don't get a shell's environment. Decrypt the sops files here before handing off.
  wrapper = pkgs.writeShellApplication {
    name = "bruno";
    runtimeInputs = [pkgs.sops];
    text = ''
      if [ -f "$HOME/.config/sops/age/keys.txt" ]; then
        set -a
        for f in "${config.dotfiles.repo}/secrets/secrets.yaml" "$HOME"/.config/secrets.d/*.yaml; do
          [ -e "$f" ] || continue
          # shellcheck disable=SC1090
          . <(sops -d --output-type dotenv "$f")
        done
        set +a
      fi
      exec ${lib.getExe bruno} "$@"
    '';
  };
in {
  options.dotfiles.bruno.enable = lib.mkEnableOption "Bruno API client and its CLI";

  config = lib.mkIf config.dotfiles.bruno.enable {
    home.packages = [wrapper pkgs.unstable.bruno-cli];

    xdg.desktopEntries.bruno = {
      name = "Bruno";
      exec = "${lib.getExe wrapper} %U";
      icon = "${bruno}/share/icons/hicolor/256x256/apps/bruno.png";
      comment = "Opensource API Client for Exploring and Testing APIs";
      categories = ["Development"];
      mimeType = ["x-scheme-handler/bruno"];
      settings.StartupWMClass = "Bruno";
    };
  };
}
