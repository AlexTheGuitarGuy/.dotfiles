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

  defaultApps = {
    "org.gnome.Loupe.desktop" = ["image/png" "image/jpeg" "image/gif" "image/webp" "image/bmp" "image/tiff" "image/avif" "image/heic" "image/svg+xml"];
    "org.pwmt.zathura.desktop" = ["application/pdf"];
    "nvim.desktop" = ["text/plain" "text/markdown" "text/x-csrc" "text/x-python" "application/json" "application/x-yaml" "application/toml" "application/x-shellscript"];
    "yazi.desktop" = ["inode/directory"];
    "org.gnome.FileRoller.desktop" = ["application/zip" "application/x-tar" "application/gzip" "application/x-compressed-tar" "application/x-xz-compressed-tar" "application/x-7z-compressed" "application/vnd.rar" "application/x-rar"];
    "mpv.desktop" = ["video/mp4" "video/x-matroska" "video/webm" "video/quicktime" "video/x-msvideo" "video/mpeg" "video/ogg" "audio/mpeg" "audio/flac" "audio/ogg" "audio/x-wav" "audio/mp4" "audio/aac" "audio/opus"];
    "writer.desktop" = ["application/vnd.oasis.opendocument.text" "application/msword" "application/vnd.openxmlformats-officedocument.wordprocessingml.document" "application/rtf"];
    "calc.desktop" = ["application/vnd.oasis.opendocument.spreadsheet" "application/vnd.ms-excel" "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" "text/csv"];
    "impress.desktop" = ["application/vnd.oasis.opendocument.presentation" "application/vnd.ms-powerpoint" "application/vnd.openxmlformats-officedocument.presentationml.presentation"];
    "draw.desktop" = ["application/vnd.oasis.opendocument.graphics"];
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
        loupe
        zathura
        file-roller
        xdg-terminal-exec
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

    xdg.configFile."xdg-terminals.list".text = "com.mitchellh.ghostty.desktop\n";

    home.activation.defaultApps = lib.hm.dag.entryAfter ["writeBoundary"] (lib.concatStrings (lib.mapAttrsToList (app: types: ''
        run ${pkgs.xdg-utils}/bin/xdg-mime default ${app} ${lib.concatStringsSep " " types}
      '')
      defaultApps));
  };
}
