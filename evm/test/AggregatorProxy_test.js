import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'
const personas = h.personas

contract('AggregatorProxy', () => {
  const SOURCE_PATH = 'AggregatorProxy.sol'
  const jobId1 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
  const deposit = h.toWei('100')
  const basePayment = h.toWei('1')
  let link, aggregator, oc1, proxy
  let jobIds = [jobId1]

  beforeEach(async () => {
    link = await h.linkContract()
    oc1 = await h.deploy('Oracle.sol', link.address)
    aggregator = await h.deploy(
      'ConversionRate.sol',
      link.address,
      basePayment,
      1,
      [oc1.address],
      [jobId1]
    )
    await link.transfer(aggregator.address, deposit)

    proxy = await h.deploy(SOURCE_PATH, aggregator.address)
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(artifacts.require(SOURCE_PATH), [
      'aggregator',
      'currentAnswer',
      'updateAggregator',
      // Ownable
      'owner',
      'renounceOwnership',
      'transferOwnership'
    ])
  })

  describe('#currentAnswer', () => {
    const response = 54321

    beforeEach(async () => {
      const requestTx = await aggregator.requestRateUpdate()
      const request = h.decodeRunRequest(requestTx.receipt.rawLogs[3])
      await h.fulfillOracleRequest(oc1, request, response)
      assertBigNum(response, await aggregator.currentAnswer.call())
    })

    it('pulls the rate from the aggregator', async () => {
      assertBigNum(response, await proxy.currentAnswer.call())
    })
  })
})
