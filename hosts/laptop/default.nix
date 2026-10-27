{
  pkgs,
  user,
  ...
}: {
  imports = [./hardware-configuration.nix];

  dotfiles.hyprland.enable = true;

  networking.hostName = "laptop";
  networking.networkmanager.enable = true;

  boot.loader.systemd-boot.enable = true;
  boot.loader.efi.canTouchEfiVariables = true;

  time.timeZone = "Europe/Chisinau";
  i18n.defaultLocale = "en_US.UTF-8";

  nix.settings.experimental-features = ["nix-command" "flakes"];
  nixpkgs.config.allowUnfree = true;

  users.users.${user} = {
    isNormalUser = true;
    extraGroups = ["wheel" "networkmanager"];
    shell = pkgs.nushell;
  };

  virtualisation.vmVariant = {
    users.users.${user}.initialPassword = "test";
    virtualisation = {
      memorySize = 2048;
      cores = 4;
      qemu.options = ["-device virtio-vga"];
    };
  };

  system.stateVersion = "26.05";
}
