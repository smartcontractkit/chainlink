import { Log } from 'ethers/providers'

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

export interface ChainlinkEvent {
  decode: Function
}

interface LogResult {
  log: Log
  eventInterface: ChainlinkEvent
}
