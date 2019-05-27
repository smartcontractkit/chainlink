import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'
const personas = h.personas

contract('AggregatorProxy', () => {
  const SOURCE_PATH = 'AggregatorProxy.sol'
  const jobId1 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
  const deposit = h.toWei('100')
  const basePayment = h.toWei('1')
  const response = 54321
  const response2 = 67890
  let link, aggregator, aggregator2, oc1, proxy
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
      'setAggregator',
      // Ownable
      'owner',
      'renounceOwnership',
      'transferOwnership'
    ])
  })

  describe('#currentAnswer', () => {
    beforeEach(async () => {
      const requestTx = await aggregator.requestRateUpdate()
      const request = h.decodeRunRequest(requestTx.receipt.rawLogs[3])
      await h.fulfillOracleRequest(oc1, request, response)
      assertBigNum(response, await aggregator.currentAnswer.call())
    })

    it('pulls the rate from the aggregator', async () => {
      assertBigNum(response, await proxy.currentAnswer.call())
    })

    context('after being updated to another contract', () => {
      beforeEach(async () => {
        aggregator2 = await h.deploy(
          'ConversionRate.sol',
          link.address,
          basePayment,
          1,
          [oc1.address],
          [jobId1]
        )
        await link.transfer(aggregator2.address, deposit)
        const requestTx = await aggregator2.requestRateUpdate()
        const request = h.decodeRunRequest(requestTx.receipt.rawLogs[3])
        await h.fulfillOracleRequest(oc1, request, response2)
        assertBigNum(response2, await aggregator2.currentAnswer.call())

        await proxy.setAggregator(aggregator2.address)
      })

      it('pulls the rate from the new aggregator', async () => {
        assertBigNum(response2, await proxy.currentAnswer.call())
      })
    })
  })

  describe('#setAggregator', () => {
    beforeEach(async () => {
      await proxy.transferOwnership(personas.Carol)

      aggregator2 = await h.deploy(
        'ConversionRate.sol',
        link.address,
        basePayment,
        1,
        [oc1.address],
        [jobId1]
      )

      assert.equal(aggregator.address, await proxy.aggregator.call())
    })

    context('when called by the owner', () => {
      it('pulls the rate from the new aggregator', async () => {
        await proxy.setAggregator(aggregator2.address, {
          from: personas.Carol
        })

        assert.equal(aggregator2.address, await proxy.aggregator.call())
      })
    })

    context('when called by a non-owner', () => {
      it('does not update', async () => {
        h.assertActionThrows(async () => {
          await proxy.setAggregator(aggregator2.address, {
            from: personas.Neil
          })
        })

        assert.equal(aggregator.address, await proxy.aggregator.call())
      })
    })
  })
})
