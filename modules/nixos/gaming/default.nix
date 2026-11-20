{
  lib,
  config,
  pkgs,
  user,
  ...
}: {
  options.dotfiles.gaming.enable = lib.mkEnableOption "Steam with gamemode, mangohud and Sunshine";

  config = lib.mkIf config.dotfiles.gaming.enable {
    programs.steam.enable = true;
    programs.gamemode.enable = true;
    hardware.graphics.enable32Bit = true;
    users.users.${user}.packages = [pkgs.mangohud];

    # Moonlight streaming to a TV without Tailscale. Started by hand (`systemctl --user start sunshine`),
    # so nothing listens on these ports otherwise; pairing needs a PIN from the web UI.
    services.sunshine = {
      enable = true;
      capSysAdmin = true;
      autoStart = false;
    };
    # The ports Moonlight uses, without the web UI (47990), which stays on localhost.
    networking.firewall.allowedTCPPorts = [47984 47989 48010];
    networking.firewall.allowedUDPPorts = [47998 47999 48000 48002 48010];
    # The module turns on avahi, which opens mDNS on every interface; add the host in Moonlight by IP instead.
    services.avahi.enable = false;
  };
}
