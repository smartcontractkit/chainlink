import { Command, flags } from '@oclif/command'
import * as Parser from '@oclif/parser'
import { cli } from 'cli-ux'
import { findABI, flattenContract } from '../services/utils'
import chalk from 'chalk'
// import axios from 'axios'
import { RuntimeConfig, RuntimeConfigParser } from '../services/runtimeConfig'
import Etherscan from '../services/etherscan'

const conf = new RuntimeConfigParser()

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

    // Check .beltrc exists
    let config: RuntimeConfig
    try {
      config = conf.load()
    } catch (e) {
      this.error(chalk.red(e))
    }

    // Skip if contract is already verified
    const etherscan = new Etherscan(config.chainId, config.etherscanAPIKey)
    const isVerified = await etherscan.isVerified(args.contractAddress)
    if (isVerified) {
      this.error(
        chalk.red(
          `${args.versionedContractName} at ${args.contractAddress} already verified.`,
        ),
      )
    }

    // Find contract ABI
    const { found, abi } = findABI(artifactsDir, args.versionedContractName)
    if (!found) {
      this.error(
        chalk.red(
          `${args.versionedContractName} ABI not found - Run 'belt compile'`,
        ),
      )
    }
    const abiMetadata = JSON.parse(abi['compilerOutput']['metadata'])

    // Flatten contract
    const mergedSource = await flattenContract(
      `${contractsDir}/${args.versionedContractName}`,
    )
    console.log(mergedSource.length)

    // TODO: fetch constructor values
    // Fetch the contract creation transaction to extract the input data
    // const encodedConstructorArgs = await Etherscan.fetchConstructorValues(contractAddress)

    const params = {
      apikey: config.etherscanAPIKey,
      module: 'contract',
      action: 'verifysourcecode',
      contractaddress: args.contractAddress,
      sourceCode: mergedSource,
      codeformat: 'solidity-single-file',
      contractname: abi.contractName,
      compilerversion: abiMetadata.compiler.version,
      optimizationUsed: abiMetadata.settings.optimizer.enabled,
      runs: abiMetadata.settings.optimizer.runs,
      constructorArguements: {}, // TODO
      evmversion: abiMetadata.settings.evmVersion,
      licenseType: 3, // (MIT)
    }

    // TODO: link libraries, pull info from ABI metadata

    console.log(params)

    // // TODO: call Etherscan verification API
    // // TODO: hardcoded to ropsten
    // const res = await axios.post('https://api-ropsten.etherscan.io/api', {
    //   params,
    // })
    // console.log(res)

    // TODO: check and poll for verification status
  }
}
