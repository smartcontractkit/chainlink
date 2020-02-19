import { Command, flags } from '@oclif/command'
import * as Parser from '@oclif/parser'
import { App } from '../services/config'

export default class Compile extends Command {
  static description =
    'Run various compilers and/or codegenners that target solidity smart contracts.'

  static examples = [
    `$ belt compile all

Creating directory at abi/v0.4...
Creating directory at abi/v0.5...
Creating directory at abi/v0.6...
Compiling 35 contracts...
...
...
Aggregator artifact saved!
AggregatorProxy artifact saved!
Chainlink artifact saved!
...`,
  ]

  static flags = {
    help: flags.help({ char: 'h' }),
    config: flags.string({
      char: 'c',
      default: 'app.config.json',
      description: 'Location of the configuration file',
    }),
  }

  static args: Parser.args.IArg[] = [
    {
      name: 'compiler',
      description:
        'Compile solidity smart contracts and output their artifacts',
      options: ['solc', 'ethers', 'truffle', 'all'],
    },
  ]

  async run() {
    const { flags, args, argv } = this.parse(Compile)
    if (argv.length === 0) {
      this._help()
    }

    try {
      const config = await import('../services/config')
      const conf = config.load(flags.config)
      const compilation =
        args.compiler === 'all'
          ? this.compileAll(conf)
          : this.compileSingle(args.compiler, conf)

      await compilation
    } catch (e) {
      this.error(e)
    }
  }

  private async compileSingle(compilerName: string, conf: App) {
    type Compilers = typeof import('../services/compilers')
    const compiler: Compilers[keyof Compilers] = await import(
      `../services/compilers/${compilerName}`
    )

    await compiler.compileAll(conf)
  }

  private async compileAll(conf: App) {
    const compilers = await import('../services/compilers')

    await compilers.solc.compileAll(conf)
    await Promise.all([
      compilers.truffle.compileAll(conf),
      compilers.ethers.compileAll(conf),
    ])
  }
}
