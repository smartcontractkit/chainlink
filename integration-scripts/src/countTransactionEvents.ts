import { ethers } from 'ethers'
import chalk from 'chalk'
import { registerPromiseHandler, createProvider, getArgs } from './common'

async function main() {
  registerPromiseHandler()
  // see: runlog_test for these variable names
  const args = getArgs(['RUN_LOG_ADDRESS', 'txid'])

  await countTransactionEvents({
    fromAddress: args.RUN_LOG_ADDRESS,
    txId: args.txid,
  })
}
main()

interface Args {
  txId: string
  fromAddress: string
}
async function countTransactionEvents({ fromAddress, txId }: Args) {
  const provider = createProvider()
  const receipt = await provider.getTransactionReceipt(txId).catch(e => {
    console.error('Error getting transaction receipt')
    console.error(chalk.red(e))
    throw e
  })
  const count = countLogsEmittedBy(fromAddress, receipt.logs)

  emitFulfillmentEvent({ count, fromAddress, txId })
}

/**
 * Count the number of logs emitted by a particular address
 * @param address The address to match against log emissions
 * @param logs The logs emitted by a transaction
 */
function countLogsEmittedBy(
  address: string,
  logs?: ethers.providers.Log[],
): number {
  if (!logs) {
    return 0
  }

  return logs.filter(l => l.address === address).length
}

interface FulfillmentEvent {
  txId: string
  fromAddress: string
  count: number
}

/**
 * Log a fulfillment event in a suitable format for our integration tests to capture via `grep`
 * @see NB: If you change this output format, make sure to change the corresponding searches for it in the integration tests.
 * @param param0 The fulfillment event to log
 */
function emitFulfillmentEvent({ count, fromAddress, txId }: FulfillmentEvent) {
  console.log(`Events from ${fromAddress} in ${txId}: ${count}`)
}
