{ pkgs # pkgs is a nixpkgs inheriting based on system with the gomod2nix overlay to build Go pkgs
, chainlink-dev # chainlink-dev is the devShell of the chainlink-dev CLI
}:

# exports an shell derivation to be imported by the default.nix
pkgs.mkShell {
  # based on pkgs, add those packages into env
  buildInputs = with pkgs; [
    # Go + tools
    go
    gopls
    gotools
    go-tools
    k3dnix
  ];

  # composing chainlink-dev shell to this shell enables us to have chainlink-dev CLI into env
  inputsFrom = [
    chainlink-dev
  ];
}
