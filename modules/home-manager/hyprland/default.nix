{
  lib,
  config,
  pkgs,
  ...
}: let
  lua = lib.generators.mkLuaInline;
  bind = key: dispatcher: {_args = [(lua ''mod .. " + ${key}"'') (lua dispatcher)];};
  mouseBind = key: dispatcher: {_args = [(lua ''mod .. " + ${key}"'') (lua dispatcher) {mouse = true;}];};

  # Home-row workspaces: a s d f / j k l ; = 1-8, u i o p = 9-12, m , . / = 13-16, q w e r = 17-20.
  workspaceKeys = ["A" "S" "D" "F" "J" "K" "L" "semicolon" "U" "I" "O" "P" "M" "comma" "period" "slash" "Q" "W" "E" "R"];
  workspaceBinds = lib.imap1 (n: key: bind key ''hl.dsp.focus({ workspace = "${toString n}" })'') workspaceKeys;
  moveBinds = lib.imap1 (n: key: bind "SHIFT + ${key}" ''hl.dsp.window.move({ workspace = "${toString n}", follow = true })'') (lib.take 8 workspaceKeys);
  focusBinds = lib.mapAttrsToList (key: dir: bind key ''hl.dsp.focus({ direction = "${dir}" })'') {
    LEFT = "l";
    RIGHT = "r";
    UP = "u";
    DOWN = "d";
  };
in {
  imports = [./caelestia.nix];

  options.dotfiles.hyprland.enable = lib.mkEnableOption "Hyprland user config";

  config = lib.mkIf config.dotfiles.hyprland.enable {
    wayland.windowManager.hyprland = {
      enable = true;
      package = null;
      portalPackage = null;
      settings = {
        mod._var = "SUPER";

        bind =
          [
            (bind "RETURN" ''hl.dsp.exec_cmd("${lib.getExe pkgs.ghostty}")'')
            (bind "C" "hl.dsp.window.close()")
            (bind "SHIFT + E" "hl.dsp.exit()")
            (bind "SPACE" ''hl.dsp.global("caelestia:launcher")'')
            (bind "SHIFT + SPACE" ''hl.dsp.exec_cmd("anyrun")'')
            {_args = ["Print" (lua ''hl.dsp.exec_cmd("hyprshot --mode region --clipboard-only --freeze")'')];}
            {_args = ["ALT + F" (lua ''hl.dsp.window.fullscreen({ action = "toggle" })'')];}
            (mouseBind "mouse:272" "hl.dsp.window.drag()")
            (mouseBind "mouse:273" "hl.dsp.window.resize()")
          ]
          ++ workspaceBinds
          ++ moveBinds
          ++ focusBinds;

        config = {
          general = {
            gaps_in = 5;
            gaps_out = 20;
            border_size = 5;
            resize_on_border = false;
            allow_tearing = false;
            layout = "dwindle";
          };
          decoration = {
            rounding = 10;
            active_opacity = 1.0;
            inactive_opacity = 0.75;
            blur = {
              enabled = true;
              size = 3;
              passes = 1;
              vibrancy = 0.1696;
            };
          };
          animations.enabled = true;
          input = {
            kb_layout = "us,ru,ro";
            kb_options = "grp:alt_shift_toggle";
          };
        };

        curve = {_args = ["myBezier" {type = "bezier"; points = [[0.05 0.9] [0.1 1.05]];}];};
        animation = [
          {leaf = "windows"; enabled = true; speed = 7; bezier = "myBezier";}
          {leaf = "windowsOut"; enabled = true; speed = 7; bezier = "default"; style = "popin 80%";}
          {leaf = "border"; enabled = true; speed = 10; bezier = "default";}
          {leaf = "borderangle"; enabled = true; speed = 8; bezier = "default";}
          {leaf = "fade"; enabled = true; speed = 7; bezier = "default";}
          {leaf = "workspaces"; enabled = true; speed = 6; bezier = "default";}
        ];

        window_rule = [
          {match.class = "^(spotify)$"; workspace = "19 silent";}
          {match.class = "^(org\\.telegram\\.desktop)$"; workspace = "20 silent";}
        ];

        monitor = {output = ""; mode = "preferred"; position = "auto"; scale = 1;};
      };
    };

    # Night light: warm from 20:00, neutral from 07:00.
    services.hyprsunset = {
      enable = true;
      settings.profile = [
        {time = "7:00"; identity = true;}
        {time = "20:00"; temperature = 4000;}
      ];
    };

    programs.anyrun = {
      enable = true;
      config = {
        plugins = map (p: "${pkgs.anyrun}/lib/lib${p}.so") ["applications" "symbols" "rink" "dictionary" "translate" "websearch"];
        x.fraction = 0.5;
        y.fraction = 0.3;
        width.fraction = 0.3;
        hideIcons = false;
        ignoreExclusiveZones = false;
        layer = "overlay";
        hidePluginInfo = false;
        closeOnClick = false;
        showResultsImmediately = false;
        maxEntries = null;
      };
      extraConfigFiles."websearch.ron".text = ''
        Config(
          prefix: "?",
          engines: [DuckDuckGo]
        )
      '';
    };

    dconf.settings."org/gnome/desktop/interface".color-scheme = "prefer-dark";

    programs.yazi.enable = true;
    services.cliphist.enable = true;

    home.packages = with pkgs; [ghostty hyprshot wl-clipboard pavucontrol];
  };
}
