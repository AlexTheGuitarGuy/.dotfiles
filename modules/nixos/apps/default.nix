{
  lib,
  config,
  user,
  ...
}: {
  options.dotfiles.apps.enable = lib.mkEnableOption "desktop apps that need system support";

  config = lib.mkIf config.dotfiles.apps.enable {
    home-manager.users.${user}.dotfiles.apps.enable = true;

    # JACK clients (REAPER) run on PipeWire's JACK implementation, no separate jackd.
    services.pipewire = {
      enable = true;
      alsa.enable = true;
      pulse.enable = true;
      jack.enable = true;
    };

    virtualisation.libvirtd.enable = true;
    programs.virt-manager.enable = true;
    users.users.${user}.extraGroups = ["libvirtd"];

    services.input-remapper.enable = true;

    programs._1password.enable = true;
    programs._1password-gui = {
      enable = true;
      polkitPolicyOwners = [user];
    };

    # Smart card daemon for MoldSign's eToken.
    services.pcscd.enable = true;
  };
}
