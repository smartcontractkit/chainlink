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
  flattenContract,
} from '../services/utils'
import Etherscan from '../services/etherscan'

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
    const { args, argv, flags } = this.parse(Deploy)
    const inputs = argv.slice(Object.keys(Deploy.args).length)

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

    // Load transaction overrides
    const overrides: DeployOverrides = {
      gasPrice: flags.gasPrice || config.gasPrice,
      gasLimit: flags.gasLimit || config.gasLimit,
      ...(flags.nonce && { nonce: flags.nonce }),
      ...(flags.value && { value: flags.value }),
    }

    // Initialize ethers wallet (signer + provider)
    const wallet = initWallet(config)

    const contractAddress = await this.deployContract(
      wallet,
      artifactsDir,
      args.versionedContractName,
      inputs,
      overrides,
    )
    if (!contractAddress) this.error('Deployed contract address is undefined.')

    await this.verifyContract(
      contractsDir,
      args.versionedContractName,
      contractAddress,
      config.chainId,
      config.etherscanAPIKey,
    )
  }

  /**
   * Deploys a smart contract.
   *
   * @param wallet Ethers wallet (signer + provider)
   * @param artifactsDir ABI directory e.g. 'abi'
   * @param versionedContractName Version and name of the chainlink contract e.g. v0.6/FluxAggregator
   * @param inputs Array of function inputs
   * @param overrides Contract call overrides e.g. gasLimit
   */
  private async deployContract(
    wallet: ethers.Wallet,
    artifactsDir: string,
    versionedContractName: string,
    inputs: string[],
    overrides: DeployOverrides,
  ) {
    // Find contract ABI
    const { found, abi } = findABI(artifactsDir, versionedContractName)
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
      // TODO: add numConfirmations to .beltrc
      const receipt = await contract.deployTransaction.wait(1) // wait for 1 confirmation
      cli.action.stop(`Deployed in tx ${receipt.transactionHash}`)
      return contract.address
    } catch (e) {
      this.error(chalk.red(e))
    }
    return
  }

  /**
   * Verifies a smart contract on Etherscan.
   *
   * @param contractsDir Contract directory e.g. 'src'
   * @param versionedContractName Version and name of the chainlink contract e.g. v0.6/FluxAggregator
   * @param contractAddress
   * @param chainId
   * @param etherscanAPIKey
   */
  private async verifyContract(
    contractsDir: string,
    versionedContractName: string,
    contractAddress: string,
    chainId: number,
    etherscanAPIKey: string,
  ) {
    // Skip if contract is already verified
    const etherscan = new Etherscan(chainId, etherscanAPIKey)
    const isVerified = await etherscan.isVerified(contractAddress)
    if (isVerified) {
      this.error(
        chalk.red(
          `${versionedContractName} at ${contractAddress} already verified.`,
        ),
      )
    }

    // Flatten contract
    const mergedSource = await flattenContract(
      `${contractsDir}/${versionedContractName}`,
    )
    console.log(mergedSource.length)

    // TODO: implement rest of verify (see commands/verify.ts)
  }
}
