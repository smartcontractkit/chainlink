import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { randomBytes } from 'crypto'
import { ethers } from 'ethers'
import { PrepaidAggregatorFactory } from '../src/generated'

let personas: setup.Personas
const provider = setup.provider()
const linkTokenFactory = new contract.LinkTokenFactory()
const prepaidAggregatorFactory = new PrepaidAggregatorFactory()

beforeAll(async () => {
  personas = await setup.users(provider).then(x => x.personas)
})

describe('PrepaidAggregator', () => {
  const paymentAmount = h.toWei('3')
  const deposit = h.toWei('100')
  const answer = 100
  const minAns = 1
  const maxAns = 1
  const rrDelay = 0
  const timeout = 1800
  const decimals = 18
  const description = 'LINK/USD'

  let aggregator: contract.Instance<PrepaidAggregatorFactory>
  let link: contract.Instance<contract.LinkTokenFactory>
  let nextRound: number
  let oracleAddresses: ethers.Wallet[]

  async function updateFutureRounds(
    aggregator: contract.Instance<PrepaidAggregatorFactory>,
    overrides: {
      minAnswers?: ethers.utils.BigNumberish
      maxAnswers?: ethers.utils.BigNumberish
      payment?: ethers.utils.BigNumberish
      restartDelay?: ethers.utils.BigNumberish
      timeout?: ethers.utils.BigNumberish
    } = {},
  ) {
    overrides = overrides || {}
    const round = {
      payment: overrides.payment || paymentAmount,
      minAnswers: overrides.minAnswers || minAns,
      maxAnswers: overrides.maxAnswers || maxAns,
      restartDelay: overrides.restartDelay || rrDelay,
      timeout: overrides.timeout || timeout,
    }

    return aggregator.updateFutureRounds(
      round.payment,
      round.minAnswers,
      round.maxAnswers,
      round.restartDelay,
      round.timeout,
    )
  }
  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(personas.Default).deploy()
    aggregator = await prepaidAggregatorFactory
      .connect(personas.Carol)
      .deploy(
        link.address,
        paymentAmount,
        timeout,
        decimals,
        ethers.utils.formatBytes32String(description),
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
    matchers.publicAbi(prepaidAggregatorFactory, [
      'addOracle',
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
      'removeOracle',
      'reportingRound',
      'restartDelay',
      'timeout',
      'updateAnswer',
      'updateAvailableFunds',
      'updateFutureRounds',
      'withdraw',
      'withdrawFunds',
      'withdrawable',
      // Owned methods:
      'acceptOwnership',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#constructor', () => {
    it('sets the paymentAmount', async () => {
      matchers.bigNum(
        ethers.utils.bigNumberify(paymentAmount),
        await aggregator.paymentAmount(),
      )
    })

    it('sets the timeout', async () => {
      matchers.bigNum(
        ethers.utils.bigNumberify(timeout),
        await aggregator.timeout(),
      )
    })

    it('sets the decimals', async () => {
      matchers.bigNum(
        ethers.utils.bigNumberify(decimals),
        await aggregator.decimals(),
      )
    })

    it('sets the description', async () => {
      assert.equal(
        ethers.utils.formatBytes32String(description),
        await aggregator.description(),
      )
    })
  })

  describe('#updateAnswer', () => {
    let minMax

    beforeEach(async () => {
      oracleAddresses = [personas.Neil, personas.Ned, personas.Nelly]
      for (let i = 0; i < oracleAddresses.length; i++) {
        minMax = i + 1
        await aggregator
          .connect(personas.Carol)
          .addOracle(oracleAddresses[i].address, minMax, minMax, rrDelay)
      }
    })

    it('updates the allocated and available funds counters', async () => {
      matchers.bigNum(0, await aggregator.allocatedFunds())

      const tx = await aggregator
        .connect(personas.Neil)
        .updateAnswer(nextRound, answer)
      const receipt = await tx.wait()

      matchers.bigNum(paymentAmount, await aggregator.allocatedFunds())
      const expectedAvailable = deposit.sub(paymentAmount)
      matchers.bigNum(expectedAvailable, await aggregator.availableFunds())
      const logged = ethers.utils.bigNumberify(
        receipt.logs?.[1].topics[1] ?? ethers.utils.bigNumberify(-1),
      )
      matchers.bigNum(expectedAvailable, logged)
    })

    it('updates the latest submission record for the oracle', async () => {
      let latest = await aggregator.latestSubmission(personas.Neil.address)
      assert.equal(0, latest[0].toNumber())
      assert.equal(0, latest[1].toNumber())

      const newAnswer = 427
      await aggregator.connect(personas.Neil).updateAnswer(nextRound, newAnswer)

      latest = await aggregator.latestSubmission(personas.Neil.address)
      assert.equal(newAnswer, latest[0].toNumber())
      assert.equal(nextRound, latest[1].toNumber())
    })

    describe('when the minimum oracles have not reported', () => {
      it('pays the oracles that have reported', async () => {
        matchers.bigNum(
          0,
          await aggregator.connect(personas.Neil).withdrawable(),
        )

        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)

        matchers.bigNum(
          paymentAmount,
          await aggregator.connect(personas.Neil).withdrawable(),
        )
        matchers.bigNum(
          0,
          await aggregator.connect(personas.Ned).withdrawable(),
        )
        matchers.bigNum(
          0,
          await aggregator.connect(personas.Nelly).withdrawable(),
        )
      })

      it('does not update the answer', async () => {
        matchers.bigNum(ethers.constants.Zero, await aggregator.latestAnswer())

        // Not updated because of changes by the owner setting minAnswerCount to 3
        await aggregator.connect(personas.Ned).updateAnswer(nextRound, answer)
        await aggregator.connect(personas.Nelly).updateAnswer(nextRound, answer)

        matchers.bigNum(ethers.constants.Zero, await aggregator.latestAnswer())
      })
    })

    describe('when an oracle prematurely bumps the round', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator, { minAnswers: 2, maxAnswers: 3 })
        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)
      })

      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator.updateAnswer(nextRound + 1, answer),
          'Not eligible to bump round',
        )
      })
    })

    describe('when the minimum number of oracles have reported', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator, { minAnswers: 2, maxAnswers: 3 })
        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)
      })

      it('updates the answer with the median', async () => {
        matchers.bigNum(0, await aggregator.latestAnswer())

        await aggregator.connect(personas.Ned).updateAnswer(nextRound, 99)
        matchers.bigNum(99, await aggregator.latestAnswer()) // ((100+99) / 2).to_i

        await aggregator.connect(personas.Nelly).updateAnswer(nextRound, 101)

        matchers.bigNum(100, await aggregator.latestAnswer())
      })

      it('updates the updated timestamp', async () => {
        const originalTimestamp = await aggregator.latestTimestamp()
        assert.equal(0, originalTimestamp.toNumber())

        await aggregator.connect(personas.Nelly).updateAnswer(nextRound, answer)

        const currentTimestamp = await aggregator.latestTimestamp()
        assert.isAbove(
          currentTimestamp.toNumber(),
          originalTimestamp.toNumber(),
        )
      })

      it('announces the new answer with a log event', async () => {
        const tx = await aggregator
          .connect(personas.Nelly)
          .updateAnswer(nextRound, answer)
        const receipt = await tx.wait()

        const newAnswer = ethers.utils.bigNumberify(
          receipt.logs?.[0].topics[1] ?? ethers.constants.Zero,
        )

        assert.equal(answer, newAnswer.toNumber())
      })

      it('does not set the timedout flag', async () => {
        assert.isFalse(await aggregator.getTimedOutStatus(nextRound))

        await aggregator.connect(personas.Nelly).updateAnswer(nextRound, answer)

        assert.equal(
          nextRound,
          (await aggregator.getOriginatingRoundOfAnswer(nextRound)).toNumber(),
        )
      })
    })

    describe('when an oracle submits for a round twice', () => {
      it('reverts', async () => {
        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)

        await matchers.evmRevert(
          aggregator.connect(personas.Neil).updateAnswer(nextRound, answer),
          'Cannot update round reports',
        )
      })
    })

    describe('when updated after the max answers submitted', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator)
        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)
      })

      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Ned).updateAnswer(nextRound, answer),
          'Round not currently eligible for reporting',
        )
      })
    })

    describe('when a new highest round number is passed in', () => {
      it('increments the answer round', async () => {
        matchers.bigNum(
          ethers.constants.Zero,
          await aggregator.reportingRound(),
        )

        for (const oracle of oracleAddresses) {
          await aggregator.connect(oracle).updateAnswer(nextRound, answer)
        }

        matchers.bigNum(ethers.constants.One, await aggregator.reportingRound())
      })

      it('announces a new round by emitting a log', async () => {
        const tx = await aggregator
          .connect(personas.Neil)
          .updateAnswer(nextRound, answer)
        const receipt = await tx.wait()

        const topics = receipt.logs?.[0].topics ?? []
        const roundNumber = ethers.utils.bigNumberify(topics[1])
        const startedBy = h.evmWordToAddress(topics[2])

        matchers.bigNum(nextRound, roundNumber.toNumber())
        matchers.bigNum(startedBy, personas.Neil.address)
      })
    })

    describe('when a round is passed in higher than expected', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Neil).updateAnswer(nextRound + 1, answer),
          'Must report on current round',
        )
      })
    })

    describe('when called by a non-oracle', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Carol).updateAnswer(nextRound, answer),
          'Only updatable by whitelisted oracles',
        )
      })
    })

    describe('when there are not sufficient available funds', () => {
      beforeEach(async () => {
        await aggregator
          .connect(personas.Carol)
          .withdrawFunds(personas.Carol.address, deposit)
      })

      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Neil).updateAnswer(nextRound, answer),
          'SafeMath: subtraction overflow',
        )
      })
    })

    describe('when price is updated mid-round', () => {
      const newAmount = h.toWei('50')

      it('pays the same amount to all oracles per round', async () => {
        matchers.bigNum(
          0,
          await aggregator.connect(personas.Neil).withdrawable(),
        )
        matchers.bigNum(
          0,
          await aggregator.connect(personas.Nelly).withdrawable(),
        )

        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)

        await updateFutureRounds(aggregator, { payment: newAmount })

        await aggregator.connect(personas.Nelly).updateAnswer(nextRound, answer)

        matchers.bigNum(
          paymentAmount,
          await aggregator.connect(personas.Neil).withdrawable(),
        )
        matchers.bigNum(
          paymentAmount,
          await aggregator.connect(personas.Nelly).withdrawable(),
        )
      })
    })

    describe('when delay is on', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator, {
          minAnswers: oracleAddresses.length,
          maxAnswers: oracleAddresses.length,
          restartDelay: 1,
        })
      })

      it("does not revert on the oracle's first round", async () => {
        // Since lastUpdatedRound defaults to zero and that's the only
        // indication that an oracle hasn't responded, this test guards against
        // the situation where we don't check that and no one can start a round.

        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)
      })

      it('does revert before the delay', async () => {
        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)

        nextRound++

        await matchers.evmRevert(
          aggregator.connect(personas.Neil).updateAnswer(nextRound, answer),
          'Not eligible to bump round',
        )
      })
    })

    describe('when an oracle starts a round before the restart delay is over', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator.connect(personas.Carol))

        oracleAddresses = [personas.Neil, personas.Ned, personas.Nelly]
        for (let i = 0; i < oracleAddresses.length; i++) {
          await aggregator
            .connect(oracleAddresses[i])
            .updateAnswer(nextRound, answer)
          nextRound++
        }

        const newDelay = 2
        // Since Ned and Nelly have answered recently, and we set the delay
        // to 2, only Nelly can answer as she is the only oracle that hasn't
        // started the last two rounds.
        await updateFutureRounds(aggregator, {
          maxAnswers: oracleAddresses.length,
          restartDelay: newDelay,
        })
      })

      describe('when called by an oracle who has not answered recently', () => {
        it('does not revert', async () => {
          await aggregator
            .connect(personas.Neil)
            .updateAnswer(nextRound, answer)
        })
      })

      describe('when called by an oracle who answered recently', () => {
        it('reverts', async () => {
          await matchers.evmRevert(
            aggregator.connect(personas.Ned).updateAnswer(nextRound, answer),
            'Round not currently eligible for reporting',
          )

          await matchers.evmRevert(
            aggregator.connect(personas.Nelly).updateAnswer(nextRound, answer),
            'Round not currently eligible for reporting',
          )
        })
      })
    })

    describe('when the price is not updated for a round', () => {
      // For a round to timeout, it needs a previous round to pull an answer
      // from, so the second round is the earliest round that can timeout,
      // pulling its answer from the first. The start of the third round is
      // the trigger that timesout the second round, so the start of the
      // third round is the earliest we can test a timeout.

      describe('on the third round or later', () => {
        beforeEach(async () => {
          await updateFutureRounds(aggregator, {
            minAnswers: oracleAddresses.length,
            maxAnswers: oracleAddresses.length,
            restartDelay: 1,
          })

          for (const oracle of oracleAddresses) {
            await aggregator.connect(oracle).updateAnswer(nextRound, answer)
          }
          nextRound++

          await aggregator.connect(personas.Ned).updateAnswer(nextRound, answer)
          await aggregator
            .connect(personas.Nelly)
            .updateAnswer(nextRound, answer)
          assert.equal(
            nextRound,
            (await aggregator.reportingRound()).toNumber(),
          )

          await h.increaseTimeBy(timeout + 1, provider)
          nextRound++
        })

        it('allows a new round to be started', async () => {
          await aggregator
            .connect(personas.Nelly)
            .updateAnswer(nextRound, answer)
        })

        it('sets the info for the previous round', async () => {
          const previousRound = nextRound - 1
          let updated = await aggregator.getTimestamp(previousRound)
          let ans = await aggregator.getAnswer(previousRound)
          assert.equal(0, updated.toNumber())
          assert.equal(0, ans.toNumber())

          const tx = await aggregator
            .connect(personas.Nelly)
            .updateAnswer(nextRound, answer)
          const receipt = await tx.wait()

          const block = await provider.getBlock(receipt.blockHash ?? '')

          updated = await aggregator.getTimestamp(previousRound)
          ans = await aggregator.getAnswer(previousRound)
          matchers.bigNum(ethers.utils.bigNumberify(block.timestamp), updated)
          assert.equal(answer, ans.toNumber())
        })

        it('sets the previous round as timed out', async () => {
          const previousRound = nextRound - 1
          assert.isFalse(await aggregator.getTimedOutStatus(previousRound))

          await aggregator
            .connect(personas.Nelly)
            .updateAnswer(nextRound, answer)

          assert.isTrue(await aggregator.getTimedOutStatus(previousRound))
          assert.equal(
            previousRound - 1,
            (
              await aggregator.getOriginatingRoundOfAnswer(previousRound)
            ).toNumber(),
          )
        })

        it('still respects the delay restriction', async () => {
          // expected to revert because the sender started the last round
          await matchers.evmRevert(
            aggregator.connect(personas.Ned).updateAnswer(nextRound, answer),
            'Round not currently eligible for reporting',
          )
        })

        it('uses the timeout set at the beginning of the round', async () => {
          await updateFutureRounds(aggregator, {
            timeout: timeout + 100000,
          })

          await aggregator
            .connect(personas.Nelly)
            .updateAnswer(nextRound, answer)
        })
      })

      describe('earlier than the third round', () => {
        beforeEach(async () => {
          await aggregator
            .connect(personas.Neil)
            .updateAnswer(nextRound, answer)
          await aggregator
            .connect(personas.Nelly)
            .updateAnswer(nextRound, answer)
          assert.equal(
            nextRound,
            (await aggregator.reportingRound()).toNumber(),
          )

          await h.increaseTimeBy(timeout + 1, provider)

          nextRound++
          assert.equal(2, nextRound)
        })

        it('does not allow a round to be started', async () => {
          await matchers.evmRevert(
            aggregator.connect(personas.Nelly).updateAnswer(nextRound, answer),
            'Must have a previous answer to pull from',
          )
        })
      })
    })
  })

  describe('#getAnswer', () => {
    const answers = [1, 10, 101, 1010, 10101, 101010, 1010101]

    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .addOracle(personas.Neil.address, minAns, maxAns, rrDelay)

      for (const answer of answers) {
        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)
        nextRound++
      }
    })

    it('retrieves the answer recorded for past rounds', async () => {
      for (let i = nextRound; i < nextRound; i++) {
        const answer = await aggregator.getAnswer(i)
        matchers.bigNum(ethers.utils.bigNumberify(answers[i - 1]), answer)
      }
    })
  })

  describe('#getTimestamp', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .addOracle(personas.Neil.address, minAns, maxAns, rrDelay)

      for (let i = 0; i < 10; i++) {
        await aggregator.connect(personas.Neil).updateAnswer(nextRound, i)
        nextRound++
      }
    })

    it('retrieves the answer recorded for past rounds', async () => {
      let lastTimestamp = ethers.constants.Zero

      for (let i = 1; i < nextRound; i++) {
        const currentTimestamp = await aggregator.getTimestamp(i)
        assert.isAtLeast(currentTimestamp.toNumber(), lastTimestamp.toNumber())
        lastTimestamp = currentTimestamp
      }
    })
  })

  describe('#addOracle', () => {
    it('increases the oracle count', async () => {
      const pastCount = await aggregator.oracleCount()
      await aggregator
        .connect(personas.Carol)
        .addOracle(personas.Neil.address, minAns, maxAns, rrDelay)
      const currentCount = await aggregator.oracleCount()

      matchers.bigNum(currentCount, pastCount + 1)
    })

    it('updates the round details', async () => {
      await aggregator
        .connect(personas.Carol)
        .addOracle(personas.Neil.address, 0, 1, 0)

      matchers.bigNum(ethers.constants.Zero, await aggregator.minAnswerCount())
      matchers.bigNum(
        ethers.utils.bigNumberify(1),
        await aggregator.maxAnswerCount(),
      )
      matchers.bigNum(ethers.constants.Zero, await aggregator.restartDelay())
    })

    it('emits a log', async () => {
      const tx = await aggregator
        .connect(personas.Carol)
        .addOracle(personas.Neil.address, minAns, maxAns, rrDelay)
      const receipt = await tx.wait()

      const added = h.evmWordToAddress(receipt.logs?.[0].topics[1])
      matchers.bigNum(added, personas.Neil.address)
    })

    describe('when the oracle has already been added', () => {
      beforeEach(async () => {
        await aggregator
          .connect(personas.Carol)
          .addOracle(personas.Neil.address, minAns, maxAns, rrDelay)
      })

      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Carol)
            .addOracle(personas.Neil.address, minAns, maxAns, rrDelay),
          'Address is already recorded as an oracle',
        )
      })
    })

    describe('when called by anyone but the owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Neil)
            .addOracle(personas.Neil.address, minAns, maxAns, rrDelay),
          'Only callable by owner',
        )
      })
    })

    describe('when an oracle gets added mid-round', () => {
      beforeEach(async () => {
        oracleAddresses = [personas.Neil, personas.Ned]
        for (let i = 0; i < oracleAddresses.length; i++) {
          await aggregator
            .connect(personas.Carol)
            .addOracle(oracleAddresses[i].address, i + 1, i + 1, rrDelay)
        }

        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)

        await aggregator
          .connect(personas.Carol)
          .addOracle(personas.Nelly.address, 3, 3, rrDelay)
      })

      it('does not allow the oracle to update the round', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Nelly).updateAnswer(nextRound, answer),
          'New oracles cannot participate in in-progress rounds',
        )
      })

      it('does allow the oracle to update future rounds', async () => {
        // complete round
        await aggregator.connect(personas.Ned).updateAnswer(nextRound, answer)

        // now can participate in new rounds
        await aggregator
          .connect(personas.Nelly)
          .updateAnswer(nextRound + 1, answer)
      })
    })

    describe('when an oracle is added after removed for a round', () => {
      it('allows the oracle to update', async () => {
        oracleAddresses = [personas.Neil, personas.Nelly]
        for (let i = 0; i < oracleAddresses.length; i++) {
          await aggregator
            .connect(personas.Carol)
            .addOracle(oracleAddresses[i].address, i + 1, i + 1, rrDelay)
        }

        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)
        await aggregator.connect(personas.Nelly).updateAnswer(nextRound, answer)
        nextRound++

        await aggregator
          .connect(personas.Carol)
          .removeOracle(personas.Nelly.address, 1, 1, rrDelay)

        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)
        nextRound++

        await aggregator
          .connect(personas.Carol)
          .addOracle(personas.Nelly.address, 1, 1, rrDelay)

        await aggregator.connect(personas.Nelly).updateAnswer(nextRound, answer)
      })
    })

    describe('when an oracle is added and immediately removed mid-round', () => {
      it('allows the oracle to update', async () => {
        oracleAddresses = [personas.Neil, personas.Nelly]
        for (let i = 0; i < oracleAddresses.length; i++) {
          await aggregator
            .connect(personas.Carol)
            .addOracle(oracleAddresses[i].address, i + 1, i + 1, rrDelay)
        }

        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)

        await aggregator
          .connect(personas.Carol)
          .removeOracle(personas.Nelly.address, 1, 1, rrDelay)
        await aggregator
          .connect(personas.Carol)
          .addOracle(personas.Nelly.address, 1, 1, rrDelay)

        await aggregator.connect(personas.Nelly).updateAnswer(nextRound, answer)
      })
    })

    const limit = 42
    describe(`when adding more than ${limit} oracles`, () => {
      it('reverts', async () => {
        for (let i = 0; i < limit; i++) {
          const minMax = i + 1
          const fakeAddress = h.addHexPrefix(randomBytes(20).toString('hex'))

          await aggregator
            .connect(personas.Carol)
            .addOracle(fakeAddress, minMax, minMax, rrDelay)
        }
        await matchers.evmRevert(
          aggregator
            .connect(personas.Carol)
            .addOracle(personas.Neil.address, limit + 1, limit + 1, rrDelay),
          `cannot add more than ${limit} oracles`,
        )
      })
    })
  })

  describe('#removeOracle', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .addOracle(personas.Neil.address, minAns, maxAns, rrDelay)
      await aggregator
        .connect(personas.Carol)
        .addOracle(personas.Nelly.address, 2, 2, rrDelay, {})
    })

    it('decreases the oracle count', async () => {
      const pastCount = await aggregator.oracleCount()
      await aggregator
        .connect(personas.Carol)
        .removeOracle(personas.Neil.address, minAns, maxAns, rrDelay)
      const currentCount = await aggregator.oracleCount()

      expect(currentCount).toEqual(pastCount - 1)
    })

    it('updates the round details', async () => {
      await aggregator
        .connect(personas.Carol)
        .removeOracle(personas.Neil.address, 0, 1, 0)

      matchers.bigNum(ethers.constants.Zero, await aggregator.minAnswerCount())
      matchers.bigNum(
        ethers.utils.bigNumberify(1),
        await aggregator.maxAnswerCount(),
      )
      matchers.bigNum(ethers.constants.Zero, await aggregator.restartDelay())
    })

    it('emits a log', async () => {
      const tx = await aggregator
        .connect(personas.Carol)
        .removeOracle(personas.Neil.address, minAns, maxAns, rrDelay)
      const receipt = await tx.wait()

      const added = h.evmWordToAddress(receipt.logs?.[0].topics[1])
      matchers.bigNum(added, personas.Neil.address)
    })

    describe('when the oracle is not currently added', () => {
      beforeEach(async () => {
        await aggregator
          .connect(personas.Carol)
          .removeOracle(personas.Neil.address, minAns, maxAns, rrDelay)
      })

      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Carol)
            .removeOracle(personas.Neil.address, minAns, maxAns, rrDelay),
          'Address is not a whitelisted oracle',
        )
      })
    })

    describe('when removing the last oracle', () => {
      it('does not revert', async () => {
        await aggregator
          .connect(personas.Carol)
          .removeOracle(personas.Neil.address, minAns, maxAns, rrDelay)

        await aggregator
          .connect(personas.Carol)
          .removeOracle(personas.Nelly.address, 0, 0, 0)
      })
    })

    describe('when called by anyone but the owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Ned)
            .removeOracle(personas.Neil.address, 0, 0, rrDelay),
          'Only callable by owner',
        )
      })
    })

    describe('when an oracle gets removed mid-round', () => {
      beforeEach(async () => {
        await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)

        await aggregator
          .connect(personas.Carol)
          .removeOracle(personas.Nelly.address, 1, 1, rrDelay)
      })

      it('is allowed to finish that round', async () => {
        await aggregator.connect(personas.Nelly).updateAnswer(nextRound, answer)
        nextRound++

        // cannot participate in future rounds
        await matchers.evmRevert(
          aggregator.connect(personas.Nelly).updateAnswer(nextRound, answer),
          'Oracle has been removed from whitelist',
        )
      })
    })
  })

  describe('#withdrawFunds', () => {
    describe('when called by the owner', () => {
      it('succeeds', async () => {
        await aggregator
          .connect(personas.Carol)
          .withdrawFunds(personas.Carol.address, deposit)

        matchers.bigNum(0, await aggregator.availableFunds())
        matchers.bigNum(deposit, await link.balanceOf(personas.Carol.address))
      })

      it('does not let withdrawals happen multiple times', async () => {
        await aggregator
          .connect(personas.Carol)
          .withdrawFunds(personas.Carol.address, deposit)

        await matchers.evmRevert(
          aggregator
            .connect(personas.Carol)
            .withdrawFunds(personas.Carol.address, deposit),
          'Insufficient funds',
        )
      })

      describe('with a number higher than the available LINK balance', () => {
        beforeEach(async () => {
          await aggregator
            .connect(personas.Carol)
            .addOracle(personas.Neil.address, minAns, maxAns, rrDelay)
          await aggregator
            .connect(personas.Neil)
            .updateAnswer(nextRound, answer)
        })

        it('fails', async () => {
          await matchers.evmRevert(
            aggregator
              .connect(personas.Carol)
              .withdrawFunds(personas.Carol.address, deposit),
            'Insufficient funds',
          )

          matchers.bigNum(
            deposit.sub(paymentAmount),
            await aggregator.availableFunds(),
          )
        })
      })
    })

    describe('when called by a non-owner', () => {
      it('fails', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Eddy)
            .withdrawFunds(personas.Carol.address, deposit),
          'Only callable by owner',
        )

        matchers.bigNum(deposit, await aggregator.availableFunds())
      })
    })
  })

  describe('#updateFutureRounds', () => {
    let minAnswerCount, maxAnswerCount
    const newPaymentAmount = h.toWei('2')
    const newMin = 1
    const newMax = 3
    const newDelay = 2

    beforeEach(async () => {
      oracleAddresses = [personas.Neil, personas.Ned, personas.Nelly]
      for (let i = 0; i < oracleAddresses.length; i++) {
        const minMax = i + 1
        await aggregator
          .connect(personas.Carol)
          .addOracle(oracleAddresses[i].address, minMax, minMax, rrDelay)
      }
      minAnswerCount = oracleAddresses.length
      maxAnswerCount = oracleAddresses.length

      matchers.bigNum(paymentAmount, await aggregator.paymentAmount())
      assert.equal(minAnswerCount, await aggregator.minAnswerCount())
      assert.equal(maxAnswerCount, await aggregator.maxAnswerCount())
    })

    it('updates the min and max answer counts', async () => {
      await updateFutureRounds(aggregator, {
        payment: newPaymentAmount,
        minAnswers: newMin,
        maxAnswers: newMax,
        restartDelay: newDelay,
      })

      matchers.bigNum(newPaymentAmount, await aggregator.paymentAmount())
      matchers.bigNum(
        ethers.utils.bigNumberify(newMin),
        await aggregator.minAnswerCount(),
      )
      matchers.bigNum(
        ethers.utils.bigNumberify(newMax),
        await aggregator.maxAnswerCount(),
      )
      matchers.bigNum(
        ethers.utils.bigNumberify(newDelay),
        await aggregator.restartDelay(),
      )
    })

    it('emits a log announcing the new round details', async () => {
      const tx = await updateFutureRounds(aggregator, {
        payment: newPaymentAmount,
        minAnswers: newMin,
        maxAnswers: newMax,
        restartDelay: newDelay,
        timeout: timeout + 1,
      })
      const receipt = await tx.wait()
      const round = h.eventArgs(receipt.events?.[0])

      matchers.bigNum(newPaymentAmount, round.paymentAmount)
      assert.equal(newMin, round.minAnswerCount)
      assert.equal(newMax, round.maxAnswerCount)
      assert.equal(newDelay, round.restartDelay)
      assert.equal(timeout + 1, round.timeout)
    })

    describe('when it is set to higher than the number or oracles', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          updateFutureRounds(aggregator, {
            maxAnswers: 4,
          }),
          'Cannot have the answer max higher oracle count',
        )
      })
    })

    describe('when it sets the min higher than the max', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          updateFutureRounds(aggregator, {
            minAnswers: 3,
            maxAnswers: 2,
          }),
          'Cannot have the answer minimum higher the max',
        )
      })
    })

    describe('when delay equal or greater the oracle count', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          updateFutureRounds(aggregator, {
            restartDelay: 3,
          }),
          'Restart delay must be less than oracle count',
        )
      })
    })

    describe('when called by anyone but the owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          updateFutureRounds(aggregator.connect(personas.Ned)),
          'Only callable by owner',
        )
      })
    })
  })

  describe('#updateAvailableFunds', () => {
    it('checks the LINK token to see if any additional funds are available', async () => {
      const originalBalance = await aggregator.availableFunds()

      await aggregator.updateAvailableFunds()

      matchers.bigNum(originalBalance, await aggregator.availableFunds())

      await link.transfer(aggregator.address, deposit)
      await aggregator.updateAvailableFunds()

      const newBalance = await aggregator.availableFunds()
      matchers.bigNum(originalBalance.add(deposit), newBalance)
    })

    it('removes allocated funds from the available balance', async () => {
      const originalBalance = await aggregator.availableFunds()

      await aggregator
        .connect(personas.Carol)
        .addOracle(personas.Neil.address, minAns, maxAns, rrDelay)
      await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)
      await link.transfer(aggregator.address, deposit)
      await aggregator.updateAvailableFunds()

      const expected = originalBalance.add(deposit).sub(paymentAmount)
      const newBalance = await aggregator.availableFunds()
      matchers.bigNum(expected, newBalance)
    })

    it('emits a log', async () => {
      await link.transfer(aggregator.address, deposit)

      const tx = await aggregator.updateAvailableFunds()
      const receipt = await tx.wait()

      const reportedBalance = ethers.utils.bigNumberify(
        receipt.logs?.[0].topics[1] ?? -1,
      )
      matchers.bigNum(await aggregator.availableFunds(), reportedBalance)
    })

    describe('when the available funds have not changed', () => {
      it('does not emit a log', async () => {
        const tx = await aggregator.updateAvailableFunds()
        const receipt = await tx.wait()

        assert.equal(0, receipt.logs?.length)
      })
    })
  })

  describe('#withdraw', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .addOracle(personas.Neil.address, minAns, maxAns, rrDelay)
      await aggregator.connect(personas.Neil).updateAnswer(nextRound, answer)
    })

    it('transfers LINK to the caller', async () => {
      const originalBalance = await link.balanceOf(aggregator.address)
      matchers.bigNum(0, await link.balanceOf(personas.Neil.address))

      await aggregator
        .connect(personas.Neil)
        .withdraw(personas.Neil.address, paymentAmount)

      matchers.bigNum(
        originalBalance.sub(paymentAmount),
        await link.balanceOf(aggregator.address),
      )
      matchers.bigNum(
        paymentAmount,
        await link.balanceOf(personas.Neil.address),
      )
    })

    it('decrements the allocated funds counter', async () => {
      const originalAllocation = await aggregator.allocatedFunds()

      await aggregator
        .connect(personas.Neil)
        .withdraw(personas.Neil.address, paymentAmount)

      matchers.bigNum(
        originalAllocation.sub(paymentAmount),
        await aggregator.allocatedFunds(),
      )
    })

    describe('when the caller withdraws more than they have', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Neil)
            .withdraw(
              personas.Neil.address,
              paymentAmount.add(ethers.utils.bigNumberify(1)),
            ),
          'Insufficient balance',
        )
      })
    })
  })

  describe('#onTokenTransfer', () => {
    it('updates the available balance', async () => {
      const originalBalance = await aggregator.availableFunds()

      await aggregator.updateAvailableFunds()

      matchers.bigNum(originalBalance, await aggregator.availableFunds())

      await link.transferAndCall(aggregator.address, deposit, '0x', {
        value: 0,
      })

      const newBalance = await aggregator.availableFunds()
      matchers.bigNum(originalBalance.add(deposit), newBalance)
    })
  })
})
