#!/usr/bin/env bash

set -e

# Run dockerized solc as if it were local solc. Specifically targets usage by
# abigen's --solc option. No warranty for other usages!
#
# Set solcVersion below, to change the version of solidity it targets.
#
# abigen captures the output of this script. To debug, send output to a log file
#
# This script has serious limitations, and if you can install solc locally, you
# may be better off using abigen with that:
#
# - If there is any stderr output from the solc command, docker will mix that up
#   with stdout, which confuses abigen.
#
# - It is kind of slow, due to docker startup time.
#
# - It assumes that only a single solidity file is to be sent to solc. This
#   assumption is baked into abigen, it seems: comma-separated paths passed via
#   a single --sol flag aren't recognized, and mulitple --sol flags only result
#   in compilation of only the solidity file from the last --sol flag. To work
#   around this, I recommend making a solidity file which imports all the
#   relevant files. See VRFAll.sol for an example. It should also be possible to
#   run multiple abigen commands for each file, but it's not clear at this stage
#   how the dependent golang contract abstractions should then interoperate.
#
# - It assumes that all solidity dependencies are in subdirectories of the main
#   solidity file. If you need to target other solidity files, you could add
#   them to an --allow-paths flag in the solc command below.
#
# abigen makes two calls to the executable specified by the --solc <executable>
# flag. First,
#
#     solc.sh --version
#
# which reports the version of solidity it targets. Second, to compile the
# solidity file specified by --sol <path to target solidity file>
#
#     solc.sh <solc arguments> <path to target solidity file>
#
# produces required solidity compiler artifacts. Only expects a single filepath
# argument, at the end of the argument list.
solcVersion=0.5.0
dockerImg="ethereum/solc:$solcVersion"

lastArg=${*: -1} # Whitespace is critical, here
if [ "$lastArg" == "--version" ]; then
    # Stubbed out because docker startup is s l o w.
    echo "solc, the solidity compiler commandline interface"
    echo "Version: $solcVersion"
    exit 0
fi

abspath() {
    cd "$(dirname "$1")"
    printf "%s/%s\n" "$(pwd)" "$(basename "$1")"
}

targetPath=$(abspath $lastArg)

# Get the path of the contracts directory, to map in the docker container.
contractsDir=$(dirname "$targetPath")
while [ "`basename $contractsDir`" != "contracts" ]; do
    contractsDir=$(dirname "$contractsDir")
    if [ "$contractsDir" == "/" ] ; then
        echo "target files must be in contracts directory!"
        exit 1
    fi
done

# ${@:1:$#-1} splices in all arguments but the last one
docker run -t --rm -v $contractsDir:$contractsDir $dockerImg \
       ${@:1:$#-1} $targetPath
