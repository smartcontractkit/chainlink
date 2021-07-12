import 'isomorphic-unfetch'
import { ethers } from 'ethers'
import { ContractReceipt } from 'ethers/contract'
import { assert } from 'chai'
import ChainlinkClient from './chainlinkClient'
import { EventDescription } from 'ethers/utils/interface'

const DEFAULT_TIMEOUT_MS = 30_000 // 30s

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

export function printHeading(message: string) {
  const dashCount = Math.floor((80 - message.length) / 2)
  const dashes = '-'.repeat(dashCount)
  console.log(`${dashes} ${message.toUpperCase()} ${dashes}`)
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
 * const args = getEnvVars(['ENV_1', 'ENV_2'])
 * // args is now available as { ENV_1: string, ENV_2: string }
 * foo(args.ENV_1, args.ENV_2)
 *
 * @param keys The keys of the environment variables to fetch
 */
export function getEnvVars<T extends string>(keys: T[]): { [K in T]: string } {
  return keys.reduce<{ [K in T]: string }>((prev, next) => {
    const envVar = process.env[next]
    if (!envVar) {
      throw new MissingEnvVarError(next)
    }
    prev[next] = envVar
    return prev
  }, {} as { [K in T]: string })
}

export async function wait(ms: number) {
  return new Promise((res) => {
    setTimeout(res, ms)
  })
}

/**
 * changePriceFeed makes a patch request to the external adapter in tools/external-adapter
 * and changes the value reported
 *
 * @param adapter the URL of the external adapter
 * @param value the value to set on the adapter
 */
export async function changePriceFeed(adapter: string, value: number) {
  console.log('Changing price feed', adapter, 'to', value)
  const url = new URL('result', adapter).href
  const response = await fetch(url, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ result: value }),
  })
  assert(response.ok)
}

/**
 * Makes a simple get request to an endpoint and ensures the service responds.
 * Status code doesn't matter - just ensures the service is running.
 *
 * @param endpoint the url of the service
 * @param timeout the time in milliseconds to wait before erroring
 */
export async function waitForService(
  endpoint: string,
  timeout = DEFAULT_TIMEOUT_MS,
) {
  await assertAsync(
    async () =>
      fetch(endpoint)
        .then(() => true)
        .catch(() => false),
    `${endpoint} is unreachable after ${timeout}ms`,
    timeout,
  )
}

/**
 * assertAsync asserts that a condition is eventually met, with a
 * default timeout of 30 seconds
 *
 * @param f function to run every second and check for truthy return value
 * @param errorMessage error message to print if unsuccessful
 * @param timeout timeout
 */
export async function assertAsync(
  f: () => boolean | Promise<boolean>,
  errorMessage: string,
  timeout = DEFAULT_TIMEOUT_MS,
) {
  const start = new Date().getTime()
  while (new Date().getTime() < start + timeout) {
    const result = await f()
    if (result === true) {
      return
    }
    await sleep(1000)
  }
  throw new Error(errorMessage)
}

/**
 * sleep returns a Promise that resolves after the given number of milliseconds.
 *
 * @param ms the number of milliseconds to sleep
 */
export function sleep(ms: number) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms)
  })
}

/**
 * assertJobRun continuously checks the CL node for the completion of a job
 * before resolving
 *
 * @param clClient the chainlink client instance
 * @param count the expected number of job runs
 * @param errorMessage error message to throw
 */
export async function assertJobRun(
  clClient: ChainlinkClient,
  count: number,
  errorMessage: string,
) {
  await assertAsync(async () => {
    const jobRuns = clClient.getJobRuns()
    console.log('Waiting for', errorMessage, `(${clClient.name})`)
    console.log(`JOB RUNS ${clClient.name}:`, clClient.name, jobRuns)
    const jobRun = jobRuns[jobRuns.length - 1]
    return jobRuns.length === count && jobRun.status === 'completed'
  }, `${errorMessage} : job not run on ${clClient.name} wanted ${count} runs`)
}

/**
 * forces parity to mimic geth's behavior of mining a block every two seconds, by broadcasting a transaction
 * at the same interval from the provided account
 * @param wallet the account to send transactions from
 * @param interval the target interval at which to send those transactions
 */
export function setRecurringTx(wallet: ethers.Wallet, interval = 2000): number {
  return (setInterval(async () => {
    await (
      await wallet.sendTransaction({
        to: ethers.constants.AddressZero,
        value: 0,
      })
    ).wait()
  }, interval) as unknown) as number
}

/**
 * fundAddress sends 1000 eth to the address provided from the default account
 *
 * @param to address to fund
 */
export async function fundAddress(to: string, ether = 1000) {
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
    value: ethers.utils.parseEther(ether.toString()),
  })
  await tx.wait()
}

/**
 * helper function to more seamlessly wait for transactions to confirm
 * before continuing with test execution
 *
 * @param tx transaction to wait for
 */
export async function txWait(
  tx: ethers.ContractTransaction,
): Promise<ContractReceipt> {
  return await tx.wait()
}

/**
 * adds a listener to the provided contract to watch for the provided
 * events; automatically parses the log data, logs the event emissions,
 * and logs the values of the event data
 *
 * @param contract contract to watch for events on
 * @param listenTo event or list of events to listen to
 */
export function logEvents(
  contract: ethers.Contract,
  contractName: string,
  listenTo: string | string[] = [],
) {
  const listenToArr = Array.isArray(listenTo) ? listenTo : [listenTo]
  if (listenTo.length > 0) {
    assert.containsAllKeys(
      contract.interface.events,
      listenToArr,
      'contract does not have requested event type',
    )
  }

  // holds the contract's events, keys are the sha256 of the event signature
  const eventsByTopic: {
    [key: string]: EventDescription
  } = Object.entries(contract.interface.events).reduce(
    (prev, [_, event]) => ({ [event.topic]: event, ...prev }),
    {},
  )

  // listen to all events, then filter by topic
  contract.on('*', (...args) => {
    const eventEmission: ethers.Event = args[args.length - 1]
    const topic = eventEmission.topics[0] as keyof typeof eventsByTopic
    const eventDesc = eventsByTopic[topic]
    // ignore events we aren't listening for
    if (listenToArr.length != 0 && !listenToArr.includes(eventDesc.name)) return
    // decode event log and generate list of args for test log
    const eventArgs = eventDesc.decode(eventEmission.data, eventEmission.topics)
    const eventArgNames = eventDesc.inputs.map((i) => i.name) as string[]
    const eventArgList = eventArgNames
      .map((argName: string) => `\t* ${argName}: ${eventArgs[argName]}`)
      .join('\n')
    console.log(
      `${eventDesc.name} event emitted by ${contractName} ` +
        `in block #${eventEmission.blockNumber}\n${eventArgList}`,
    )
  })
}
