{
  lib,
  config,
  pkgs,
  ...
}: let
  lua = lib.generators.mkLuaInline;
  bind = key: dispatcher: {_args = [(lua ''mod .. " + ${key}"'') (lua dispatcher)];};
  mouseBind = key: dispatcher: {_args = [(lua ''mod .. " + ${key}"'') (lua dispatcher) {mouse = true;}];};
  mediaKey = key: cmd: repeating: {_args = [key (lua ''hl.dsp.exec_cmd("${cmd}")'') {inherit repeating; locked = true;}];};

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

  screenshot = "hyprshot --mode region --raw --freeze | satty --filename - --copy-command wl-copy --early-exit --actions-on-enter save-to-clipboard --output-filename ~/Pictures/Screenshots/satty-%Y%m%d-%H%M%S.png";
  osk = pkgs.writeShellApplication {
    name = "osk";
    runtimeInputs = with pkgs; [wvkbd procps jq];
    text = ''
      if pkill -x wvkbd-mobintl; then
        read -r kb idx < <(hyprctl devices -j | jq -r '.keyboards[] | select(.main) | "\(.name) \(.active_layout_index)"')
        hyprctl switchxkblayout "$kb" "$idx" >/dev/null
      else
        exec wvkbd-mobintl -L 300
      fi
    '';
  };
  focusOnActivateLater = "sleep 60 && hyprctl eval 'hl.config({ misc = { focus_on_activate = true } })'";

  gapsOut = 8;
  borderSize = 5;

  cfg = config.dotfiles.hyprland;
in {
  imports = [./caelestia.nix];

  options.dotfiles.hyprland = {
    enable = lib.mkEnableOption "Hyprland user config";
    autostart = lib.mkOption {
      type = lib.types.listOf lib.types.str;
      default = [];
      example = ["Telegram" "spotify"];
      description = "Commands started once when Hyprland starts.";
    };
    pinnedApps = lib.mkOption {
      type = lib.types.listOf (lib.types.attrsOf lib.types.str);
      default = [];
      example = [
        {
          class = "org.telegram.desktop";
          initial_title = "Telegram";
          workspace = "20";
        }
      ];
      description = "Windows that open on their workspace and are moved back when moved away. class and the optional initial_title match exactly.";
    };
  };

  config = lib.mkIf config.dotfiles.hyprland.enable {
    wayland.windowManager.hyprland = {
      enable = true;
      package = null;
      portalPackage = null;
      extraConfig =
        lib.optionalString (config.dotfiles.hyprland.autostart != []) ''
          hl.on("hyprland.start", function()
          ${lib.concatMapStrings (cmd: "  hl.exec_cmd(${builtins.toJSON cmd})\n") (config.dotfiles.hyprland.autostart ++ [focusOnActivateLater])}end)
        ''
        + ''
          local pinned = ${lib.generators.toLua {} cfg.pinnedApps}
          hl.on("window.move_to_workspace", function(w, ws)
            for _, app in ipairs(pinned) do
              if w.class == app.class and (app.initial_title == nil or w.initial_title == app.initial_title) and ws.name ~= app.workspace then
                hl.dispatch(hl.dsp.window.move({ window = w, workspace = app.workspace, follow = false }))
              end
            end
          end)
        '';
      settings = {
        mod._var = "SUPER";

        bind =
          [
            (bind "RETURN" ''hl.dsp.exec_cmd("${lib.getExe pkgs.ghostty}")'')
            (bind "C" "hl.dsp.window.close()")
            (bind "SHIFT + E" "hl.dsp.exit()")
            (bind "SPACE" ''hl.dsp.global("caelestia:launcher")'')
            {_args = ["Print" (lua ''hl.dsp.exec_cmd("${screenshot}")'')];}
            (bind "SHIFT + Z" ''hl.dsp.exec_cmd("${screenshot}")'')
            (bind "SHIFT + O" ''hl.dsp.exec_cmd("${lib.getExe osk}")'')
            (bind "SHIFT + R" ''hl.dsp.exec_cmd("gpu-screen-recorder-gtk")'')
            (bind "V" ''hl.dsp.exec_cmd("caelestia clipboard")'')
            (bind "X" ''hl.dsp.exec_cmd("caelestia emoji -p")'')
            (bind "ESCAPE" ''hl.dsp.exec_cmd("caelestia shell lock lock")'')
            (mediaKey "XF86MonBrightnessUp" "caelestia shell brightness set +5%" true)
            (mediaKey "XF86MonBrightnessDown" "caelestia shell brightness set 5%-" true)
            (mediaKey "XF86AudioRaiseVolume" "${pkgs.wireplumber}/bin/wpctl set-volume -l 1 @DEFAULT_AUDIO_SINK@ 5%+" true)
            (mediaKey "XF86AudioLowerVolume" "${pkgs.wireplumber}/bin/wpctl set-volume @DEFAULT_AUDIO_SINK@ 5%-" true)
            (mediaKey "XF86AudioMute" "${pkgs.wireplumber}/bin/wpctl set-mute @DEFAULT_AUDIO_SINK@ toggle" false)
            (mediaKey "XF86AudioMicMute" "${pkgs.wireplumber}/bin/wpctl set-mute @DEFAULT_AUDIO_SOURCE@ toggle" false)
            (mediaKey "XF86AudioPlay" "caelestia shell mpris playPause" false)
            (mediaKey "XF86AudioPause" "caelestia shell mpris playPause" false)
            (mediaKey "XF86AudioNext" "caelestia shell mpris next" false)
            (mediaKey "XF86AudioPrev" "caelestia shell mpris previous" false)
            {_args = ["ALT + F" (lua ''hl.dsp.window.fullscreen({ action = "toggle" })'')];}
            (mouseBind "mouse:272" "hl.dsp.window.drag()")
            (mouseBind "mouse:273" "hl.dsp.window.resize()")
          ]
          ++ workspaceBinds
          ++ moveBinds
          ++ focusBinds;

        config = {
          general = {
            gaps_in = 3;
            gaps_out = gapsOut;
            border_size = borderSize;
            resize_on_border = false;
            allow_tearing = false;
            layout = "dwindle";
          };
          decoration = {
            rounding = config.programs.caelestia.settings.border.rounding - gapsOut - borderSize;
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
          misc.focus_on_activate = config.dotfiles.hyprland.autostart == [];
          input = {
            kb_layout = "us,ru,ro";
            kb_options = "grp:alt_shift_toggle";
            touchpad.natural_scroll = true;
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

        window_rule =
          map (app: {
            match = lib.mapAttrs (_: v: "^(${lib.escapeRegex v})$") (removeAttrs app ["workspace"]);
            workspace = "${app.workspace} silent";
          })
          cfg.pinnedApps
          ++ [
            {
              match.class = "^(org\\.telegram\\.desktop)$";
              suppress_event = "activatefocus";
            }
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

    dconf.settings."org/gnome/desktop/interface".color-scheme = "prefer-dark";

    programs.yazi = {
      enable = true;
      plugins.smart-enter = pkgs.yaziPlugins.smart-enter;
      keymap.mgr.prepend_keymap = [
        {
          on = "<Enter>";
          run = "plugin smart-enter";
          desc = "Enter the folder or open the file";
        }
      ];
    };
    services.cliphist.enable = true;

    home.pointerCursor = {
      package = pkgs.bibata-cursors;
      name = "Bibata-Modern-Ice";
      size = 24;
      gtk.enable = true;
      hyprcursor.enable = true;
    };

    gtk.font.size = lib.mkForce 10;

    xdg.desktopEntries =
      {
        osk = {
          name = "On-screen keyboard";
          exec = lib.getExe osk;
          icon = "input-keyboard";
          categories = ["Utility"];
        };
      }
      // lib.mapAttrs (exec: name: {
        inherit name exec;
        noDisplay = true;
      }) {
        kvantummanager = "Kvantum Manager";
        qt5ct = "Qt5 Settings";
        qt6ct = "Qt6 Settings";
      };
    dconf.settings."org/gnome/desktop/interface".text-scaling-factor = 1.25;

    programs.fuzzel = {
      enable = true;
      settings.main = {
        width = 60;
        lines = 15;
      };
    };

    home.packages = with pkgs; [osk ghostty hyprshot satty wl-clipboard pavucontrol brightnessctl brotab gpu-screen-recorder-gtk];
  };
}
