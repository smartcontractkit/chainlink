{
  description = "Chainlink development shell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    foundry.url = "github:shazow/foundry.nix/monthly";
    flake-utils.url = "github:numtide/flake-utils";
    foundry.inputs.flake-utils.follows = "flake-utils";
    nur.url = "github:nix-community/NUR";
    goreleaser-nur.url = "github:goreleaser/nur";
  };

  outputs = inputs @ {
    self,
    nixpkgs,
    flake-utils,
    foundry,
    nur,
    goreleaser-nur,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let

        pkgs = import nixpkgs { inherit system; 
          config = { allowUnfree = true; }; 
          overlays = [
            (final: prev: {
              nur = import nur
                {
                  pkgs = prev;
                  repoOverrides = {
                    goreleaser = import goreleaser-nur { pkgs = prev; };
                  };
                };
            })
            foundry.overlay 
          ];
         };
    in rec {
      devShell = pkgs.callPackage ./shell.nix {
        inherit pkgs;
      };
      formatter = pkgs.nixpkgs-fmt;
    });
}
