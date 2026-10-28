{
  lib,
  config,
  pkgs,
  ...
}: let
  appimage = pname: version: url: hash:
    pkgs.appimageTools.wrapType2 {
      inherit pname version;
      src = pkgs.fetchurl {inherit url hash;};
    };

  apps2samsung = appimage "apps2samsung" "3.0.0" "https://github.com/Apps2Samsung/Apps2Samsung/releases/download/v3.0.0/Apps2Samsung-v3.0.0-linux-x64.AppImage" "sha256-lHzltMTC4UdtU1nmbLjrWyh27i3yJSVwgj/CeKxUhqw=";
  tizenbrew-device-manager = appimage "tizenbrew-device-manager" "1.0.0" "https://github.com/reisxd/tizenbrew-device-manager/releases/download/v1.0.0/tizenbrew-device-manager_1.0.0_amd64.AppImage" "sha256-GQMMsxflW0AQas0pIBGN6A2BNri6vVNz8YHHJNd1ui4=";

  # A Node binary packed with its payload; patchelf would break it, so it runs in an FHS env.
  tizenbrew-installer = pkgs.buildFHSEnv {
    name = "tizenbrew-installer";
    runScript = pkgs.fetchurl {
      url = "https://github.com/reisxd/TizenBrewInstaller/releases/download/v1.1.3/TizenBrewInstaller-linux-x64";
      hash = "sha256-DcgNH8X5FInFJX1guGCnwIKqPyxLvD7+ymw/jqEkYKE=";
      executable = true;
    };
  };

  # MoldSign has no public download; its folder (bundled JRE, writes its data next to itself)
  # is copied to ~/Apps/MoldSign by hand. The FHS env gives the JRE and the eToken library
  # the usual /lib layout plus the X11 and smart card libraries.
  moldsign = pkgs.buildFHSEnv {
    name = "moldsign";
    targetPkgs = p:
      with p; [pcsclite zlib freetype fontconfig alsa-lib libx11 libxext libxrender libxtst libxi];
    runScript = "sh ${config.home.homeDirectory}/Apps/MoldSign/MoldSign-Desk.sh";
  };

  godot-nvim = pkgs.writeShellApplication {
    name = "godot-nvim";
    runtimeInputs = with pkgs; [jq python3 procps hostname];
    runtimeEnv.GODOT_NVIM_EXEC = "${config.home.profileDirectory}/bin/godot-nvim";
    text = builtins.readFile ./godot-nvim.sh;
  };
in {
  options.dotfiles.apps.enable = lib.mkEnableOption "desktop apps";

  config = lib.mkIf config.dotfiles.apps.enable {
    home.packages =
      (with pkgs; [
        telegram-desktop
        slack
        teams-for-linux
        viber
        spotify

        obs-studio
        kdePackages.kdenlive
        gimp
        mpv
        snapshot
        reaper
        qpwgraph
        aseprite
        godot_4
        renpy
        anki
        prismlauncher
        wineWow64Packages.stable

        libreoffice-fresh
        qbittorrent
        qalculate-gtk
        proton-vpn
        gnome-network-displays
        wvkbd
      ])
      # Far behind on stable, and these move fast.
      ++ (with pkgs.unstable; [
        vscode
        jetbrains.idea
        jetbrains.datagrip
        mongodb-compass
      ])
      ++ [apps2samsung tizenbrew-device-manager tizenbrew-installer moldsign godot-nvim];

    programs.yazi.enable = true;
  };
}
