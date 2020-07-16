import { Command, flags } from '@oclif/command'
import * as Parser from '@oclif/parser'
import { cli } from 'cli-ux'
import { findABI } from '../services/utils'
import chalk from 'chalk'

export default class Inspect extends Command {
  static description = 'Inspects the API of a chainlink smart contract.'

  static examples = [
    'belt inspect [<options>] <version/contract>',
    'belt inspect v0.6/AccessControlledAggregator',
  ]

  static flags = {
    ...cli.table.flags(),
    config: flags.string({
      char: 'c',
      default: 'app.config.json',
      description: 'Location of the configuration file',
    }),
  }

  static args: Parser.args.IArg[] = [
    {
      name: 'versionedContractName',
      description:
        'Version and name of the chainlink contract e.g. v0.6/FluxAggregator',
    },
  ]

  async run() {
    const { args, flags } = this.parse(Inspect)

    // Load app.config.json
    const appConfig = await import('../services/config')
    const { artifactsDir } = appConfig.load(flags.config)

    // Find contract ABI
    const { found, abi } = findABI(artifactsDir, args.versionedContractName)
    if (!found) {
      this.error(
        chalk.red(
          `${args.versionedContractName} ABI not found - Run 'belt compile'`,
        ),
      )
    }

    // Parse ABI for function APIs
    const devdoc = abi['compilerOutput']['devdoc']
    const abiMethods = Object.keys(devdoc['methods'])
    const data = abiMethods.map(m => {
      return { name: m, details: devdoc['methods'][m]['details'] || '' }
    })

    // Render table
    cli.table(
      data,
      {
        name: {
          header: 'Function Name',
        },
        details: {
          header: 'Function Description',
        },
      },
      {
        printLine: this.log,
        ...flags, // parsed flags
      },
    )
  }
}
