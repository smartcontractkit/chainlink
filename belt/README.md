# @chainlink/belt

A toolbelt for performing various commands on chainlink smart contracts.
This cli tool is currently used within `@chainlink/contracts` for the usage of running
build and development tools across multiple solidity contract revisions.

<!-- toc -->
* [@chainlink/belt](#chainlinkbelt)
* [Usage](#usage)
* [Commands](#commands)
<!-- tocstop -->

# Usage

<!-- usage -->
```sh-session
$ npm install -g @chainlink/belt
$ belt COMMAND
running command...
$ belt (-v|--version|version)
@chainlink/belt/0.0.1 linux-x64 node-v10.16.3
$ belt --help [COMMAND]
USAGE
  $ belt COMMAND
...
```
<!-- usagestop -->

# Commands

<!-- commands -->
* [`belt box [PATH]`](#belt-box-path)
* [`belt compile [COMPILER]`](#belt-compile-compiler)
* [`belt help [COMMAND]`](#belt-help-command)

## `belt box [PATH]`

Modify a truffle box to a specified solidity version

```
USAGE
  $ belt box [PATH]

ARGUMENTS
  PATH  the path to the truffle box

OPTIONS
  -d, --dryRun         output the replaced strings, but dont change them
  -h, --help           show CLI help
  -i, --interactive    run this command in interactive mode
  -l, --list           list the available solidity versions

  -s, --solVer=solVer  the solidity version to change the truffle box to
                       either a solidity version alias "v0.6" | "0.6" or its full version "0.6.2"

EXAMPLES
  belt box --solVer 0.6 -d path/to/box
  belt box --interactive path/to/box
  belt box -l
```

_See code: [src/src/commands/box.ts](https://github.com/smartcontractkit/chainlink/blob/v0.0.1/src/src/commands/box.ts)_

## `belt compile [COMPILER]`

Run various compilers and/or codegenners that target solidity smart contracts.

```
USAGE
  $ belt compile [COMPILER]

ARGUMENTS
  COMPILER  (solc|ethers|truffle|all) Compile solidity smart contracts and output their artifacts

OPTIONS
  -c, --config=config  [default: app.config.json] Location of the configuration file
  -h, --help           show CLI help

EXAMPLE
  $ belt compile all

  Creating directory at abi/v0.4...
  Creating directory at abi/v0.5...
  Creating directory at abi/v0.6...
  Compiling 35 contracts...
  ...
  ...
  Aggregator artifact saved!
  AggregatorProxy artifact saved!
  Chainlink artifact saved!
  ...
```

_See code: [src/src/commands/compile.ts](https://github.com/smartcontractkit/chainlink/blob/v0.0.1/src/src/commands/compile.ts)_

## `belt help [COMMAND]`

display help for belt

```
USAGE
  $ belt help [COMMAND]

ARGUMENTS
  COMMAND  command to show help for

OPTIONS
  --all  see all commands in CLI
```

_See code: [@oclif/plugin-help](https://github.com/oclif/plugin-help/blob/v2.2.3/src/commands/help.ts)_
<!-- commandsstop -->
