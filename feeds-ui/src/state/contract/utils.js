import { provider } from './contract'
import { ethers } from 'ethers'

export function removeListener(filter, eventListener) {
  try {
    provider.removeListener(filter, eventListener)
  } catch (error) {}
}

export function formatEthPrice(value) {
  return ethers.utils.formatEther(value.mul(10000000000), {
    commify: true,
    pad: true
  })
}

export async function getLogs({ name, filter, eventInterface, cb = () => {} }) {
  const logs = await provider.getLogs(filter)
  const blockTimePromises = []

  for (let i = 0; i < logs.length; i++) {
    blockTimePromises.push(provider.getBlock(logs[i].blockNumber))
  }
  const blockTimes = await Promise.all(blockTimePromises)

  return blockTimes
    .filter(block => {
      return block !== null
    })
    .map((block, i) => {
      const decodedLog = eventInterface.decode(logs[i].data, logs[i].topics)
      const meta = {
        name,
        timestamp: block.timestamp,
        blockNumber: logs[i].blockNumber,
        transactionHash: logs[i].transactionHash
      }

      return {
        ...{ meta },
        ...{ rawLog: logs[i] },
        ...{ decodedLog },
        ...cb(decodedLog)
      }
    })
}

export async function getLogsFromEvent({
  name,
  log,
  eventInterface,
  cb = () => {}
}) {
  const block = await provider.getBlock(log.blockNumber)
  const decodedLog = eventInterface.decode(log.data, log.topics)
  const meta = {
    name,
    timestamp: block.timestamp,
    blockNumber: log.blockNumber,
    transactionHash: log.transactionHash
  }
  return {
    ...{ meta },
    ...{ rawLog: log },
    ...{ decodedLog },
    ...cb(decodedLog)
  }
}

export async function getLogsFromEventWithoutTimestamp({
  name,
  log,
  eventInterface,
  cb = () => {}
}) {
  const decodedLog = eventInterface.decode(log.data, log.topics)
  const meta = {
    name,
    blockNumber: log.blockNumber,
    transactionHash: log.transactionHash
  }
  return {
    ...{ meta },
    ...{ rawLog: log },
    ...{ decodedLog },
    ...cb(decodedLog)
  }
}

export async function getLogsWithoutTimestamp({
  name,
  filter,
  eventInterface,
  cb = () => {}
}) {
  const logs = await provider.getLogs(filter)

  return logs
    .filter(log => {
      return log !== null
    })
    .map((log, i) => {
      const decodedLog = eventInterface.decode(logs[i].data, logs[i].topics)
      const meta = {
        name,
        blockNumber: logs[i].blockNumber,
        transactionHash: logs[i].transactionHash
      }

      return {
        ...{ meta },
        ...{ rawLog: logs[i] },
        ...{ decodedLog },
        ...cb(decodedLog)
      }
    })
}
