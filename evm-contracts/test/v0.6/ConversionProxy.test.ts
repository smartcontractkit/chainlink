import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { MockAggregatorFactory } from '../../ethers/v0.6/MockAggregatorFactory'
import { ConversionProxyFactory } from '../../ethers/v0.6/ConversionProxyFactory'

let personas: setup.Personas
let defaultAccount: ethers.Wallet

const provider = setup.provider()
const aggregatorFactory = new MockAggregatorFactory()
const conversionProxyFactory = new ConversionProxyFactory()

beforeAll(async () => {
  const users = await setup.users(provider)
  personas = users.personas
  defaultAccount = users.roles.defaultAccount
})

describe('ConversionProxy', () => {
  const asset = h.numToBytes32(13240400000) // The asset represented in some fiat currency
  const fiatAnswer = h.numToBytes32(124330000) // The fiat currency to USD
  const ethAnswer = h.numToBytes32(186090000000000) // The asset represented in ETH
  const convertedFiat = h.numToBytes32(16461789320) // The asset converted to USD
  const convertedEth = h.numToBytes32(2463906) // The asset converted to ETH
  const fiatDecimals = 8
  const ethDecimals = 18

  let aggregator: contract.Instance<MockAggregatorFactory>
  let aggregatorFiat: contract.Instance<MockAggregatorFactory>
  let aggregatorEth: contract.Instance<MockAggregatorFactory>
  let proxy: contract.CallableOverrideInstance<ConversionProxyFactory>
  const deployment = setup.snapshot(provider, async () => {
    aggregator = await aggregatorFactory
      .connect(defaultAccount)
      .deploy(fiatDecimals, asset)
    aggregatorFiat = await aggregatorFactory
      .connect(defaultAccount)
      .deploy(fiatDecimals, fiatAnswer)
    aggregatorEth = await aggregatorFactory
      .connect(defaultAccount)
      .deploy(ethDecimals, ethAnswer)
    proxy = contract.callableAggregator(
      await conversionProxyFactory
        .connect(defaultAccount)
        .deploy(aggregator.address, aggregatorFiat.address),
    )
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(conversionProxyFactory, [
      'decimals',
      'from',
      'getAnswer',
      'getRoundData',
      'getTimestamp',
      'latestAnswer',
      'latestRound',
      'latestRoundData',
      'latestTimestamp',
      'setAddresses',
      'to',
      // Owned methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  it('deploys with the given parameters stored', async () => {
    assert.equal(aggregator.address, await proxy.from())
    assert.equal(aggregatorFiat.address, await proxy.to())
  })

  describe('#setAddresses', () => {
    let newAggregator: contract.Instance<MockAggregatorFactory>
    let newAggregatorFiat: contract.Instance<MockAggregatorFactory>

    beforeEach(async () => {
      newAggregator = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(fiatDecimals, asset)
      newAggregatorFiat = await aggregatorFactory
        .connect(defaultAccount)
        .deploy(fiatDecimals, fiatAnswer)
    })

    describe('when called by a stranger', () => {
      it('reverts', async () => {
        await matchers.evmRevert(async () => {
          await proxy
            .connect(personas.Carol)
            .setAddresses(newAggregator.address, newAggregatorFiat.address)
        })
      })
    })

    describe('when called by the owner', () => {
      it('updates the addresses and decimals', async () => {
        await proxy.setAddresses(
          newAggregator.address,
          newAggregatorFiat.address,
        )
        matchers.bigNum(
          ethers.utils.bigNumberify(fiatDecimals),
          await proxy.decimals(),
        )
        assert.equal(newAggregator.address, await proxy.from())
        assert.equal(newAggregatorFiat.address, await proxy.to())
      })
    })
  })

  describe('#latestAnswer', () => {
    describe('when converting from ETH to fiat', () => {
      beforeEach(async () => {
        await proxy.setAddresses(aggregator.address, aggregatorEth.address)
        matchers.bigNum(
          ethers.utils.bigNumberify(asset),
          await aggregator.latestAnswer(),
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
        matchers.bigNum(
          ethers.utils.bigNumberify(asset),
          await aggregator.latestAnswer(),
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
      const height = await aggregator.latestTimestamp()
      assert.notEqual('0', height.toString())
    })

    it('pulls the timestamp from the proxy', async () => {
      matchers.bigNum(await aggregator.latestRound(), await proxy.latestRound())
    })
  })

  describe('#getAnswer', () => {
    const newAnswer = h.numToBytes32(13340400000)
    const newFiatAnswer = h.numToBytes32(125330000)
    const expectedRates = [
      h.numToBytes32(16594193320), // (13240400000 * 125330000) / (10 ** 8)
      h.numToBytes32(16719523320), // (13340400000 * 125330000) / (10 ** 8)
    ]

    beforeEach(async () => {
      await aggregator.updateAnswer(newAnswer)
      await aggregatorFiat.updateAnswer(newFiatAnswer)
    })

    it('returns the rate of the latest conversion rate', async () => {
      for (let i = 2; i > 0; i--) {
        matchers.bigNum(
          ethers.utils.bigNumberify(expectedRates[i - 1]),
          await proxy.getAnswer(i),
        )
      }
    })
  })
})
