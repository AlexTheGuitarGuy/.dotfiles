{
  pkgs,
  hostModules,
}:
pkgs.testers.runNixOSTest {
  name = "desktop";

  node.pkgsReadOnly = false;

  nodes.machine = {
    imports = hostModules;
    users.users.alex.initialPassword = "test";
    virtualisation = {
      memorySize = 2048;
      qemu.options = ["-vga none" "-device virtio-vga"];
    };
  };

  testScript = ''
    machine.wait_for_unit("greetd.service")
    machine.wait_until_tty_matches("1", "Username")
    machine.send_chars("alex\n")
    machine.wait_until_tty_matches("1", "Password")
    machine.send_chars("test\n")

    hyprctl = "su -s /bin/sh alex -c 'XDG_RUNTIME_DIR=/run/user/1000 HYPRLAND_INSTANCE_SIGNATURE=$(ls /run/user/1000/hypr) hyprctl %s'"
    machine.wait_until_succeeds("ls /run/user/1000/hypr/*/.socket.sock", timeout=60)
    machine.wait_until_succeeds(hyprctl % "version", timeout=10)
    machine.sleep(3)
    machine.screenshot("hyprland")

    machine.fail(hyprctl % "clients | grep -iq class:.*ghostty")
    machine.send_key("meta_l-ret")
    machine.wait_until_succeeds(hyprctl % "clients | grep -iq class:.*ghostty", timeout=60)
    machine.sleep(3)
    machine.screenshot("ghostty")

    import shlex
    def nu(code):
        return machine.succeed("su -s /bin/sh alex -c " + shlex.quote("nu -l -c " + shlex.quote(code))).strip()

    assert nu("git config user.name") == "AlexSeremet"
    assert "delta" in nu("git config pager.diff")
    assert ".excluded/" in nu("open ~/.config/git/ignore")
    assert nu("$env.config.edit_mode") == "vi"
    assert nu("which gst | get type.0") == "alias"
    assert nu("which zvim | get type.0") == "custom"
    assert nu("which cd | get definition.0") == "__zoxide_z"
    assert nu("$env.STARSHIP_SHELL") == "nu"
  '';
}
