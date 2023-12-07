{ pkgs }:

let
  scriptName = "test-run-smoke";
  go = pkgs.go_1_21;
  scriptBuildInputs = with pkgs; [ go go-ethereum ]; # add geth to handle missing underlying dependencies
  testDir = toString ./.;

  wrapperScript = pkgs.writeShellScriptBin scriptName ''
    #!/usr/bin/env bash
    cd ${testDir}
    ${testDir}/run_test.sh "$@"
  '';

  symlink = pkgs.symlinkJoin {
    name = scriptName;
    paths = [ wrapperScript ] ++ scriptBuildInputs;
    buildInputs = [ pkgs.makeWrapper ];
    postBuild = ''
      wrapProgram $out/bin/${scriptName} \
        --prefix PATH : $out/bin \
        --set GOROOT "${go}/share/go" \
        --run 'export GOPATH=$HOME/go' \
        --run 'export PATH=$GOPATH/bin:$PATH'
    '';
  };

  package = pkgs.stdenv.mkDerivation {
    name = scriptName;
    buildCommand = "ln -s ${symlink} $out";
    installPhase = ''
      mkdir -p $out/bin
      cp ${wrapperScript}/bin/${scriptName} $out/bin/
      chmod +x $out/bin/${scriptName}
    '';
    meta = {
      description = "Main script to run ethereum smoke tests";
      longDescription = ''
        TBD
      '';
      homepage = "TBD";
      maintainers = [ "QA, TT" ];
    };
  };
in
{
  package = package;
  devShells.default = pkgs.mkShell {
    # based on pkgs, add those packages into env
    buildInputs = with pkgs; [
      # Go + tools
      go
      gopls
      gotools
      go-tools
    ];

    # import any additional build inputs from goPackage
    inputsFrom = [ package ];
  };
}
