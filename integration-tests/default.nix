{ pkgs # pkgs is a nixpkgs inheriting the system with the gomod2nix overlay to build Go pkgs
}:


let
  # imports ./ethereum.print-chain/default.nix inheriting pkgs and chainlink-dev args
  ethereum_test-run-smoke = pkgs.callPackage ./smoke { inherit pkgs; };
in
{
  # exports a package.<package-name> to be consumed by the flake.nix
  packages = {
    chainlink-dev_ethereum_test-run-smoke = ethereum_test-run-smoke.package;
  };
}
