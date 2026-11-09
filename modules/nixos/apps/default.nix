{
  lib,
  config,
  pkgs,
  user,
  ...
}: let
  officeMimes = [
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
    "application/vnd.openxmlformats-officedocument.presentationml.presentation"
    "application/msword"
    "application/vnd.ms-excel"
    "application/vnd.ms-powerpoint"
    "application/rtf"
  ];

  # Own profile dir so it never clashes with a running office instance.
  officeThumbnail = pkgs.writeShellApplication {
    name = "office-thumbnail";
    runtimeInputs = with pkgs; [libreoffice imagemagick];
    text = ''
      tmp=$(mktemp -d)
      trap 'rm -rf "$tmp"' EXIT
      soffice --headless -env:UserInstallation=file://"$tmp"/profile --convert-to png --outdir "$tmp" "$1"
      magick "$tmp"/*.png -thumbnail "''${2}x''${2}" "$3"
    '';
  };

  officeThumbnailer = pkgs.writeTextDir "share/thumbnailers/libreoffice.thumbnailer" ''
    [Thumbnailer Entry]
    TryExec=${lib.getExe officeThumbnail}
    Exec=${lib.getExe officeThumbnail} %i %s %o
    MimeType=${lib.concatStringsSep ";" officeMimes};
  '';
in {
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

    services.input-remapper = {
      enable = true;
      enableUdevRules = true;
    };

    programs.thunar.enable = true;
    services.tumbler.enable = true;
    programs.gdk-pixbuf.modulePackages = [pkgs.webp-pixbuf-loader];
    environment.systemPackages = [officeThumbnailer];

    programs._1password.enable = true;
    programs._1password-gui = {
      enable = true;
      polkitPolicyOwners = [user];
    };
    # 1Password only talks to browsers it knows; LibreWolf is not on its built-in list.
    environment.etc."1password/custom_allowed_browsers" = {
      text = "librewolf\n";
      mode = "0755";
    };
  };
}
