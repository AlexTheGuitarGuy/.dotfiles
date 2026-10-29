{
  pkgs,
  user,
  ...
}: {
  imports = [./hardware-configuration.nix];

  dotfiles.hyprland.enable = true;
  dotfiles.apps.enable = true;
  dotfiles.docker.enable = true;

  # Intel UHD 630 drives the desktop and the laptop screen; the RTX 2060 (external HDMI/DP ports)
  # runs apps started with `nvidia-offload`.
  services.xserver.videoDrivers = ["modesetting" "nvidia"];
  hardware.graphics.enable = true;
  hardware.nvidia = {
    open = true;
    modesetting.enable = true;
    powerManagement = {
      enable = true;
      finegrained = true;
    };
    prime = {
      offload = {
        enable = true;
        enableOffloadCmd = true;
      };
      intelBusId = "PCI:0@0:2:0";
      nvidiaBusId = "PCI:1@0:0:0";
    };
  };

  networking.hostName = "dell-g5";
  networking.networkmanager.enable = true;

  boot.loader.systemd-boot = {
    enable = true;
    configurationLimit = 10;
  };
  boot.loader.efi.canTouchEfiVariables = true;
  zramSwap.enable = true;

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
