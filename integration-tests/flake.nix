{
  description = "Integration tests";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    # define gomod2nix which is a dependency to build Go lang plugins
    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.utils.follows = "flake-utils";
    };

    # define chainlink-dev flake as a input
    chainlink-dev.url = "git+ssh://git@github.com/smartcontractkit/chainlink-project-nix-poc?ref=main&dir=chainlink-dev"; # this could also be a project URL
    chainlink.url = "git+ssh://git@github.com/smartcontractkit/chainlink-project-nix-poc?ref=main&dir=chainlink";

    # defines plugins module (default.nix) that aggregates all plugins of this project
    plugins = {
      url = "./";
      flake = false;
    };
  };
  outputs = { self, nixpkgs, flake-utils, gomod2nix, ... }@inputs:
      # it enables us to generate all outputs for each default system
      (flake-utils.lib.eachDefaultSystem
        (system:
          let
            pkgs = import nixpkgs {
              inherit system;

              # add gomod2nix in pkgs as a Go lang dependecy for plugins
              overlays = [
                gomod2nix.overlays.default
              ];
            };

            # wrap together dependencies for plugins
            commonArgs = {
              inherit pkgs;
            };

            # import custom plugins with required args
            plugins = pkgs.callPackage ./. commonArgs;

            # it flats the set tree of the plugins for packages and shell.
            # (e.g. [ packages.pkg1, packages.pkg2 ] -> [ pkg1, pkg2 ] )
            pluginsPackage = flake-utils.lib.flattenTree plugins.packages;
          in
          {
            # it outputs packages all packages defined in plugins
            packages = pluginsPackage;

            # it outputs the default shell
            devShells.default =
              pkgs.mkShell {
                buildInputs = with pkgs; [
                   # Go + tools
                   go
                   gopls
                   gotools
                   go-tools

                    # k8s
                    kube3d
                    kubectl
                    k9s
                    kubernetes-helm

                    # NOTE: cannot import all packages through chainlink-dev
                    # nested relative path import of chainlink breaks nix (cannot dynamically calculate absoluate path)
                    inputs.chainlink-dev.packages.${system}.chainlink-dev
                    inputs.chainlink.packages.${system}.chainlink

                    # add all local plugins of this project
                    self.packages.${system}.integration-tests_run-smoke
                  ]++ pkgs.lib.optionals pkgs.stdenv.isDarwin [
                         # Additional darwin specific inputs can be set here
                         pkgs.libiconv
                         pkgs.darwin.apple_sdk.frameworks.Security
                       ];
              };

            apps.integration-tests_run-smoke = {
              type = "app";
              program = "${self.packages.${system}.integration-tests_run-smoke}/bin/integration-tests_run-smoke";
            };
          })
      );

}
