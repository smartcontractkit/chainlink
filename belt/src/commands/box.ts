import { Command, flags } from '@oclif/command'
import * as Parser from '@oclif/parser'
import chalk from 'chalk'
import ux from 'cli-ux'
import * as cli from 'inquirer'
import {
  getSolidityVersionBy,
  getSolidityVersions,
  modifyTruffleBoxWith,
} from '../services/truffle-box'
export default class Box extends Command {
  static description = 'Modify a truffle box to a specified solidity version'

  static examples = [
    'belt box --solVer 0.6 -d path/to/box',
    'belt box --interactive path/to/box',
    'belt box -l',
  ]

  static flags = {
    help: flags.help({ char: 'h' }),
    interactive: flags.boolean({
      char: 'i',
      description: 'run this command in interactive mode',
    }),
    solVer: flags.string({
      char: 's',
      description:
        'the solidity version to change the truffle box to\neither a solidity version alias "v0.6" | "0.6" or its full version "0.6.2"',
    }),
    list: flags.boolean({
      char: 'l',
      description: 'list the available solidity versions',
    }),
    dryRun: flags.boolean({
      char: 'd',
      description: 'output the replaced strings, but dont change them',
    }),
  }

  static args: Parser.args.IArg[] = [
    {
      name: 'path',
      default: '.',
      description: 'the path to the truffle box',
    },
  ]

  async run() {
    const { flags, args } = this.parse(Box)
    if (flags.list) {
      return this.handleList()
    }

    if (flags.interactive) {
      return await this.handleInteractive(args.path, flags.dryRun)
    }
    if (flags.solVer) {
      return this.handleNonInteractive(args.path, flags.dryRun, flags.solVer)
    }

    this._help()
  }

  /**
   * Handle printing out a list of available solidity versions
   */
  private handleList() {
    const versions = getSolidityVersions().map(([alias, full]) => ({
      alias,
      full,
    }))

    this.log(chalk.greenBright('Available Solidity Versions'))
    ux.table(versions, {
      alias: {},
      full: {},
    })
    this.log('')
  }

  /**
   * Handle interactive mode.
   * Prompts user for a solidity version number then proceeds to
   * do a find-replace within their box for the selected version
   *
   * @param path The path to the truffle box
   * @param dryRun Dont replace the file contents, print the diff instead
   */
  private async handleInteractive(path: string, dryRun: boolean) {
    const solidityVersions = getSolidityVersions()
    const { solcVersion } = await cli.prompt([
      {
        name: 'solcVersion',
        type: 'list',
        choices: solidityVersions.map(([, version]) => version),
        message:
          'What version of solidity do you want to use with your smart contracts?',
      },
    ])

    const fullVersion = this.getFullVersion(solcVersion)

    modifyTruffleBoxWith(fullVersion, path, dryRun)
    this.log(
      chalk.greenBright(
        `Done!\nPlease run "npm i" to install the new changes made.`,
      ),
    )
  }

  /**
   * Handle non-interactive mode "--solVer".
   * solidity version number then proceeds to
   * do a find-replace within their box for the selected version
   *
   * @param path The path to the truffle box
   * @param dryRun Dont replace the file contents, print the diff instead
   * @param versionAliasOrVersion Either a solidity version alias "v0.6" | "0.6" or its full version "0.6.2"
   */
  private handleNonInteractive(
    path: string,
    dryRun: boolean,
    versionAliasOrVersion: string,
  ) {
    const fullVersion = this.getFullVersion(versionAliasOrVersion)

    modifyTruffleBoxWith(fullVersion, path, dryRun)
  }

  private getFullVersion(versionAliasOrVersion: string) {
    let fullVersion: ReturnType<typeof getSolidityVersionBy>
    try {
      fullVersion = getSolidityVersionBy(versionAliasOrVersion)
    } catch {
      const error = chalk.red('Could not find given solidity version\n')
      this.log(error)
      this.handleList()
      this.exit(1)
    }

    return fullVersion
  }
}
