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
      matchBlocks = {
        "github.com" = {
          hostname = "github.com";
          identityFile = "${keys}/github-personal";
          identitiesOnly = true;
        };
        "gitlab.com" = {
          hostname = "gitlab.com";
          identityFile = "${keys}/gitlab-professional";
          identitiesOnly = true;
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
