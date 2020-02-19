import { ethers } from 'ethers'
import { createProvider, getArgs, registerPromiseHandler } from './common'
import { PrepaidAggregatorFactory } from '@chainlink/contracts/ethers/v0.5/PrepaidAggregatorFactory'

// TESTING ONLY
const privKey =
  'e08e067c08a8b145704cad7cd2129a8fbb753cd513b53cca7ff89dca724ef0a9'
// const oracle2 = '0x6c1676F8492119C16eA62E07351E0094e93c6667'

async function main() {
  registerPromiseHandler()
  const [functionName, ...functionArgs] = process.argv.slice(2)
  const { PREPAID_AGGREGATOR_ADDRESS } = getArgs(['PREPAID_AGGREGATOR_ADDRESS'])
  const provider = createProvider()
  const wallet = new ethers.Wallet(privKey, provider)
  const prepaidAggregatorFactory = new PrepaidAggregatorFactory(wallet)
  const prepaidAggregator = await prepaidAggregatorFactory.attach(
    PREPAID_AGGREGATOR_ADDRESS,
  )

  const result = await prepaidAggregator[functionName](...functionArgs)
  console.log(`${functionName}:`, formatResult(result))
}
main()

function formatResult(result: any): string {
  if (result instanceof ethers.utils.BigNumber) {
    result = result.toString()
  }
  return result.toString()
}
