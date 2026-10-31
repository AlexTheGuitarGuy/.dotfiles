{
  lib,
  config,
  pkgs,
  ...
}: let
  cfg = config.dotfiles.librewolf;
  live = config.lib.file.mkOutOfStoreSymlink;
  repo = "${config.dotfiles.repo}/modules/home-manager/librewolf";
  configPath = ".config/librewolf/librewolf";
in {
  options.dotfiles.librewolf = {
    enable = lib.mkEnableOption "LibreWolf";
    profilePath = lib.mkOption {
      type = lib.types.str;
      readOnly = true;
      default = "${configPath}/default";
      description = "The LibreWolf profile directory, relative to home. Other modules link profile files here.";
    };
  };

  config = lib.mkIf cfg.enable {
    programs.librewolf = {
      enable = true;
      inherit configPath;
      profiles.default = {
        id = 0;
        path = "default";
        isDefault = true;
      };
      policies.ExtensionSettings."{d634138d-c276-4fc8-924b-40a0ea21d284}".installation_mode = "blocked";
    };

    mozilla.firefoxNativeMessagingHosts = [pkgs.tridactyl-native];

    home.file = {
      "${cfg.profilePath}/chrome".source = live "${repo}/config/chrome";
      "${configPath}/librewolf.overrides.cfg".source = live "${repo}/config/librewolf.overrides.cfg";
      ".config/librewolf/darkreader.json".source = live "${repo}/config/darkreader.json";
      ".config/librewolf/search-shortcuts.html".source = live "${repo}/config/search-shortcuts.html";
      ".config/librewolf/native-messaging-hosts".source = live "${config.home.homeDirectory}/.mozilla/native-messaging-hosts";
    };

    xdg.configFile.tridactyl.source = live "${repo}/tridactyl";

    home.activation.librewolfDefaultBrowser = lib.hm.dag.entryAfter ["writeBoundary"] ''
      run ${pkgs.xdg-utils}/bin/xdg-mime default librewolf.desktop x-scheme-handler/http x-scheme-handler/https text/html application/xhtml+xml
    '';
  };
}
