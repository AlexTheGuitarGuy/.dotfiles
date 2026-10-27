{user, ...}: {
  home.username = user;
  home.homeDirectory = "/home/${user}";
  home.stateVersion = "26.05";

  dotfiles = {
    git.enable = true;
    nushell.enable = true;
    starship.enable = true;
  };
}
