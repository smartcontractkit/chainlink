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
        nur.repos.goreleaser.goreleaser-pro
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

    shellHook = ''
        ${if stdenv.isDarwin then "
          source ./nix-darwin-shell-hook.sh
        " else ""}
        if [ -z "$GORELEASER_KEY" ]; then
          if [ "$IS_CRIB" ]; then
            echo 'GORELEASER_KEY must be set in CRIB environments. You can find it in our 1p vault under "goreleaser-pro-license".'
            exit 1
          fi
          echo 'If you plan on using goreleaser, make sure the env var GORELEASER_KEY is set, you can find it in our 1p vault under "goreleaser-pro-license"'
        fi
      '';

    GOROOT = "${go}/share/go";
    PGDATA = "db";
    CL_DATABASE_URL = "postgresql://chainlink:chainlink@localhost:5432/chainlink_test?sslmode=disable";
  }
