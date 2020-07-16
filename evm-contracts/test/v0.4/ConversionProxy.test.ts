import {
  contract,
  helpers as h,
  matchers,
  oracle,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { AggregatorFactory } from '../../ethers/v0.4/AggregatorFactory'
import { ConversionProxyFactory } from '../../ethers/v0.4/ConversionProxyFactory'
import { OracleFactory } from '../../ethers/v0.4/OracleFactory'

let personas: setup.Personas
let defaultAccount: ethers.Wallet

const provider = setup.provider()
const linkTokenFactory = new contract.LinkTokenFactory()
const aggregatorFactory = new AggregatorFactory()
const oracleFactory = new OracleFactory()
const conversionProxyFactory = new ConversionProxyFactory()

beforeAll(async () => {
  const users = await setup.users(provider)
  personas = users.personas
  defaultAccount = users.roles.defaultAccount
})

describe('ConversionProxy', () => {
  const jobId1 =
    '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
  const deposit = h.toWei('100')
  const basePayment = h.toWei('1')
  const asset = h.numToBytes32(13240400000) // The asset represented in some fiat currency
  const fiatAnswer = h.numToBytes32(124330000) // The fiat currency to USD
  const ethAnswer = h.numToBytes32(186090000000000) // The asset represented in ETH
  const convertedFiat = h.numToBytes32(16461789320) // The asset converted to USD
  const convertedEth = h.numToBytes32(2463906) // The asset converted to ETH
  const fiatDecimals = 8
  const ethDecimals = 18

  let link: contract.Instance<contract.LinkTokenFactory>
  let aggregator: contract.CallableOverrideInstance<AggregatorFactory>
  let aggregatorFiat: contract.CallableOverrideInstance<AggregatorFactory>
  let aggregatorEth: contract.CallableOverrideInstance<AggregatorFactory>
  let oc1: contract.Instance<OracleFactory>
  let proxy: contract.CallableOverrideInstance<ConversionProxyFactory>
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(defaultAccount).deploy()
    oc1 = await oracleFactory.connect(defaultAccount).deploy(link.address)
    aggregator = contract.callableAggregator(
      await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, basePayment, 1, [oc1.address], [jobId1]),
    )
    aggregatorFiat = contract.callableAggregator(
      await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, basePayment, 1, [oc1.address], [jobId1]),
    )
    aggregatorEth = contract.callableAggregator(
      await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, basePayment, 1, [oc1.address], [jobId1]),
    )
    await link.transfer(aggregator.address, deposit)
    await link.transfer(aggregatorFiat.address, deposit)
    await link.transfer(aggregatorEth.address, deposit)
    proxy = contract.callableAggregator(
      await conversionProxyFactory
        .connect(defaultAccount)
        .deploy(fiatDecimals, aggregator.address, aggregatorFiat.address),
    )
  }) // 690110508578

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(conversionProxyFactory, [
      'decimals',
      'from',
      'to',
      'latestAnswer',
      'latestRound',
      'getAnswer',
      'setAddresses',
      'latestTimestamp',
      'getTimestamp',
      // Ownable methods:
      'owner',
      'renounceOwnership',
      'transferOwnership',
    ])
  })

  it('deploys with the given parameters stored', async () => {
    assert.equal(aggregator.address, await proxy.from())
    matchers.bigNum(
      ethers.utils.bigNumberify(fiatDecimals),
      await proxy.decimals(),
    )
    assert.equal(aggregatorFiat.address, await proxy.to())
  })

  describe('#setAddresses', () => {
    let newAggregator: contract.Instance<AggregatorFactory>
    let newAggregatorFiat: contract.Instance<AggregatorFactory>

    beforeEach(async () => {
      newAggregator = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])
      newAggregatorFiat = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(link.address, basePayment, 1, [oc1.address], [jobId1])
    })

    describe('when called by a stranger', () => {
      it('reverts', async () => {
        await matchers.evmRevert(async () => {
          await proxy
            .connect(personas.Carol)
            .setAddresses(
              ethDecimals,
              newAggregator.address,
              newAggregatorFiat.address,
            )
        })
      })
    })

    describe('when called by the owner', () => {
      it('updates the addresses and decimals', async () => {
        await proxy.setAddresses(
          ethDecimals,
          newAggregator.address,
          newAggregatorFiat.address,
        )
        assert.equal(newAggregator.address, await proxy.from())
        matchers.bigNum(
          ethers.utils.bigNumberify(ethDecimals),
          await proxy.decimals(),
        )
        assert.equal(newAggregatorFiat.address, await proxy.to())
      })
    })
  })

  describe('#latestAnswer', () => {
    describe('when converting from ETH to fiat', () => {
      beforeEach(async () => {
        await proxy.setAddresses(
          ethDecimals,
          aggregator.address,
          aggregatorEth.address,
        )

        const requestTx = await aggregator.requestRateUpdate()
        const receipt = await requestTx.wait()

        const request = oracle.decodeRunRequest(receipt.logs?.[3])
        await oc1.fulfillOracleRequest(
          ...oracle.convertFufillParams(request, asset),
        )
        matchers.bigNum(
          ethers.utils.bigNumberify(asset),
          await aggregator.latestAnswer(),
        )

        const requestTx2 = await aggregatorEth.requestRateUpdate()
        const receipt2 = await requestTx2.wait()

        const request2 = oracle.decodeRunRequest(receipt2.logs?.[3])
        await oc1.fulfillOracleRequest(
          ...oracle.convertFufillParams(request2, ethAnswer),
        )
        matchers.bigNum(
          ethers.utils.bigNumberify(ethAnswer),
          await aggregatorEth.latestAnswer(),
        )
      })

      it('pulls the converted rate from the proxy', async () => {
        matchers.bigNum(
          ethers.utils.bigNumberify(convertedEth),
          await proxy.latestAnswer(),
        )
        const latestRound = await proxy.latestRound()
        matchers.bigNum(
          ethers.utils.bigNumberify(convertedEth),
          await proxy.getAnswer(latestRound),
        )
      })
    })

    describe('when converting from fiat to fiat', () => {
      beforeEach(async () => {
        const requestTx = await aggregator.requestRateUpdate()
        const receipt = await requestTx.wait()

        const request = oracle.decodeRunRequest(receipt.logs?.[3])
        await oc1.fulfillOracleRequest(
          ...oracle.convertFufillParams(request, asset),
        )
        matchers.bigNum(
          ethers.utils.bigNumberify(asset),
          await aggregator.latestAnswer(),
        )

        const requestTx2 = await aggregatorFiat.requestRateUpdate()
        const receipt2 = await requestTx2.wait()

        const request2 = oracle.decodeRunRequest(receipt2.logs?.[3])
        await oc1.fulfillOracleRequest(
          ...oracle.convertFufillParams(request2, fiatAnswer),
        )
        matchers.bigNum(
          ethers.utils.bigNumberify(fiatAnswer),
          await aggregatorFiat.latestAnswer(),
        )
      })

      it('pulls the converted rate from the proxy', async () => {
        matchers.bigNum(
          ethers.utils.bigNumberify(convertedFiat),
          await proxy.latestAnswer(),
        )
        const latestRound = await proxy.latestRound()
        matchers.bigNum(
          ethers.utils.bigNumberify(convertedFiat),
          await proxy.getAnswer(latestRound),
        )
      })
    })
  })

  describe('#latestTimestamp', () => {
    beforeEach(async () => {
      const requestTx = await aggregator.requestRateUpdate()
      const receipt = await requestTx.wait()
      const request = oracle.decodeRunRequest(receipt.logs?.[3])

      await oc1.fulfillOracleRequest(
        ...oracle.convertFufillParams(request, fiatAnswer),
      )
      const height = await aggregator.latestTimestamp()
      assert.notEqual('0', height.toString())
    })

    it('pulls the timestamp from the proxy', async () => {
      matchers.bigNum(
        await aggregator.latestTimestamp(),
        await proxy.latestTimestamp(),
      )
      const latestRound = await proxy.latestRound()
      matchers.bigNum(
        await aggregator.latestTimestamp(),
        await proxy.getTimestamp(latestRound),
      )
      matchers.bigNum(
        await aggregator.getTimestamp(latestRound),
        await proxy.getTimestamp(latestRound),
      )
    })
  })

  describe('#latestRound', () => {
    beforeEach(async () => {
      const requestTx = await aggregator.requestRateUpdate()
      const receipt = await requestTx.wait()
      const request = oracle.decodeRunRequest(receipt.logs?.[3])

      await oc1.fulfillOracleRequest(
        ...oracle.convertFufillParams(request, fiatAnswer),
      )
      const height = await aggregator.latestTimestamp()
      assert.notEqual('0', height.toString())
    })

    it('pulls the timestamp from the proxy', async () => {
      matchers.bigNum(await aggregator.latestRound(), await proxy.latestRound())
    })
  })
})
