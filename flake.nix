{
  description = "Chainlink development shell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    dapp.url = "github:dapphub/dapptools";
    foundry.url = "github:shazow/foundry.nix";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    dapp,
    foundry,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {inherit system;};
      inherit (pkgs) lib;
    in {
      devShell = pkgs.callPackage ./shell.nix {
        foundry-bin = foundry.defaultPackage.${system};
        dapp =
          dapp.packages.${system}
          // {
            solc-static-versions = with lib; filterAttrs (n: v: hasPrefix "solc" n) dapp.packages.${system};
          };
      };

      formatter = pkgs.alejandra;
    });
}
