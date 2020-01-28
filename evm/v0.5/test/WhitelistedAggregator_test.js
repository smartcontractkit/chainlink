import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'

import { expectRevert } from 'openzeppelin-test-helpers'

contract('WhitelistedAggregator', () => {
  const Aggregator = artifacts.require('WhitelistedAggregator.sol')
  const personas = h.personas
  const paymentAmount = h.toWei('3')
  const deposit = h.toWei('100')
  const answer = 100
  const minAns = 1
  const maxAns = 1
  const rrDelay = 0
  const timeout = 1800
  const decimals = 18
  const description = 'LINK/USD'

  let aggregator, link, nextRound

  beforeEach(async () => {
    link = await h.linkContract(personas.defaultAccount)
    aggregator = await Aggregator.new(
      link.address,
      paymentAmount,
      timeout,
      decimals,
      h.toHex(description),
      {
        from: personas.Carol,
      },
    )
    await link.transfer(aggregator.address, deposit)
    await aggregator.updateAvailableFunds()
    assertBigNum(deposit, await link.balanceOf.call(aggregator.address))
    nextRound = 1
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(Aggregator, [
      'addOracle',
      'addToWhitelist',
      'allocatedFunds',
      'availableFunds',
      'description',
      'getAnswer',
      'getOriginatingRoundOfAnswer',
      'getTimedOutStatus',
      'getTimestamp',
      'latestAnswer',
      'latestRound',
      'latestSubmission',
      'latestTimestamp',
      'maxAnswerCount',
      'minAnswerCount',
      'onTokenTransfer',
      'oracleCount',
      'paymentAmount',
      'decimals',
      'removeFromWhitelist',
      'removeOracle',
      'reportingRound',
      'restartDelay',
      'timeout',
      'updateAnswer',
      'updateAvailableFunds',
      'updateFutureRounds',
      'whitelisted',
      'withdraw',
      'withdrawFunds',
      'withdrawable',
      // Owned methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#constructor', async () => {
    it('sets the paymentAmount', async () => {
      assertBigNum(
        h.bigNum(paymentAmount),
        await aggregator.paymentAmount.call(),
      )
    })

    it('sets the timeout', async () => {
      assertBigNum(h.bigNum(timeout), await aggregator.timeout.call())
    })

    it('sets the decimals', async () => {
      assertBigNum(h.bigNum(decimals), await aggregator.decimals.call())
    })

    it('sets the description', async () => {
      assert.equal(
        description,
        web3.utils.toUtf8(await aggregator.description.call()),
      )
    })
  })

  describe('#getAnswer', async () => {
    beforeEach(async () => {
      await aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
        from: personas.Carol,
      })
      await aggregator.updateAnswer(nextRound, answer, {
        from: personas.Neil,
      })
    })

    context('when the reader is not whitelisted', () => {
      it('does not allow getAnswer to be called', async () => {
        const round = await aggregator.latestRound.call()
        await expectRevert(
          aggregator.getAnswer.call(round, {
            from: personas.Eddy,
          }),
          'Not whitelisted',
        )
      })
    })

    context('when the reader is whitelisted', () => {
      beforeEach(async () => {
        await aggregator.addToWhitelist(personas.Eddy, { from: personas.Carol })
      })

      it('allows getAnswer to be called', async () => {
        const round = await aggregator.latestRound.call()
        const answer = await aggregator.getAnswer.call(round, {
          from: personas.Eddy,
        })
        assertBigNum(h.bigNum(answer), answer)
      })
    })
  })

  describe('#getTimestamp', async () => {
    beforeEach(async () => {
      await aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
        from: personas.Carol,
      })
      await aggregator.updateAnswer(nextRound, answer, {
        from: personas.Neil,
      })
    })

    context('when the reader is not whitelisted', () => {
      it('does not allow getTimestamp to be called', async () => {
        const round = await aggregator.latestRound.call()
        await expectRevert(
          aggregator.getTimestamp.call(round, {
            from: personas.Eddy,
          }),
          'Not whitelisted',
        )
      })
    })

    context('when the reader is whitelisted', () => {
      beforeEach(async () => {
        await aggregator.addToWhitelist(personas.Eddy, { from: personas.Carol })
      })

      it('allows getTimestamp to be called', async () => {
        const round = await aggregator.latestRound.call()
        const currentTimestamp = await aggregator.getTimestamp.call(round, {
          from: personas.Eddy,
        })
        assert.isAbove(currentTimestamp.toNumber(), 0)
      })
    })
  })

  describe('#latestAnswer', async () => {
    beforeEach(async () => {
      await aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
        from: personas.Carol,
      })
      await aggregator.updateAnswer(nextRound, answer, {
        from: personas.Neil,
      })
    })

    context('when the reader is not whitelisted', () => {
      it('does not allow latestAnswer to be called', async () => {
        await expectRevert(
          aggregator.latestAnswer.call({
            from: personas.Eddy,
          }),
          'Not whitelisted',
        )
      })
    })

    context('when the reader is whitelisted', () => {
      beforeEach(async () => {
        await aggregator.addToWhitelist(personas.Eddy, { from: personas.Carol })
      })

      it('allows latestAnswer to be called', async () => {
        const answer = await aggregator.latestAnswer.call({
          from: personas.Eddy,
        })
        assertBigNum(h.bigNum(answer), answer)
      })
    })
  })

  describe('#latestTimestamp', async () => {
    beforeEach(async () => {
      await aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
        from: personas.Carol,
      })
      await aggregator.updateAnswer(nextRound, answer, {
        from: personas.Neil,
      })
    })

    context('when the reader is not whitelisted', () => {
      it('does not allow latestTimestamp to be called', async () => {
        await expectRevert(
          aggregator.latestTimestamp.call({
            from: personas.Eddy,
          }),
          'Not whitelisted',
        )
      })
    })

    context('when the reader is whitelisted', () => {
      beforeEach(async () => {
        await aggregator.addToWhitelist(personas.Eddy, { from: personas.Carol })
      })

      it('allows latestTimestamp to be called', async () => {
        const currentTimestamp = await aggregator.latestTimestamp.call({
          from: personas.Eddy,
        })
        assert.isAbove(currentTimestamp.toNumber(), 0)
      })
    })
  })
})
