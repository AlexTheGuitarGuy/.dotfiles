{
  lib,
  config,
  pkgs,
  user,
  ...
}: let
  theme = config.dotfiles.theme;
  flavor =
    if theme.polarity == "light"
    then "latte"
    else "mocha";
  ansiAccent =
    {
      rosewater = "lightred";
      flamingo = "lightred";
      pink = "lightmagenta";
      mauve = "magenta";
      red = "red";
      maroon = "lightred";
      peach = "yellow";
      yellow = "lightyellow";
      green = "green";
      teal = "cyan";
      sky = "lightcyan";
      sapphire = "cyan";
      blue = "blue";
      lavender = "lightblue";
    }.${
      theme.accent
    };
  switchDark = "/run/dotfiles-theme-base/bin/switch-to-configuration";
  switchLight = "/run/current-system/specialisation/light/bin/switch-to-configuration";
  themeToggle = pkgs.writeShellApplication {
    name = "theme-toggle";
    runtimeInputs = [pkgs.coreutils pkgs.util-linux];
    text = ''
      state="''${XDG_STATE_HOME:-$HOME/.local/state}/dotfiles/polarity"
      if [ "''${1:-}" = restore ]; then
        want="$(cat "$state" 2>/dev/null || echo dark)"
      elif [ -e ${switchLight} ]; then
        want=light
      else
        want=dark
      fi
      mkdir -p "$(dirname "$state")"
      echo "$want" > "$state"

      lock=/run/nixos/switch-to-configuration.lock
      if [ -e "$lock" ]; then flock "$lock" true; fi
      if [ "$want" = light ] && [ -e ${switchLight} ]; then
        exec /run/wrappers/bin/sudo -n ${switchLight} test
      elif [ "$want" = dark ] && [ "''${1:-}" != restore ] && [ ! -e ${switchLight} ]; then
        exec /run/wrappers/bin/sudo -n ${switchDark} test
      fi
    '';
  };
in {
  options.dotfiles.hyprland.enable = lib.mkEnableOption "Hyprland desktop";
  options.dotfiles.hyprland.startCommand = lib.mkOption {
    type = lib.types.str;
    default = lib.getExe' config.programs.hyprland.package "start-hyprland";
    description = "Command the login screen runs to start the Hyprland session.";
  };

  config = lib.mkIf config.dotfiles.hyprland.enable {
    home-manager.users.${user} = {lib, ...}: {
      dotfiles.hyprland.enable = true;
      systemd.user.services.theme-restore = {
        Unit = {
          Description = "Reapply the light theme when it was last chosen";
          After = ["graphical-session.target"];
        };
        Service = {
          Type = "simple";
          ExecStart = "${lib.getExe themeToggle} restore";
        };
        Install.WantedBy = ["graphical-session.target"];
      };
      home.activation.themeRestore = lib.hm.dag.entryAfter ["reloadSystemd"] ''
        run ${pkgs.systemd}/bin/systemctl --user start --no-block theme-restore.service || true
      '';
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
        command = "${lib.getExe pkgs.tuigreet} --time --theme 'border=${ansiAccent};text=white;prompt=${ansiAccent};time=cyan;action=${ansiAccent};button=${ansiAccent};input=white' --cmd ${config.dotfiles.hyprland.startCommand}";
        user = "greeter";
      };
    };

    environment.sessionVariables.NIXOS_OZONE_WL = "1";

    programs.gpu-screen-recorder.enable = true;

    services.gnome.gnome-keyring.enable = true;
    services.gnome.gcr-ssh-agent.enable = true;
    programs.ssh.startAgent = false;
    security.pam.services.greetd.enableGnomeKeyring = true;
    services.fwupd.enable = true;

    xdg.portal = {
      enable = true;
      extraPortals = [pkgs.xdg-desktop-portal-gtk];
    };
    environment.systemPackages = [
      (lib.hiPrio (pkgs.writeShellScriptBin "xdg-open" ''exec ${lib.getExe' pkgs.glib "gio"} open "$@"''))
      themeToggle
    ];

    specialisation.light.configuration.dotfiles.theme.polarity = lib.mkForce "light";
    system.activationScripts.dotfilesThemeBase = lib.mkIf (theme.polarity == "dark") ''
      ln -sfn "$systemConfig" /run/dotfiles-theme-base
    '';

    security.sudo.extraRules = [
      {
        users = [user];
        commands = map (command: {
          command = "${command} test";
          options = ["NOPASSWD"];
        }) [switchDark switchLight];
      }
    ];

    stylix = {
      enable = true;
      autoEnable = true;
      base16Scheme = "${pkgs.base16-schemes}/share/themes/catppuccin-${flavor}.yaml";
      inherit (theme) polarity;
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
          inherit flavor;
          inherit (theme) accent;
        };
        dark = "Papirus-Dark";
        light = "Papirus-Light";
      };
    };

    fonts = {
      enableDefaultPackages = true;
      packages = [
        pkgs.noto-fonts
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
      # Helvetica falls back to TeX Gyre Heros, whose i blurs into an l at small sizes.
      fontconfig.localConf = ''
        <?xml version="1.0"?>
        <!DOCTYPE fontconfig SYSTEM "urn:fontconfig:fonts.dtd">
        <fontconfig>
          <alias binding="same">
            <family>Helvetica</family>
            <prefer><family>Liberation Sans</family></prefer>
          </alias>
        </fontconfig>
      '';
    };
  };
}
