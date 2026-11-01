{
  lib,
  config,
  pkgs,
  ...
}: {
  options.dotfiles.cli.enable = lib.mkEnableOption "everyday command line tools";

  config = lib.mkIf config.dotfiles.cli.enable {
    programs.bat.enable = true;
    programs.fd.enable = true;
    programs.tealdeer = {
      enable = true;
      settings.updates.auto_update = true;
    };

    home.packages = with pkgs; [python3 lm_sensors nvtopPackages.full];
  };
}
