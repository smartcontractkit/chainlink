/* eslint-disable @typescript-eslint/no-use-before-define */
import { Command, flags } from '@oclif/command'
import * as Parser from '@oclif/parser'
import cli from 'cli-ux'
import chalk from 'chalk'
import { ethers } from 'ethers'
import { RuntimeConfigParser, RuntimeConfig } from '../services/runtimeConfig'
import {
  findABI,
  parseArrayInputs,
  initWallet,
  getConstructorABI,
} from '../services/utils'

const conf = new RuntimeConfigParser()

export interface DeployOverrides {
  gasPrice?: number
  gasLimit?: number
  nonce?: number
  value?: number
}

export default class Deploy extends Command {
  static description = 'Deploys a chainlink smart contract.'

  static examples = [
    'belt deploy [<options>] <version/contract> [<args>]',
    'belt deploy v0.6/AccessControlledAggregator 0x01be23585060835e02b77ef475b0cc51aa1e0709 160000000000000000 300 1 1000000000 18 LINK/USD',
  ]
  static strict = false

  static flags = {
    help: flags.help({ char: 'h' }),
    gasPrice: flags.integer({
      char: 'g',
      description: 'Gas price',
    }),
    gasLimit: flags.integer({
      char: 'l',
      description: 'Gas limit',
    }),
    nonce: flags.integer({
      char: 'n',
      description: 'Nonce',
    }),
    value: flags.integer({
      char: 'v',
      description: 'Value',
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
    const { args, argv, flags } = this.parse(Deploy)
    const inputs = argv.slice(Object.keys(Deploy.args).length)

    // Check .beltrc exists
    let config: RuntimeConfig
    try {
      config = conf.load()
    } catch (e) {
      this.error(chalk.red(e))
    }

    // Load transaction overrides
    const overrides: DeployOverrides = {
      gasPrice: flags.gasPrice || config.gasPrice,
      gasLimit: flags.gasLimit || config.gasLimit,
      ...(flags.nonce && { nonce: flags.nonce }),
      ...(flags.value && { value: flags.value }),
    }

    // Initialize ethers wallet (signer + provider)
    const wallet = initWallet(config)

    await this.deployContract(
      wallet,
      args.versionedContractName,
      inputs,
      overrides,
    )
  }

  /**
   * Deploys a smart contract.
   *
   * @param wallet Ethers wallet (signer + provider)
   * @param versionedContractName Version and name of the chainlink contract e.g. v0.6/FluxAggregator
   * @param inputs Array of function inputs
   * @param overrides Contract call overrides e.g. gasLimit
   */
  private async deployContract(
    wallet: ethers.Wallet,
    versionedContractName: string,
    inputs: string[],
    overrides: DeployOverrides,
  ) {
    // Find contract ABI
    const { found, abi } = findABI(versionedContractName)
    if (!found) {
      this.error(
        chalk.red(
          `${versionedContractName} ABI not found - Run 'belt compile'`,
        ),
      )
    }

    // Validate command inputs against constructor inputs
    const constructorABI = getConstructorABI(abi)
    const numConstructorInputs = constructorABI['inputs'].length
    if (numConstructorInputs !== inputs.length) {
      this.error(
        chalk.red(
          `Received ${inputs.length} arguments, constructor expected ${numConstructorInputs}`,
        ),
      )
    }

    // Transforms string arrays to arrays
    const parsedInputs = parseArrayInputs(inputs)

    // Intialize ethers contract factory
    const factory = new ethers.ContractFactory(
      abi['compilerOutput']['abi'],
      abi['compilerOutput']['evm']['bytecode'],
      wallet,
    )

    // Deploy contract
    let contract: ethers.Contract
    try {
      contract = await factory.deploy(...parsedInputs, overrides)
      cli.action.start(
        `Deploying ${versionedContractName} to ${contract.address} `,
      )
      const receipt = await contract.deployTransaction.wait() // defaults to 1 confirmation
      cli.action.stop(`Deployed in tx ${receipt.transactionHash}`)
      this.log(contract.address)
    } catch (e) {
      this.error(chalk.red(e))
    }
    return
  }
}