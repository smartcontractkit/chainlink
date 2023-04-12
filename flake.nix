{
  description = "Chainlink development shell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    dapp.url = "github:dapphub/dapptools";
    foundry.url = "github:shazow/foundry.nix";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    dapp,
    foundry,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {inherit system;};
      inherit (pkgs) lib;
    in rec {
      devShell = pkgs.callPackage ./shell.nix {
        foundry-bin = foundry.defaultPackage.${system};
        dapp =
          dapp.packages.${system}
          // {
            solc-static-versions = with lib; filterAttrs (n: v: hasPrefix "solc" n) dapp.packages.${system};
          };
      };

      formatter = pkgs.alejandra;

      packages = rec {
        chainlink = let
          ui-tag = lib.fileContents ./operator_ui/TAG;
          operator-ui = builtins.fetchTarball {
            name = "smartcontractkit-operator-ui-${ui-tag}";
            url = "https://github.com/smartcontractkit/operator-ui/releases/download/${ui-tag}/smartcontractkit-operator-ui-${builtins.substring 1 99 ui-tag}.tgz";
            sha256 = "sha256:17qlf5s7pzvgibkzmf4cav3lb8ycqk793c1b1mwx4hsi39afnasr"; # needs update when operator-ui gets updated
          };
        in
          pkgs.buildGoModule rec {
            pname = "chainlink";
            version = lib.fileContents ./VERSION;
            src = ./.;
            vendorHash = "sha256-oNOLWVvRA1bAfdGgbDtawcK37WKRf43j1QhQWhBdTJk="; # needs update when go.mod gets updated
            subPackages = ["."];
            doCheck = false;
            ldflags = let
              prefix = "github.com/smartcontractkit/chainlink/v2/core/static";
            in [
              "-s"
              "-w"
              "-X ${prefix}.Version=${version}"
              "-X ${prefix}.Sha=${self.rev or "dirty"}"
              "-X ${prefix}.BuildUser=nix"
              "-X ${prefix}.BuildDate=1980-01-01T00:00:00Z"
            ];

            preBuild = ''
              cp -r ${operator-ui}/artifacts ./core/web/assets
            '';
            postInstall = ''
              mkdir -p $out/lib
              cp --reflink=auto vendor/github.com/CosmWasm/wasmvm/internal/api/libwasmvm.* $out/lib/

              patchelf --print-rpath $out/bin/chainlink \
              | sed "s|$(pwd)/vendor/github.com/CosmWasm/wasmvm/internal/api|$out/lib|" \
              | xargs patchelf $out/bin/chainlink --set-rpath
            '';
          };

        default = chainlink;
      };

      apps = rec {
        chainlink = {
          type = "app";
          program = "${packages.chainlink}/bin/chainlink";
        };
        default = chainlink;
      };
    });
}
