{user, ...}: {
  home.username = user;
  home.homeDirectory = "/home/${user}";
  home.stateVersion = "26.05";

  programs.nushell.enable = true;
}
