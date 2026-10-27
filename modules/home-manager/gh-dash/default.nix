{
  lib,
  config,
  ...
}: {
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
              command = "cd {{.RepoPath}}; lazygit\n";
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
        pager = {diff = "diffnav";};
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
            background = {selected = "#313244";};
            border = {
              faint = "#181825";
              primary = "#89b4fa";
              secondary = "#585b70";
            };
            text = {
              actor = "#b4befe";
              faint = "#6c7086";
              inverted = "#1e1e2e";
              primary = "#cdd6f4";
              secondary = "#a6adc8";
              success = "#a6e3a1";
              warning = "#fab387";
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
