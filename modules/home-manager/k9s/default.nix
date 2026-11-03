{
  lib,
  config,
  ...
}: let
  p = config.dotfiles.theme.palette;
in {
  options.dotfiles.k9s.enable = lib.mkEnableOption "k9s";

  config = lib.mkIf config.dotfiles.k9s.enable {
    programs.k9s = {
      enable = true;
      settings = {
        k9s = {
          apiServerTimeout = "15s";
          defaultView = "";
          disablePodCounting = false;
          imageScans = {
            enable = false;
            exclusions = {
              labels = {};
              namespaces = [];
            };
          };
          liveViewAutoRefresh = false;
          logger = {
            buffer = 5000;
            disableAutoscroll = false;
            showTime = false;
            sinceSeconds = -1;
            tail = 100;
            textWrap = false;
          };
          maxConnRetry = 5;
          noExitOnCtrlC = false;
          portForwardAddress = "localhost";
          readOnly = false;
          refreshRate = 2;
          screenDumpDir = "/home/alex/.local/state/k9s/screen-dumps";
          shellPod = {
            image = "busybox:1.35.0";
            limits = {
              cpu = "100m";
              memory = "100Mi";
            };
            namespace = "default";
          };
          skipLatestRevCheck = false;
          thresholds = {
            cpu = {
              critical = 90;
              warn = 70;
            };
            memory = {
              critical = 90;
              warn = 70;
            };
          };
          ui = {
            crumbsless = false;
            defaultsToFullScreen = false;
            enableMouse = false;
            headless = false;
            logoless = false;
            noIcons = false;
            reactive = false;
            skin = "catppuccin";
            splashless = false;
            useFullGVRTitle = false;
          };
        };
      };
      aliases = {
        cr = "clusterroles";
        crb = "clusterrolebindings";
        dp = "deployments";
        jo = "jobs";
        np = "networkpolicies";
        rb = "rolebindings";
        ro = "roles";
        sec = "v1/secrets";
      };
      skins.catppuccin = {
        k9s = {
          body = {
            bgColor = p.base;
            fgColor = p.text;
            logoColor = p.accent;
          };
          dialog = {
            bgColor = p.surface0;
            buttonBgColor = p.overlay1;
            buttonFgColor = p.base;
            buttonFocusBgColor = p.pink;
            buttonFocusFgColor = p.base;
            fgColor = p.yellow;
            fieldFgColor = p.text;
            labelFgColor = p.rosewater;
          };
          frame = {
            border = {
              fgColor = p.accent;
              focusColor = p.lavender;
            };
            crumbs = {
              activeColor = p.flamingo;
              bgColor = p.maroon;
              fgColor = p.base;
            };
            menu = {
              fgColor = p.text;
              keyColor = p.accent;
              numKeyColor = p.maroon;
            };
            status = {
              addColor = p.green;
              completedColor = p.overlay0;
              errorColor = p.red;
              highlightColor = p.sky;
              killColor = p.accent;
              modifyColor = p.lavender;
              newColor = p.accent;
              pendingColor = p.peach;
            };
            title = {
              bgColor = p.base;
              counterColor = p.yellow;
              fgColor = p.teal;
              filterColor = p.green;
              highlightColor = p.pink;
            };
          };
          help = {
            bgColor = p.base;
            fgColor = p.text;
            keyColor = p.accent;
            numKeyColor = p.maroon;
            sectionColor = p.green;
          };
          info = {
            fgColor = p.peach;
            sectionColor = p.text;
          };
          prompt = {
            bgColor = p.mantle;
            fgColor = p.text;
            suggestColor = p.accent;
          };
          views = {
            charts = {
              bgColor = p.base;
              chartBgColor = p.base;
              defaultChartColors = [p.green p.red];
              defaultDialColors = [p.green p.red];
              dialBgColor = p.base;
              resourceColors = {
                cpu = [p.accent p.accent];
                mem = [p.yellow p.peach];
              };
            };
            logs = {
              bgColor = p.base;
              fgColor = p.text;
              indicator = {
                bgColor = p.base;
                fgColor = p.lavender;
                toggleOffColor = p.subtext0;
                toggleOnColor = p.green;
              };
            };
            table = {
              bgColor = p.base;
              cursorBgColor = p.surface1;
              cursorFgColor = p.surface0;
              fgColor = p.text;
              header = {
                bgColor = p.base;
                fgColor = p.yellow;
                sorterColor = p.sky;
              };
              markColor = p.rosewater;
            };
            xray = {
              bgColor = p.base;
              cursorColor = p.surface1;
              cursorTextColor = p.base;
              fgColor = p.text;
              graphicColor = p.pink;
            };
            yaml = {
              colonColor = p.subtext0;
              keyColor = p.accent;
              valueColor = p.text;
            };
          };
        };
      };
    };
  };
}
