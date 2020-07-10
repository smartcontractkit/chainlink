import {
  contract,
  helpers as h,
  matchers,
  setup,
  wallet,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { FluxAggregatorFactory } from '../../ethers/v0.6/FluxAggregatorFactory'
import { FluxAggregatorTestHelperFactory } from '../../ethers/v0.6/FluxAggregatorTestHelperFactory'
import { AggregatorValidatorMockFactory } from '../../ethers/v0.6/AggregatorValidatorMockFactory'
import { GasGuzzlerFactory } from '../../ethers/v0.6/GasGuzzlerFactory'
import { DeviationFlaggingValidatorFactory } from '../../ethers/v0.6/DeviationFlaggingValidatorFactory'
import { FlagsFactory } from '../../ethers/v0.6/FlagsFactory'
import { SimpleWriteAccessControllerFactory } from '../../ethers/v0.6/SimpleWriteAccessControllerFactory'

let personas: setup.Personas
const provider = setup.provider()
const linkTokenFactory = new contract.LinkTokenFactory()
const fluxAggregatorFactory = new FluxAggregatorFactory()
const validatorMockFactory = new AggregatorValidatorMockFactory()
const testHelperFactory = new FluxAggregatorTestHelperFactory()
const validatorFactory = new DeviationFlaggingValidatorFactory()
const flagsFactory = new FlagsFactory()
const acFactory = new SimpleWriteAccessControllerFactory()
const gasGuzzlerFactory = new GasGuzzlerFactory()
const emptyAddress = '0x0000000000000000000000000000000000000000'

beforeAll(async () => {
  personas = await setup.users(provider).then(x => x.personas)
})

describe('FluxAggregator', () => {
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

  let aggregator: contract.Instance<FluxAggregatorFactory>
  let link: contract.Instance<contract.LinkTokenFactory>
  let testHelper: contract.Instance<FluxAggregatorTestHelperFactory>
  let validator: contract.Instance<AggregatorValidatorMockFactory>
  let gasGuzzler: contract.Instance<GasGuzzlerFactory>
  let nextRound: number
  let oracles: ethers.Wallet[]

  async function updateFutureRounds(
    aggregator: contract.Instance<FluxAggregatorFactory>,
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
    aggregator: contract.Instance<FluxAggregatorFactory>,
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
    aggregator: contract.Instance<FluxAggregatorFactory>,
    submitters: ethers.Wallet[],
    currentSubmission: number = answer,
  ): Promise<number> {
    for (const submitter of submitters) {
      await aggregator.connect(submitter).submit(nextRound, currentSubmission)
    }
    nextRound++
    return nextRound
  }

  const ShouldBeSet = 'expects it to be different'
  const ShouldNotBeSet = 'expects it to equal'
  let startingState: any

  async function checkOracleRoundState(
    state: any,
    want: {
      eligibleToSubmit: boolean
      roundId: ethers.utils.BigNumberish
      latestSubmission: ethers.utils.BigNumberish
      startedAt: string
      timeout: ethers.utils.BigNumberish
      availableFunds: ethers.utils.BigNumberish
      oracleCount: ethers.utils.BigNumberish
      paymentAmount: ethers.utils.BigNumberish
    },
  ) {
    assert.equal(
      want.eligibleToSubmit,
      state._eligibleToSubmit,
      'round state: unexecpted eligibility',
    )
    matchers.bigNum(
      want.roundId,
      state._roundId,
      'round state: unexpected Round ID',
    )
    matchers.bigNum(
      want.latestSubmission,
      state._latestSubmission,
      'round state: unexpected latest submission',
    )
    if (want.startedAt === ShouldBeSet) {
      assert.isAbove(
        state._startedAt.toNumber(),
        startingState._startedAt.toNumber(),
        'round state: expected the started at to be the same as previous',
      )
    } else {
      matchers.bigNum(
        0,
        state._startedAt,
        'round state: expected the started at not to be updated',
      )
    }
    matchers.bigNum(
      want.timeout,
      state._timeout.toNumber(),
      'round state: unexepcted timeout',
    )
    matchers.bigNum(
      want.availableFunds,
      state._availableFunds,
      'round state: unexepected funds',
    )
    matchers.bigNum(
      want.oracleCount,
      state._oracleCount,
      'round state: unexpected oracle count',
    )
    matchers.bigNum(
      want.paymentAmount,
      state._paymentAmount,
      'round state: unexpected paymentamount',
    )
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
      'setValidator',
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
        matchers.evmRevert(
          aggregator.getRoundData(nextRound),
          'No data present',
        )

        await aggregator.connect(personas.Nelly).submit(nextRound, answer)

        const round = await aggregator.getRoundData(nextRound)
        assert.equal(nextRound, round.answeredInRound.toNumber())
      })

      it('updates the round details', async () => {
        matchers.evmRevert(aggregator.latestRoundData(), 'No data present')

        h.increaseTimeBy(15, provider)
        await aggregator.connect(personas.Nelly).submit(nextRound, answer)

        const roundAfter = await aggregator.getRoundData(nextRound)
        matchers.bigNum(nextRound, roundAfter.roundId)
        matchers.bigNum(answer, roundAfter.answer)
        assert.isFalse(roundAfter.startedAt.isZero())
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
        matchers.evmRevert(
          aggregator.getRoundData(nextRound),
          'No data present',
        )

        await aggregator.connect(oracles[0]).submit(nextRound, answer)

        matchers.evmRevert(
          aggregator.getRoundData(nextRound),
          'No data present',
        )
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
        matchers.evmRevert(
          aggregator.getRoundData(previousRound),
          'No data present',
        )

        await aggregator.connect(personas.Nelly).submit(nextRound, answer)

        const round = await aggregator.getRoundData(previousRound)
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

    describe('when a validator is set', () => {
      beforeEach(async () => {
        await updateFutureRounds(aggregator, { minAnswers: 1, maxAnswers: 1 })
        oracles = [personas.Nelly]

        validator = await validatorMockFactory.connect(personas.Carol).deploy()
        await aggregator.connect(personas.Carol).setValidator(validator.address)
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
          .setValidator(gasGuzzler.address)
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

  describe('#getAnswer', () => {
    const answers = [1, 10, 101, 1010, 10101, 101010, 1010101]

    beforeEach(async () => {
      await addOracles(aggregator, [personas.Neil], minAns, maxAns, rrDelay)

      for (const answer of answers) {
        await aggregator.connect(personas.Neil).submit(nextRound, answer)
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
      await addOracles(aggregator, [personas.Neil], minAns, maxAns, rrDelay)

      for (let i = 0; i < 10; i++) {
        await aggregator.connect(personas.Neil).submit(nextRound, i + 1)
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

  describe('#addOracles', () => {
    it('increases the oracle count', async () => {
      const pastCount = await aggregator.oracleCount()
      await addOracles(aggregator, [personas.Neil], minAns, maxAns, rrDelay)
      const currentCount = await aggregator.oracleCount()

      matchers.bigNum(currentCount, pastCount + 1)
    })

    it('adds the address in getOracles', async () => {
      await addOracles(aggregator, [personas.Neil], minAns, maxAns, rrDelay)
      assert.deepEqual([personas.Neil.address], await aggregator.getOracles())
    })

    it('updates the round details', async () => {
      await addOracles(
        aggregator,
        [personas.Neil, personas.Ned, personas.Nelly],
        1,
        3,
        2,
      )
      matchers.bigNum(1, await aggregator.minSubmissionCount())
      matchers.bigNum(3, await aggregator.maxSubmissionCount())
      matchers.bigNum(2, await aggregator.restartDelay())
    })

    it('emits a log', async () => {
      const tx = await aggregator
        .connect(personas.Carol)
        .addOracles([personas.Ned.address], [personas.Neil.address], 1, 1, 0)
      const receipt = await tx.wait()

      const oracleAddedEvent = h.eventArgs(receipt.events?.[0])
      assert.equal(oracleAddedEvent.oracle, personas.Ned.address)
      assert.isTrue(oracleAddedEvent.whitelisted)
      const oracleAdminUpdatedEvent = h.eventArgs(receipt.events?.[1])
      assert.equal(oracleAdminUpdatedEvent.oracle, personas.Ned.address)
      assert.equal(oracleAdminUpdatedEvent.newAdmin, personas.Neil.address)
    })

    describe('when the oracle has already been added', () => {
      beforeEach(async () => {
        await addOracles(aggregator, [personas.Neil], minAns, maxAns, rrDelay)
      })

      it('reverts', async () => {
        await matchers.evmRevert(
          addOracles(aggregator, [personas.Neil], minAns, maxAns, rrDelay),
          'oracle already enabled',
        )
      })
    })

    describe('when called by anyone but the owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Neil)
            .addOracles(
              [personas.Neil.address],
              [personas.Neil.address],
              minAns,
              maxAns,
              rrDelay,
            ),
          'Only callable by owner',
        )
      })
    })

    describe('when an oracle gets added mid-round', () => {
      beforeEach(async () => {
        oracles = [personas.Neil, personas.Ned]
        await addOracles(
          aggregator,
          oracles,
          oracles.length,
          oracles.length,
          rrDelay,
        )

        await aggregator.connect(personas.Neil).submit(nextRound, answer)

        await addOracles(
          aggregator,
          [personas.Nelly],
          oracles.length + 1,
          oracles.length + 1,
          rrDelay,
        )
      })

      it('does not allow the oracle to update the round', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Nelly).submit(nextRound, answer),
          'not yet enabled oracle',
        )
      })

      it('does allow the oracle to update future rounds', async () => {
        // complete round
        await aggregator.connect(personas.Ned).submit(nextRound, answer)

        // now can participate in new rounds
        await aggregator.connect(personas.Nelly).submit(nextRound + 1, answer)
      })
    })

    describe('when an oracle is added after removed for a round', () => {
      it('allows the oracle to update', async () => {
        oracles = [personas.Neil, personas.Nelly]
        await addOracles(
          aggregator,
          oracles,
          oracles.length,
          oracles.length,
          rrDelay,
        )

        await aggregator.connect(personas.Neil).submit(nextRound, answer)
        await aggregator.connect(personas.Nelly).submit(nextRound, answer)
        nextRound++

        await aggregator
          .connect(personas.Carol)
          .removeOracles([personas.Nelly.address], 1, 1, rrDelay)

        await aggregator.connect(personas.Neil).submit(nextRound, answer)
        nextRound++

        await addOracles(aggregator, [personas.Nelly], 1, 1, rrDelay)

        await aggregator.connect(personas.Nelly).submit(nextRound, answer)
      })
    })

    describe('when an oracle is added and immediately removed mid-round', () => {
      it('allows the oracle to update', async () => {
        await addOracles(
          aggregator,
          oracles,
          oracles.length,
          oracles.length,
          rrDelay,
        )

        await aggregator.connect(personas.Neil).submit(nextRound, answer)
        await aggregator.connect(personas.Nelly).submit(nextRound, answer)
        nextRound++

        await aggregator
          .connect(personas.Carol)
          .removeOracles([personas.Nelly.address], 1, 1, rrDelay)

        await aggregator.connect(personas.Neil).submit(nextRound, answer)
        nextRound++

        await addOracles(aggregator, [personas.Nelly], 1, 1, rrDelay)

        await aggregator.connect(personas.Nelly).submit(nextRound, answer)
      })
    })

    describe('when an oracle is re-added with a different admin address', () => {
      it('reverts', async () => {
        await addOracles(
          aggregator,
          oracles,
          oracles.length,
          oracles.length,
          rrDelay,
        )

        await aggregator.connect(personas.Neil).submit(nextRound, answer)

        await aggregator
          .connect(personas.Carol)
          .removeOracles([personas.Nelly.address], 1, 1, rrDelay)

        await matchers.evmRevert(
          aggregator
            .connect(personas.Carol)
            .addOracles(
              [personas.Nelly.address],
              [personas.Carol.address],
              1,
              1,
              rrDelay,
            ),
          'owner cannot overwrite admin',
        )
      })
    })

    const limit = 77
    describe(`when adding more than ${limit} oracles`, () => {
      let oracles: ethers.Wallet[]

      beforeEach(async () => {
        oracles = []
        for (let i = 0; i < limit; i++) {
          const account = await wallet.createWallet(provider, i + 100)
          await personas.Default.sendTransaction({
            to: account.address,
            value: h.toWei('0.01'),
          })
          oracles.push(account)
        }

        await link.transfer(
          aggregator.address,
          paymentAmount.mul(limit).mul(reserveRounds),
        )
        await aggregator.updateAvailableFunds()

        let addresses = oracles.slice(0, 50).map(o => o.address)
        await aggregator
          .connect(personas.Carol)
          .addOracles(addresses, addresses, 1, 50, rrDelay)
        // add in two transactions to avoid gas limit issues
        addresses = oracles.slice(50, 100).map(o => o.address)
        await aggregator
          .connect(personas.Carol)
          .addOracles(addresses, addresses, 1, oracles.length, rrDelay)
      })

      it('not use too much gas', async () => {
        let tx: any
        assert.deepEqual(
          // test adveserial quickselect algo
          [2, 4, 6, 8, 10, 12, 14, 16, 1, 9, 5, 11, 3, 13, 7, 15],
          adverserialQuickselectList(16),
        )
        const inputs = adverserialQuickselectList(limit)
        for (let i = 0; i < limit; i++) {
          tx = await aggregator.connect(oracles[i]).submit(nextRound, inputs[i])
        }
        assert(!!tx)
        if (tx) {
          const receipt = await tx.wait()
          assert.isAbove(400_000, receipt.gasUsed.toNumber())
        }
      })

      function adverserialQuickselectList(len: number): number[] {
        const xs: number[] = []
        const pi: number[] = []
        for (let i = 0; i < len; i++) {
          pi[i] = i
          xs[i] = 0
        }

        for (let l = len; l > 0; l--) {
          const pivot = Math.floor((l - 1) / 2)
          xs[pi[pivot]] = l
          const temp = pi[l - 1]
          pi[l - 1] = pi[pivot]
          pi[pivot] = temp
        }
        return xs
      }

      it('reverts when another oracle is added', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Carol)
            .addOracles(
              [personas.Neil.address],
              [personas.Neil.address],
              limit + 1,
              limit + 1,
              rrDelay,
            ),
          'max oracles allowed',
        )
      })
    })

    it('reverts when minSubmissions is set to 0', async () => {
      await matchers.evmRevert(
        addOracles(aggregator, [personas.Neil], 0, 0, 0),
        'min must be greater than 0',
      )
    })
  })

  describe('#removeOracles', () => {
    beforeEach(async () => {
      oracles = [personas.Neil, personas.Nelly]
      await addOracles(
        aggregator,
        oracles,
        oracles.length,
        oracles.length,
        rrDelay,
      )
    })

    it('decreases the oracle count', async () => {
      const pastCount = await aggregator.oracleCount()
      await aggregator
        .connect(personas.Carol)
        .removeOracles([personas.Neil.address], minAns, maxAns, rrDelay)
      const currentCount = await aggregator.oracleCount()

      expect(currentCount).toEqual(pastCount - 1)
    })

    it('updates the round details', async () => {
      await aggregator
        .connect(personas.Carol)
        .removeOracles([personas.Neil.address], 1, 1, 0)

      matchers.bigNum(1, await aggregator.minSubmissionCount())
      matchers.bigNum(1, await aggregator.maxSubmissionCount())
      matchers.bigNum(ethers.constants.Zero, await aggregator.restartDelay())
    })

    it('emits a log', async () => {
      const tx = await aggregator
        .connect(personas.Carol)
        .removeOracles([personas.Neil.address], minAns, maxAns, rrDelay)
      const receipt = await tx.wait()

      const oracleRemovedEvent = h.eventArgs(receipt.events?.[0])
      assert.equal(oracleRemovedEvent.oracle, personas.Neil.address)
      assert.isFalse(oracleRemovedEvent.whitelisted)
    })

    it('removes the address in getOracles', async () => {
      await aggregator
        .connect(personas.Carol)
        .removeOracles([personas.Neil.address], minAns, maxAns, rrDelay)
      assert.deepEqual([personas.Nelly.address], await aggregator.getOracles())
    })

    describe('when the oracle is not currently added', () => {
      beforeEach(async () => {
        await aggregator
          .connect(personas.Carol)
          .removeOracles([personas.Neil.address], minAns, maxAns, rrDelay)
      })

      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Carol)
            .removeOracles([personas.Neil.address], minAns, maxAns, rrDelay),
          'oracle not enabled',
        )
      })
    })

    describe('when removing the last oracle', () => {
      it('does not revert', async () => {
        await aggregator
          .connect(personas.Carol)
          .removeOracles([personas.Neil.address], minAns, maxAns, rrDelay)

        await aggregator
          .connect(personas.Carol)
          .removeOracles([personas.Nelly.address], 0, 0, 0)
      })
    })

    describe('when called by anyone but the owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Ned)
            .removeOracles([personas.Neil.address], 0, 0, rrDelay),
          'Only callable by owner',
        )
      })
    })

    describe('when an oracle gets removed', () => {
      beforeEach(async () => {
        await aggregator
          .connect(personas.Carol)
          .removeOracles([personas.Nelly.address], 1, 1, rrDelay)
      })

      it('is allowed to report on one more round', async () => {
        // next round
        await advanceRound(aggregator, [personas.Nelly])
        // finish round
        await advanceRound(aggregator, [personas.Neil])

        // cannot participate in future rounds
        await matchers.evmRevert(
          aggregator.connect(personas.Nelly).submit(nextRound, answer),
          'no longer allowed oracle',
        )
      })
    })

    describe('when an oracle gets removed mid-round', () => {
      beforeEach(async () => {
        await aggregator.connect(personas.Neil).submit(nextRound, answer)

        await aggregator
          .connect(personas.Carol)
          .removeOracles([personas.Nelly.address], 1, 1, rrDelay)
      })

      it('is allowed to finish that round and one more round', async () => {
        await advanceRound(aggregator, [personas.Nelly]) // finish round

        await advanceRound(aggregator, [personas.Nelly]) // next round

        // cannot participate in future rounds
        await matchers.evmRevert(
          aggregator.connect(personas.Nelly).submit(nextRound, answer),
          'no longer allowed oracle',
        )
      })
    })

    it('reverts when minSubmissions is set to 0', async () => {
      await matchers.evmRevert(
        aggregator
          .connect(personas.Carol)
          .removeOracles([personas.Nelly.address], 0, 0, 0),
        'min must be greater than 0',
      )
    })
  })

  describe('#getOracles', () => {
    describe('after adding oracles', () => {
      beforeEach(async () => {
        await addOracles(aggregator, [personas.Neil], minAns, maxAns, rrDelay)

        assert.deepEqual([personas.Neil.address], await aggregator.getOracles())
      })

      it('returns the addresses of added oracles', async () => {
        await addOracles(aggregator, [personas.Ned], minAns, maxAns, rrDelay)

        assert.deepEqual(
          [personas.Neil.address, personas.Ned.address],
          await aggregator.getOracles(),
        )

        await addOracles(aggregator, [personas.Nelly], minAns, maxAns, rrDelay)
        assert.deepEqual(
          [personas.Neil.address, personas.Ned.address, personas.Nelly.address],
          await aggregator.getOracles(),
        )
      })
    })

    describe('after removing oracles', () => {
      beforeEach(async () => {
        await addOracles(
          aggregator,
          [personas.Neil, personas.Ned, personas.Nelly],
          minAns,
          maxAns,
          rrDelay,
        )

        assert.deepEqual(
          [personas.Neil.address, personas.Ned.address, personas.Nelly.address],
          await aggregator.getOracles(),
        )
      })

      it('reorders when removing from the beginning', async () => {
        await aggregator
          .connect(personas.Carol)
          .removeOracles([personas.Neil.address], minAns, maxAns, rrDelay)
        assert.deepEqual(
          [personas.Nelly.address, personas.Ned.address],
          await aggregator.getOracles(),
        )
      })

      it('reorders when removing from the middle', async () => {
        await aggregator
          .connect(personas.Carol)
          .removeOracles([personas.Ned.address], minAns, maxAns, rrDelay)
        assert.deepEqual(
          [personas.Neil.address, personas.Nelly.address],
          await aggregator.getOracles(),
        )
      })

      it('pops the last node off at the end', async () => {
        await aggregator
          .connect(personas.Carol)
          .removeOracles([personas.Nelly.address], minAns, maxAns, rrDelay)
        assert.deepEqual(
          [personas.Neil.address, personas.Ned.address],
          await aggregator.getOracles(),
        )
      })
    })
  })

  describe('#withdrawFunds', () => {
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
        'insufficient reserve funds',
      )
    })

    describe('with a number higher than the available LINK balance', () => {
      beforeEach(async () => {
        await addOracles(aggregator, [personas.Neil], minAns, maxAns, rrDelay)

        await aggregator.connect(personas.Neil).submit(nextRound, answer)
      })

      it('fails', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Carol)
            .withdrawFunds(personas.Carol.address, deposit),
          'insufficient reserve funds',
        )

        matchers.bigNum(
          deposit.sub(paymentAmount),
          await aggregator.availableFunds(),
        )
      })
    })

    describe('with oracles still present', () => {
      beforeEach(async () => {
        oracles = [personas.Neil, personas.Ned, personas.Nelly]
        await addOracles(aggregator, oracles, 1, 1, rrDelay)

        matchers.bigNum(deposit, await aggregator.availableFunds())
      })

      it('does not allow withdrawal with less than 2x rounds of payments', async () => {
        const oracleReserve = paymentAmount
          .mul(oracles.length)
          .mul(reserveRounds)
        const allowed = deposit.sub(oracleReserve)

        //one more than the allowed amount cannot be withdrawn
        await matchers.evmRevert(
          aggregator
            .connect(personas.Carol)
            .withdrawFunds(personas.Carol.address, allowed.add(1)),
          'insufficient reserve funds',
        )

        // the allowed amount can be withdrawn
        await aggregator
          .connect(personas.Carol)
          .withdrawFunds(personas.Carol.address, allowed)
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
    let minSubmissionCount, maxSubmissionCount
    const newPaymentAmount = h.toWei('2')
    const newMin = 1
    const newMax = 3
    const newDelay = 2

    beforeEach(async () => {
      oracles = [personas.Neil, personas.Ned, personas.Nelly]
      minSubmissionCount = oracles.length
      maxSubmissionCount = oracles.length
      await addOracles(
        aggregator,
        oracles,
        minSubmissionCount,
        maxSubmissionCount,
        rrDelay,
      )

      matchers.bigNum(paymentAmount, await aggregator.paymentAmount())
      assert.equal(minSubmissionCount, await aggregator.minSubmissionCount())
      assert.equal(maxSubmissionCount, await aggregator.maxSubmissionCount())
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
        await aggregator.minSubmissionCount(),
      )
      matchers.bigNum(
        ethers.utils.bigNumberify(newMax),
        await aggregator.maxSubmissionCount(),
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
      assert.equal(newMin, round.minSubmissionCount)
      assert.equal(newMax, round.maxSubmissionCount)
      assert.equal(newDelay, round.restartDelay)
      assert.equal(timeout + 1, round.timeout)
    })

    describe('when it is set to higher than the number or oracles', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          updateFutureRounds(aggregator, {
            maxAnswers: 4,
          }),
          'max cannot exceed total',
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
          'max must equal/exceed min',
        )
      })
    })

    describe('when delay equal or greater the oracle count', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          updateFutureRounds(aggregator, {
            restartDelay: 3,
          }),
          'revert delay cannot exceed total',
        )
      })
    })

    describe('when the payment amount does not cover reserve rounds', () => {
      beforeEach(async () => {})

      it('reverts', async () => {
        const most = deposit.div(oracles.length * reserveRounds)

        await matchers.evmRevert(
          updateFutureRounds(aggregator, {
            payment: most.add(1),
          }),
          'insufficient funds for payment',
        )

        await updateFutureRounds(aggregator, {
          payment: most,
        })
      })
    })

    describe('min oracles is set to 0', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator.updateFutureRounds(paymentAmount, 0, 0, rrDelay, timeout),
          'min must be greater than 0',
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

      await addOracles(aggregator, [personas.Neil], minAns, maxAns, rrDelay)
      await aggregator.connect(personas.Neil).submit(nextRound, answer)
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

  describe('#withdrawPayment', () => {
    beforeEach(async () => {
      await addOracles(aggregator, [personas.Neil], minAns, maxAns, rrDelay)
      await aggregator.connect(personas.Neil).submit(nextRound, answer)
    })

    it('transfers LINK to the recipient', async () => {
      const originalBalance = await link.balanceOf(aggregator.address)
      matchers.bigNum(0, await link.balanceOf(personas.Neil.address))

      await aggregator
        .connect(personas.Neil)
        .withdrawPayment(
          personas.Neil.address,
          personas.Neil.address,
          paymentAmount,
        )

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
        .withdrawPayment(
          personas.Neil.address,
          personas.Neil.address,
          paymentAmount,
        )

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
            .withdrawPayment(
              personas.Neil.address,
              personas.Neil.address,
              paymentAmount.add(ethers.utils.bigNumberify(1)),
            ),
          'revert insufficient withdrawable funds',
        )
      })
    })

    describe('when the caller is not the admin', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Nelly)
            .withdrawPayment(
              personas.Neil.address,
              personas.Nelly.address,
              ethers.utils.bigNumberify(1),
            ),
          'only callable by admin',
        )
      })
    })
  })

  describe('#transferAdmin', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .addOracles(
          [personas.Ned.address],
          [personas.Neil.address],
          minAns,
          maxAns,
          rrDelay,
        )
    })

    describe('when the admin tries to transfer the admin', () => {
      it('works', async () => {
        const tx = await aggregator
          .connect(personas.Neil)
          .transferAdmin(personas.Ned.address, personas.Nelly.address)
        const receipt = await tx.wait()
        assert.equal(
          personas.Neil.address,
          await aggregator.getAdmin(personas.Ned.address),
        )
        const event = h.eventArgs(receipt.events?.[0])
        assert.equal(event.oracle, personas.Ned.address)
        assert.equal(event.admin, personas.Neil.address)
        assert.equal(event.newAdmin, personas.Nelly.address)
      })
    })

    describe('when the non-admin owner tries to update the admin', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Carol)
            .transferAdmin(personas.Ned.address, personas.Nelly.address),
          'revert only callable by admin',
        )
      })
    })

    describe('when the non-admin oracle tries to update the admin', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Ned)
            .transferAdmin(personas.Ned.address, personas.Nelly.address),
          'revert only callable by admin',
        )
      })
    })
  })

  describe('#acceptAdmin', () => {
    beforeEach(async () => {
      await aggregator
        .connect(personas.Carol)
        .addOracles(
          [personas.Ned.address],
          [personas.Neil.address],
          minAns,
          maxAns,
          rrDelay,
        )
      const tx = await aggregator
        .connect(personas.Neil)
        .transferAdmin(personas.Ned.address, personas.Nelly.address)
      await tx.wait()
    })

    describe('when the new admin tries to accept', () => {
      it('works', async () => {
        const tx = await aggregator
          .connect(personas.Nelly)
          .acceptAdmin(personas.Ned.address)
        const receipt = await tx.wait()
        assert.equal(
          personas.Nelly.address,
          await aggregator.getAdmin(personas.Ned.address),
        )
        const event = h.eventArgs(receipt.events?.[0])
        assert.equal(event.oracle, personas.Ned.address)
        assert.equal(event.newAdmin, personas.Nelly.address)
      })
    })

    describe('when someone other than the new admin tries to accept', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Ned).acceptAdmin(personas.Ned.address),
          'only callable by pending admin',
        )
        await matchers.evmRevert(
          aggregator.connect(personas.Neil).acceptAdmin(personas.Ned.address),
          'only callable by pending admin',
        )
      })
    })
  })

  describe('#onTokenTransfer', () => {
    it('updates the available balance', async () => {
      const originalBalance = await aggregator.availableFunds()

      await aggregator.updateAvailableFunds()

      matchers.bigNum(originalBalance, await aggregator.availableFunds())

      await link.transferAndCall(aggregator.address, deposit, '0x')

      const newBalance = await aggregator.availableFunds()
      matchers.bigNum(originalBalance.add(deposit), newBalance)
    })

    it('reverts given calldata', async () => {
      await matchers.evmRevert(
        // error message is not bubbled up by link token
        link.transferAndCall(aggregator.address, deposit, '0x12345678'),
      )
    })
  })

  describe('#requestNewRound', () => {
    beforeEach(async () => {
      await addOracles(aggregator, [personas.Neil], 1, 1, 0)

      await aggregator.connect(personas.Neil).submit(nextRound, answer)
      nextRound = nextRound + 1

      await aggregator.setRequesterPermissions(personas.Carol.address, true, 0)
    })

    it('announces a new round via log event', async () => {
      const tx = await aggregator.requestNewRound()
      const receipt = await tx.wait()
      const event = matchers.eventExists(
        receipt,
        aggregator.interface.events.NewRound,
      )

      matchers.bigNum(nextRound, h.eventArgs(event).roundId)
    })

    describe('when there is a round in progress', () => {
      beforeEach(async () => {
        await aggregator.requestNewRound()
      })

      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator.requestNewRound(),
          'prev round must be supersedable',
        )
      })

      describe('when that round has timed out', () => {
        beforeEach(async () => {
          await h.increaseTimeBy(timeout + 1, provider)
          await h.mineBlock(provider)
        })

        it('starts a new round', async () => {
          const tx = await aggregator.requestNewRound()
          const receipt = await tx.wait()
          const event = matchers.eventExists(
            receipt,
            aggregator.interface.events.NewRound,
          )
          matchers.bigNum(nextRound + 1, h.eventArgs(event).roundId)
        })
      })
    })

    describe('when there is a restart delay set', () => {
      beforeEach(async () => {
        await aggregator.setRequesterPermissions(personas.Eddy.address, true, 1)
      })

      it('reverts if a round is started before the delay', async () => {
        await aggregator.connect(personas.Eddy).requestNewRound()

        await aggregator.connect(personas.Neil).submit(nextRound, answer)
        nextRound = nextRound + 1

        // Eddy can't start because of the delay
        await matchers.evmRevert(
          aggregator.connect(personas.Eddy).requestNewRound(),
          'must delay requests',
        )
        // Carol starts a new round instead
        await aggregator.connect(personas.Carol).requestNewRound()

        // round completes
        await aggregator.connect(personas.Neil).submit(nextRound, answer)
        nextRound = nextRound + 1

        // now Eddy can start again
        await aggregator.connect(personas.Eddy).requestNewRound()
      })
    })

    describe('when all oracles have been removed and then re-added', () => {
      it('does not get stuck', async () => {
        await aggregator
          .connect(personas.Carol)
          .removeOracles([personas.Neil.address], 0, 0, 0)

        // advance a few rounds
        for (let i = 0; i < 7; i++) {
          await aggregator.requestNewRound()
          nextRound = nextRound + 1
          await h.increaseTimeBy(timeout + 1, provider)
          await h.mineBlock(provider)
        }

        await addOracles(aggregator, [personas.Neil], 1, 1, 0)
        await aggregator.connect(personas.Neil).submit(nextRound, answer)
      })
    })
  })

  describe('#setRequesterPermissions', () => {
    beforeEach(async () => {
      await addOracles(aggregator, [personas.Neil], 1, 1, 0)

      await aggregator.connect(personas.Neil).submit(nextRound, answer)
      nextRound = nextRound + 1
    })

    describe('when called by the owner', () => {
      it('allows the specified address to start new rounds', async () => {
        await aggregator.setRequesterPermissions(personas.Neil.address, true, 0)

        await aggregator.connect(personas.Neil).requestNewRound()
      })

      it('emits a log announcing the update', async () => {
        const tx = await aggregator.setRequesterPermissions(
          personas.Neil.address,
          true,
          0,
        )
        const receipt = await tx.wait()
        const event = matchers.eventExists(
          receipt,
          aggregator.interface.events.RequesterPermissionsSet,
        )
        const args = h.eventArgs(event)

        assert.equal(args.requester, personas.Neil.address)
        assert.equal(args.authorized, true)
      })

      describe('when the address is already authorized', () => {
        beforeEach(async () => {
          await aggregator.setRequesterPermissions(
            personas.Neil.address,
            true,
            0,
          )
        })

        it('does not emit a log for already authorized accounts', async () => {
          const tx = await aggregator.setRequesterPermissions(
            personas.Neil.address,
            true,
            0,
          )
          const receipt = await tx.wait()
          assert.equal(0, receipt?.logs?.length)
        })
      })

      describe('when permission is removed by the owner', () => {
        beforeEach(async () => {
          await aggregator.setRequesterPermissions(
            personas.Neil.address,
            true,
            0,
          )
        })

        it('does not allow the specified address to start new rounds', async () => {
          await aggregator.setRequesterPermissions(
            personas.Neil.address,
            false,
            0,
          )

          await matchers.evmRevert(
            aggregator.connect(personas.Neil).requestNewRound(),
            'not authorized requester',
          )
        })

        it('emits a log announcing the update', async () => {
          const tx = await aggregator.setRequesterPermissions(
            personas.Neil.address,
            false,
            0,
          )
          const receipt = await tx.wait()
          const event = matchers.eventExists(
            receipt,
            aggregator.interface.events.RequesterPermissionsSet,
          )
          const args = h.eventArgs(event)

          assert.equal(args.requester, personas.Neil.address)
          assert.equal(args.authorized, false)
        })

        it('does not emit a log for accounts without authorization', async () => {
          const tx = await aggregator.setRequesterPermissions(
            personas.Ned.address,
            false,
            0,
          )
          const receipt = await tx.wait()
          assert.equal(0, receipt?.logs?.length)
        })
      })
    })

    describe('when called by a stranger', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator
            .connect(personas.Neil)
            .setRequesterPermissions(personas.Neil.address, true, 0),
          'Only callable by owner',
        )

        await matchers.evmRevert(
          aggregator.connect(personas.Neil).requestNewRound(),
          'not authorized requester',
        )
      })
    })
  })

  describe('#oracleRoundState', () => {
    describe('when round ID 0 is passed in', () => {
      const previousSubmission = 42
      let baseFunds: any
      let minAnswers: number
      let maxAnswers: number
      let submitters: ethers.Wallet[]

      beforeEach(async () => {
        oracles = [
          personas.Neil,
          personas.Ned,
          personas.Nelly,
          personas.Nancy,
          personas.Norbert,
        ]
        minAnswers = 3
        maxAnswers = 4

        await addOracles(aggregator, oracles, minAnswers, maxAnswers, rrDelay)
        submitters = [
          personas.Nelly,
          personas.Ned,
          personas.Neil,
          personas.Nancy,
        ]
        await advanceRound(aggregator, submitters, previousSubmission)
        baseFunds = h.bigNum(deposit).sub(paymentAmount.mul(submitters.length))
        startingState = await aggregator.oracleRoundState(
          personas.Nelly.address,
          0,
        )
      })

      it('returns all of the important round information', async () => {
        const state = await aggregator.oracleRoundState(
          personas.Nelly.address,
          0,
        )

        await checkOracleRoundState(state, {
          eligibleToSubmit: true,
          roundId: 2,
          latestSubmission: previousSubmission,
          startedAt: ShouldNotBeSet,
          timeout: 0,
          availableFunds: baseFunds,
          oracleCount: oracles.length,
          paymentAmount,
        })
      })

      it('reverts if called by a contract', async () => {
        testHelper = await testHelperFactory.connect(personas.Carol).deploy()
        await matchers.evmRevert(
          testHelper.readOracleRoundState(
            aggregator.address,
            personas.Neil.address,
          ),
          'off-chain reading only',
        )
      })

      describe('when the restart delay is not enforced', () => {
        beforeEach(async () => {
          await updateFutureRounds(aggregator, {
            minAnswers,
            maxAnswers,
            restartDelay: 0,
          })
        })

        describe('< min submissions and oracle not included', () => {
          beforeEach(async () => {
            await advanceRound(aggregator, [personas.Neil])
          })

          it('is eligible to submit', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              0,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: true,
              roundId: 2,
              latestSubmission: previousSubmission,
              startedAt: ShouldBeSet,
              timeout,
              availableFunds: baseFunds.sub(paymentAmount),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })
        })

        describe('< min submissions and oracle included', () => {
          beforeEach(async () => {
            await advanceRound(aggregator, [personas.Nelly])
          })

          it('is not eligible to submit', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              0,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: false,
              roundId: 2,
              latestSubmission: answer,
              startedAt: ShouldBeSet,
              timeout,
              availableFunds: baseFunds.sub(paymentAmount),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })

          describe('and timed out', () => {
            beforeEach(async () => {
              await h.increaseTimeBy(timeout + 1, provider)
              await h.mineBlock(provider)
            })

            it('is eligible to submit', async () => {
              const state = await aggregator.oracleRoundState(
                personas.Nelly.address,
                0,
              )

              await checkOracleRoundState(state, {
                eligibleToSubmit: true,
                roundId: 3,
                latestSubmission: answer,
                startedAt: ShouldNotBeSet,
                timeout: 0,
                availableFunds: baseFunds.sub(paymentAmount),
                oracleCount: oracles.length,
                paymentAmount,
              })
            })
          })
        })

        describe('>= min sumbissions and oracle not included', () => {
          beforeEach(async () => {
            await advanceRound(aggregator, [
              personas.Neil,
              personas.Nancy,
              personas.Ned,
            ])
          })

          it('is eligible to submit', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              0,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: true,
              roundId: 2,
              latestSubmission: previousSubmission,
              startedAt: ShouldBeSet,
              timeout,
              availableFunds: baseFunds.sub(paymentAmount.mul(3)),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })
        })

        describe('>= min submissions and oracle included', () => {
          beforeEach(async () => {
            await advanceRound(aggregator, [
              personas.Neil,
              personas.Nelly,
              personas.Ned,
            ])
          })

          it('is eligible to submit', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              0,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: true,
              roundId: 3,
              latestSubmission: answer,
              startedAt: ShouldNotBeSet,
              timeout: 0,
              availableFunds: baseFunds.sub(paymentAmount.mul(3)),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })

          describe('and timed out', () => {
            beforeEach(async () => {
              await h.increaseTimeBy(timeout + 1, provider)
              await h.mineBlock(provider)
            })

            it('is eligible to submit', async () => {
              const state = await aggregator.oracleRoundState(
                personas.Nelly.address,
                0,
              )

              await checkOracleRoundState(state, {
                eligibleToSubmit: true,
                roundId: 3,
                latestSubmission: answer,
                startedAt: ShouldNotBeSet,
                timeout: 0,
                availableFunds: baseFunds.sub(paymentAmount.mul(3)),
                oracleCount: oracles.length,
                paymentAmount,
              })
            })
          })
        })

        describe('max submissions and oracle not included', () => {
          beforeEach(async () => {
            submitters = [
              personas.Neil,
              personas.Ned,
              personas.Nancy,
              personas.Norbert,
            ]
            assert.equal(
              submitters.length,
              maxAnswers,
              'precondition, please update submitters if maxAnswers changes',
            )
            await advanceRound(aggregator, submitters)
          })

          it('is eligible to submit', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              0,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: true,
              roundId: 3,
              latestSubmission: previousSubmission,
              startedAt: ShouldNotBeSet,
              timeout: 0,
              availableFunds: baseFunds.sub(paymentAmount.mul(4)),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })
        })

        describe('max submissions and oracle included', () => {
          beforeEach(async () => {
            submitters = [
              personas.Neil,
              personas.Ned,
              personas.Nelly,
              personas.Nancy,
            ]
            assert.equal(
              submitters.length,
              maxAnswers,
              'precondition, please update submitters if maxAnswers changes',
            )
            await advanceRound(aggregator, submitters)
          })

          it('is eligible to submit', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              0,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: true,
              roundId: 3,
              latestSubmission: answer,
              startedAt: ShouldNotBeSet,
              timeout: 0,
              availableFunds: baseFunds.sub(paymentAmount.mul(4)),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })
        })
      })

      describe('when the restart delay is enforced', () => {
        beforeEach(async () => {
          await updateFutureRounds(aggregator, {
            minAnswers,
            maxAnswers,
            restartDelay: maxAnswers - 1,
          })
        })

        describe('< min submissions and oracle not included', () => {
          beforeEach(async () => {
            await advanceRound(aggregator, [personas.Neil, personas.Ned])
          })

          it('is eligible to submit', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              0,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: true,
              roundId: 2,
              latestSubmission: previousSubmission,
              startedAt: ShouldBeSet,
              timeout,
              availableFunds: baseFunds.sub(paymentAmount.mul(2)),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })
        })

        describe('< min submissions and oracle included', () => {
          beforeEach(async () => {
            await advanceRound(aggregator, [personas.Neil, personas.Nelly])
          })

          it('is not eligible to submit', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              0,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: false,
              roundId: 2,
              latestSubmission: answer,
              startedAt: ShouldBeSet,
              timeout,
              availableFunds: baseFunds.sub(paymentAmount.mul(2)),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })

          describe('and timed out', () => {
            beforeEach(async () => {
              await h.increaseTimeBy(timeout + 1, provider)
              await h.mineBlock(provider)
            })

            it('is eligible to submit', async () => {
              const state = await aggregator.oracleRoundState(
                personas.Nelly.address,
                0,
              )

              await checkOracleRoundState(state, {
                eligibleToSubmit: false,
                roundId: 3,
                latestSubmission: answer,
                startedAt: ShouldNotBeSet,
                timeout: 0,
                availableFunds: baseFunds.sub(paymentAmount.mul(2)),
                oracleCount: oracles.length,
                paymentAmount,
              })
            })
          })
        })

        describe('>= min sumbissions and oracle not included', () => {
          beforeEach(async () => {
            await advanceRound(aggregator, [
              personas.Neil,
              personas.Ned,
              personas.Nancy,
            ])
          })

          it('is eligible to submit', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              0,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: true,
              roundId: 2,
              latestSubmission: previousSubmission,
              startedAt: ShouldBeSet,
              timeout,
              availableFunds: baseFunds.sub(paymentAmount.mul(3)),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })
        })

        describe('>= min submissions and oracle included', () => {
          beforeEach(async () => {
            await advanceRound(aggregator, [
              personas.Neil,
              personas.Ned,
              personas.Nelly,
            ])
          })

          it('is eligible to submit', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              0,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: false,
              roundId: 3,
              latestSubmission: answer,
              startedAt: ShouldNotBeSet,
              timeout: 0,
              availableFunds: baseFunds.sub(paymentAmount.mul(3)),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })

          describe('and timed out', () => {
            beforeEach(async () => {
              await h.increaseTimeBy(timeout + 1, provider)
              await h.mineBlock(provider)
            })

            it('is eligible to submit', async () => {
              const state = await aggregator.oracleRoundState(
                personas.Nelly.address,
                0,
              )

              await checkOracleRoundState(state, {
                eligibleToSubmit: false, // restart delay enforced
                roundId: 3,
                latestSubmission: answer,
                startedAt: ShouldNotBeSet,
                timeout: 0,
                availableFunds: baseFunds.sub(paymentAmount.mul(3)),
                oracleCount: oracles.length,
                paymentAmount,
              })
            })
          })
        })

        describe('max submissions and oracle not included', () => {
          beforeEach(async () => {
            submitters = [
              personas.Neil,
              personas.Ned,
              personas.Nancy,
              personas.Norbert,
            ]
            assert.equal(
              submitters.length,
              maxAnswers,
              'precondition, please update submitters if maxAnswers changes',
            )
            await advanceRound(aggregator, submitters, answer)
          })

          it('is not eligible to submit', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              0,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: false,
              roundId: 3,
              latestSubmission: previousSubmission,
              startedAt: ShouldNotBeSet,
              timeout: 0, // details have been deleted
              availableFunds: baseFunds.sub(paymentAmount.mul(4)),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })
        })

        describe('max submissions and oracle included', () => {
          beforeEach(async () => {
            submitters = [
              personas.Neil,
              personas.Ned,
              personas.Nelly,
              personas.Nancy,
            ]
            assert.equal(
              submitters.length,
              maxAnswers,
              'precondition, please update submitters if maxAnswers changes',
            )
            await advanceRound(aggregator, submitters, answer)
          })

          it('is not eligible to submit', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              0,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: false,
              roundId: 3,
              latestSubmission: answer,
              startedAt: ShouldNotBeSet,
              timeout: 0,
              availableFunds: baseFunds.sub(paymentAmount.mul(4)),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })
        })
      })
    })

    describe('when non-zero round ID 0 is passed in', () => {
      const answers = [0, 42, 47, 52, 57]
      let currentFunds: any

      beforeEach(async () => {
        oracles = [personas.Neil, personas.Ned, personas.Nelly]

        await addOracles(aggregator, oracles, 2, 3, rrDelay)
        startingState = await aggregator.oracleRoundState(
          personas.Nelly.address,
          0,
        )
        await advanceRound(aggregator, oracles, answers[1])
        await advanceRound(
          aggregator,
          [personas.Neil, personas.Ned],
          answers[2],
        )
        await advanceRound(aggregator, oracles, answers[3])
        await advanceRound(aggregator, [personas.Neil], answers[4])
        const submissionsSoFar = 9
        currentFunds = h
          .bigNum(deposit)
          .sub(paymentAmount.mul(submissionsSoFar))
      })

      it('returns info about previous rounds', async () => {
        const state = await aggregator.oracleRoundState(
          personas.Nelly.address,
          1,
        )

        await checkOracleRoundState(state, {
          eligibleToSubmit: false,
          roundId: 1,
          latestSubmission: answers[3],
          startedAt: ShouldBeSet,
          timeout: 0,
          availableFunds: currentFunds,
          oracleCount: oracles.length,
          paymentAmount: 0,
        })
      })

      it('returns info about previous rounds that were not submitted to', async () => {
        const state = await aggregator.oracleRoundState(
          personas.Nelly.address,
          2,
        )

        await checkOracleRoundState(state, {
          eligibleToSubmit: false,
          roundId: 2,
          latestSubmission: answers[3],
          startedAt: ShouldBeSet,
          timeout,
          availableFunds: currentFunds,
          oracleCount: oracles.length,
          paymentAmount,
        })
      })

      describe('for the current round', () => {
        describe('which has not been submitted to', () => {
          it("returns info about the current round that hasn't been submitted to", async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              4,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: true,
              roundId: 4,
              latestSubmission: answers[3],
              startedAt: ShouldBeSet,
              timeout,
              availableFunds: currentFunds,
              oracleCount: oracles.length,
              paymentAmount,
            })
          })

          it('returns info about the subsequent round', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              5,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: false,
              roundId: 5,
              latestSubmission: answers[3],
              startedAt: ShouldNotBeSet,
              timeout: 0,
              availableFunds: currentFunds,
              oracleCount: oracles.length,
              paymentAmount,
            })
          })
        })

        describe('which has been submitted to', () => {
          beforeEach(async () => {
            await aggregator.connect(personas.Nelly).submit(4, answers[4])
          })

          it("returns info about the current round that hasn't been submitted to", async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              4,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: false,
              roundId: 4,
              latestSubmission: answers[4],
              startedAt: ShouldBeSet,
              timeout,
              availableFunds: currentFunds.sub(paymentAmount),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })

          it('returns info about the subsequent round', async () => {
            const state = await aggregator.oracleRoundState(
              personas.Nelly.address,
              5,
            )

            await checkOracleRoundState(state, {
              eligibleToSubmit: true,
              roundId: 5,
              latestSubmission: answers[4],
              startedAt: ShouldNotBeSet,
              timeout: 0,
              availableFunds: currentFunds.sub(paymentAmount),
              oracleCount: oracles.length,
              paymentAmount,
            })
          })
        })
      })

      it('returns speculative info about future rounds', async () => {
        const state = await aggregator.oracleRoundState(
          personas.Nelly.address,
          6,
        )

        await checkOracleRoundState(state, {
          eligibleToSubmit: false,
          roundId: 6,
          latestSubmission: answers[3],
          startedAt: ShouldNotBeSet,
          timeout: 0,
          availableFunds: currentFunds,
          oracleCount: oracles.length,
          paymentAmount,
        })
      })
    })
  })

  describe('#getRoundData', () => {
    let latestRoundId: any
    beforeEach(async () => {
      oracles = [personas.Nelly]
      const minMax = oracles.length
      await addOracles(aggregator, oracles, minMax, minMax, rrDelay)
      await advanceRound(aggregator, oracles, answer)
      latestRoundId = await aggregator.latestRound()
    })

    it('returns the relevant round information', async () => {
      const round = await aggregator.getRoundData(latestRoundId)
      matchers.bigNum(latestRoundId, round.roundId)
      matchers.bigNum(answer, round.answer)
      const nowSeconds = new Date().valueOf() / 1000
      assert.isAbove(round.updatedAt.toNumber(), nowSeconds - 120)
      matchers.bigNum(round.updatedAt, round.startedAt)
      matchers.bigNum(latestRoundId, round.answeredInRound)
    })

    it('reverts if a round is not present', async () => {
      await matchers.evmRevert(
        aggregator.getRoundData(latestRoundId.add(1)),
        'No data present',
      )
    })
  })

  describe('#latestRoundData', () => {
    beforeEach(async () => {
      oracles = [personas.Nelly]
      const minMax = oracles.length
      await addOracles(aggregator, oracles, minMax, minMax, rrDelay)
    })

    describe('when an answer has already been received', () => {
      beforeEach(async () => {
        await advanceRound(aggregator, oracles, answer)
      })

      it('reverts if there is no latest round', async () => {
        const round = await aggregator.latestRoundData()
        const latestRoundId = await aggregator.latestRound()

        matchers.bigNum(latestRoundId, round.roundId)
        matchers.bigNum(answer, round.answer)
        const nowSeconds = new Date().valueOf() / 1000
        assert.isAbove(round.updatedAt.toNumber(), nowSeconds - 120)
        matchers.bigNum(round.updatedAt, round.startedAt)
        matchers.bigNum(latestRoundId, round.answeredInRound)
      })
    })

    it('reverts if a round is not present', async () => {
      await matchers.evmRevert(aggregator.latestRoundData(), 'No data present')
    })
  })

  describe('#setValidator', () => {
    beforeEach(async () => {
      validator = await validatorMockFactory.connect(personas.Carol).deploy()
    })

    it('emits a log event showing the validator was changed', async () => {
      const tx = await aggregator
        .connect(personas.Carol)
        .setValidator(validator.address)
      const receipt = await tx.wait()
      const eventLog = matchers.eventExists(
        receipt,
        aggregator.interface.events.ValidatorUpdated,
      )

      assert.equal(emptyAddress, h.eventArgs(eventLog).previous)
      assert.equal(validator.address, h.eventArgs(eventLog).current)

      const sameChangeTx = await aggregator
        .connect(personas.Carol)
        .setValidator(validator.address)
      const sameChangeReceipt = await sameChangeTx.wait()
      assert.equal(0, sameChangeReceipt.events?.length)
      matchers.eventDoesNotExist(
        sameChangeReceipt,
        aggregator.interface.events.ValidatorUpdated,
      )
    })

    describe('when called by a non-owner', () => {
      it('reverts', async () => {
        await matchers.evmRevert(
          aggregator.connect(personas.Neil).setValidator(validator.address),
          'Only callable by owner',
        )
      })
    })
  })

  describe('integrating with historic deviation checker', () => {
    let validator: contract.Instance<DeviationFlaggingValidatorFactory>
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

      await aggregator.connect(personas.Carol).setValidator(validator.address)

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
