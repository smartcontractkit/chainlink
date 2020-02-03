import { Command, flags } from '@oclif/command'
import * as Parser from '@oclif/parser'
import * as compilers from '../services/compilers'
import * as config from '../services/config'

type Compiler = keyof typeof compilers

export default class Compile extends Command {
  static description =
    'Run various compilers and/or codegenners that target solidity smart contracts.'

  static examples = [
    `$ linkbelt compile all

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
    console.log({ flags, args })
    if (argv.length === 0) {
      this._help()
    }

    try {
      const conf = config.load(flags.config)
      const compilation =
        args.compiler === 'all'
          ? this.compileAll(conf)
          : this.compileSingle(args.compiler, conf)

      await compilation

      this.log('Done!')
    } catch (e) {
      this.error(e)
    }
  }

  private async compileSingle(compiler: Compiler, conf: config.App) {
    await compilers[compiler].compileAll(conf)
  }

  private async compileAll(conf: config.App) {
    await compilers.solc.compileAll(conf)
    await Promise.all([
      compilers.truffle.compileAll(conf),
      compilers.ethers.compileAll(conf),
    ])
  }
}
