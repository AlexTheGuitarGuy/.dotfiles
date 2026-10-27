{
  lib,
  config,
  ...
}: {
  options.dotfiles.gh.enable = lib.mkEnableOption "GitHub CLI";

  config = lib.mkIf config.dotfiles.gh.enable {
    programs.gh = {
      enable = true;
      gitCredentialHelper.enable = false;
    };

    xdg.configFile."gh/config.yml".source = lib.mkForce (config.lib.file.mkOutOfStoreSymlink "${config.dotfiles.repo}/modules/home-manager/gh/config.yml");
  };
}
