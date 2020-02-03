# @chainlink/linkbelt

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
$ npm install -g @chainlink/linkbelt
$ linkbelt COMMAND
running command...
$ linkbelt (-v|--version|version)
@chainlink/linkbelt/0.0.1 linux-x64 node-v10.16.3
$ linkbelt --help [COMMAND]
USAGE
  $ linkbelt COMMAND
...
```

<!-- usagestop -->

# Commands

<!-- commands -->

- [`linkbelt compile [COMPILER]`](#linkbelt-compile-compiler)
- [`linkbelt help [COMMAND]`](#linkbelt-help-command)

## `linkbelt compile [COMPILER]`

Run various compilers and/or codegenners that target solidity smart contracts.

```
USAGE
  $ linkbelt compile [COMPILER]

ARGUMENTS
  COMPILER  (solc|ethers|truffle|all) Compile solidity smart contracts and output their artifacts

OPTIONS
  -c, --config=config  [default: app.config.json] Location of the configuration file
  -h, --help           show CLI help

EXAMPLE
  $ linkbelt compile all

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

## `linkbelt help [COMMAND]`

display help for linkbelt

```
USAGE
  $ linkbelt help [COMMAND]

ARGUMENTS
  COMMAND  command to show help for

OPTIONS
  --all  see all commands in CLI
```

_See code: [@oclif/plugin-help](https://github.com/oclif/plugin-help/blob/v2.2.3/src/commands/help.ts)_

<!-- commandsstop -->
