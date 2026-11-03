{
  lib,
  config,
  ...
}: {
  options.dotfiles.laptop.enable = lib.mkEnableOption "laptop services: battery, power profiles, Bluetooth";

  config = lib.mkIf config.dotfiles.laptop.enable {
    services.upower.enable = true;
    services.power-profiles-daemon.enable = true;
    hardware.bluetooth.enable = true;
    services.blueman.enable = true;
  };
}
