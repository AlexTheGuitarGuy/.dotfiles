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

  slimblade = "Kensington SlimBlade Pro(2.4GHz Receiver) Kensington SlimBlade Pro Trackball(2.4GHz Receiver)";

  defaultApps = {
    "org.gnome.Loupe.desktop" = ["image/png" "image/jpeg" "image/gif" "image/webp" "image/bmp" "image/tiff" "image/avif" "image/heic" "image/svg+xml"];
    "org.pwmt.zathura.desktop" = ["application/pdf"];
    "nvim.desktop" = ["text/plain" "text/markdown" "text/x-csrc" "text/x-python" "application/json" "application/x-yaml" "application/toml" "application/x-shellscript"];
    "yazi.desktop" = ["inode/directory"];
    "org.gnome.FileRoller.desktop" = ["application/zip" "application/x-tar" "application/gzip" "application/x-compressed-tar" "application/x-xz-compressed-tar" "application/x-7z-compressed" "application/vnd.rar" "application/x-rar"];
    "mpv.desktop" = ["video/mp4" "video/x-matroska" "video/webm" "video/quicktime" "video/x-msvideo" "video/mpeg" "video/ogg" "audio/mpeg" "audio/flac" "audio/ogg" "audio/x-wav" "audio/mp4" "audio/aac" "audio/opus"];
    "onlyoffice-desktopeditors.desktop" = ["application/vnd.oasis.opendocument.text" "application/msword" "application/vnd.openxmlformats-officedocument.wordprocessingml.document" "application/rtf" "application/vnd.oasis.opendocument.spreadsheet" "application/vnd.ms-excel" "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" "text/csv" "application/vnd.oasis.opendocument.presentation" "application/vnd.ms-powerpoint" "application/vnd.openxmlformats-officedocument.presentationml.presentation" "application/vnd.ms-word.document.macroEnabled.12" "application/vnd.openxmlformats-officedocument.wordprocessingml.template" "application/vnd.ms-word.template.macroEnabled.12" "application/vnd.oasis.opendocument.text-template" "application/vnd.ms-excel.sheet.macroEnabled.12" "application/vnd.openxmlformats-officedocument.spreadsheetml.template" "application/vnd.ms-excel.template.macroEnabled.12" "application/vnd.oasis.opendocument.spreadsheet-template" "application/vnd.ms-powerpoint.presentation.macroEnabled.12" "application/vnd.openxmlformats-officedocument.presentationml.slideshow" "application/vnd.openxmlformats-officedocument.presentationml.template" "application/vnd.oasis.opendocument.presentation-template"];
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

        obs-studio
        gimp
        mpv
        snapshot
        reaper
        postgresql_18
        aseprite
        renpy
        anki
        prismlauncher
        wineWow64Packages.stable

        onlyoffice-desktopeditors
        loupe
        zathura
        file-roller
        xdg-terminal-exec
        qbittorrent
        qalculate-gtk
        proton-vpn
        wvkbd
      ])
      # Far behind on stable, and these move fast.
      ++ (with pkgs.unstable; [
        vscode
        jetbrains.idea
        jetbrains.datagrip
        mongodb-compass
        podman-desktop
        godot_4
      ])
      ++ [apps2samsung tizenbrew-device-manager tizenbrew-installer godot-nvim];

    xdg.configFile."xdg-terminals.list".text = "com.mitchellh.ghostty.desktop\n";

    xdg.dataFile."PrismLauncher/themes/catppuccin-mocha-mauve".source = "${pkgs.fetchFromGitHub {
      owner = "catppuccin";
      repo = "prismlauncher";
      rev = "1ef0e745286ce32750abc3bc5d3ca8073a623b60";
      hash = "sha256-RN1VM29OH75GHSHgpu3HW8YSonIEnX7MZzCObJ6BwtQ=";
    }}/themes/mocha/mauve";

    home.activation.defaultApps = lib.hm.dag.entryAfter ["writeBoundary"] (lib.concatStrings (lib.mapAttrsToList (app: types: ''
        run ${pkgs.xdg-utils}/bin/xdg-mime default ${app} ${lib.concatStringsSep " " types}
      '')
      defaultApps));

    home.file.".jdks/jdk-17".source = "${pkgs.jdk17}/lib/openjdk";
    home.file.".jdks/jdk-25".source = "${pkgs.jdk25}/lib/openjdk";

    xdg.configFile."input-remapper-2/presets/${slimblade}/top-buttons-swapped.json".source = ./input-remapper/slimblade.json;
    home.activation.inputRemapperConfig = lib.hm.dag.entryAfter ["writeBoundary"] (let
      cfg = pkgs.writeText "input-remapper-config.json" (builtins.toJSON {
        version = "2.0.1";
        autoload.${slimblade} = "top-buttons-swapped";
      });
    in ''
      run rm -f ${config.xdg.configHome}/input-remapper-2/config.json
      run install -Dm644 ${cfg} ${config.xdg.configHome}/input-remapper-2/config.json
    '');

    dotfiles.hyprland.autostart = ["input-remapper-control --command autoload"];

    home.activation.godotNvim = lib.hm.dag.entryAfter ["writeBoundary"] ''
      run ${lib.getExe godot-nvim} setup || true
    '';
  };
}
