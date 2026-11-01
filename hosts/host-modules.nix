inputs: {
  hostName,
  extraModules ? [],
  extraHomeModules ? [],
}: let
  user = "alex";
in
  [
    {_module.args = {inherit inputs user;};}
    ({lib, ...}: {
      system.configurationRevision = lib.mkDefault "dotfiles ${inputs.self.shortRev or inputs.self.dirtyShortRev or "unknown"}";
    })
    {
      nixpkgs.overlays = [
        (final: prev: {
          evidence = final.callPackage ../packages/evidence {};
          gitnexus = final.callPackage ../packages/gitnexus {};
          unstable = import inputs.nixpkgs-unstable {
            inherit (prev.stdenv.hostPlatform) system;
            inherit (prev) config;
          };
          inherit (final.unstable) ghostty lazygit yazi;
        })
      ];
    }
    inputs.home-manager.nixosModules.home-manager
    inputs.stylix.nixosModules.stylix
    {
      home-manager = {
        useGlobalPkgs = true;
        useUserPackages = true;
        backupFileExtension = "hm-backup";
        extraSpecialArgs = {inherit inputs user;};
        users.${user}.imports =
          [
            ../modules/home-manager
            ./${hostName}/home.nix
          ]
          ++ extraHomeModules;
      };
    }
    ../modules/nixos
    ./${hostName}
  ]
  ++ extraModules
