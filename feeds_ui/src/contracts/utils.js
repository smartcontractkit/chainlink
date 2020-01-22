import { ethers } from 'ethers'
import { networkName, MAINNET_ID } from 'utils/'

export function createContract(address, provider, abi) {
  return new ethers.Contract(address, abi, provider)
}

export function createInfuraProvider(networkId = MAINNET_ID) {
  const provider = new ethers.providers.JsonRpcProvider(
    `https://${networkName(networkId)}.infura.io/v3/${
      process.env.REACT_APP_INFURA_KEY
    }`,
  )
  provider.pollingInterval = 8000

  return provider
}

/**
 * @dev Format an aggregator answer
 * @param value The Big Number to format
 * @param multiply The number to divide the result by. See Multiply adapter in Chainlink Job Specification -  https://docs.chain.link/docs/job-specifications
 * @param decimalPlaces The number to show decimal places
 */

export function formatAnswer(value, multiply, decimalPlaces) {
  const decimals = 10 ** decimalPlaces
  const divided = value.mul(decimals).div(multiply)
  const formatted = ethers.utils.formatUnits(divided, decimalPlaces)
  return formatted.toString()
}

export async function getLogs(
  { provider, filter, eventInterface },
  /* eslint-disable-next-line @typescript-eslint/no-empty-function */
  cb = () => {},
) {
  const logs = await provider.getLogs(filter)
  const result = logs
    .filter(l => {
      return l !== null
    })
    .map(log => decodeLog({ log, eventInterface }, cb))
  return result
}

/* eslint-disable-next-line @typescript-eslint/no-empty-function */
export function decodeLog({ log, eventInterface }, cb = () => {}) {
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
