{ stdenv, pkgs }:

pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    go_1_17

    postgresql_13
    python3
    python3Packages.pip
    curl
    nodejs-16_x
    (yarn.override { nodejs = nodejs-16_x; })
    # TODO: compiler / gcc for secp compilation
    nodePackages.ganache-cli
    # py3: web3 slither-analyzer crytic-compile
    # echidna
    # go-ethereum # geth
    # parity # openethereum
    # go-mockery

    # tooling
    goimports
    gopls
    delve
    golangci-lint

    # gofuzz
  ];
  LD_LIBRARY_PATH="${stdenv.cc.cc.lib}/lib64:$LD_LIBRARY_PATH";
  GOROOT="${pkgs.go_1_17}/share/go";

  PGDATA="db";
}

