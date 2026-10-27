{
  lib,
  config,
  ...
}: {
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
            skin = "catppuccin-mocha";
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
      skins.catppuccin-mocha = {
        k9s = {
          body = {
            bgColor = "#1e1e2e";
            fgColor = "#cdd6f4";
            logoColor = "#cba6f7";
          };
          dialog = {
            bgColor = "#313244";
            buttonBgColor = "#7f849c";
            buttonFgColor = "#1e1e2e";
            buttonFocusBgColor = "#f5c2e7";
            buttonFocusFgColor = "#1e1e2e";
            fgColor = "#f9e2af";
            fieldFgColor = "#cdd6f4";
            labelFgColor = "#f5e0dc";
          };
          frame = {
            border = {
              fgColor = "#cba6f7";
              focusColor = "#b4befe";
            };
            crumbs = {
              activeColor = "#f2cdcd";
              bgColor = "#eba0ac";
              fgColor = "#1e1e2e";
            };
            menu = {
              fgColor = "#cdd6f4";
              keyColor = "#89b4fa";
              numKeyColor = "#eba0ac";
            };
            status = {
              addColor = "#a6e3a1";
              completedColor = "#6c7086";
              errorColor = "#f38ba8";
              highlightColor = "#89dceb";
              killColor = "#cba6f7";
              modifyColor = "#b4befe";
              newColor = "#89b4fa";
              pendingColor = "#fab387";
            };
            title = {
              bgColor = "#1e1e2e";
              counterColor = "#f9e2af";
              fgColor = "#94e2d5";
              filterColor = "#a6e3a1";
              highlightColor = "#f5c2e7";
            };
          };
          help = {
            bgColor = "#1e1e2e";
            fgColor = "#cdd6f4";
            keyColor = "#89b4fa";
            numKeyColor = "#eba0ac";
            sectionColor = "#a6e3a1";
          };
          info = {
            fgColor = "#fab387";
            sectionColor = "#cdd6f4";
          };
          prompt = {
            bgColor = "#181825";
            fgColor = "#cdd6f4";
            suggestColor = "#89b4fa";
          };
          views = {
            charts = {
              bgColor = "#1e1e2e";
              chartBgColor = "#1e1e2e";
              defaultChartColors = ["#a6e3a1" "#f38ba8"];
              defaultDialColors = ["#a6e3a1" "#f38ba8"];
              dialBgColor = "#1e1e2e";
              resourceColors = {
                cpu = ["#cba6f7" "#89b4fa"];
                mem = ["#f9e2af" "#fab387"];
              };
            };
            logs = {
              bgColor = "#1e1e2e";
              fgColor = "#cdd6f4";
              indicator = {
                bgColor = "#1e1e2e";
                fgColor = "#b4befe";
                toggleOffColor = "#a6adc8";
                toggleOnColor = "#a6e3a1";
              };
            };
            table = {
              bgColor = "#1e1e2e";
              cursorBgColor = "#45475a";
              cursorFgColor = "#313244";
              fgColor = "#cdd6f4";
              header = {
                bgColor = "#1e1e2e";
                fgColor = "#f9e2af";
                sorterColor = "#89dceb";
              };
              markColor = "#f5e0dc";
            };
            xray = {
              bgColor = "#1e1e2e";
              cursorColor = "#45475a";
              cursorTextColor = "#1e1e2e";
              fgColor = "#cdd6f4";
              graphicColor = "#f5c2e7";
            };
            yaml = {
              colonColor = "#a6adc8";
              keyColor = "#89b4fa";
              valueColor = "#cdd6f4";
            };
          };
        };
      };
    };
  };
}
