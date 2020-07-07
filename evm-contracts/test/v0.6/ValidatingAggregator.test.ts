import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { ValidatingAggregatorFactory } from '../../ethers/v0.6/ValidatingAggregatorFactory'
import { AnswerValidatorTestHelperFactory } from '../../ethers/v0.6/AnswerValidatorTestHelperFactory'
import { GasGuzzlerFactory } from '../../ethers/v0.6/GasGuzzlerFactory'
import { HistoricDeviationValidatorFactory } from '../../ethers/v0.6/HistoricDeviationValidatorFactory'
import { FlagsFactory } from '../../ethers/v0.6/FlagsFactory'
import { SimpleWriteAccessControllerFactory } from '../../ethers/v0.6/SimpleWriteAccessControllerFactory'

let personas: setup.Personas
const provider = setup.provider()
const linkTokenFactory = new contract.LinkTokenFactory()
const fluxAggregatorFactory = new ValidatingAggregatorFactory()
const answerValidatorFactory = new AnswerValidatorTestHelperFactory()
const validatorFactory = new HistoricDeviationValidatorFactory()
const flagsFactory = new FlagsFactory()
const acFactory = new SimpleWriteAccessControllerFactory()
const gasGuzzlerFactory = new GasGuzzlerFactory()
const emptyAddress = '0x0000000000000000000000000000000000000000'

beforeAll(async () => {
  personas = await setup.users(provider).then(x => x.personas)
})

describe('ValidatingAggregator', () => {
  const paymentAmount = h.toWei('3')
  const deposit = h.toWei('100')
  const answer = 100
  const minAns = 1
  const maxAns = 1
  const rrDelay = 0
  const timeout = 1800
  const decimals = 18
  const description = 'LINK/USD'
  const reserveRounds = 2
  const minSubmissionValue = h.bigNum('1')
  const maxSubmissionValue = h.bigNum('100000000000000000000')

  let aggregator: contract.Instance<ValidatingAggregatorFactory>
  let link: contract.Instance<contract.LinkTokenFactory>
  let validator: contract.Instance<AnswerValidatorTestHelperFactory>
  let gasGuzzler: contract.Instance<GasGuzzlerFactory>
  let nextRound: number
  let oracles: ethers.Wallet[]

  async function updateFutureRounds(
    aggregator: contract.Instance<ValidatingAggregatorFactory>,
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

  async function addOracles(
    aggregator: contract.Instance<ValidatingAggregatorFactory>,
    oraclesAndAdmin: ethers.Wallet[],
    minAnswers: number,
    maxAnswers: number,
    restartDelay: number,
  ): Promise<ethers.ContractTransaction> {
    return aggregator.connect(personas.Carol).addOracles(
      oraclesAndAdmin.map(oracle => oracle.address),
      oraclesAndAdmin.map(admin => admin.address),
      minAnswers,
      maxAnswers,
      restartDelay,
    )
  }

  async function advanceRound(
    aggregator: contract.Instance<ValidatingAggregatorFactory>,
    submitters: ethers.Wallet[],
    currentSubmission: number = answer,
  ): Promise<number> {
    for (const submitter of submitters) {
      await aggregator.connect(submitter).submit(nextRound, currentSubmission)
    }
    nextRound++
    return nextRound
  }

  const deployment = setup.snapshot(provider, async () => {
    link = await linkTokenFactory.connect(personas.Default).deploy()
    aggregator = await fluxAggregatorFactory
      .connect(personas.Carol)
      .deploy(
        link.address,
        paymentAmount,
        timeout,
        emptyAddress,
        minSubmissionValue,
        maxSubmissionValue,
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
    matchers.publicAbi(fluxAggregatorFactory, [
      'acceptAdmin',
      'addOracles',
      'allocatedFunds',
      'answerValidator',
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
      'maxSubmissionValue',
      'minSubmissionCount',
      'minSubmissionValue',
      'onTokenTransfer',
      'oracleCount',
      'oracleRoundState',
      'paymentAmount',
      'removeOracles',
      'reportingRound',
      'requestNewRound',
      'restartDelay',
      'setRequesterPermissions',
      'setAnswerValidator',
      'submit',
      'timeout',
      'transferAdmin',
      'updateAvailableFunds',
      'updateFutureRounds',
      'withdrawFunds',
      'withdrawPayment',
      'withdrawablePayment',
      'version',
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

    it('sets the version to 3', async () => {
      matchers.bigNum(3, await aggregator.version())
    })
  })

  describe('#submit', () => {
    let minMax

    beforeEach(async () => {
      oracles = [personas.Neil, personas.Ned, personas.Nelly]
      minMax = oracles.length
      await addOracles(aggregator, oracles, minMax, minMax, rrDelay)
    })

    it('updates the allocated and available funds counters', async () => {
      matchers.bigNum(0, await aggregator.allocatedFunds())

      const tx = await aggregator
        .connect(personas.Neil)
        .submit(nextRound, answer)
      const receipt = await tx.wait()

      matchers.bigNum(paymentAmount, await aggregator.allocatedFunds())
      const expectedAvailable = deposit.sub(paymentAmount)
      matchers.bigNum(expectedAvailable, await aggregator.availableFunds())
      const logged = ethers.utils.bigNumberify(
        receipt.logs?.[2].topics[1] ?? ethers.utils.bigNumberify(-1),
      )
      matchers.bigNum(expectedAvailable, logged)
    })

    it('updates the latest submission record for the oracle', async () => {
      let latest = await aggregator.latestSubmission(personas.Neil.address)
      assert.equal(0, latest[0].toNumber())
      assert.equal(0, latest[1].toNumber())

      const newAnswer = 427
      await aggregator.connect(personas.Neil).submit(nextRound, newAnswer)

      latest = await aggregator.latestSubmission(personas.Neil.address)
      assert.equal(newAnswer, latest[0].toNumber())
      assert.equal(nextRound, latest[1].toNumber())
    })

    it('emits a log event announcing submission details', async () => {
      const tx = await aggregator
        .connect(personas.Nelly)
        .submit(nextRound, answer)
      const receipt = await tx.wait()
      const round = h.eventArgs(receipt.events?.[1])

      assert.equal(answer, round.submission)
      assert.equal(nextRound, round.round)
      assert.equal(personas.Nelly.address, round.oracle)
    })

    describe('when the minimum oracles have not reported', () => {
      it('pays the oracles that have reported', async () => {
        matchers.bigNum(
          0,
          await aggregator
            .connect(personas.Neil)
            .withdrawablePayment(personas.Neil.address),
        )

        await aggregator.connect(personas.Neil).submit(nextRound, answer)

        matchers.bigNum(
          paymentAmount,
          await aggregator
            .connect(personas.Neil)
            .withdrawablePayment(personas.Neil.address),
        )
        matchers.bigNum(
          0,
          await aggregator
            .connect(personas.Ned)
            .withdrawablePayment(personas.Ned.address),
        )
        matchers.bigNum(
          0,
          await aggregator
            .connect(personas.Nelly)
            .withdrawablePayment(personas.Nelly.address),
        )
      })

      it('does not update the answer', async () => {
        matchers.bigNum(ethers.constants.Zero, await aggregator.latestAnswer())

        // Not updated because of changes by the owner setting minSubmissionCount to 3
        await aggregator.connect(personas.Ned).submit(nextRound, answer)
        await aggregator.connect(personas.Nelly).submit(nextRound, answer)

        matchers.bigNum(ethers.constants.Zero, await aggregator.latestAnswer())
      })
    })

    describe('when an oracle prematurely bumps the round', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator, { minAnswers: 2, maxAnswers: 3 })
        await aggregator.connect(personas.Neil).submit(nextRound, answer)
      })

      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Neil).submit(nextRound + 1, answer),
          'previous round not supersedable',
        )
      })
    })

    describe('when the minimum number of oracles have reported', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator, { minAnswers: 2, maxAnswers: 3 })
        await aggregator.connect(personas.Neil).submit(nextRound, answer)
      })

      it('updates the answer with the median', async () => {
        matchers.bigNum(0, await aggregator.latestAnswer())

        await aggregator.connect(personas.Ned).submit(nextRound, 99)
        matchers.bigNum(99, await aggregator.latestAnswer()) // ((100+99) / 2).to_i

        await aggregator.connect(personas.Nelly).submit(nextRound, 101)

        matchers.bigNum(100, await aggregator.latestAnswer())
      })

      it('updates the updated timestamp', async () => {
        const originalTimestamp = await aggregator.latestTimestamp()
        assert.isAbove(originalTimestamp.toNumber(), 0)

        await aggregator.connect(personas.Nelly).submit(nextRound, answer)

        const currentTimestamp = await aggregator.latestTimestamp()
        assert.isAbove(
          currentTimestamp.toNumber(),
          originalTimestamp.toNumber(),
        )
      })

      it('announces the new answer with a log event', async () => {
        const tx = await aggregator
          .connect(personas.Nelly)
          .submit(nextRound, answer)
        const receipt = await tx.wait()

        const newAnswer = ethers.utils.bigNumberify(
          receipt.logs?.[0].topics[1] ?? ethers.constants.Zero,
        )

        assert.equal(answer, newAnswer.toNumber())
      })

      it('does not set the timedout flag', async () => {
        let round = await aggregator.getRoundData(nextRound)
        assert.notEqual(round.roundId, round.answeredInRound)

        await aggregator.connect(personas.Nelly).submit(nextRound, answer)

        round = await aggregator.getRoundData(nextRound)
        assert.equal(nextRound, round.answeredInRound.toNumber())
      })

      it('updates the round details', async () => {
        const roundBefore = await aggregator.getRoundData(nextRound)
        matchers.bigNum(nextRound, roundBefore.roundId)
        matchers.bigNum(0, roundBefore.answer)
        assert.isFalse(roundBefore.startedAt.isZero())
        matchers.bigNum(0, roundBefore.updatedAt)
        matchers.bigNum(0, roundBefore.answeredInRound)

        const roundBeforeLatest = await aggregator.latestRoundData()
        matchers.bigNum(nextRound - 1, roundBeforeLatest.roundId)

        h.increaseTimeBy(15, provider)
        await aggregator.connect(personas.Nelly).submit(nextRound, answer)

        const roundAfter = await aggregator.getRoundData(nextRound)
        matchers.bigNum(nextRound, roundAfter.roundId)
        matchers.bigNum(answer, roundAfter.answer)
        matchers.bigNum(roundBefore.startedAt, roundAfter.startedAt)
        matchers.bigNum(
          await aggregator.getTimestamp(nextRound),
          roundAfter.updatedAt,
        )
        matchers.bigNum(nextRound, roundAfter.answeredInRound)

        assert.isBelow(
          roundAfter.startedAt.toNumber(),
          roundAfter.updatedAt.toNumber(),
        )

        const roundAfterLatest = await aggregator.latestRoundData()
        matchers.bigNum(roundAfter.roundId, roundAfterLatest.roundId)
        matchers.bigNum(roundAfter.answer, roundAfterLatest.answer)
        matchers.bigNum(roundAfter.startedAt, roundAfterLatest.startedAt)
        matchers.bigNum(roundAfter.updatedAt, roundAfterLatest.updatedAt)
        matchers.bigNum(
          roundAfter.answeredInRound,
          roundAfterLatest.answeredInRound,
        )
      })
    })

    describe('when an oracle submits for a round twice', () => {
      it('reverts', async () => {
        await aggregator.connect(personas.Neil).submit(nextRound, answer)

        await matchers.evmRevert(
          aggregator.connect(personas.Neil).submit(nextRound, answer),
          'cannot report on previous rounds',
        )
      })
    })

    describe('when updated after the max answers submitted', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator)
        await aggregator.connect(personas.Neil).submit(nextRound, answer)
      })

      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Ned).submit(nextRound, answer),
          'round not accepting submissions',
        )
      })
    })

    describe('when a new highest round number is passed in', () => {
      it('increments the answer round', async () => {
        matchers.bigNum(
          ethers.constants.Zero,
          await aggregator.reportingRound(),
        )

        await advanceRound(aggregator, oracles)

        matchers.bigNum(ethers.constants.One, await aggregator.reportingRound())
      })

      it('sets the startedAt time for the reportingRound', async () => {
        let round = await aggregator.getRoundData(nextRound)
        matchers.bigNum(ethers.constants.Zero, round.startedAt)

        await aggregator.connect(oracles[0]).submit(nextRound, answer)

        round = await aggregator.getRoundData(nextRound)

        expect(round.startedAt).not.toBe(0)
      })

      it('announces a new round by emitting a log', async () => {
        const tx = await aggregator
          .connect(personas.Neil)
          .submit(nextRound, answer)
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
          aggregator.connect(personas.Neil).submit(nextRound + 1, answer),
          'invalid round to report',
        )
      })
    })

    describe('when called by a non-oracle', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Carol).submit(nextRound, answer),
          'not enabled oracle',
        )
      })
    })

    describe('when there are not sufficient available funds', () => {
      beforeEach(async () => {
        await aggregator
          .connect(personas.Carol)
          .withdrawFunds(
            personas.Carol.address,
            deposit.sub(paymentAmount.mul(oracles.length).mul(reserveRounds)),
          )

        // drain remaining funds
        await advanceRound(aggregator, oracles)
        await advanceRound(aggregator, oracles)
      })

      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Neil).submit(nextRound, answer),
          'SafeMath: subtraction overflow',
        )
      })
    })

    describe('when a new round opens before the previous rounds closes', () => {
      beforeEach(async () => {
        oracles = [personas.Nancy, personas.Norbert]
        await addOracles(aggregator, oracles, 3, 4, rrDelay)
        await advanceRound(aggregator, [
          personas.Nelly,
          personas.Neil,
          personas.Nancy,
        ])

        // start the next round
        await aggregator.connect(personas.Nelly).submit(nextRound, answer)
      })

      it('still allows the previous round to be answered', async () => {
        await aggregator.connect(personas.Ned).submit(nextRound - 1, answer)
      })

      describe('once the current round is answered', () => {
        beforeEach(async () => {
          oracles = [personas.Neil, personas.Nancy]
          for (let i = 0; i < oracles.length; i++) {
            await aggregator.connect(oracles[i]).submit(nextRound, answer)
          }
        })

        it('does not allow reports for the previous round', async () => {
          await matchers.evmRevert(
            aggregator.connect(personas.Ned).submit(nextRound - 1, answer),
            'invalid round to report',
          )
        })
      })

      describe('when the previous round has finished', () => {
        beforeEach(async () => {
          await aggregator
            .connect(personas.Norbert)
            .submit(nextRound - 1, answer)
        })

        it('does not allow reports for the previous round', async () => {
          await matchers.evmRevert(
            aggregator.connect(personas.Ned).submit(nextRound - 1, answer),
            'round not accepting submissions',
          )
        })
      })
    })

    describe('when price is updated mid-round', () => {
      const newAmount = h.toWei('50')

      it('pays the same amount to all oracles per round', async () => {
        await link.transfer(
          aggregator.address,
          newAmount.mul(oracles.length).mul(reserveRounds),
        )
        await aggregator.updateAvailableFunds()

        matchers.bigNum(
          0,
          await aggregator
            .connect(personas.Neil)
            .withdrawablePayment(personas.Neil.address),
        )
        matchers.bigNum(
          0,
          await aggregator
            .connect(personas.Nelly)
            .withdrawablePayment(personas.Nelly.address),
        )

        await aggregator.connect(personas.Neil).submit(nextRound, answer)

        await updateFutureRounds(aggregator, { payment: newAmount })

        await aggregator.connect(personas.Nelly).submit(nextRound, answer)

        matchers.bigNum(
          paymentAmount,
          await aggregator
            .connect(personas.Neil)
            .withdrawablePayment(personas.Neil.address),
        )
        matchers.bigNum(
          paymentAmount,
          await aggregator
            .connect(personas.Nelly)
            .withdrawablePayment(personas.Nelly.address),
        )
      })
    })

    describe('when delay is on', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator, {
          minAnswers: oracles.length,
          maxAnswers: oracles.length,
          restartDelay: 1,
        })
      })

      it("does not revert on the oracle's first round", async () => {
        // Since lastUpdatedRound defaults to zero and that's the only
        // indication that an oracle hasn't responded, this test guards against
        // the situation where we don't check that and no one can start a round.

        await aggregator.connect(personas.Neil).submit(nextRound, answer)
      })

      it('does revert before the delay', async () => {
        await aggregator.connect(personas.Neil).submit(nextRound, answer)

        nextRound++

        await matchers.evmRevert(
          aggregator.connect(personas.Neil).submit(nextRound, answer),
          'previous round not supersedable',
        )
      })
    })

    describe('when an oracle starts a round before the restart delay is over', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator.connect(personas.Carol), {
          minAnswers: 1,
          maxAnswers: 1,
        })

        oracles = [personas.Neil, personas.Ned, personas.Nelly]
        for (let i = 0; i < oracles.length; i++) {
          await aggregator.connect(oracles[i]).submit(nextRound, answer)
          nextRound++
        }

        const newDelay = 2
        // Since Ned and Nelly have answered recently, and we set the delay
        // to 2, only Nelly can answer as she is the only oracle that hasn't
        // started the last two rounds.
        await updateFutureRounds(aggregator, {
          maxAnswers: oracles.length,
          restartDelay: newDelay,
        })
      })

      describe('when called by an oracle who has not answered recently', () => {
        it('does not revert', async () => {
          await aggregator.connect(personas.Neil).submit(nextRound, answer)
        })
      })

      describe('when called by an oracle who answered recently', () => {
        it('reverts', async () => {
          await matchers.evmRevert(
            aggregator.connect(personas.Ned).submit(nextRound, answer),
            'round not accepting submissions',
          )

          await matchers.evmRevert(
            aggregator.connect(personas.Nelly).submit(nextRound, answer),
            'round not accepting submissions',
          )
        })
      })
    })

    describe('when the price is not updated for a round', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator, {
          minAnswers: oracles.length,
          maxAnswers: oracles.length,
          restartDelay: 1,
        })

        for (const oracle of oracles) {
          await aggregator.connect(oracle).submit(nextRound, answer)
        }
        nextRound++

        await aggregator.connect(personas.Ned).submit(nextRound, answer)
        await aggregator.connect(personas.Nelly).submit(nextRound, answer)
        assert.equal(nextRound, (await aggregator.reportingRound()).toNumber())

        await h.increaseTimeBy(timeout + 1, provider)
        nextRound++
      })

      it('allows a new round to be started', async () => {
        await aggregator.connect(personas.Nelly).submit(nextRound, answer)
      })

      it('sets the info for the previous round', async () => {
        const previousRound = nextRound - 1
        let updated = await aggregator.getTimestamp(previousRound)
        let ans = await aggregator.getAnswer(previousRound)
        assert.equal(0, updated.toNumber())
        assert.equal(0, ans.toNumber())

        const tx = await aggregator
          .connect(personas.Nelly)
          .submit(nextRound, answer)
        const receipt = await tx.wait()

        const block = await provider.getBlock(receipt.blockHash ?? '')

        updated = await aggregator.getTimestamp(previousRound)
        ans = await aggregator.getAnswer(previousRound)
        matchers.bigNum(ethers.utils.bigNumberify(block.timestamp), updated)
        assert.equal(answer, ans.toNumber())

        const round = await aggregator.getRoundData(previousRound)
        matchers.bigNum(previousRound, round.roundId)
        matchers.bigNum(ans, round.answer)
        matchers.bigNum(updated, round.updatedAt)
        matchers.bigNum(previousRound - 1, round.answeredInRound)
      })

      it('sets the previous round as timed out', async () => {
        const previousRound = nextRound - 1
        let round = await aggregator.getRoundData(previousRound)
        matchers.bigNum(0, round.answeredInRound)

        await aggregator.connect(personas.Nelly).submit(nextRound, answer)

        round = await aggregator.getRoundData(previousRound)
        assert.notEqual(round.roundId, round.answeredInRound)
        matchers.bigNum(previousRound - 1, round.answeredInRound)
      })

      it('still respects the delay restriction', async () => {
        // expected to revert because the sender started the last round
        await matchers.evmRevert(
          aggregator.connect(personas.Ned).submit(nextRound, answer),
        )
      })

      it('uses the timeout set at the beginning of the round', async () => {
        await updateFutureRounds(aggregator, {
          timeout: timeout + 100000,
        })

        await aggregator.connect(personas.Nelly).submit(nextRound, answer)
      })
    })

    describe('submitting values near the edges of allowed values', () => {
      it('rejects values below the submission value range', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Neil)
            .submit(nextRound, minSubmissionValue.sub(1)),
          'value below minSubmissionValue',
        )
      })

      it('accepts submissions equal to the min submission value', async () => {
        await aggregator
          .connect(personas.Neil)
          .submit(nextRound, minSubmissionValue)
      })

      it('accepts submissions equal to the max submission value', async () => {
        await aggregator
          .connect(personas.Neil)
          .submit(nextRound, maxSubmissionValue)
      })

      it('rejects submissions equal to the max submission value', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Neil)
            .submit(nextRound, maxSubmissionValue.add(1)),
          'value above maxSubmissionValue',
        )
      })
    })

    describe('when an answer validator is set', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator, { minAnswers: 1, maxAnswers: 1 })
        oracles = [personas.Nelly]

        validator = await answerValidatorFactory
          .connect(personas.Carol)
          .deploy()
        await aggregator
          .connect(personas.Carol)
          .setAnswerValidator(validator.address)
        assert.equal(validator.address, await aggregator.answerValidator())
      })

      it('calls out to the validator', async () => {
        const tx = await aggregator
          .connect(personas.Nelly)
          .submit(nextRound, answer)
        const receipt = await tx.wait()

        const event = matchers.eventExists(
          receipt,
          validator.interface.events.Validated,
        )
        matchers.bigNum(0, h.bigNum(event.topics[1]))
        matchers.bigNum(answer, h.bigNum(event.topics[2]))
      })
    })

    describe('when the answer validator eats all gas', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator, { minAnswers: 1, maxAnswers: 1 })
        oracles = [personas.Nelly]

        gasGuzzler = await gasGuzzlerFactory.connect(personas.Carol).deploy()
        await aggregator
          .connect(personas.Carol)
          .setAnswerValidator(gasGuzzler.address)
        assert.equal(gasGuzzler.address, await aggregator.answerValidator())
      })

      it('still updates', async () => {
        matchers.bigNum(0, await aggregator.latestAnswer())

        await aggregator
          .connect(personas.Nelly)
          .submit(nextRound, answer, { gasLimit: 500000 })

        matchers.bigNum(answer, await aggregator.latestAnswer())
      })
    })
  })

  describe('#setAnswerValidator', () => {
    beforeEach(async () => {
      validator = await answerValidatorFactory.connect(personas.Carol).deploy()
    })

    it('changes the answer validator', async () => {
      assert.equal(emptyAddress, await aggregator.answerValidator())

      await aggregator
        .connect(personas.Carol)
        .setAnswerValidator(validator.address)

      assert.equal(validator.address, await aggregator.answerValidator())
    })

    it('emits a log event', async () => {
      const tx = await aggregator
        .connect(personas.Carol)
        .setAnswerValidator(validator.address)
      const receipt = await tx.wait()
      const eventLog = matchers.eventExists(
        receipt,
        aggregator.interface.events.AnswerValidatorUpdated,
      )

      assert.equal(emptyAddress, h.eventArgs(eventLog).previous)
      assert.equal(validator.address, h.eventArgs(eventLog).current)

      const sameChangeTx = await aggregator
        .connect(personas.Carol)
        .setAnswerValidator(validator.address)
      const sameChangeReceipt = await sameChangeTx.wait()
      assert.equal(0, sameChangeReceipt.events?.length)
      matchers.eventDoesNotExist(
        sameChangeReceipt,
        aggregator.interface.events.AnswerValidatorUpdated,
      )
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Neil)
            .setAnswerValidator(validator.address),
          'Only callable by owner',
        )
      })
    })
  })

  describe('integrating with historic deviation checker', () => {
    let validator: contract.Instance<HistoricDeviationValidatorFactory>
    let flags: contract.Instance<FlagsFactory>
    let ac: contract.Instance<SimpleWriteAccessControllerFactory>
    const flaggingThreshold = 1000 // 1%

    beforeEach(async () => {
      ac = await acFactory.connect(personas.Carol).deploy()
      flags = await flagsFactory.connect(personas.Carol).deploy(ac.address)
      validator = await validatorFactory
        .connect(personas.Carol)
        .deploy(flags.address, flaggingThreshold)
      await ac.connect(personas.Carol).addAccess(validator.address)

      await aggregator
        .connect(personas.Carol)
        .setAnswerValidator(validator.address)

      oracles = [personas.Nelly]
      const minMax = oracles.length
      await addOracles(aggregator, oracles, minMax, minMax, rrDelay)
    })

    it('raises a flag on with high enough deviation', async () => {
      await aggregator.connect(personas.Nelly).submit(nextRound, 100)
      nextRound++

      const tx = await aggregator.connect(personas.Nelly).submit(nextRound, 102)
      const receipt = await tx.wait()
      const event = matchers.eventExists(receipt, flags.interface.events.FlagOn)

      assert.equal(flags.address, event.address)
      assert.equal(aggregator.address, h.evmWordToAddress(event.topics[1]))
    })

    it('does not raise a flag with low enough deviation', async () => {
      await aggregator.connect(personas.Nelly).submit(nextRound, 100)
      nextRound++

      const tx = await aggregator.connect(personas.Nelly).submit(nextRound, 101)
      const receipt = await tx.wait()
      matchers.eventDoesNotExist(receipt, flags.interface.events.FlagOn)
    })
  })
})
