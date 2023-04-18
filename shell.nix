{
  pkgs ? import <nixpkgs> {},
  dapp ? let
    lock = (builtins.fromJSON (builtins.readFile ./flake.lock)).nodes.dapp.locked;
  in
    import (fetchTarball {
      url = "https://github.com/dapphub/dapptools/archive/${lock.rev}.tar.gz";
      sha256 = lock.narHash;
    }) {inherit (pkgs) system;},
  foundry-bin ? let
    lock = (builtins.fromJSON (builtins.readFile ./flake.lock)).nodes.foundry.locked;
  in
    import (fetchTarball {
        url = "https://github.com/shazow/foundry.nix/archive/${lock.rev}.tar.gz";
        sha256 = lock.narHash;
      }
      + "/foundry-bin") {inherit pkgs;},
}: let
  go = pkgs.go_1_20;
  postgresql = pkgs.postgresql_14;
  nodejs = pkgs.nodejs-16_x;
  nodePackages = pkgs.nodePackages.override {inherit nodejs;};

  solcs = with dapp.solc-static-versions; [
    solc_0_6_6
    solc_0_7_6
    solc_0_8_6
    solc_0_8_15
    solc_0_8_16
  ];
in
  pkgs.mkShell {
    nativeBuildInputs = with pkgs;
      [
        go
        postgresql
        curl
        cacert
        which
        git

        nodejs
        nodePackages.pnpm
        # TODO: compiler / gcc for secp compilation
        nodePackages.ganache
        go-ethereum # geth

        # tooling
        gotools
        gopls
        delve
        golangci-lint
        go-mockery

        github-cli
        kubernetes

        # gofuzz
        foundry-bin
        lcov
      ]
      ++ solcs;

    # env vars
    CL_DATABASE_URL = "postgresql://chainlink:chainlink@localhost:5432/chainlink_test?sslmode=disable";
  }
