{
  description = "Alex's NixOS and Home Manager config";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";

    home-manager = {
      url = "github:nix-community/home-manager/release-26.05";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs: {
    lib.mkHost = import ./hosts/mk-host.nix inputs;

    nixosConfigurations.dell-g5 = inputs.self.lib.mkHost {hostName = "dell-g5";};
  };
}
