{
  description = "Chainlink development shell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    foundry.url = "github:shazow/foundry.nix";
    flake-utils.url = "github:numtide/flake-utils";
    foundry.inputs.flake-utils.follows = "flake-utils";
  };

  outputs = inputs@{ self, nixpkgs, flake-utils, foundry, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; overlays = [ foundry.overlay ]; };
      in
      rec {
        devShell = pkgs.callPackage ./shell.nix {
          inherit pkgs;
        };
        formatter = pkgs.nixpkgs-fmt;
      });
}
