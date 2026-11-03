{
  lib,
  config,
  pkgs,
  ...
}: let
  p = config.dotfiles.theme.palette;
in {
  options.dotfiles.lazygit.enable = lib.mkEnableOption "lazygit";

  config = lib.mkIf config.dotfiles.lazygit.enable {
    programs.lazygit = {
      enable = true;
      settings = {
        customCommands = [
          {
            command = "f=$(git rev-parse --git-path COMMIT_EDITMSG) && echo \"{{.Form.Msg}}\" > \"$f\" && nvim \"$f\" && [ -s \"$f\" ] && git commit -F \"$f\" || echo \"Commit message is empty, commit aborted.\"";
            context = "files";
            description = "Pick AI commit";
            key = "<c-z>";
            output = "terminal";
            prompts = [
              {
                command = "ai-commit-msgs";
                filter = "^(?P<number>\\d+)\\.\\s(?P<message>.+)$";
                key = "Msg";
                labelFormat = "{{ .number }}: {{ .message | green }}";
                title = "AI Commits";
                type = "menuFromCommand";
                valueFormat = "{{ .message }}";
              }
            ];
          }
          {
            command = "git diff HEAD | ${pkgs.diffnav}/bin/diffnav";
            context = "files";
            description = "Open working tree diff in diffnav";
            key = "<c-t>";
            output = "terminal";
          }
          {
            command = "git show --format= {{.SelectedCommit.Hash}} | ${pkgs.diffnav}/bin/diffnav";
            context = "commits";
            description = "Open commit diff in diffnav";
            key = "<c-t>";
            output = "terminal";
          }
          {
            command = "b={{.SelectedLocalBranch.Name}}; base=$(gh pr view \"$b\" --json baseRefName -q .baseRefName 2>/dev/null) || base=$(git symbolic-ref --short refs/remotes/origin/HEAD | sed 's@^origin/@@'); git diff \"origin/$base...$b\" | ${pkgs.diffnav}/bin/diffnav";
            context = "localBranches";
            description = "Open branch diff vs parent (PR base or default branch) in diffnav";
            key = "<c-t>";
            output = "terminal";
          }
        ];
        git = {
          diffRenderers = [
            {
              colorArg = "always";
              command = "delta --paging=never";
            }
          ];
        };
        gui = {
          theme = {
            activeBorderColor = [p.accent "bold"];
            cherryPickedCommitBgColor = [p.surface1];
            cherryPickedCommitFgColor = [p.accent];
            defaultFgColor = [p.text];
            inactiveBorderColor = [p.subtext0];
            inactiveViewSelectedLineBgColor = [p.overlay0];
            markedBaseCommitBgColor = [p.yellow];
            markedBaseCommitFgColor = [p.accent];
            optionsTextColor = [p.accent];
            searchingActiveBorderColor = [p.yellow];
            selectedLineBgColor = [p.surface0];
            unstagedChangesColor = [p.red];
          };
        };
        keybinding = {universal = {suspendApp = "<disabled>";};};
        os = {
          editPreset = "nvim";
          open = "xdg-open {{filename}}";
          openLink = "xdg-open {{link}}";
        };
      };
    };
  };
}
