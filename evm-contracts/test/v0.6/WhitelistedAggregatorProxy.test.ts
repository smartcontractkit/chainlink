import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { MockAggregatorFactory } from '../../ethers/v0.6/MockAggregatorFactory'
import { WhitelistedAggregatorProxyFactory } from '../../ethers/v0.6/WhitelistedAggregatorProxyFactory'

let personas: setup.Personas
let defaultAccount: ethers.Wallet

const provider = setup.provider()
const linkTokenFactory = new contract.LinkTokenFactory()
const aggregatorFactory = new MockAggregatorFactory()
const whitelistedAggregatorProxyFactory = new WhitelistedAggregatorProxyFactory()

beforeAll(async () => {
  const users = await setup.users(provider)

  personas = users.personas
  defaultAccount = users.roles.defaultAccount
})

describe('WhitelistedAggregatorProxy', () => {
  const deposit = h.toWei('100')
  const answer = h.numToBytes32(54321)
  const roundId = 17
  const decimals = 18
  const timestamp = 678
  const startedAt = 677

  let link: contract.Instance<contract.LinkTokenFactory>
  let aggregator: contract.Instance<MockAggregatorFactory>
  let proxy: contract.CallableOverrideInstance<WhitelistedAggregatorProxyFactory>

  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(defaultAccount).deploy()
    aggregator = await aggregatorFactory
      .connect(defaultAccount)
      .deploy(decimals, 0)
    await aggregator.updateRoundData(roundId, answer, timestamp, startedAt)
    await link.transfer(aggregator.address, deposit)
    proxy = contract.callableAggregator(
      await whitelistedAggregatorProxyFactory
        .connect(defaultAccount)
        .deploy(aggregator.address),
    )
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(whitelistedAggregatorProxyFactory, [
      'aggregator',
      'decimals',
      'getAnswer',
      'getRoundData',
      'getTimestamp',
      'latestAnswer',
      'latestRound',
      'latestRoundData',
      'latestTimestamp',
      'setAggregator',
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
      await matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).latestAnswer()
      }, 'Not whitelisted')
    })

    it('latestTimestamp reverts', async () => {
      await matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).latestTimestamp()
      }, 'Not whitelisted')
    })

    it('getAnswer reverts', async () => {
      await matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).getAnswer(1)
      }, 'Not whitelisted')
    })

    it('getTimestamp reverts', async () => {
      await matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).getTimestamp(1)
      }, 'Not whitelisted')
    })

    it('latestRound reverts', async () => {
      await matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).latestRound()
      }, 'Not whitelisted')
    })

    it('getRoundData reverts', async () => {
      await matchers.evmRevert(async () => {
        await proxy.connect(personas.Carol).getRoundData(1)
      }, 'Not whitelisted')
    })
  })

  describe('if the caller is whitelisted', () => {
    beforeEach(async () => {
      await proxy.addToWhitelist(defaultAccount.address)

      matchers.bigNum(
        ethers.utils.bigNumberify(answer),
        await aggregator.latestAnswer(),
      )
      const height = await aggregator.latestTimestamp()
      assert.notEqual('0', height.toString())
    })

    it('pulls the rate from the aggregator', async () => {
      matchers.bigNum(answer, await proxy.latestAnswer())
      const latestRound = await proxy.latestRound()
      matchers.bigNum(answer, await proxy.getAnswer(latestRound))
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

    it('getRoundData works', async () => {
      const latestRound = await proxy.latestRound()
      await proxy.latestRound()
      const round = await proxy.getRoundData(latestRound)
      await proxy.getRoundData(latestRound)
      matchers.bigNum(roundId, round.roundId)
      matchers.bigNum(answer, round.answer)
      matchers.bigNum(startedAt, round.startedAt)
      matchers.bigNum(timestamp, round.updatedAt)
    })
  })
})
