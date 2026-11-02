{user, ...}: {
  home.username = user;
  home.homeDirectory = "/home/${user}";
  home.stateVersion = "26.05";

  dotfiles = {
    btop.enable = true;
    bruno.enable = true;
    claude.enable = true;
    cli.enable = true;
    fastfetch.enable = true;
    gh.enable = true;
    gh-dash.enable = true;
    herdr.enable = true;
    journal.enable = true;
    ghostty.enable = true;
    git.enable = true;
    gpg.enable = true;
    k9s.enable = true;
    lazygit.enable = true;
    librewolf.enable = true;
    nushell.enable = true;
    nvim.enable = true;
    opencode.enable = true;
    playwright.enable = true;
    posting.enable = true;
    secrets.enable = true;
    shots.enable = true;
    spotify.enable = true;
    ssh.enable = true;
    starship.enable = true;

    hyprland = {
      autostart = [
        "ghostty --class=com.mitchellh.ghostty.home -e nu -e herdr"
        "sleep 5 && librewolf"
        "bruno"
        "podman-desktop"
        "teams-for-linux"
        "slack"
        "Telegram"
        "spotify"
        "1password --silent"
      ];
      pinnedApps = [
        {
          class = "com.mitchellh.ghostty.home";
          workspace = "1";
        }
        {
          class = "librewolf";
          workspace = "2";
        }
        {
          class = "bruno";
          workspace = "5";
        }
        {
          class = "podman-desktop";
          workspace = "6";
        }
        {
          class = "teams-for-linux";
          workspace = "8";
        }
        {
          class = "slack";
          workspace = "18";
        }
        {
          class = "spotify";
          workspace = "19";
        }
        {
          class = "org.telegram.desktop";
          initial_title = "Telegram";
          workspace = "20";
        }
      ];
    };
  };

  programs.caelestia.settings.services.gpuType = "None";
}
