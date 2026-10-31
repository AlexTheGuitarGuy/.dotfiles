{
  lib,
  config,
  pkgs,
  ...
}: {
  options.dotfiles.gpg.enable = lib.mkEnableOption "GnuPG with its agent";

  config = lib.mkIf config.dotfiles.gpg.enable {
    programs.gpg.enable = true;

    services.gpg-agent = {
      enable = true;
      pinentry.package = pkgs.pinentry-qt;
    };
  };
}
