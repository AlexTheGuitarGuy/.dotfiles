{
  description = "Alex's NixOS and Home Manager config";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";

    home-manager = {
      url = "github:nix-community/home-manager/release-26.05";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    sops-nix = {
      url = "github:Mic92/sops-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs: let
    pkgs = inputs.nixpkgs.legacyPackages.x86_64-linux;
  in {
    lib = {
      hostModules = import ./hosts/host-modules.nix inputs;
      mkHost = args:
        inputs.nixpkgs.lib.nixosSystem {
          modules = [{nixpkgs.hostPlatform = "x86_64-linux";}] ++ inputs.self.lib.hostModules args;
        };
    };

    nixosConfigurations.dell-g5 = inputs.self.lib.mkHost {hostName = "dell-g5";};

    checks.x86_64-linux.desktop = import ./checks/desktop.nix {
      inherit pkgs;
      hostModules = inputs.self.lib.hostModules {hostName = "dell-g5";};
    };
  };
}
