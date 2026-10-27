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
    home-manager.users.alex = {
      dotfiles.nvim.configSource = ../modules/home-manager/nvim/config;
      sops.defaultSopsFile = pkgs.lib.mkForce ./test-secrets.yaml;
    };
    systemd.tmpfiles.rules = [
      "d /home/alex/.config 0755 alex users"
      "d /home/alex/.config/sops 0700 alex users"
      "d /home/alex/.config/sops/age 0700 alex users"
      "C /home/alex/.config/sops/age/keys.txt 0600 alex users - ${./test-age-key.txt}"
    ];
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
    machine.succeed("su -s /bin/sh alex -c 'mkdir -p ~/hooktest && cd ~/hooktest && git init -q && git commit -q --allow-empty -m init && git worktree add -q ../hooktest-wt'")
    todos = machine.succeed("cat /home/alex/hooktest-wt/.excluded/todos.md")
    assert todos == "# TODOs\n", repr(todos)
    assert "dry run" in nu("cd ~/hooktest; ^git wt-clean")
    assert nu("$env.config.edit_mode") == "vi"
    assert nu("which gst | get type.0") == "alias"
    assert nu("which zvim | get type.0") == "custom"
    assert nu("which cd | get definition.0") == "__zoxide_z"
    assert nu("$env.STARSHIP_SHELL") == "nu"

    assert nu("$env.EDITOR") == "nvim"
    machine.succeed("test -f /home/alex/.config/nvim/lua/plugins/lspconfig.lua")
    machine.fail("test -e /home/alex/.config/nvim/lua/plugins/mason.lua")
    for tool in ["lua-language-server", "typescript-language-server", "vscode-eslint-language-server", "sqls", "ngserver", "prettierd", "stylua", "js-debug", "rg", "fd", "tree-sitter", "cc", "gio"]:
        path = nu(f"^nvim --headless --clean -c 'lua io.write(vim.fn.exepath(\"{tool}\"))' -c q")
        assert path.startswith("/nix/store/"), f"{tool} not on nvim PATH: {path!r}"

    for cmd in ["btop --version", "fastfetch --version", "gh --version", "gh dash --help", "ghostty --version", "k9s version --short", "lazygit --version", "posting --help"]:
        nu(f"^{cmd} | complete | if $in.exit_code != 0 {{ error make {{msg: $in.stderr}} }}")
    font = machine.succeed("su -s /bin/sh alex -c 'fc-match \"JetBrainsMono Nerd Font\" file'")
    assert "JetBrainsMonoNerdFont" in font, font
    assert nu("^readlink -m ~/.config/btop") == "/home/alex/.dotfiles/modules/home-manager/btop/config"
    assert nu("^readlink -m ~/.config/gh/config.yml") == "/home/alex/.dotfiles/modules/home-manager/gh/config.yml"

    nu("mkdir ~/.config/secrets.d")
    nu("{TEST_KEY: hello, OTHER: x} | to yaml | ^sops -e --age (^age-keygen -y ~/.config/sops/age/keys.txt) --input-type yaml --output-type yaml /dev/stdin | save ~/.config/secrets.d/test.yaml")
    machine.fail("grep -q hello /home/alex/.config/secrets.d/test.yaml")
    assert nu("secrets TEST_KEY | to nuon") == "{TEST_KEY: hello}"

    assert nu("open ~/.ssh/keys/github-personal") == "test-github-personal"
    assert nu("open ~/.ssh/keys/gitlab-professional") == "test-gitlab-professional"
    assert nu("^stat -L -c %a ~/.ssh/keys/github-personal") == "400"
    assert "identityfile /home/alex/.ssh/keys/github-personal" in nu("^ssh -G github.com")
    assert "identitiesonly yes" in nu("^ssh -G gitlab.com")
    machine.succeed("test -f /home/alex/.ssh/pubkeys/github-personal.pub")
  '';
}
