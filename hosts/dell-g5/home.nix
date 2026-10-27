{user, ...}: {
  home.username = user;
  home.homeDirectory = "/home/${user}";
  home.stateVersion = "26.05";

  dotfiles = {
    btop.enable = true;
    claude.enable = true;
    fastfetch.enable = true;
    gh.enable = true;
    gh-dash.enable = true;
    herdr.enable = true;
    ghostty.enable = true;
    git.enable = true;
    k9s.enable = true;
    lazygit.enable = true;
    nushell.enable = true;
    nvim.enable = true;
    opencode.enable = true;
    playwright.enable = true;
    posting.enable = true;
    secrets.enable = true;
    shots.enable = true;
    ssh.enable = true;
    starship.enable = true;
  };
}
