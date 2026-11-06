{
  lib,
  config,
  pkgs,
  ...
}: {
  options.dotfiles.laptop.enable = lib.mkEnableOption "laptop services: battery, power profiles, Bluetooth";

  config = lib.mkIf config.dotfiles.laptop.enable {
    services.upower.enable = true;
    services.power-profiles-daemon.enable = true;
    hardware.bluetooth.enable = true;
    services.blueman.enable = true;
    environment.systemPackages = [
      (pkgs.writeShellApplication {
        name = "bt-sync-windows-key";
        runtimeInputs = with pkgs; [chntpw ntfs3g];
        text = builtins.readFile ./bt-sync-windows-key.sh;
      })
    ];
  };
}
