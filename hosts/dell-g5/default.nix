{
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
  environment.sessionVariables = {
    LIBVA_DRIVER_NAME = "iHD";
    __EGL_VENDOR_LIBRARY_FILENAMES = "/run/opengl-driver/share/glvnd/egl_vendor.d/50_mesa.json";
  };
  environment.systemPackages = [
    (pkgs.writeShellScriptBin "nvidia-offload" ''
      export __NV_PRIME_RENDER_OFFLOAD=1
      export __NV_PRIME_RENDER_OFFLOAD_PROVIDER=NVIDIA-G0
      export __GLX_VENDOR_LIBRARY_NAME=nvidia
      export __VK_LAYER_NV_optimus=NVIDIA_only
      export __EGL_VENDOR_LIBRARY_FILENAMES=/run/opengl-driver/share/glvnd/egl_vendor.d/10_nvidia.json
      exec "$@"
    '')
  ];
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
        enableOffloadCmd = false;
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
