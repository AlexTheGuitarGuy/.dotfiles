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

  # 1Password writes this file itself on other distros. Here ~/.mozilla/native-messaging-hosts
  # is a read-only link, so it is declared, pointing at the NixOS setgid wrapper.
  onePasswordHost = pkgs.writeTextDir "lib/mozilla/native-messaging-hosts/com.1password.1password.json" (builtins.toJSON {
    name = "com.1password.1password";
    description = "1Password BrowserSupport";
    path = "/run/wrappers/bin/1Password-BrowserSupport";
    type = "stdio";
    allowed_extensions = [
      "{0a75d802-9aed-41e7-8daa-24c067386e82}"
      "{25fc87fa-4d31-4fee-b5c1-c32a7844c063}"
      "{d634138d-c276-4fc8-924b-40a0ea21d284}"
    ];
  });

  brotabHost = pkgs.writeTextDir "lib/mozilla/native-messaging-hosts/brotab_mediator.json" (builtins.toJSON {
    name = "brotab_mediator";
    description = "brotab tab control for the caelestia launcher";
    path = "${pkgs.brotab}/bin/bt_mediator";
    type = "stdio";
    allowed_extensions = ["brotab_mediator@example.org"];
  });
in {
  options.dotfiles.librewolf = {
    enable = lib.mkEnableOption "LibreWolf";
    profilePath = lib.mkOption {
      type = lib.types.str;
      readOnly = true;
      default = "${configPath}/default";
      description = "The LibreWolf profile directory, relative to home. Other modules link profile files here.";
    };
    fingerprintingExemptions = lib.mkOption {
      type = lib.types.listOf lib.types.str;
      default = [];
      description = "Base domains that get no fingerprinting protection (privacy.fingerprintingProtection.granularOverrides).";
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
        settings."privacy.fingerprintingProtection.granularOverrides" = builtins.toJSON (map (domain: {
            firstPartyDomain = domain;
            overrides = "-AllTargets";
          })
          cfg.fingerprintingExemptions);
      };
      policies.ExtensionSettings."brotab_mediator@example.org" = {
        installation_mode = "normal_installed";
        install_url = "https://addons.mozilla.org/firefox/downloads/latest/brotab/latest.xpi";
      };
    };

    dotfiles.librewolf.fingerprintingExemptions = ["google.com"];

    mozilla.firefoxNativeMessagingHosts = [pkgs.tridactyl-native onePasswordHost brotabHost];

    home.file = {
      "${cfg.profilePath}/chrome/userChrome.css".text = config.dotfiles.theme.recolor (builtins.readFile ./config/userChrome.css);
      "${cfg.profilePath}/chrome/userContent.css".text = config.dotfiles.theme.recolor (builtins.readFile ./config/userContent.css);
      "${configPath}/librewolf.overrides.cfg".source = live "${repo}/config/librewolf.overrides.cfg";
      ".config/librewolf/darkreader.json".source = live "${repo}/config/darkreader.json";
      ".config/librewolf/search-shortcuts.html".source = live "${repo}/config/search-shortcuts.html";
      ".config/librewolf/native-messaging-hosts".source = live "${config.home.homeDirectory}/.mozilla/native-messaging-hosts";
    };

    xdg.configFile."tridactyl/tridactylrc".source = live "${repo}/config/tridactylrc";
    xdg.configFile."tridactyl/themes/catppuccin.css".text = config.dotfiles.theme.recolor (builtins.readFile ./config/tridactyl.css);

    home.activation.librewolfDefaultBrowser = lib.hm.dag.entryAfter ["writeBoundary"] ''
      run ${pkgs.xdg-utils}/bin/xdg-mime default librewolf.desktop x-scheme-handler/http x-scheme-handler/https text/html application/xhtml+xml
    '';
  };
}
