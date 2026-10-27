{
  inputs,
  lib,
  config,
  pkgs,
  ...
}: {
  imports = [inputs.sops-nix.homeManagerModules.sops];

  options.dotfiles.secrets.enable = lib.mkEnableOption "sops secrets";

  config = lib.mkIf config.dotfiles.secrets.enable {
    home.packages = [pkgs.sops pkgs.age];

    sops = {
      defaultSopsFile = ../../../secrets/secrets.yaml;
      age.keyFile = "${config.home.homeDirectory}/.config/sops/age/keys.txt";
    };

    programs.nushell.extraConfig = ''
      # Secrets (sops-encrypted, decrypted per command, never exported to the shell)
      def secrets [...keys: string] {
          ["${config.dotfiles.repo}/secrets/secrets.yaml"] ++ (glob ~/.config/secrets.d/*.yaml)
          | path expand
          | where { path exists }
          | each {|f| sops -d $f | from yaml }
          | reduce {|it, acc| $acc | merge $it }
          | select ...$keys
      }
    '';
  };
}
