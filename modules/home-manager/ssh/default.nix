{
  lib,
  config,
  ...
}: let
  keys = "${config.home.homeDirectory}/.ssh/keys";
in {
  options.dotfiles.ssh.enable = lib.mkEnableOption "SSH config and keys";

  config = lib.mkIf config.dotfiles.ssh.enable {
    programs.ssh = {
      enable = true;
      enableDefaultConfig = false;
      settings = {
        "github.com" = {
          HostName = "github.com";
          IdentityFile = "${keys}/github-personal";
          IdentitiesOnly = true;
        };
        "gitlab.com" = {
          HostName = "gitlab.com";
          IdentityFile = "${keys}/gitlab-professional";
          IdentitiesOnly = true;
        };
      };
    };

    home.file.".ssh/pubkeys" = {
      source = ./pubkeys;
      recursive = true;
    };

    sops.secrets = {
      SSH_KEY_GITHUB_PERSONAL.path = "${keys}/github-personal";
      SSH_KEY_GITLAB_PROFESSIONAL.path = "${keys}/gitlab-professional";
    };
  };
}
