{
  lib,
  config,
  pkgs,
  ...
}: let
  cfg = config.dotfiles.nvim;
in {
  options.dotfiles.nvim = {
    enable = lib.mkEnableOption "Neovim";
    configSource = lib.mkOption {
      type = lib.types.path;
      default = config.lib.file.mkOutOfStoreSymlink "${config.dotfiles.repo}/modules/home-manager/nvim/config";
      description = "Neovim config directory, linked to ~/.config/nvim. Defaults to the live repo checkout so edits and lazy-lock.json updates land in git.";
    };
  };

  config = lib.mkIf cfg.enable {
    systemd.user.sessionVariables.EDITOR = "nvim";

    programs.neovim = {
      enable = true;
      defaultEditor = true;
      package = pkgs.unstable.neovim-unwrapped;
      sideloadInitLua = true;
      extraPackages =
        (with pkgs.unstable; [
          angular-language-server
          bash-language-server
          buf
          dockerfile-language-server
          gopls
          graphql-language-service-cli
          lua-language-server
          rust-analyzer
          sqls
          tailwindcss-language-server
          typescript
          typescript-language-server
          vscode-langservers-extracted
          yaml-language-server

          prettierd
          stylua

          vscode-js-debug
          tree-sitter
        ])
        ++ (with pkgs; [
          nodejs
          fd
          gcc
          git
          glib
          ripgrep
        ]);
    };

    xdg.configFile.nvim.source = cfg.configSource;
  };
}
