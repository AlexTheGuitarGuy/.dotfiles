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
      default = config.lib.file.mkOutOfStoreSymlink "${config.home.homeDirectory}/.dotfiles/modules/home-manager/nvim/config";
      description = "Neovim config directory, linked to ~/.config/nvim. Defaults to the live repo checkout so edits and lazy-lock.json updates land in git.";
    };
  };

  config = lib.mkIf cfg.enable {
    programs.neovim = {
      enable = true;
      defaultEditor = true;
      sideloadInitLua = true;
      extraPackages = with pkgs; [
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
        nodejs

        fd
        gcc
        git
        glib
        ripgrep
        tree-sitter
      ];
    };

    xdg.configFile.nvim.source = cfg.configSource;
  };
}
