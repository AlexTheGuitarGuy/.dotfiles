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
    };
  };
}
