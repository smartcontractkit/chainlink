import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { MockAggregatorFactory } from '../../ethers/v0.6/MockAggregatorFactory'
import { WhitelistedConversionProxyFactory } from '../../ethers/v0.6/WhitelistedConversionProxyFactory'

let personas: setup.Personas
let defaultAccount: ethers.Wallet

const provider = setup.provider()
const aggregatorFactory = new MockAggregatorFactory()
const whitelistedConversionProxyFactory = new WhitelistedConversionProxyFactory()

beforeAll(async () => {
  const users = await setup.users(provider)

  personas = users.personas
  defaultAccount = users.roles.defaultAccount
})

describe('WhitelistedConversionProxy', () => {
  const response = h.numToBytes32(13240400000)
  const fiatAnswer = h.numToBytes32(124330000)
  const convertedFiat = h.numToBytes32(16461789320)
  const decimals = 8

  let aggregator: contract.Instance<MockAggregatorFactory>
  let aggregator2: contract.Instance<MockAggregatorFactory>
  let proxy: contract.CallableOverrideInstance<WhitelistedConversionProxyFactory>

  const deployment = setup.snapshot(provider, async () => {
    aggregator = await aggregatorFactory
      .connect(defaultAccount)
      .deploy(decimals, response)
    aggregator2 = await aggregatorFactory
      .connect(defaultAccount)
      .deploy(decimals, fiatAnswer)
    proxy = contract.callableAggregator(
      await whitelistedConversionProxyFactory
        .connect(defaultAccount)
        .deploy(aggregator.address, aggregator2.address),
    )
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(whitelistedConversionProxyFactory, [
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
      // Ownable methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
      // Whitelisted methods:
      'addToWhitelist',
      'disableWhitelist',
      'enableWhitelist',
      'removeFromWhitelist',
      'whitelistEnabled',
      'whitelisted',
    ])
  })

  describe('if the caller is not whitelisted', () => {
    it('latestAnswer reverts', async () => {
      matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).latestAnswer()
      }, 'Not whitelisted')
    })

    it('latestTimestamp reverts', async () => {
      matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).latestTimestamp()
      }, 'Not whitelisted')
    })

    it('getAnswer reverts', async () => {
      matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).getAnswer(1)
      }, 'Not whitelisted')
    })

    it('getTimestamp reverts', async () => {
      matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).getTimestamp(1)
      }, 'Not whitelisted')
    })

    it('latestRound reverts', async () => {
      matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).latestRound()
      }, 'Not whitelisted')
    })
  })

  describe('if the caller is whitelisted', () => {
    beforeEach(async () => {
      await proxy.addToWhitelist(defaultAccount.address)

      matchers.bigNum(
        ethers.utils.bigNumberify(response),
        await aggregator.latestAnswer(),
      )
      const height = await aggregator.latestTimestamp()
      assert.notEqual('0', height.toString())
    })

    it('pulls the rate from the aggregator', async () => {
      matchers.bigNum(convertedFiat, await proxy.latestAnswer())
      const latestRound = await proxy.latestRound()
      matchers.bigNum(convertedFiat, await proxy.getAnswer(latestRound))
    })

    it('pulls the timestamp from the aggregator', async () => {
      matchers.bigNum(
        await aggregator.latestTimestamp(),
        await proxy.latestTimestamp(),
      )
      const latestRound = await proxy.latestRound()
      matchers.bigNum(
        await aggregator.latestTimestamp(),
        await proxy.getTimestamp(latestRound),
      )
    })
  })
})
