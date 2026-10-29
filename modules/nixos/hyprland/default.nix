{
  lib,
  config,
  pkgs,
  user,
  ...
}: {
  options.dotfiles.hyprland.enable = lib.mkEnableOption "Hyprland desktop";

  config = lib.mkIf config.dotfiles.hyprland.enable {
    home-manager.users.${user} = {
      dotfiles.hyprland.enable = true;
      stylix.targets = {
        gtk.enable = true;
        qt.enable = true;
        hyprland.enable = true;
        yazi.enable = true;
      };
    };

    programs.hyprland = {
      enable = true;
      xwayland.enable = true;
    };

    services.greetd = {
      enable = true;
      settings.default_session = {
        command = "${lib.getExe pkgs.tuigreet} --time --cmd ${lib.getExe' config.programs.hyprland.package "start-hyprland"}";
        user = "greeter";
      };
    };

    environment.sessionVariables.NIXOS_OZONE_WL = "1";

    programs.gpu-screen-recorder.enable = true;

    xdg.portal = {
      enable = true;
      xdgOpenUsePortal = true;
      extraPortals = [pkgs.xdg-desktop-portal-gtk];
    };

    # Only GTK, Qt, Hyprland borders and yazi; the other tools carry their own Catppuccin themes.
    stylix = {
      enable = true;
      autoEnable = false;
      base16Scheme = "${pkgs.base16-schemes}/share/themes/catppuccin-mocha.yaml";
      polarity = "dark";
      fonts = {
        monospace = {
          package = pkgs.nerd-fonts.jetbrains-mono;
          name = "JetBrainsMono Nerd Font Mono";
        };
        sansSerif = {
          package = pkgs.dejavu_fonts;
          name = "DejaVu Sans";
        };
        serif = {
          package = pkgs.dejavu_fonts;
          name = "DejaVu Serif";
        };
      };
      targets = {
        gtk.enable = true;
        qt.enable = true;
        fontconfig.enable = true;
        font-packages.enable = true;
      };
    };

    fonts = {
      enableDefaultPackages = true;
      packages = [pkgs.nerd-fonts.jetbrains-mono];
    };
  };
}
