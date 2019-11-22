import { ethers } from 'ethers'

export function formatEthPrice(value) {
  return ethers.utils.formatEther(value.mul(10000000000), {
    commify: true,
    pad: true,
  })
}

export async function getLogs(
  { provider, filter, eventInterface },
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
