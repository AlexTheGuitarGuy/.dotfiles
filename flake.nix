{
  description = "Alex's NixOS and Home Manager config";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";
    nixpkgs-unstable.url = "github:NixOS/nixpkgs/nixos-unstable";

    home-manager = {
      url = "github:nix-community/home-manager/release-26.05";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    herdr.url = "github:ogulcancelik/herdr/v0.9.5";

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

    packages.x86_64-linux.evidence = pkgs.callPackage ./packages/evidence {};
    packages.x86_64-linux.gitnexus = pkgs.callPackage ./packages/gitnexus {};

    checks.x86_64-linux.url-bar = pkgs.runCommand "url-bar-test" {nativeBuildInputs = [pkgs.nodejs];} ''
      node --test ${./modules/home-manager/playwright/config}/url-bar.test.mjs
      touch $out
    '';

    checks.x86_64-linux.desktop = import ./checks/desktop.nix {
      inherit pkgs;
      hostModules = inputs.self.lib.hostModules {hostName = "dell-g5";};
    };
  };
}
