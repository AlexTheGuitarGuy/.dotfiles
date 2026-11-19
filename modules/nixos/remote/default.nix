{
  lib,
  config,
  pkgs,
  user,
  ...
}: let
  cfg = config.dotfiles.remote;
in {
  options.dotfiles.remote = {
    enable = lib.mkEnableOption "ssh, mosh and KDE Connect over Tailscale only";
    authorizedKeys = lib.mkOption {
      type = lib.types.listOf lib.types.str;
      default = [];
      description = "SSH public keys allowed to log in as the main user.";
    };
  };

  config = lib.mkIf cfg.enable {
    services.tailscale.enable = true;
    # sshd, mosh-server and KDE Connect stay closed on every other interface.
    networking.firewall.trustedInterfaces = ["tailscale0"];

    services.openssh = {
      enable = true;
      openFirewall = false;
      settings = {
        PasswordAuthentication = false;
        KbdInteractiveAuthentication = false;
        PermitRootLogin = "no";
        X11Forwarding = false;
      };
    };
    users.users.${user}.openssh.authorizedKeys.keys = cfg.authorizedKeys;

    environment.systemPackages = [pkgs.mosh];

    # Not programs.kdeconnect: it opens its ports on every interface.
    home-manager.users.${user}.services.kdeconnect = {
      enable = true;
      indicator = true;
    };
  };
}
