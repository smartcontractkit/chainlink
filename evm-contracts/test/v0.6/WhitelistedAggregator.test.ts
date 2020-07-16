import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { WhitelistedAggregatorFactory } from '../../ethers/v0.6/WhitelistedAggregatorFactory'

const aggregatorFactory = new WhitelistedAggregatorFactory()
const linkTokenFactory = new contract.LinkTokenFactory()
const provider = setup.provider()
let personas: setup.Personas
beforeAll(async () => {
  await setup.users(provider).then(u => (personas = u.personas))
})

describe('WhitelistedAggregator', () => {
  const paymentAmount = h.toWei('3')
  const deposit = h.toWei('100')
  const answer = 100
  const minAns = 1
  const maxAns = 1
  const rrDelay = 0
  const timeout = 1800
  const decimals = 18
  const description = 'LINK/USD'

  let link: contract.Instance<contract.LinkTokenFactory>
  let aggregator: contract.CallableOverrideInstance<WhitelistedAggregatorFactory>
  let nextRound: number

  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(personas.Default).deploy()
    aggregator = contract.callableAggregator(
      await (aggregatorFactory as any).connect(personas.Carol).deploy(
        link.address,
        paymentAmount,
        timeout,
        decimals,
        h.toBytes32String(description),
        // Remove when this PR gets merged:
        // https://github.com/ethereum-ts/TypeChain/pull/218
        { gasLimit: 8_000_000 },
      ),
    )
    await link.transfer(aggregator.address, deposit)
    await aggregator.updateAvailableFunds()
    matchers.bigNum(deposit, await link.balanceOf(aggregator.address))
  })

  beforeEach(async () => {
    await deployment()
    nextRound = 1
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(aggregatorFactory, [
      'acceptAdmin',
      'addOracles',
      'allocatedFunds',
      'availableFunds',
      'decimals',
      'description',
      'getAdmin',
      'getAnswer',
      'getOracles',
      'getRoundData',
      'getTimestamp',
      'latestAnswer',
      'latestRound',
      'latestRoundData',
      'latestSubmission',
      'latestTimestamp',
      'linkToken',
      'maxSubmissionCount',
      'minSubmissionCount',
      'onTokenTransfer',
      'oracleCount',
      'oracleRoundState',
      'paymentAmount',
      'removeOracles',
      'reportingRound',
      'requestNewRound',
      'restartDelay',
      'setRequesterPermissions',
      'submit',
      'timeout',
      'transferAdmin',
      'updateAvailableFunds',
      'updateFutureRounds',
      'withdrawFunds',
      'withdrawPayment',
      'withdrawablePayment',
      'VERSION',
      // Owned methods:
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

  describe('#constructor', () => {
    it('sets the paymentAmount', async () => {
      matchers.bigNum(h.bigNum(paymentAmount), await aggregator.paymentAmount())
    })

    it('sets the timeout', async () => {
      matchers.bigNum(h.bigNum(timeout), await aggregator.timeout())
    })

    it('sets the decimals', async () => {
      matchers.bigNum(h.bigNum(decimals), await aggregator.decimals())
    })

    it('sets the description', async () => {
      assert.equal(
        description,
        h.parseBytes32String(await aggregator.description()),
      )
    })
  })

  describe('#getAnswer', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .addOracles(
          [personas.Neil.address],
          [personas.Neil.address],
          minAns,
          maxAns,
          rrDelay,
        )
      await aggregator.connect(personas.Neil).submit(nextRound, answer)
    })

    describe('when the reader is not whitelisted', () => {
      it('does not allow getAnswer to be called', async () => {
        const round = await aggregator.latestRound()
        await matchers.evmRevert(
          aggregator.connect(personas.Eddy).getAnswer(round),
          'Not whitelisted',
        )
      })
    })

    describe('when the reader is whitelisted', () => {
      beforeEach(async () => {
        await aggregator
          .connect(personas.Carol)
          .addToWhitelist(personas.Eddy.address)
      })

      it('allows getAnswer to be called', async () => {
        const round = await aggregator.latestRound()
        const answer = await aggregator.connect(personas.Eddy).getAnswer(round)
        matchers.bigNum(h.bigNum(answer), answer)
      })
    })
  })

  describe('#getTimestamp', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .addOracles(
          [personas.Neil.address],
          [personas.Neil.address],
          minAns,
          maxAns,
          rrDelay,
        )
      await aggregator.connect(personas.Neil).submit(nextRound, answer)
    })

    describe('when the reader is not whitelisted', () => {
      it('does not allow getTimestamp to be called', async () => {
        const round = await aggregator.latestRound()
        await matchers.evmRevert(
          aggregator.connect(personas.Eddy).getTimestamp(round),
          'Not whitelisted',
        )
      })
    })

    describe('when the reader is whitelisted', () => {
      beforeEach(async () => {
        await aggregator
          .connect(personas.Carol)
          .addToWhitelist(personas.Eddy.address)
      })

      it('allows getTimestamp to be called', async () => {
        const round = await aggregator.latestRound()
        const currentTimestamp = await aggregator
          .connect(personas.Eddy)
          .getTimestamp(round)
        assert.isAbove(currentTimestamp.toNumber(), 0)
      })
    })
  })

  describe('#latestAnswer', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .addOracles(
          [personas.Neil.address],
          [personas.Neil.address],
          minAns,
          maxAns,
          rrDelay,
        )
      await aggregator.connect(personas.Neil).submit(nextRound, answer)
    })

    describe('when the reader is not whitelisted', () => {
      it('does not allow latestAnswer to be called', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Eddy).latestAnswer(),
          'Not whitelisted',
        )
      })
    })

    describe('when the reader is whitelisted', () => {
      beforeEach(async () => {
        await aggregator
          .connect(personas.Carol)
          .addToWhitelist(personas.Eddy.address)
      })

      it('allows latestAnswer to be called', async () => {
        const answer = await aggregator.connect(personas.Eddy).latestAnswer()
        matchers.bigNum(h.bigNum(answer), answer)
      })
    })
  })

  describe('#latestTimestamp', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .addOracles(
          [personas.Neil.address],
          [personas.Neil.address],
          minAns,
          maxAns,
          rrDelay,
        )
      await aggregator.connect(personas.Neil).submit(nextRound, answer)
    })

    describe('when the reader is not whitelisted', () => {
      it('does not allow latestTimestamp to be called', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Eddy).latestTimestamp(),
          'Not whitelisted',
        )
      })
    })

    describe('when the reader is whitelisted', () => {
      beforeEach(async () => {
        await aggregator
          .connect(personas.Carol)
          .addToWhitelist(personas.Eddy.address)
      })

      it('allows latestTimestamp to be called', async () => {
        const currentTimestamp = await aggregator
          .connect(personas.Eddy)
          .latestTimestamp()
        assert.isAbove(currentTimestamp.toNumber(), 0)
      })
    })
  })
})
