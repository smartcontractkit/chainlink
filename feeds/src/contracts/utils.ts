import { ethers } from 'ethers'
import { FunctionFragment } from 'ethers/utils'
import { JsonRpcProvider, Log, Filter } from 'ethers/providers'
import { Config } from '../config'
import { networkName, Networks } from '../utils'

/**
 * Connect to a deployed contract
 *
 * @param address Deployed address of the contract
 * @param provider Network to connect to
 * @param contractInterface ABI of the contract
 */
export function createContract(
  address: string,
  provider: JsonRpcProvider,
  contractInterface: FunctionFragment[],
) {
  return new ethers.Contract(address, contractInterface, provider)
}

/**
 * Initialize the infura provider for the given network
 *
 * @param networkId The network id from the whitelist
 */
export function createInfuraProvider(
  networkId: Networks = Networks.MAINNET,
): JsonRpcProvider {
  const provider = new ethers.providers.JsonRpcProvider(
    Config.devProvider() ??
      `https://${networkName(networkId)}.infura.io/v3/${Config.infuraKey()}`,
  )
  provider.pollingInterval = 8000

  return provider
}

/**
 * Format an aggregator answer
 *
 * @param value The Big Number to format
 * @param multiply The number to divide the result by. See Multiply adapter in Chainlink Job Specification -  https://docs.chain.link/docs/job-specifications
 * @param decimalPlaces The number to show decimal places
 */
export function formatAnswer(
  value: any,
  multiply: string,
  decimalPlaces: number,
  formatDecimalPlaces: number,
): string {
  try {
    const decimals = 10 ** decimalPlaces
    const divided = value.mul(decimals).div(multiply)
    const formatted = ethers.utils.formatUnits(
      divided,
      decimalPlaces + formatDecimalPlaces,
    )

    return formatted.toString()
  } catch {
    return value
  }
}

interface ChainlinkEvent {
  decode: Function
}

interface Query {
  provider: JsonRpcProvider
  filter: Filter
  eventInterface: ChainlinkEvent
}

/**
 * Retrieve the event logs for the matching filter
 *
 * @param query
 * @param cb
 */
export async function getLogs(
  { provider, filter, eventInterface }: Query,
  /* eslint-disable-next-line @typescript-eslint/no-empty-function */
  cb: any = () => {},
): Promise<any[]> {
  const logs = await provider.getLogs(filter)
  const result = logs
    .filter(l => l !== null)
    .map(log => decodeLog({ log, eventInterface }, cb))
  return result
}

interface LogResult {
  log: Log
  eventInterface: ChainlinkEvent
}

/**
 * Decode data & topics into args
 *
 * @param logResult
 * @param cb
 */
/* eslint-disable-next-line @typescript-eslint/no-empty-function */
export function decodeLog(
  { log, eventInterface }: LogResult,
  cb: Function = () => {},
) {
  const decodedLog = eventInterface.decode(log.data, log.topics)
  const meta = {
    blockNumber: log.blockNumber,
    transactionHash: log.transactionHash,
  }
  const result = {
    ...{ meta },
    ...cb(decodedLog),
  }

  return result
}
