import { ethers } from 'ethers'
import { JsonRpcProvider, Log } from 'ethers/providers'
import { networkName, Networks } from '../utils'

/**
 * TODO
 *
 * @param address TODO
 * @param provider TODO
 * @param abi TODO
 */
export function createContract(
  address: string,
  provider: JsonRpcProvider,
  abi: any,
) {
  return new ethers.Contract(address, abi, provider)
}

const REACT_APP_INFURA_KEY = process.env.REACT_APP_INFURA_KEY

/**
 * Initialize the infura provider for the given network
 *
 * @param networkId The network id from the whitelist
 */
export function createInfuraProvider(
  networkId: Networks = Networks.MAINNET,
): JsonRpcProvider {
  const provider = new ethers.providers.JsonRpcProvider(
    `https://${networkName(networkId)}.infura.io/v3/${REACT_APP_INFURA_KEY}`,
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
  multiply?: string,
  decimalPlaces = 2,
): string {
  const decimals = 10 ** decimalPlaces
  const divided = value.mul(decimals).div(multiply)
  const formatted = ethers.utils.formatUnits(divided, decimalPlaces)
  return formatted.toString()
}

interface Filter {
  fromBlock: any
  toBlock: any
}

// event ChainlinkRequested(bytes32 indexed id);
// event ChainlinkFulfilled(bytes32 indexed id);
// event ChainlinkCancelled(bytes32 indexed id);
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
 * @param query TODO
 * @param cb TODO
 */
export async function getLogs(
  { provider, filter, eventInterface }: Query,
  /* eslint-disable-next-line @typescript-eslint/no-empty-function */
  cb = () => {},
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
 * TODO
 *
 * @param todo TODO
 * @param cb TODO
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
