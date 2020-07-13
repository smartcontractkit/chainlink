import { Command, flags } from '@oclif/command'
import * as Parser from '@oclif/parser'
import { cli } from 'cli-ux'
import { findABI } from '../services/utils'
import chalk from 'chalk'
import { join } from 'path'
import process from 'process'
import { merge } from 'sol-merger'

export default class Verify extends Command {
  static description = 'Verifies a chainlink smart contract on Etherscan.'

  static examples = [
    'belt verify [<options>] <version/contract> <address>',
    'belt verify v0.6/AccessControlledAggregator 0xe47D8b2CC42F07cdf05ca791bab47bc47Ed8B5CD',
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
    {
      name: 'contractAddress',
      description: 'Address of the chainlink contract',
    },
  ]

  async run() {
    const { args, flags } = this.parse(Verify)

    // Load app.config.json
    const appConfig = await import('../services/config')
    const { contractsDir, artifactsDir } = appConfig.load(flags.config)

    // Find contract ABI
    const { found, abi } = findABI(artifactsDir, args.versionedContractName)
    if (!found) {
      this.error(
        chalk.red(
          `${args.versionedContractName} ABI not found - Run 'belt compile'`,
        ),
      )
    }
    console.log(abi) // May be needed by etherscan

    // TODO: first, call etherscan to check if contract is already verified

    // https://api.etherscan.io/api?module=contract&action=getsourcecode&address=0xBB9bc244D798123fDe783fCc1C72d3Bb8C189413&apikey=YourApiKeyToken

    // Flatten contract
    // TODO: refactor to utils
    const cwd = process.cwd()
    const contractPath = join(
      cwd,
      contractsDir,
      `${args.versionedContractName}.sol`,
    )
    const mergedSource = await merge(contractPath, { removeComments: false })
    console.log(mergedSource)

    // TODO: call Etherscan API

    // axios...

    // {
    //   apikey: $('#apikey').val(),                     //A valid API-Key is required        
    //   module: 'contract',                             //Do not change
    //   action: 'verifysourcecode',                     //Do not change
    //   contractaddress: $('#contractaddress').val(),   //Contract Address starts with 0x...     
    //   sourceCode: $('#sourceCode').val(),             //Contract Source Code (Flattened if necessary)
    //   codeformat: $('#codeformat').val(),             //solidity-single-file (default) or solidity-standard-json-input (for std-input-json-format support
    //   contractname: $('#contractname').val(),         //ContractName (if codeformat=solidity-standard-json-input, then enter contractname as ex: erc20.sol:erc20)
    //   compilerversion: $('#compilerversion').val(),   // see http://etherscan.io/solcversions for list of support versions
    //   optimizationUsed: $('#optimizationUsed').val(), //0 = No Optimization, 1 = Optimization used (applicable when codeformat=solidity-single-file)
    //   runs: 200,                                      //set to 200 as default unless otherwise  (applicable when codeformat=solidity-single-file)        
    //   constructorArguements: $('#constructorArguements').val(),   //if applicable
    //   evmversion: $('#evmVersion').val(),             //leave blank for compiler default, homestead, tangerineWhistle, spuriousDragon, byzantium, constantinople, petersburg, istanbul (applicable when codeformat=solidity-single-file)
    //   licenseType: $('#licenseType').val(),
    // }
  }
}
