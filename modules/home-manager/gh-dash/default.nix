{
  lib,
  config,
  pkgs,
  ...
}: let
  p = config.dotfiles.theme.palette;
in {
  options.dotfiles.gh-dash.enable = lib.mkEnableOption "gh-dash";

  config = lib.mkIf config.dotfiles.gh-dash.enable {
    programs.gh-dash = {
      enable = true;
      settings = {
        confirmQuit = false;
        defaults = {
          issuesLimit = 20;
          layout = {
            issues = {
              assignees = {
                hidden = true;
                width = 20;
              };
              createdAt = {width = 5;};
              creator = {width = 10;};
              creatorIcon = {hidden = false;};
              repo = {width = 15;};
              updatedAt = {width = 5;};
            };
            prs = {
              assignees = {
                hidden = true;
                width = 20;
              };
              author = {width = 15;};
              authorIcon = {hidden = false;};
              base = {
                hidden = true;
                width = 15;
              };
              createdAt = {width = 5;};
              lines = {width = 15;};
              repo = {width = 20;};
              updatedAt = {width = 5;};
            };
          };
          notificationsLimit = 20;
          prApproveComment = "LGTM";
          preview = {
            open = true;
            width = 0.45;
          };
          prsLimit = 20;
          refetchIntervalMinutes = 30;
          view = "prs";
        };
        includeReadNotifications = true;
        issuesSections = [
          {
            filters = "is:open author:@me";
            title = "My Issues";
          }
          {
            filters = "is:open assignee:@me";
            title = "Assigned";
          }
          {
            filters = "is:open involves:@me -author:@me";
            title = "Involved";
          }
        ];
        keybindings = {
          prs = [
            {
              command = "if [ -n \"{{.RepoPath}}\" ]; then cd \"{{.RepoPath}}\"; fi; lazygit";
              key = "g";
              name = "lazygit";
            }
            {
              command = "gh enhance -R {{.RepoName}} {{.PrNumber}}";
              key = "T";
              name = "enhance";
            }
          ];
        };
        notificationsSections = [
          {
            filters = "";
            title = "All";
          }
          {
            filters = "reason:author";
            title = "Created";
          }
          {
            filters = "reason:participating";
            title = "Participating";
          }
          {
            filters = "reason:mention";
            title = "Mentioned";
          }
          {
            filters = "reason:review-requested";
            title = "Review Requested";
          }
          {
            filters = "reason:assign";
            title = "Assigned";
          }
          {
            filters = "reason:subscribed";
            title = "Subscribed";
          }
          {
            filters = "reason:team-mention";
            title = "Team Mentioned";
          }
        ];
        pager = {diff = "${pkgs.diffnav}/bin/diffnav";};
        prSections = [
          {
            filters = "is:open author:@me";
            title = "My Pull Requests";
          }
          {
            filters = "is:open review-requested:@me";
            title = "Needs My Review";
          }
          {
            filters = "is:open involves:@me -author:@me";
            title = "Involved";
          }
        ];
        repo = {
          branchesRefetchIntervalSeconds = 30;
          prsRefetchIntervalSeconds = 60;
        };
        showAuthorIcons = true;
        smartFilteringAtLaunch = true;
        theme = {
          colors = {
            background = {selected = p.surface0;};
            border = {
              faint = p.mantle;
              primary = p.accent;
              secondary = p.surface2;
            };
            text = {
              actor = p.lavender;
              faint = p.overlay0;
              inverted = p.base;
              primary = p.text;
              secondary = p.subtext0;
              success = p.green;
              warning = p.peach;
            };
          };
          ui = {
            sectionsShowCount = true;
            table = {
              compact = false;
              showSeparator = true;
            };
          };
        };
      };
    };
  };
}
