{ pkgs ? import <nixpkgs> { } }:
with pkgs;
let
  go = go_1_19;
  postgresql = postgresql_14;
  nodejs = nodejs-16_x;
  nodePackages = pkgs.nodePackages.override { inherit nodejs; };
in
mkShell {
  nativeBuildInputs = [
    go

    postgresql
    python3
    python3Packages.pip
    curl
    nodejs
    nodePackages.pnpm
    # TODO: compiler / gcc for secp compilation
    nodePackages.ganache
    # py3: web3 slither-analyzer crytic-compile
    # echidna
    go-ethereum # geth
    # parity # openethereum
    go-mockery

    # tooling
    gotools
    gopls
    delve
    golangci-lint
    github-cli

    # gofuzz
  ] ++ lib.optionals stdenv.isLinux [
    # some dependencies needed for node-gyp on pnpm install
    pkg-config
    libudev-zero
    libusb1
  ];
  LD_LIBRARY_PATH = "${stdenv.cc.cc.lib}/lib64:$LD_LIBRARY_PATH";
  GOROOT = "${go}/share/go";

  PGDATA = "db";
  CL_DATABASE_URL = "postgresql://chainlink:chainlink@localhost:5432/chainlink_test?sslmode=disable";
  shellHook = ''
    export GOPATH=$HOME/go
    export PATH=$GOPATH/bin:$PATH
  '';
}
