{pkgs, isCrib}:
with pkgs; let
  go = go_1_23;
  postgresql = postgresql_15;
  nodejs = nodejs-18_x;
  nodePackages = pkgs.nodePackages.override {inherit nodejs;};
  pnpm = pnpm_10;

  version = "v2.0.18";
  getBinDerivation =
    {
      name,
      filename,
      sha256,
    }:
    pkgs.stdenv.mkDerivation rec {
      inherit name;
      url = "https://github.com/anza-xyz/agave/releases/download/${version}/${filename}";

      nativeBuildInputs = lib.optionals stdenv.isLinux [ autoPatchelfHook ];

      autoPatchelfIgnoreMissingDeps = stdenv.isLinux;

      buildInputs = with pkgs; [stdenv.cc.cc.lib] ++ lib.optionals stdenv.isLinux [ stdenv.cc.cc.libgcc libudev-zero ];

      src = pkgs.fetchzip {
        inherit url sha256;
      };

      installPhase = ''
        mkdir -p $out/bin
        ls -lah $src
        cp -r $src/bin/* $out/bin
      '';
    };

  solanaBinaries = {
      x86_64-linux = getBinDerivation {
        name = "solana-cli-x86_64-linux";
        filename = "solana-release-x86_64-unknown-linux-gnu.tar.bz2";
        ### BEGIN_LINUX_SHA256 ###
        sha256 = "sha256-3FW6IMZeDtyU4GTsRIwT9BFLNzLPEuP+oiQdur7P13s=";
        ### END_LINUX_SHA256 ###
      };
      aarch64-apple-darwin = getBinDerivation {
        name = "solana-cli-aarch64-apple-darwin";
        filename = "solana-release-aarch64-apple-darwin.tar.bz2";
        ### BEGIN_DARWIN_SHA256 ###
        sha256 = "sha256-6VjycYU0NU0evXoqtGAZMYGHQEKijofnFQnBJNVsb6Q=";
        ### END_DARWIN_SHA256 ###
      };
  };

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
        solanaBinaries.x86_64-linux
      ] ++ lib.optionals isCrib [
        patchelf
      ] ++ pkgs.lib.optionals (pkgs.stdenv.isDarwin && pkgs.stdenv.hostPlatform.isAarch64 && !isCrib) [
        solanaBinaries.aarch64-apple-darwin
      ];

    shellHook = ''
      ${if !isCrib then "" else ''
        ${if stdenv.isDarwin then "source $(git rev-parse --show-toplevel)/nix-darwin-shell-hook.sh" else ""}
      ''}
    '';

    PGDATA = "db";
    CL_DATABASE_URL = "postgresql://chainlink:chainlink@localhost:5432/chainlink_test?sslmode=disable";
  }
