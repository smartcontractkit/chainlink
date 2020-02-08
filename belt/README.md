# @chainlink/belt

A toolbelt for performing various commands on chainlink smart contracts.
This cli tool is currently used within `@chainlink/contracts` for the usage of running
build and development tools across multiple solidity contract verisions.

<!-- toc -->

- [Usage](#usage)
- [Commands](#commands)
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

- [`belt compile [COMPILER]`](#belt-compile-compiler)
- [`belt help [COMMAND]`](#belt-help-command)

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
