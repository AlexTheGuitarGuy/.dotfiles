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
    programs.fzf.enable = true;
    programs.direnv = {
      enable = true;
      nix-direnv.enable = true;
      config.global = {
        hide_env_diff = true;
        log_filter = "^$";
      };
    };
    programs.tealdeer = {
      enable = true;
      settings.updates.auto_update = true;
    };

    home.packages = with pkgs; [python3 go gcc lm_sensors nvtopPackages.full dnsutils whois bluetuith p7zip wget terraform minikube];
  };
}
