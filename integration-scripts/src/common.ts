import { ethers } from 'ethers'
import chalk from 'chalk'

/**
 * Devnet miner address
 */
export const DEVNET_ADDRESS = '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f'

/**
 * Default credentials for testing node
 */
export const credentials = {
  email: 'notreal@fakeemail.ch',
  password: 'twochains',
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
  return keys.reduce<{ [K in T]: string }>(
    (prev, next) => {
      const envVar = process.env[next]
      if (!envVar) {
        throw new MissingEnvVarError(next)
      }

      prev[next] = envVar
      return prev
    },
    {} as { [K in T]: string },
  )
}

/**
 * Registers a global promise handler that will exit the currently
 * running process if an unhandled promise rejection is caught
 */
export function registerPromiseHandler() {
  process.on('unhandledRejection', e => {
    console.error(chalk.red(e as any))
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
