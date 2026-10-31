{
  lib,
  config,
  user,
  ...
}: {
  options.dotfiles.docker.enable = lib.mkEnableOption "Docker engine with compose";

  config = lib.mkIf config.dotfiles.docker.enable {
    virtualisation.docker = {
      enable = true;
      daemon.settings.default-address-pools = [
        {
          base = "198.18.0.0/15";
          size = 24;
        }
      ];
    };
    users.users.${user}.extraGroups = ["docker"];
  };
}
