{pkgs}:
with pkgs; let
  go = go_1_21;
  postgresql = postgresql_14;
  nodejs = nodejs-18_x;
  nodePackages = pkgs.nodePackages.override {inherit nodejs;};
  pnpm = pnpm_9;

  mkShell' = mkShell.override {
    # The current nix default sdk for macOS fails to compile go projects, so we use a newer one for now.
    stdenv =
      if stdenv.isDarwin
      then overrideSDK stdenv "11.0"
      else stdenv;
  };
in
  mkShell' {
    nativeBuildInputs =
      [
        go
        goreleaser
        postgresql

        python3
        python3Packages.pip
        protobuf
        protoc-gen-go
        protoc-gen-go-grpc

        foundry-bin

        curl
        nodejs
        pnpm
        # TODO: compiler / gcc for secp compilation
        go-ethereum # geth
        go-mockery

        # tooling
        gotools
        gopls
        delve
        golangci-lint
        github-cli
        jq

        # gofuzz
      ]
      ++ lib.optionals stdenv.isLinux [
        # some dependencies needed for node-gyp on pnpm install
        pkg-config
        libudev-zero
        libusb1
      ];
    GOROOT = "${go}/share/go";

    PGDATA = "db";
    CL_DATABASE_URL = "postgresql://chainlink:chainlink@localhost:5432/chainlink_test?sslmode=disable";
  }
