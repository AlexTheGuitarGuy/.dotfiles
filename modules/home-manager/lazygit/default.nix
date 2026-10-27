{
  lib,
  config,
  ...
}: {
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
            command = "git diff HEAD | diffnav";
            context = "files";
            description = "Open working tree diff in diffnav";
            key = "<c-t>";
            output = "terminal";
          }
          {
            command = "git show --format= {{.SelectedCommit.Hash}} | diffnav";
            context = "commits";
            description = "Open commit diff in diffnav";
            key = "<c-t>";
            output = "terminal";
          }
          {
            command = "b={{.SelectedLocalBranch.Name}}; base=$(gh pr view \"$b\" --json baseRefName -q .baseRefName 2>/dev/null) || base=$(git symbolic-ref --short refs/remotes/origin/HEAD | sed 's@^origin/@@'); git diff \"origin/$base...$b\" | diffnav";
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
            activeBorderColor = ["#89b4fa" "bold"];
            cherryPickedCommitBgColor = ["#45475a"];
            cherryPickedCommitFgColor = ["#89b4fa"];
            defaultFgColor = ["#cdd6f4"];
            inactiveBorderColor = ["#a6adc8"];
            inactiveViewSelectedLineBgColor = ["#6c7086"];
            markedBaseCommitBgColor = ["#f9e2af"];
            markedBaseCommitFgColor = ["#89b4fa"];
            optionsTextColor = ["#89b4fa"];
            searchingActiveBorderColor = ["#f9e2af"];
            selectedLineBgColor = ["#313244"];
            unstagedChangesColor = ["#f38ba8"];
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
