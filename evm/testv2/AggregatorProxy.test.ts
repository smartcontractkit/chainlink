import * as h from '../src/helpersV2'
import { assertBigNum } from '../src/matchersV2'
import { ethers } from 'ethers'
import { Instance } from '../src/contract'
import { LinkTokenFactory } from '../src/generated/LinkTokenFactory'
import { AggregatorFactory } from '../src/generated/AggregatorFactory'
import { OracleFactory } from '../src/generated/OracleFactory'
import { AggregatorProxyFactory } from '../src/generated/AggregatorProxyFactory'
import { assert } from 'chai'
import ganache from 'ganache-core'

let personas: h.Personas
let defaultAccount: ethers.Wallet

const provider = new ethers.providers.Web3Provider(ganache.provider() as any)
const linkTokenFactory = new LinkTokenFactory()
const aggregatorFactory = new AggregatorFactory()
const oracleFactory = new OracleFactory()
const aggregatorProxyFactory = new AggregatorProxyFactory()

beforeAll(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(provider)

  personas = rolesAndPersonas.personas
  defaultAccount = rolesAndPersonas.roles.defaultAccount
})

describe('AggregatorProxy', () => {
  const jobId1 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
  const deposit = h.toWei('100')
  const basePayment = h.toWei('1')
  const response = h.numToBytes32(54321)
  const response2 = h.numToBytes32(67890)

  let link: Instance<LinkTokenFactory>
  let aggregator: Instance<AggregatorFactory>
  let aggregator2: Instance<AggregatorFactory>
  let oc1: Instance<OracleFactory>
  let proxy: Instance<AggregatorProxyFactory>
  const deployment = h.useSnapshot(provider, async () => {
    link = await linkTokenFactory.connect(defaultAccount).deploy()
    oc1 = await oracleFactory.connect(defaultAccount).deploy(link.address)
    aggregator = await aggregatorFactory
      .connect(defaultAccount)
      .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])
    await link.transfer(aggregator.address, deposit)
    proxy = await aggregatorProxyFactory
      .connect(defaultAccount)
      .deploy(aggregator.address)
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(aggregatorProxyFactory, [
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
      const receipt = await requestTx.wait()

      const request = h.decodeRunRequest(receipt.logs![3])
      await h.fulfillOracleRequest(oc1, request, response)
      assertBigNum(
        ethers.utils.bigNumberify(response),
        await aggregator.currentAnswer(),
      )
    })

    it('pulls the rate from the aggregator', async () => {
      assertBigNum(response, await proxy.currentAnswer())
    })

    describe('after being updated to another contract', () => {
      beforeEach(async () => {
        aggregator2 = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])
        await link.transfer(aggregator2.address, deposit)
        const requestTx = await aggregator2.requestRateUpdate()
        const receipt = await requestTx.wait()
        const request = h.decodeRunRequest(receipt.logs![3])

        await h.fulfillOracleRequest(oc1, request, response2)
        assertBigNum(response2, await aggregator2.currentAnswer())

        await proxy.setAggregator(aggregator2.address)
      })

      it('pulls the rate from the new aggregator', async () => {
        assertBigNum(response2, await proxy.currentAnswer())
      })
    })
  })

  describe('#updatedHeight', () => {
    beforeEach(async () => {
      const requestTx = await aggregator.requestRateUpdate()
      const receipt = await requestTx.wait()
      const request = h.decodeRunRequest(receipt.logs![3])

      await h.fulfillOracleRequest(oc1, request, response)
      const height = await aggregator.updatedHeight()
      assert.notEqual('0', height.toString())
    })

    it('pulls the height from the aggregator', async () => {
      assertBigNum(
        await aggregator.updatedHeight(),
        await proxy.updatedHeight(),
      )
    })

    describe('after being updated to another contract', () => {
      beforeEach(async () => {
        aggregator2 = await aggregatorFactory
          .connect(defaultAccount)
          .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])
        await link.transfer(aggregator2.address, deposit)

        const requestTx = await aggregator2.requestRateUpdate()
        const receipt = await requestTx.wait()
        const request = h.decodeRunRequest(receipt.logs![3])

        await h.fulfillOracleRequest(oc1, request, response2)
        const height2 = await aggregator2.updatedHeight()
        assert.notEqual('0', height2.toString())

        const height1 = await aggregator.updatedHeight()
        assert.notEqual(height1.toString(), height2.toString())

        await proxy.setAggregator(aggregator2.address)
      })

      it('pulls the height from the new aggregator', async () => {
        assertBigNum(
          await aggregator2.updatedHeight(),
          await proxy.updatedHeight(),
        )
      })
    })
  })

  describe('#setAggregator', () => {
    beforeEach(async () => {
      await proxy.transferOwnership(personas.Carol.address)

      aggregator2 = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])

      assert.equal(aggregator.address, await proxy.aggregator())
    })

    describe('when called by the owner', () => {
      it('sets the address of the new aggregator', async () => {
        await proxy.connect(personas.Carol).setAggregator(aggregator2.address)

        assert.equal(aggregator2.address, await proxy.aggregator())
      })
    })

    describe('when called by a non-owner', () => {
      it('does not update', async () => {
        h.assertActionThrows(async () => {
          await proxy.connect(personas.Neil).setAggregator(aggregator2.address)
        })

        assert.equal(aggregator.address, await proxy.aggregator())
      })
    })
  })

  describe('#destroy', () => {
    beforeEach(async () => {
      await proxy.transferOwnership(personas.Carol.address)
    })

    describe('when called by the owner', () => {
      it('succeeds', async () => {
        await proxy.connect(personas.Carol).destroy()

        assert.equal('0x', await provider.getCode(proxy.address))
      })
    })

    describe('when called by a non-owner', () => {
      it('fails', async () => {
        await h.assertActionThrows(async () => {
          await proxy.connect(personas.Eddy).destroy()
        })

        assert.notEqual('0x', await provider.getCode(proxy.address))
      })
    })
  })
})
