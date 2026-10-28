{
  pkgs,
  hostModules,
}:
pkgs.testers.runNixOSTest {
  name = "apps";

  node.pkgsReadOnly = false;

  nodes.machine = {
    imports = hostModules;
    dotfiles.hyprland.enable = pkgs.lib.mkForce false;
    virtualisation.memorySize = 2048;
  };

  testScript = ''
    machine.wait_for_unit("home-manager-alex.service")

    for app in ["Telegram", "slack", "teams-for-linux", "viber", "spotify", "obs", "kdenlive", "gimp", "mpv", "snapshot", "reaper", "qpwgraph", "aseprite", "godot", "renpy", "anki", "prismlauncher", "wine", "libreoffice", "qbittorrent", "qalculate-gtk", "protonvpn-app", "gnome-network-displays", "wvkbd-mobintl", "code", "intellij-idea", "datagrip", "mongodb-compass", "apps2samsung", "tizenbrew-device-manager", "tizenbrew-installer", "moldsign", "yazi", "virt-manager", "1password", "input-remapper-gtk"]:
        machine.succeed(f"su -s /bin/sh alex -c 'command -v {app}'")
    machine.succeed("systemctl start input-remapper.service")
    machine.wait_for_unit("input-remapper.service")
    machine.succeed("systemctl is-enabled libvirtd.socket pcscd.socket && test -d /etc/pipewire/jack.conf.d")
    assert "MoldSign-Desk.sh" in machine.fail("su -s /bin/sh alex -c 'moldsign 2>&1'")
  '';
}
