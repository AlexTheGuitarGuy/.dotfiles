{
  lib,
  config,
  user,
  ...
}: {
  options.dotfiles.docker.enable = lib.mkEnableOption "Docker engine with compose";

  config = lib.mkIf config.dotfiles.docker.enable {
    virtualisation.docker.enable = true;
    users.users.${user}.extraGroups = ["docker"];
  };
}
