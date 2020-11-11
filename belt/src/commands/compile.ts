import { Command, flags } from '@oclif/command'
import * as Parser from '@oclif/parser'
import { App } from '../services/config'

type Compilers = typeof import('../services/compilers')
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
      options: ['solc', 'solc-ovm', 'ethers', 'truffle', 'all', 'all-ovm'],
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

      switch (args.compiler) {
        case 'all':
          return await this.compileAll('solc', conf)
        case 'all-ovm':
          return await this.compileAll('solc-ovm', conf)
        default:
          return await this.compile(args.compiler, conf)
      }
    } catch (e) {
      this.error(e)
    }
  }

  private async getCompiler(name: string): Promise<Compilers[keyof Compilers]> {
    return await import(`../services/compilers/${name}`)
  }

  private async compile(compilerName: string, conf: App) {
    const compiler = await this.getCompiler(compilerName)
    await compiler.compileAll(conf)
  }

  private async compileAll(compilerName: string, conf: App) {
    const compilers = await import('../services/compilers')

    await this.compile(compilerName, conf)
    await Promise.all([
      compilers.truffle.compileAll(conf),
      compilers.ethers.compileAll(conf),
    ])
  }
}
