import { ethers } from 'ethers'
import chalk from 'chalk'
import { Compiler } from '@0x/sol-compiler'
import {
  SolCompilerArtifactAdapter,
  Web3ProviderEngine,
  RevertTraceSubprovider,
} from '@0x/sol-trace'
import 'source-map-support/register'
import {
  FakeGasEstimateSubprovider,
  GanacheSubprovider,
} from '@0x/subproviders'
import { resolve, join } from 'path'
import { rm, cp } from 'shelljs'

/**
 * Devnet miner address
 * FIXME: duplicated
 */
export const DEVNET_ADDRESS = '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f'

/**
 * Default credentials for testing node
 * FIXME: duplicated
 */
export const credentials = {
  email: 'notreal@fakeemail.ch',
  password: 'twochains',
}

export async function createTraceProvider() {
  const mnemonic =
    'dose weasel clever culture letter volume endorse used harvest ripple circle install'
  const accountIndex = 0
  const path = `m/44'/60'/${accountIndex}'/0/0`
  const root = resolve(__dirname, '../')
  const evmv05Contracts = join(
    resolve(__dirname, '../../'),
    'evm',
    'v0.5',
    'contracts',
  )
  const contracts = join(root, 'contracts', 'v0.5')
  const artifacts = join(root, 'artifacts')
  console.warn(chalk.yellow('Removing contracts/v0.5 dir:', contracts))
  rm('-rf', contracts)
  console.warn(chalk.yellow('Removing artifacts dir:', artifacts))
  rm('-r', artifacts)
  console.log(
    chalk.green(`Copying contracts from ${evmv05Contracts} to ${contracts}`),
  )
  cp('-r', evmv05Contracts, contracts)

  const compiler = new Compiler({
    artifactsDir: artifacts,
    contracts: '*',
    contractsDir: contracts,
    solcVersion: '0.5.0',
    useDockerisedSolc: false,
    compilerSettings: {
      outputSelection: {
        '*': {
          '*': [
            'abi',
            'evm.bytecode.object',
            'evm.bytecode.sourceMap',
            'evm.deployedBytecode.object',
            'evm.deployedBytecode.sourceMap',
          ],
        },
      },
    },
  })

  console.log(chalk.green('Compiling contracts in:', contracts))
  console.log(chalk.green('Outputting artifacts to:', artifacts))
  await compiler.compileAsync()

  const defaultFromAddress = await ethers.Wallet.fromMnemonic(
    mnemonic,
    path,
  ).getAddress()
  console.log(
    chalk.green(
      'Default from address derived from mnemonic:',
      defaultFromAddress,
    ),
  )

  const artifactAdapter = new SolCompilerArtifactAdapter(artifacts, contracts)
  const revertTraceSubprovider = new RevertTraceSubprovider(
    artifactAdapter,
    defaultFromAddress,
    true,
  )

  const providerEngine = new Web3ProviderEngine()
  providerEngine.addProvider(new FakeGasEstimateSubprovider(4 * 10 ** 6)) // Ganache does a poor job of estimating gas, so just crank it up for testing.
  providerEngine.addProvider(revertTraceSubprovider)
  providerEngine.addProvider(
    // Start an in-process ganache instance
    new GanacheSubprovider({
      mnemonic,
      hdPath: path,
      vmErrorsOnRPCResponse: true,
    }),
  )
  providerEngine.start()

  const provider = new ethers.providers.Web3Provider(providerEngine)
  const accounts = await provider.listAccounts()
  console.log(chalk.green(`Accounts from provider: ${accounts}`))

  return { provider, defaultFromAddress }
}

export function createProvider(): ethers.providers.JsonRpcProvider {
  const port = process.env.ETH_HTTP_PORT || `18545`
  const providerURL = process.env.ETH_HTTP_URL || `http://localhost:${port}`

  return new ethers.providers.JsonRpcProvider(providerURL)
}

/**
 * MissingEnvVarError occurs when an expected environment variable does not exist.
 */
class MissingEnvVarError extends Error {
  constructor(envKey: string) {
    super()
    this.name = 'MissingEnvVarError'
    this.message = this.formErrorMsg(envKey)
  }

  private formErrorMsg(envKey: string) {
    const errMsg = `Not enough arguments supplied. 
      Expected "${envKey}" to be supplied as environment variable.`

    return errMsg
  }
}

/**
 * Get environment variables in a friendly object format
 *
 * @example
 * const args = getArgs(['ENV_1', 'ENV_2'])
 * // args is now available as { ENV_1: string, ENV_2: string }
 * foo(args.ENV_1, args.ENV_2)
 *
 * @param keys The keys of the environment variables to fetch
 */
export function getArgs<T extends string>(keys: T[]): { [K in T]: string } {
  return keys.reduce<{ [K in T]: string }>((prev, next) => {
    const envVar = process.env[next]
    if (!envVar) {
      throw new MissingEnvVarError(next)
    }
    prev[next] = envVar
    return prev
  }, {} as { [K in T]: string })
}

/**
 * Registers a global promise handler that will exit the currently
 * running process if an unhandled promise rejection is caught
 */
export function registerPromiseHandler() {
  process.on('unhandledRejection', e => {
    console.error(e)
    console.error(chalk.red('Exiting due to promise rejection'))
    process.exit(1)
  })
}

interface DeployFactory {
  new (signer?: ethers.Signer): ethers.ContractFactory
}

interface DeployContractArgs<T extends DeployFactory> {
  Factory: T
  name: string
  signer: ethers.Signer
}

export async function deployContract<T extends DeployFactory>(
  { Factory, name, signer }: DeployContractArgs<T>,
  ...deployArgs: Parameters<InstanceType<T>['deploy']>
): Promise<ReturnType<InstanceType<T>['deploy']>> {
  const contractFactory = new Factory(signer)
  const contract = await contractFactory.deploy(...deployArgs)
  await contract.deployed()
  console.log(`Deployed ${name} at: ${contract.address}`)

  return contract as any
}
