import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'

import { expectRevert, time } from 'openzeppelin-test-helpers'

contract('PrepaidAggregator', () => {
  const Aggregator = artifacts.require('PrepaidAggregator.sol')
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

  let aggregator, link, nextRound, oracles

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
      // Ownable methods:
      'isOwner',
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

  describe('#updateAnswer', async () => {
    let minMax

    beforeEach(async () => {
      oracles = [personas.Neil, personas.Ned, personas.Nelly]
      for (let i = 0; i < oracles.length; i++) {
        minMax = i + 1
        await aggregator.addOracle(oracles[i], minMax, minMax, rrDelay, {
          from: personas.Carol,
        })
      }
    })

    it('updates the allocated and available funds counters', async () => {
      assertBigNum(0, await aggregator.allocatedFunds.call())

      const tx = await aggregator.updateAnswer(nextRound, answer, {
        from: personas.Neil,
      })

      assertBigNum(paymentAmount, await aggregator.allocatedFunds.call())
      const expectedAvailable = deposit.sub(paymentAmount)
      assertBigNum(expectedAvailable, await aggregator.availableFunds.call())
      const logged = h.bigNum(tx.receipt.rawLogs[1].topics[1])
      assertBigNum(expectedAvailable, logged)
    })

    it('updates the latest submission record for the oracle', async () => {
      let latest = await aggregator.latestSubmission.call(personas.Neil)
      assert.equal(0, latest[0])
      assert.equal(0, latest[1])

      const newAnswer = 427
      await aggregator.updateAnswer(nextRound, newAnswer, {
        from: personas.Neil,
      })

      latest = await aggregator.latestSubmission.call(personas.Neil)
      assert.equal(newAnswer, latest[0])
      assert.equal(nextRound, latest[1])
    })

    context('when the minimum oracles have not reported', async () => {
      it('pays the oracles that have reported', async () => {
        assertBigNum(
          0,
          await aggregator.withdrawable.call({ from: personas.Neil }),
        )

        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })

        assertBigNum(
          paymentAmount,
          await aggregator.withdrawable.call({ from: personas.Neil }),
        )
        assertBigNum(
          0,
          await aggregator.withdrawable.call({ from: personas.Ned }),
        )
        assertBigNum(
          0,
          await aggregator.withdrawable.call({ from: personas.Nelly }),
        )
      })

      it('does not update the answer', async () => {
        assert.equal(0, await aggregator.latestAnswer.call())

        // Not updated because of changes by the owner setting minAnswerCount to 3
        await aggregator.updateAnswer(nextRound, answer, { from: personas.Ned })
        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Nelly,
        })

        assert.equal(0, await aggregator.latestAnswer.call())
      })
    })

    context('when an oracle prematurely bumps the round', async () => {
      beforeEach(async () => {
        await aggregator.updateFutureRounds(paymentAmount, 2, 3, 0, timeout, {
          from: personas.Carol,
        })
        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })
      })

      it('reverts', async () => {
        await expectRevert(
          aggregator.updateAnswer(nextRound + 1, answer, {
            from: personas.Neil,
          }),
          'Not eligible to bump round',
        )
      })
    })

    context('when the minimum number of oracles have reported', async () => {
      beforeEach(async () => {
        await aggregator.updateFutureRounds(paymentAmount, 2, 3, 0, timeout, {
          from: personas.Carol,
        })
        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })
      })

      it('updates the answer with the median', async () => {
        assert.equal(0, await aggregator.latestAnswer.call())

        await aggregator.updateAnswer(nextRound, 99, { from: personas.Ned })
        assert.equal(99, await aggregator.latestAnswer.call()) // ((100+99) / 2).to_i

        await aggregator.updateAnswer(nextRound, 101, {
          from: personas.Nelly,
        })

        assert.equal(100, await aggregator.latestAnswer.call())
      })

      it('updates the updated timestamp', async () => {
        const originalTimestamp = await aggregator.latestTimestamp.call()
        assert.equal(0, originalTimestamp.toNumber())

        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Nelly,
        })

        const currentTimestamp = await aggregator.latestTimestamp.call()
        assert.isAbove(
          currentTimestamp.toNumber(),
          originalTimestamp.toNumber(),
        )
      })

      it('announces the new answer with a log event', async () => {
        const tx = await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Nelly,
        })
        const log = tx.receipt.rawLogs[0]
        const newAnswer = h.bigNum(log.topics[1])

        assert.equal(answer, newAnswer.toNumber())
      })

      it('does not set the timedout flag', async () => {
        assert.isFalse(await aggregator.getTimedOutStatus.call(nextRound))

        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Nelly,
        })

        assert.equal(
          nextRound,
          await aggregator.getOriginatingRoundOfAnswer.call(nextRound),
        )
      })
    })

    context('when an oracle submits for a round twice', async () => {
      it('reverts', async () => {
        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })

        await expectRevert(
          aggregator.updateAnswer(nextRound, answer, {
            from: personas.Neil,
          }),
          'Cannot update round reports',
        )
      })
    })

    context('when updated after the max answers submitted', async () => {
      beforeEach(async () => {
        await aggregator.updateFutureRounds(
          paymentAmount,
          1,
          1,
          rrDelay,
          timeout,
          {
            from: personas.Carol,
          },
        )
        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })
      })

      it('reverts', async () => {
        await expectRevert(
          aggregator.updateAnswer(nextRound, answer, {
            from: personas.Ned,
          }),
          'Round not currently eligible for reporting',
        )
      })
    })

    context('when a new highest round number is passed in', async () => {
      it('increments the answer round', async () => {
        assert.equal(0, await aggregator.reportingRound.call())

        for (const oracle of oracles) {
          await aggregator.updateAnswer(nextRound, answer, { from: oracle })
        }

        assert.equal(1, await aggregator.reportingRound.call())
      })

      it('announces a new round by emitting a log', async () => {
        const tx = await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })
        const log = tx.receipt.rawLogs[0]
        const roundNumber = h.bigNum(log.topics[1])
        const startedBy = h.evmWordToAddress(log.topics[2])

        assert.equal(nextRound, roundNumber.toNumber())
        assert.equal(startedBy, personas.Neil)
      })
    })

    context('when a round is passed in higher than expected', async () => {
      it('reverts', async () => {
        await expectRevert(
          aggregator.updateAnswer(nextRound + 1, answer, {
            from: personas.Neil,
          }),
          'Must report on current round',
        )
      })
    })

    context('when called by a non-oracle', async () => {
      it('reverts', async () => {
        await expectRevert(
          aggregator.updateAnswer(nextRound, answer, {
            from: personas.Carol,
          }),
          'Only updatable by whitelisted oracles',
        )
      })
    })

    context('when there are not sufficient available funds', async () => {
      beforeEach(async () => {
        await aggregator.withdrawFunds(personas.Carol, deposit, {
          from: personas.Carol,
        })
      })

      it('reverts', async () => {
        await expectRevert(
          aggregator.updateAnswer(nextRound, answer, {
            from: personas.Neil,
          }),
          'SafeMath: subtraction overflow',
        )
      })
    })

    context('when price is updated mid-round', async () => {
      const newAmount = h.toWei('50')

      it('pays the same amount to all oracles per round', async () => {
        assertBigNum(
          0,
          await aggregator.withdrawable.call({ from: personas.Neil }),
        )
        assertBigNum(
          0,
          await aggregator.withdrawable.call({ from: personas.Nelly }),
        )

        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })

        await aggregator.updateFutureRounds(
          newAmount,
          minMax,
          minMax,
          rrDelay,
          timeout,
          { from: personas.Carol },
        )

        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Nelly,
        })

        assertBigNum(
          paymentAmount,
          await aggregator.withdrawable.call({ from: personas.Neil }),
        )
        assertBigNum(
          paymentAmount,
          await aggregator.withdrawable.call({ from: personas.Nelly }),
        )
      })
    })

    context('when delay is on', async () => {
      beforeEach(async () => {
        const minMax = oracles.length
        const delay = 1

        await aggregator.updateFutureRounds(
          paymentAmount,
          minMax,
          minMax,
          delay,
          timeout,
          {
            from: personas.Carol,
          },
        )
      })

      it("does not revert on the oracle's first round", async () => {
        // Since lastUpdatedRound defaults to zero and that's the only
        // indication that an oracle hasn't responded, this test guards against
        // the situation where we don't check that and no one can start a round.
        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })
      })

      it('does revert before the delay', async () => {
        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })

        nextRound++

        await expectRevert(
          aggregator.updateAnswer(nextRound, answer, {
            from: personas.Neil,
          }),
          'Not eligible to bump round',
        )
      })
    })

    context(
      'when an oracle starts a round before the restart delay is over',
      async () => {
        beforeEach(async () => {
          await aggregator.updateFutureRounds(paymentAmount, 1, 1, 0, timeout, {
            from: personas.Carol,
          })

          oracles = [personas.Neil, personas.Ned, personas.Nelly]
          for (let i = 0; i < oracles.length; i++) {
            await aggregator.updateAnswer(nextRound, answer, {
              from: oracles[i],
            })
            nextRound++
          }

          const newDelay = 2
          // Since Ned and Nelly have answered recently, and we set the delay
          // to 2, only Nelly can answer as she is the only oracle that hasn't
          // started the last two rounds.
          await aggregator.updateFutureRounds(
            paymentAmount,
            1,
            oracles.length,
            newDelay,
            timeout,
            {
              from: personas.Carol,
            },
          )
        })

        context(
          'when called by an oracle who has not answered recently',
          async () => {
            it('does not revert', async () => {
              await aggregator.updateAnswer(nextRound, answer, {
                from: personas.Neil,
              })
            })
          },
        )

        context('when called by an oracle who answered recently', async () => {
          it('reverts', async () => {
            await expectRevert(
              aggregator.updateAnswer(nextRound, answer, {
                from: personas.Ned,
              }),
              'Round not currently eligible for reporting',
            )

            await expectRevert(
              aggregator.updateAnswer(nextRound, answer, {
                from: personas.Nelly,
              }),
              'Round not currently eligible for reporting',
            )
          })
        })
      },
    )

    context('when the price is not updated for a round', async () => {
      // For a round to timeout, it needs a previous round to pull an answer
      // from, so the second round is the earliest round that can timeout,
      // pulling its answer from the first. The start of the third round is
      // the trigger that timesout the second round, so the start of the
      // third round is the earliest we can test a timeout.

      context('on the third round or later', async () => {
        const delay = 1

        beforeEach(async () => {
          await aggregator.updateFutureRounds(
            paymentAmount,
            oracles.length,
            oracles.length,
            delay,
            timeout,
            {
              from: personas.Carol,
            },
          )

          for (const oracle of oracles) {
            await aggregator.updateAnswer(nextRound, answer, { from: oracle })
          }
          nextRound++

          await aggregator.updateAnswer(nextRound, answer, {
            from: personas.Ned,
          })
          await aggregator.updateAnswer(nextRound, answer, {
            from: personas.Nelly,
          })
          assert.equal(nextRound, await aggregator.reportingRound.call())

          await time.increase(time.duration.seconds(timeout + 1))
          nextRound++
        })

        it('allows a new round to be started', async () => {
          await aggregator.updateAnswer(nextRound, answer, {
            from: personas.Nelly,
          })
        })

        it('sets the info for the previous round', async () => {
          const previousRound = nextRound - 1
          let updated = await aggregator.getTimestamp.call(previousRound)
          let ans = await aggregator.getAnswer.call(previousRound)
          assert.equal(0, updated)
          assert.equal(0, ans)

          const tx = await aggregator.updateAnswer(nextRound, answer, {
            from: personas.Nelly,
          })
          const block = await web3.eth.getBlock(tx.receipt.blockHash)

          updated = await aggregator.getTimestamp.call(previousRound)
          ans = await aggregator.getAnswer.call(previousRound)
          assertBigNum(h.bigNum(block.timestamp), updated)
          assert.equal(answer, ans)
        })

        it('sets the previous round as timed out', async () => {
          const previousRound = nextRound - 1
          assert.isFalse(await aggregator.getTimedOutStatus.call(previousRound))

          await aggregator.updateAnswer(nextRound, answer, {
            from: personas.Nelly,
          })

          assert.isTrue(await aggregator.getTimedOutStatus.call(previousRound))
          assert.equal(
            previousRound - 1,
            await aggregator.getOriginatingRoundOfAnswer.call(previousRound),
          )
        })

        it('still respects the delay restriction', async () => {
          // expected to revert because the sender started the last round
          await expectRevert(
            aggregator.updateAnswer(nextRound, answer, {
              from: personas.Ned,
            }),
            'Round not currently eligible for reporting',
          )
        })

        it('uses the timeout set at the beginning of the round', async () => {
          await aggregator.updateFutureRounds(
            paymentAmount,
            oracles.length,
            oracles.length,
            delay,
            timeout + 100000,
            {
              from: personas.Carol,
            },
          )

          await aggregator.updateAnswer(nextRound, answer, {
            from: personas.Nelly,
          })
        })
      })

      context('earlier than the third round', async () => {
        beforeEach(async () => {
          await aggregator.updateAnswer(nextRound, answer, {
            from: personas.Neil,
          })
          await aggregator.updateAnswer(nextRound, answer, {
            from: personas.Nelly,
          })
          assert.equal(nextRound, await aggregator.reportingRound.call())

          await time.increase(time.duration.seconds(timeout + 1))

          nextRound++
          assert.equal(2, nextRound)
        })

        it('does not allow a round to be started', async () => {
          await expectRevert(
            aggregator.updateAnswer(nextRound, answer, {
              from: personas.Nelly,
            }),
            'Must have a previous answer to pull from',
          )
        })
      })
    })
  })

  describe('#getAnswer', async () => {
    const answers = [1, 10, 101, 1010, 10101, 101010, 1010101]

    beforeEach(async () => {
      await aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
        from: personas.Carol,
      })

      for (const answer of answers) {
        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })
        nextRound++
      }
    })

    it('retrieves the answer recorded for past rounds', async () => {
      for (let i = nextRound; i < nextRound; i++) {
        const answer = await aggregator.getAnswer.call(i)
        assertBigNum(h.bigNum(answers[i - 1]), answer)
      }
    })
  })

  describe('#getTimestamp', async () => {
    beforeEach(async () => {
      await aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
        from: personas.Carol,
      })

      for (let i = 0; i < 10; i++) {
        await aggregator.updateAnswer(nextRound, i, {
          from: personas.Neil,
        })
        nextRound++
      }
    })

    it('retrieves the answer recorded for past rounds', async () => {
      let lastTimestamp = h.bigNum(0)

      for (let i = 1; i < nextRound; i++) {
        const currentTimestamp = await aggregator.getTimestamp.call(i)
        assert.isAtLeast(currentTimestamp.toNumber(), lastTimestamp.toNumber())
        lastTimestamp = currentTimestamp
      }
    })
  })

  describe('#addOracle', async () => {
    it('increases the oracle count', async () => {
      const pastCount = await aggregator.oracleCount.call()
      await aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
        from: personas.Carol,
      })
      const currentCount = await aggregator.oracleCount.call()

      assertBigNum(currentCount, pastCount.add(h.bigNum(1)))
    })

    it('updates the round details', async () => {
      await aggregator.addOracle(personas.Neil, 0, 1, 0, {
        from: personas.Carol,
      })

      assertBigNum(h.bigNum(0), await aggregator.minAnswerCount.call())
      assertBigNum(h.bigNum(1), await aggregator.maxAnswerCount.call())
      assertBigNum(h.bigNum(0), await aggregator.restartDelay.call())
    })

    it('emits a log', async () => {
      const tx = await aggregator.addOracle(
        personas.Neil,
        minAns,
        maxAns,
        rrDelay,
        {
          from: personas.Carol,
        },
      )

      const added = h.evmWordToAddress(tx.receipt.rawLogs[0].topics[1])
      assertBigNum(added, personas.Neil)
    })

    context('when the oracle has already been added', async () => {
      beforeEach(async () => {
        await aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
          from: personas.Carol,
        })
      })

      it('reverts', async () => {
        await expectRevert(
          aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
            from: personas.Carol,
          }),
          'Address is already recorded as an oracle',
        )
      })
    })

    context('when called by anyone but the owner', async () => {
      it('reverts', async () => {
        await expectRevert(
          aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
            from: personas.Neil,
          }),
          'Ownable: caller is not the owner',
        )
      })
    })

    context('when an oracle gets added mid-round', async () => {
      beforeEach(async () => {
        oracles = [personas.Neil, personas.Ned]
        for (let i = 0; i < oracles.length; i++) {
          await aggregator.addOracle(oracles[i], i + 1, i + 1, rrDelay, {
            from: personas.Carol,
          })
        }

        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })

        await aggregator.addOracle(personas.Nelly, 3, 3, rrDelay, {
          from: personas.Carol,
        })
      })

      it('does not allow the oracle to update the round', async () => {
        await expectRevert(
          aggregator.updateAnswer(nextRound, answer, {
            from: personas.Nelly,
          }),
          'New oracles cannot participate in in-progress rounds',
        )
      })

      it('does allow the oracle to update future rounds', async () => {
        // complete round
        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Ned,
        })

        // now can participate in new rounds
        await aggregator.updateAnswer(nextRound + 1, answer, {
          from: personas.Nelly,
        })
      })
    })

    context('when an oracle is added after removed for a round', async () => {
      it('allows the oracle to update', async () => {
        oracles = [personas.Neil, personas.Nelly]
        for (let i = 0; i < oracles.length; i++) {
          await aggregator.addOracle(oracles[i], i + 1, i + 1, rrDelay, {
            from: personas.Carol,
          })
        }

        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })
        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Nelly,
        })
        nextRound++

        await aggregator.removeOracle(personas.Nelly, 1, 1, rrDelay, {
          from: personas.Carol,
        })

        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })
        nextRound++

        await aggregator.addOracle(personas.Nelly, 1, 1, rrDelay, {
          from: personas.Carol,
        })

        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Nelly,
        })
      })
    })

    context(
      'when an oracle is added and immediately removed mid-round',
      async () => {
        it('allows the oracle to update', async () => {
          oracles = [personas.Neil, personas.Nelly]
          for (let i = 0; i < oracles.length; i++) {
            await aggregator.addOracle(oracles[i], i + 1, i + 1, rrDelay, {
              from: personas.Carol,
            })
          }

          await aggregator.updateAnswer(nextRound, answer, {
            from: personas.Neil,
          })

          await aggregator.removeOracle(personas.Nelly, 1, 1, rrDelay, {
            from: personas.Carol,
          })
          await aggregator.addOracle(personas.Nelly, 1, 1, rrDelay, {
            from: personas.Carol,
          })

          await aggregator.updateAnswer(nextRound, answer, {
            from: personas.Nelly,
          })
        })
      },
    )

    const limit = 42
    context(`when adding more than ${limit} oracles`, async () => {
      it('reverts', async () => {
        for (let i = 0; i < limit; i++) {
          const minMax = i + 1
          const fakeAddress = web3.utils.randomHex(20)
          await aggregator.addOracle(fakeAddress, minMax, minMax, rrDelay, {
            from: personas.Carol,
          })
        }
        await expectRevert(
          aggregator.addOracle(personas.Neil, limit + 1, limit + 1, rrDelay, {
            from: personas.Carol,
          }),
          `cannot add more than ${limit} oracles`,
        )
      })
    })
  })

  describe('#removeOracle', async () => {
    beforeEach(async () => {
      await aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
        from: personas.Carol,
      })
      await aggregator.addOracle(personas.Nelly, 2, 2, rrDelay, {
        from: personas.Carol,
      })
    })

    it('decreases the oracle count', async () => {
      const pastCount = await aggregator.oracleCount.call()
      await aggregator.removeOracle(personas.Neil, minAns, maxAns, rrDelay, {
        from: personas.Carol,
      })
      const currentCount = await aggregator.oracleCount.call()

      assertBigNum(currentCount, pastCount.sub(h.bigNum(1)))
    })

    it('updates the round details', async () => {
      await aggregator.removeOracle(personas.Neil, 0, 1, 0, {
        from: personas.Carol,
      })

      assertBigNum(h.bigNum(0), await aggregator.minAnswerCount.call())
      assertBigNum(h.bigNum(1), await aggregator.maxAnswerCount.call())
      assertBigNum(h.bigNum(0), await aggregator.restartDelay.call())
    })

    it('emits a log', async () => {
      const tx = await aggregator.removeOracle(
        personas.Neil,
        minAns,
        maxAns,
        rrDelay,
        {
          from: personas.Carol,
        },
      )

      const added = h.evmWordToAddress(tx.receipt.rawLogs[0].topics[1])
      assertBigNum(added, personas.Neil)
    })

    context('when the oracle is not currently added', async () => {
      beforeEach(async () => {
        await aggregator.removeOracle(personas.Neil, minAns, maxAns, rrDelay, {
          from: personas.Carol,
        })
      })

      it('reverts', async () => {
        await expectRevert(
          aggregator.removeOracle(personas.Neil, minAns, maxAns, rrDelay, {
            from: personas.Carol,
          }),
          'Address is not a whitelisted oracle',
        )
      })
    })

    context('when removing the last oracle', async () => {
      it('does not revert', async () => {
        await aggregator.removeOracle(personas.Neil, minAns, maxAns, rrDelay, {
          from: personas.Carol,
        })

        await aggregator.removeOracle(personas.Nelly, 0, 0, 0, {
          from: personas.Carol,
        })
      })
    })

    context('when called by anyone but the owner', async () => {
      it('reverts', async () => {
        await expectRevert(
          aggregator.removeOracle(personas.Neil, 0, 0, rrDelay, {
            from: personas.Ned,
          }),
          'Ownable: caller is not the owner',
        )
      })
    })

    context('when an oracle gets removed mid-round', async () => {
      beforeEach(async () => {
        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Neil,
        })

        await aggregator.removeOracle(personas.Nelly, 1, 1, rrDelay, {
          from: personas.Carol,
        })
      })

      it('is allowed to finish that round', async () => {
        await aggregator.updateAnswer(nextRound, answer, {
          from: personas.Nelly,
        })
        nextRound++

        // cannot participate in future rounds
        await expectRevert(
          aggregator.updateAnswer(nextRound, answer, {
            from: personas.Nelly,
          }),
          'Oracle has been removed from whitelist',
        )
      })
    })
  })

  describe('#withdrawFunds', () => {
    context('when called by the owner', () => {
      it('succeeds', async () => {
        await aggregator.withdrawFunds(personas.Carol, deposit, {
          from: personas.Carol,
        })

        assertBigNum(0, await aggregator.availableFunds.call())
        assertBigNum(deposit, await link.balanceOf.call(personas.Carol))
      })

      it('does not let withdrawals happen multiple times', async () => {
        await aggregator.withdrawFunds(personas.Carol, deposit, {
          from: personas.Carol,
        })

        await expectRevert(
          aggregator.withdrawFunds(personas.Carol, deposit, {
            from: personas.Carol,
          }),
          'Insufficient funds',
        )
      })

      context('with a number higher than the available LINK balance', () => {
        beforeEach(async () => {
          await aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
            from: personas.Carol,
          })
          await aggregator.updateAnswer(nextRound, answer, {
            from: personas.Neil,
          })
        })

        it('fails', async () => {
          await expectRevert(
            aggregator.withdrawFunds(personas.Carol, deposit, {
              from: personas.Carol,
            }),
            'Insufficient funds',
          )

          assertBigNum(
            deposit.sub(paymentAmount),
            await aggregator.availableFunds.call(),
          )
        })
      })
    })

    context('when called by a non-owner', () => {
      it('fails', async () => {
        await expectRevert(
          aggregator.withdrawFunds(personas.Carol, deposit, {
            from: personas.Eddy,
          }),
          'Ownable: caller is not the owner',
        )

        assertBigNum(deposit, await aggregator.availableFunds.call())
      })
    })
  })

  describe('#updateFutureRounds', async () => {
    let minAnswerCount, maxAnswerCount
    const newPaymentAmount = h.toWei('2')
    const newMin = 1
    const newMax = 3
    const newDelay = 2

    beforeEach(async () => {
      oracles = [personas.Neil, personas.Ned, personas.Nelly]
      for (let i = 0; i < oracles.length; i++) {
        const minMax = i + 1
        await aggregator.addOracle(oracles[i], minMax, minMax, rrDelay, {
          from: personas.Carol,
        })
      }
      minAnswerCount = oracles.length
      maxAnswerCount = oracles.length

      assertBigNum(paymentAmount, await aggregator.paymentAmount.call())
      assert.equal(minAnswerCount, await aggregator.minAnswerCount.call())
      assert.equal(maxAnswerCount, await aggregator.maxAnswerCount.call())
    })

    it('updates the min and max answer counts', async () => {
      await aggregator.updateFutureRounds(
        newPaymentAmount,
        newMin,
        newMax,
        newDelay,
        timeout,
        {
          from: personas.Carol,
        },
      )

      assertBigNum(newPaymentAmount, await aggregator.paymentAmount.call())
      assertBigNum(h.bigNum(newMin), await aggregator.minAnswerCount.call())
      assertBigNum(h.bigNum(newMax), await aggregator.maxAnswerCount.call())
      assertBigNum(h.bigNum(newDelay), await aggregator.restartDelay.call())
    })

    it('emits a log announcing the new round details', async () => {
      const tx = await aggregator.updateFutureRounds(
        paymentAmount,
        newMin,
        newMax,
        newDelay,
        timeout,
        {
          from: personas.Carol,
        },
      )

      const round = h.parseAggregatorRoundLog(tx.receipt.rawLogs[0])

      assertBigNum(paymentAmount, round.paymentAmount)
      assertBigNum(h.bigNum(newMin), round.minAnswerCount)
      assertBigNum(h.bigNum(newMax), round.maxAnswerCount)
      assertBigNum(h.bigNum(newDelay), round.restartDelay)
      assertBigNum(h.bigNum(timeout), round.timeout)
    })

    context('when it is set to higher than the number or oracles', async () => {
      it('reverts', async () => {
        await expectRevert(
          aggregator.updateFutureRounds(
            paymentAmount,
            minAnswerCount,
            4,
            rrDelay,
            timeout,
            {
              from: personas.Carol,
            },
          ),
          'Cannot have the answer max higher oracle count',
        )
      })
    })

    context('when it sets the min higher than the max', async () => {
      it('reverts', async () => {
        await expectRevert(
          aggregator.updateFutureRounds(paymentAmount, 3, 2, rrDelay, timeout, {
            from: personas.Carol,
          }),
          'Cannot have the answer minimum higher the max',
        )
      })
    })

    context('when delay equal or greater the oracle count', async () => {
      it('reverts', async () => {
        await expectRevert(
          aggregator.updateFutureRounds(paymentAmount, 1, 1, 3, timeout, {
            from: personas.Carol,
          }),
          'Restart delay must be less than oracle count',
        )
      })
    })

    context('when called by anyone but the owner', async () => {
      it('reverts', async () => {
        await expectRevert(
          aggregator.updateFutureRounds(paymentAmount, 1, 3, rrDelay, timeout, {
            from: personas.Ned,
          }),
          'caller is not the owner',
        )
      })
    })
  })

  describe('#updateAvailableFunds', async () => {
    it('checks the LINK token to see if any additional funds are available', async () => {
      const originalBalance = await aggregator.availableFunds.call()

      await aggregator.updateAvailableFunds()

      assertBigNum(originalBalance, await aggregator.availableFunds.call())

      await link.transfer(aggregator.address, deposit)
      await aggregator.updateAvailableFunds()

      const newBalance = await aggregator.availableFunds.call()
      assertBigNum(originalBalance.add(deposit), newBalance)
    })

    it('removes allocated funds from the available balance', async () => {
      const originalBalance = await aggregator.availableFunds.call()

      await aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
        from: personas.Carol,
      })
      await aggregator.updateAnswer(nextRound, answer, {
        from: personas.Neil,
      })
      await link.transfer(aggregator.address, deposit)
      await aggregator.updateAvailableFunds()

      const expected = originalBalance.add(deposit).sub(paymentAmount)
      const newBalance = await aggregator.availableFunds.call()
      assertBigNum(expected, newBalance)
    })

    it('emits a log', async () => {
      await link.transfer(aggregator.address, deposit)

      const tx = await aggregator.updateAvailableFunds()

      const reportedBalance = h.bigNum(tx.receipt.rawLogs[0].topics[1])
      assertBigNum(await aggregator.availableFunds.call(), reportedBalance)
    })
  })

  describe('#withdraw', async () => {
    beforeEach(async () => {
      await aggregator.addOracle(personas.Neil, minAns, maxAns, rrDelay, {
        from: personas.Carol,
      })
      await aggregator.updateAnswer(nextRound, answer, {
        from: personas.Neil,
      })
    })

    it('transfers LINK to the caller', async () => {
      const originalBalance = await link.balanceOf.call(aggregator.address)
      assertBigNum(0, await link.balanceOf.call(personas.Neil))

      await aggregator.withdraw(personas.Neil, paymentAmount, {
        from: personas.Neil,
      })

      assertBigNum(
        originalBalance.sub(paymentAmount),
        await link.balanceOf.call(aggregator.address),
      )
      assertBigNum(paymentAmount, await link.balanceOf.call(personas.Neil))
    })

    it('decrements the allocated funds counter', async () => {
      const originalAllocation = await aggregator.allocatedFunds.call()

      await aggregator.withdraw(personas.Neil, paymentAmount, {
        from: personas.Neil,
      })

      assertBigNum(
        originalAllocation.sub(paymentAmount),
        await aggregator.allocatedFunds.call(),
      )
    })

    context('when the caller withdraws more than they have', async () => {
      it('reverts', async () => {
        await expectRevert(
          aggregator.withdraw(personas.Neil, paymentAmount.add(h.bigNum(1)), {
            from: personas.Neil,
          }),
          'Insufficient balance',
        )
      })
    })
  })
})
