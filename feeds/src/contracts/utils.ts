import { ethers } from 'ethers'
import { FunctionFragment } from 'ethers/utils'
import { JsonRpcProvider, Log } from 'ethers/providers'
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

const REACT_APP_INFURA_KEY = process.env.REACT_APP_INFURA_KEY
const REACT_APP_DEV_PROVIDER = process.env.REACT_APP_DEV_PROVIDER

/**
 * Initialize the infura provider for the given network
 *
 * @param networkId The network id from the whitelist
 */
export function createInfuraProvider(
  networkId: Networks = Networks.MAINNET,
): JsonRpcProvider {
  const provider = new ethers.providers.JsonRpcProvider(
    REACT_APP_DEV_PROVIDER ??
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
