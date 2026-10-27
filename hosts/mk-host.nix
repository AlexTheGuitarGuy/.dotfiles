inputs: {
  hostName,
  extraModules ? [],
  extraHomeModules ? [],
}: let
  user = "alex";
in
  inputs.nixpkgs.lib.nixosSystem {
    specialArgs = {inherit inputs user;};
    modules =
      [
        {nixpkgs.hostPlatform = "x86_64-linux";}
        inputs.home-manager.nixosModules.home-manager
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
      ++ extraModules;
  }
