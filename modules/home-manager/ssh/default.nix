{
  lib,
  config,
  ...
}: let
  ssh = "${config.home.homeDirectory}/.ssh";
in {
  options.dotfiles.ssh.enable = lib.mkEnableOption "SSH config and keys";

  config = lib.mkIf config.dotfiles.ssh.enable {
    sshAuthSock.initialization = {
      bash = ''export SSH_AUTH_SOCK="$XDG_RUNTIME_DIR/ssh-agent"'';
      fish = ''set -x SSH_AUTH_SOCK "$XDG_RUNTIME_DIR/ssh-agent"'';
      nushell = ''$env.SSH_AUTH_SOCK = $"($env.XDG_RUNTIME_DIR)/ssh-agent"'';
    };

    programs.ssh = {
      enable = true;
      enableDefaultConfig = false;
      settings = {
        "github.com" = {
          HostName = "github.com";
          IdentityFile = "${ssh}/github-personal";
          IdentitiesOnly = true;
        };
        "gitlab.com" = {
          HostName = "gitlab.com";
          IdentityFile = "${ssh}/gitlab-professional";
          IdentitiesOnly = true;
        };
        "*".AddKeysToAgent = "yes";
      };
    };

    home.file = {
      ".ssh/github-personal.pub".source = ./pubkeys/github-personal.pub;
      ".ssh/gitlab-professional.pub".source = ./pubkeys/gitlab-professional.pub;
    };

    sops.secrets = {
      SSH_KEY_GITHUB_PERSONAL.path = "${ssh}/github-personal";
      SSH_KEY_GITLAB_PROFESSIONAL.path = "${ssh}/gitlab-professional";
    };
  };
}
