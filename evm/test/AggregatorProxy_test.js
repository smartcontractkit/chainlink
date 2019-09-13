import * as h from '../src/helpers'
import { assertBigNum } from '../src/matchers'
const personas = h.personas
const Aggregator = artifacts.require('Aggregator.sol')
const AggregatorProxy = artifacts.require('AggregatorProxy.sol')
const Oracle = artifacts.require('Oracle.sol')

contract('AggregatorProxy', () => {
  const jobId1 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
  const deposit = h.toWei('100')
  const basePayment = h.toWei('1')
  const response = 54321
  const response2 = 67890
  let link, aggregator, aggregator2, oc1, proxy

  beforeEach(async () => {
    link = await h.linkContract()
    oc1 = await Oracle.new(link.address)
    aggregator = await Aggregator.new(
      link.address,
      basePayment,
      1,
      [oc1.address],
      [jobId1],
    )
    await link.transfer(aggregator.address, deposit)

    proxy = await AggregatorProxy.new(aggregator.address)
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(AggregatorProxy, [
      'aggregator',
      'currentAnswer',
      'destroy',
      'setAggregator',
      'updatedHeight',
      // Ownable methods:
      'owner',
      'renounceOwnership',
      'transferOwnership',
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
        aggregator2 = await Aggregator.new(
          link.address,
          basePayment,
          1,
          [oc1.address],
          [jobId1],
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

  describe('#updatedHeight', () => {
    beforeEach(async () => {
      const requestTx = await aggregator.requestRateUpdate()
      const request = h.decodeRunRequest(requestTx.receipt.rawLogs[3])
      await h.fulfillOracleRequest(oc1, request, response)
      const height = await aggregator.updatedHeight.call()
      assert.notEqual('0', height.toString())
    })

    it('pulls the height from the aggregator', async () => {
      assertBigNum(
        await aggregator.updatedHeight.call(),
        await proxy.updatedHeight.call(),
      )
    })

    context('after being updated to another contract', () => {
      beforeEach(async () => {
        aggregator2 = await Aggregator.new(
          link.address,
          basePayment,
          1,
          [oc1.address],
          [jobId1],
        )
        await link.transfer(aggregator2.address, deposit)
        const requestTx = await aggregator2.requestRateUpdate()
        const request = h.decodeRunRequest(requestTx.receipt.rawLogs[3])
        await h.fulfillOracleRequest(oc1, request, response2)
        const height2 = await aggregator2.updatedHeight.call()
        assert.notEqual('0', height2.toString())
        const height1 = await aggregator.updatedHeight.call()
        assert.notEqual(height1.toString(), height2.toString())

        await proxy.setAggregator(aggregator2.address)
      })

      it('pulls the height from the new aggregator', async () => {
        assertBigNum(
          await aggregator2.updatedHeight.call(),
          await proxy.updatedHeight.call(),
        )
      })
    })
  })

  describe('#setAggregator', () => {
    beforeEach(async () => {
      await proxy.transferOwnership(personas.Carol)

      aggregator2 = await Aggregator.new(
        link.address,
        basePayment,
        1,
        [oc1.address],
        [jobId1],
      )

      assert.equal(aggregator.address, await proxy.aggregator.call())
    })

    context('when called by the owner', () => {
      it('sets the address of the new aggregator', async () => {
        await proxy.setAggregator(aggregator2.address, {
          from: personas.Carol,
        })

        assert.equal(aggregator2.address, await proxy.aggregator.call())
      })
    })

    context('when called by a non-owner', () => {
      it('does not update', async () => {
        h.assertActionThrows(async () => {
          await proxy.setAggregator(aggregator2.address, {
            from: personas.Neil,
          })
        })

        assert.equal(aggregator.address, await proxy.aggregator.call())
      })
    })
  })

  describe('#destroy', () => {
    beforeEach(async () => {
      await proxy.transferOwnership(personas.Carol)
    })

    context('when called by the owner', () => {
      it('succeeds', async () => {
        await proxy.destroy({ from: personas.Carol })

        assert.equal('0x', await web3.eth.getCode(proxy.address))
      })
    })

    context('when called by a non-owner', () => {
      it('fails', async () => {
        await h.assertActionThrows(async () => {
          await proxy.destroy({ from: personas.Eddy })
        })

        assert.notEqual('0x', await web3.eth.getCode(proxy.address))
      })
    })
  })
})
