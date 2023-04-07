#!/usr/bin/env bash

set -ex

optimize_runs="$1"
solpath="$2"
solcoptions=("--optimize" "--optimize-runs" "$optimize_runs" "--metadata-hash" "none")

basefilename="$(basename "$solpath" .sol)"
pkgname="$(echo $basefilename | tr '[:upper:]' '[:lower:]')"

here="$(dirname $0)"
pkgdir="${here}/${pkgname}"
mkdir -p "$pkgdir"
outpath="${pkgdir}/${pkgname}.go"
abi="${pkgdir}/${basefilename}.abi"
bin="${pkgdir}/${basefilename}.bin"

version="0.7.6"
solc="solc-${version}"
if ! [ -x "$(command -v ${solc})" ]; then
  # TODO: reuse contract's get_solc script
  solc-select install "${version}"
  solc-select use "${version}"
  solc="solc"
fi
"${solc}" --version | grep "${version}" || ( echo "You need solc version ${version}" && exit 1 )

# FIXME: solc seems to find and compile every .sol file in this path, so invoking this once for every file produces n*3 artifacts
"${solc}" "$solpath" ${solcoptions[@]} --abi --bin --combined-json bin,bin-runtime,srcmap-runtime --overwrite -o "$(dirname $outpath)"

go run wrap.go "$abi" "$bin" "$basefilename" "$pkgname"
