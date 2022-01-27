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

solc --version | grep 0.7.6 || ( echo "You need solc version 0.7.6" && exit 1 )

# FIXME: solc seems to find and compile every .sol file in this path, so invoking this once for every file produces n*3 artifacts
solc "$solpath" ${solcoptions[@]} --abi --bin --combined-json bin,bin-runtime,srcmap-runtime --overwrite -o "$(dirname $outpath)"

go run wrap.go "$abi" "$bin" "$basefilename" "$pkgname"
