{ pkgs # pkgs is a nixpkgs inheriting the system with the gomod2nix overlay to build Go pkgs
}:


let
  # imports ./ethereum.print-chain/default.nix inheriting pkgs and chainlink-dev args
  test-run-smoke = pkgs.callPackage ./smoke { inherit pkgs; };
in
{
  # exports a package.<package-name> to be consumed by the flake.nix
  packages = {
    integration-tests_run-smoke = test-run-smoke.package;
  };
}
