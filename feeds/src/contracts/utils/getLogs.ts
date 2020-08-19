import { JsonRpcProvider, Filter } from 'ethers/providers'
import { decodeLog, ChainlinkEvent } from './decodeLog'

/**
 * Retrieve the event logs for the matching filter
 *
 * @param query
 * @param cb
 */
export async function getLogs(
  { provider, filter, eventInterface }: Query,
  /* eslint-disable-next-line @typescript-eslint/no-empty-function */
  cb: Function = () => {},
): Promise<any[]> {
  const logs = await provider.getLogs(filter)
  const result = logs
    .filter(l => l !== null)
    .map(log => decodeLog({ log, eventInterface }, cb))
  return result
}

export interface Query {
  provider: JsonRpcProvider
  filter: Filter
  eventInterface: ChainlinkEvent
}
