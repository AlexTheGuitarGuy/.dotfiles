inputs: {
  hostName,
  extraModules ? [],
  extraHomeModules ? [],
}: let
  user = "alex";
in
  [
    {_module.args = {inherit inputs user;};}
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
  ++ extraModules
