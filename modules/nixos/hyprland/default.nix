{
  lib,
  config,
  pkgs,
  user,
  ...
}: {
  options.dotfiles.hyprland.enable = lib.mkEnableOption "Hyprland desktop";
  options.dotfiles.hyprland.startCommand = lib.mkOption {
    type = lib.types.str;
    default = lib.getExe' config.programs.hyprland.package "start-hyprland";
    description = "Command the login screen runs to start the Hyprland session.";
  };

  config = lib.mkIf config.dotfiles.hyprland.enable {
    home-manager.users.${user} = {
      dotfiles.hyprland.enable = true;
      stylix.targets = {
        neovim.enable = false;
        firefox.enable = false;
        librewolf.enable = false;
        gnome.enable = false;
        kde.enable = false;
        ghostty.enable = false;
        k9s.enable = false;
        lazygit.enable = false;
        starship.enable = false;
        nushell.enable = false;
      };
    };

    programs.hyprland = {
      enable = true;
      xwayland.enable = true;
    };

    services.greetd = {
      enable = true;
      settings.default_session = {
        command = "${lib.getExe pkgs.tuigreet} --time --theme 'border=magenta;text=white;prompt=magenta;time=cyan;action=magenta;button=magenta;input=white' --cmd ${config.dotfiles.hyprland.startCommand}";
        user = "greeter";
      };
    };

    environment.sessionVariables.NIXOS_OZONE_WL = "1";

    programs.gpu-screen-recorder.enable = true;

    services.gnome.gnome-keyring.enable = true;
    services.gnome.gcr-ssh-agent.enable = false;
    programs.ssh.startAgent = true;
    security.pam.services.greetd.enableGnomeKeyring = true;
    services.fwupd.enable = true;

    xdg.portal = {
      enable = true;
      extraPortals = [pkgs.xdg-desktop-portal-gtk];
    };
    environment.systemPackages = [
      (lib.hiPrio (pkgs.writeShellScriptBin "xdg-open" ''exec ${lib.getExe' pkgs.glib "gio"} open "$@"''))
    ];

    stylix = {
      enable = true;
      autoEnable = true;
      base16Scheme = "${pkgs.base16-schemes}/share/themes/catppuccin-mocha.yaml";
      polarity = "dark";
      fonts = {
        monospace = {
          package = pkgs.nerd-fonts.jetbrains-mono;
          name = "JetBrainsMono Nerd Font Mono";
        };
      };
      fonts.sizes.popups = 10;
      icons = {
        enable = true;
        package = pkgs.catppuccin-papirus-folders.override {
          flavor = "mocha";
          accent = "mauve";
        };
        dark = "Papirus-Dark";
        light = "Papirus-Light";
      };
    };

    fonts = {
      enableDefaultPackages = true;
      packages = [
        pkgs.nerd-fonts.jetbrains-mono
        pkgs.nerd-fonts.fantasque-sans-mono
        pkgs.fantasque-sans-mono
        pkgs.victor-mono
        pkgs.poppins
        (pkgs.runCommand "apple-logo-font" {nativeBuildInputs = [pkgs.fontforge];} ''
          fontforge -lang=py -c '
          import fontforge
          src = fontforge.open("${pkgs.font-awesome}/share/fonts/opentype/Font Awesome 7 Brands-Regular-400.otf")
          src.selection.select(0xf179)
          src.copy()
          out = fontforge.font()
          out.encoding = "UnicodeFull"
          out.em = src.em
          out.ascent = src.ascent
          out.descent = src.descent
          out.fontname = "AppleLogo"
          out.familyname = "Apple Logo"
          out.fullname = "Apple Logo"
          out.selection.select(0xf8ff)
          out.paste()
          out.generate("AppleLogo.ttf")
          '
          install -Dm444 AppleLogo.ttf $out/share/fonts/truetype/AppleLogo.ttf
        '')
        (pkgs.runCommand "google-sans-flex" {} ''
          install -Dm444 ${pkgs.fetchurl {
            name = "GoogleSansFlex.ttf";
            url = "https://github.com/google/fonts/raw/6c1dbe27dcf7d817d4b62cec033ba77fa8e2a5ec/ofl/googlesansflex/GoogleSansFlex%5BROND,opsz,slnt,wdth,wght%5D.ttf";
            hash = "sha256-qax3B5iFfV6KxZPO1Csi5Ap8CM7Jgxn5Kw7Mwng6d7o=";
          }} $out/share/fonts/truetype/GoogleSansFlex.ttf
        '')
      ];
    };
  };
}
