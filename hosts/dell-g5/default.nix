{
  lib,
  config,
  pkgs,
  user,
  ...
}: {
  imports = [./hardware-configuration.nix];

  dotfiles.hyprland.enable = true;
  dotfiles.apps.enable = true;
  dotfiles.docker.enable = true;
  dotfiles.laptop.enable = true;

  # Intel UHD 630 drives the desktop and the laptop screen; the RTX 2060 (external HDMI/DP ports)
  # runs apps started with `nvidia-offload`.
  services.xserver.videoDrivers = ["modesetting" "nvidia"];
  hardware.graphics = {
    enable = true;
    extraPackages = [pkgs.intel-media-driver];
  };
  services.udev.extraRules = ''
    SUBSYSTEM=="drm", KERNEL=="card*", KERNELS=="0000:01:00.0", DRIVERS=="nvidia", SYMLINK+="dri/dgpu"
    SUBSYSTEM=="drm", KERNEL=="card*", KERNELS=="0000:00:02.0", DRIVERS=="i915", SYMLINK+="dri/igpu"
  '';
  dotfiles.hyprland.startCommand = lib.getExe (pkgs.writeShellApplication {
    name = "start-hyprland-gpu";
    text = ''
      dgpu="$(basename "$(readlink -f /dev/dri/dgpu)")"
      if grep -qx connected /sys/class/drm/"$dgpu"-*/status 2>/dev/null; then
        export AQ_DRM_DEVICES=/dev/dri/dgpu:/dev/dri/igpu
      fi
      exec ${lib.getExe' config.programs.hyprland.package "start-hyprland"}
    '';
  });

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
  networking.networkmanager.wifi.powersave = false;

  services.printing.enable = true;

  services.pipewire.alsa.support32Bit = true;
  security.rtkit.enable = true;

  boot.loader.systemd-boot = {
    enable = true;
    configurationLimit = 10;
  };
  boot.loader.efi.canTouchEfiVariables = true;
  zramSwap.enable = true;

  time.timeZone = "Europe/Chisinau";
  i18n.defaultLocale = "en_US.UTF-8";
  i18n.extraLocaleSettings = {
    LC_ADDRESS = "ro_RO.UTF-8";
    LC_IDENTIFICATION = "ro_RO.UTF-8";
    LC_MEASUREMENT = "ro_RO.UTF-8";
    LC_MONETARY = "ro_RO.UTF-8";
    LC_NAME = "ro_RO.UTF-8";
    LC_NUMERIC = "en_US.UTF-8";
    LC_PAPER = "ro_RO.UTF-8";
    LC_TELEPHONE = "ro_RO.UTF-8";
    LC_TIME = "ro_RO.UTF-8";
  };

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
