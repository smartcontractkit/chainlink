import chalk from 'chalk'
import { ethers } from 'ethers'
import 'source-map-support/register'

/**
 * Devnet miner address
 * FIXME: duplicated
 */
export const DEVNET_ADDRESS = '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f'
export const DEVNET_PRIVATE_KEY =
  '34d2ee6c703f755f9a205e322c68b8ff3425d915072ca7483190ac69684e548c'

/**
 * Default credentials for testing node
 * FIXME: duplicated
 */
export const credentials = {
  email: 'notreal@fakeemail.ch',
  password: 'twochains',
}

export const GETH_DEV_ADDRESS = '0x7db75251a74f40b15631109ba44d33283ed48528'

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

export async function wait(ms: number) {
  return new Promise(res => {
    setTimeout(res, ms)
  })
}

export async function waitForService(endpoint: string, timeout = 20000) {
  await assertAsync(
    async () => {
      return fetch(endpoint)
        .then(() => true)
        .catch(() => false)
    },
    `${endpoint} is unreachable after ${timeout}ms`,
    timeout,
  )
}

/**
 * assertAsync asserts that a condition is evantually met, with a
 * default timeout of 20 seconds
 * @param f function to run every second and check for truthy return value
 * @param errorMessage error message to print if unseccessful
 * @param timeout timeout
 */
export async function assertAsync(
  f: () => boolean | Promise<boolean>,
  errorMessage: string,
  timeout = 20000,
) {
  return new Promise((res, rej) => {
    // eslint-disable-next-line
    let interval: NodeJS.Timeout, timer: NodeJS.Timeout

    function resolveIfFulfilled(fulfilled: boolean) {
      if (fulfilled === true) {
        clearTimeout(timer)
        clearInterval(interval)
        res()
      }
    }

    timer = setTimeout(() => {
      clearInterval(interval)
      rej(errorMessage)
    }, timeout)

    interval = setInterval(() => {
      const result = f()
      if (result instanceof Promise) {
        result.then(resolveIfFulfilled)
      } else {
        resolveIfFulfilled(result)
      }
    }, 1000)
  })
}

export async function fundAddress(to: string) {
  const gethMode = !!process.env.GETH_MODE || false
  const provider = createProvider()
  let signer: ethers.Signer
  if (gethMode) {
    signer = provider.getSigner(GETH_DEV_ADDRESS)
  } else {
    signer = new ethers.Wallet(DEVNET_PRIVATE_KEY).connect(provider)
  }
  const tx = await signer.sendTransaction({
    to,
    value: ethers.utils.parseEther('1000'),
  })
  await tx.wait()
}
