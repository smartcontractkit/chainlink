import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'

contract('PushAggregator', () => {
  const Aggregator = artifacts.require('PushAggregator.sol')
  const personas = h.personas
  const paymentAmount = h.toWei('3')
  const deposit = h.toWei('100')

  let aggregator, link, nextRound, oracles

  beforeEach(async () => {
    link = await h.linkContract(personas.defaultAccount)
    aggregator = await Aggregator.new(link.address, paymentAmount, {
      from: personas.Carol,
    })
    await link.transfer(aggregator.address, deposit)
    assertBigNum(deposit, await link.balanceOf.call(aggregator.address))
    nextRound = 1
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(Aggregator, [
      'addOracle',
      'currentAnswer',
      'currentRound',
      'maxAnswerCount',
      'minAnswerCount',
      'oracleCount',
      'paymentAmount',
      'removeOracle',
      'setAnswerCountRange',
      'setPaymentAmount',
      'transferLINK',
      'updateAnswer',
      // Ownable methods:
      'isOwner',
      'owner',
      'transferOwnership',
    ])
  })

  describe('#updateAnswer', async () => {
    beforeEach(async () => {
      oracles = [personas.Neil, personas.Ned, personas.Nelly]
      for (const oracle of oracles) {
        await aggregator.addOracle(oracle, { from: personas.Carol })
      }
    })

    context('when the minimum number of oracles have NOT reported', async () => {
      it('pays the oracles that have reported', async () => {
        assertBigNum(0, await link.balanceOf.call(personas.Neil))

        await aggregator.updateAnswer(100, nextRound, { from: personas.Neil })

        assertBigNum(paymentAmount, await link.balanceOf.call(personas.Neil))
        assertBigNum(0, await link.balanceOf.call(personas.Ned))
        assertBigNum(0, await link.balanceOf.call(personas.Nelly))
      })

      it('updates the answer', async () => {
        const answer = 100

        assert.equal(0, await aggregator.currentAnswer.call())

        await aggregator.updateAnswer(answer, nextRound, { from: personas.Ned })
        await aggregator.updateAnswer(answer, nextRound, {
          from: personas.Nelly,
        })

        assert.equal(0, await aggregator.currentAnswer.call())
      })
    })

    context('when the minimum number of oracles have reported', async () => {
      it('updates the answer', async () => {
        assert.equal(0, await aggregator.currentAnswer.call())

        await aggregator.updateAnswer(99, nextRound, { from: personas.Neil })
        await aggregator.updateAnswer(100, nextRound, { from: personas.Ned })
        await aggregator.updateAnswer(101, nextRound, { from: personas.Nelly })

        assert.equal(100, await aggregator.currentAnswer.call())
      })

      it('announces the new answer with a log event', async () => {
        assert.equal(0, await aggregator.currentAnswer.call())

        await aggregator.updateAnswer(99, nextRound, { from: personas.Neil })
        await aggregator.updateAnswer(100, nextRound, { from: personas.Ned })
        const tx = await aggregator.updateAnswer(101, nextRound, {
          from: personas.Nelly,
        })
        const log = tx.receipt.rawLogs[0]
        const newAnswer = web3.utils.toBN(log.topics[1])

        assert.equal(100, newAnswer.toNumber())
      })
    })

    context('when an oracle submits for a round twice', async () => {
      it('reverts', async () => {
        await aggregator.updateAnswer(100, nextRound, { from: personas.Neil })

        await h.assertActionThrows(async () => {
          await aggregator.updateAnswer(100, nextRound, { from: personas.Neil })
        })
      })
    })

    context('when updated after the max answers submitted', async () => {
      beforeEach(async () => {
        await aggregator.setAnswerCountRange(1, 1, { from: personas.Carol })
        await aggregator.updateAnswer(100, nextRound, { from: personas.Neil })
      })

      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.updateAnswer(100, nextRound, { from: personas.Ned })
        })
      })
    })

    context('when a new highest round number is passed in', async () => {
      it('increments the answer round', async () => {
        assert.equal(0, await aggregator.currentRound.call())

        for (const oracle of oracles) {
          await aggregator.updateAnswer(100, nextRound, { from: oracle })
        }

        assert.equal(1, await aggregator.currentRound.call())
      })

      it('announces a new round by emitting a log', async () => {
        const tx = await aggregator.updateAnswer(100, nextRound, { from: personas.Neil })
        const log = tx.receipt.rawLogs[0]
        const roundNumber = web3.utils.toBN(log.topics[1])

        assert.equal(nextRound, roundNumber.toNumber())
      })
    })

    context('when a round is passed in higher than expected', async () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.updateAnswer(100, nextRound + 1, { from: personas.Neil })
        })
      })
    })

    context('when called by a non-oracle', async () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.updateAnswer(100, nextRound, { from: personas.Carol })
        })
      })
    })
  })

  describe('#addOracle', async () => {
    it('increases the oracle count', async () => {
      const pastCount = await aggregator.oracleCount.call()
      await aggregator.addOracle(personas.Neil, { from: personas.Carol })
      const currentCount = await aggregator.oracleCount.call()

      assertBigNum(currentCount, pastCount.add(web3.utils.toBN('1')))
    })

    it('updates the answer range', async () => {
      const pastMin = await aggregator.minAnswerCount.call()
      const pastMax = await aggregator.maxAnswerCount.call()

      await aggregator.addOracle(personas.Neil, { from: personas.Carol })

      assertBigNum(
        pastMin.add(web3.utils.toBN('1')),
        await aggregator.minAnswerCount.call(),
      )
      assertBigNum(
        pastMax.add(web3.utils.toBN('1')),
        await aggregator.maxAnswerCount.call(),
      )
    })

    context('when the oracle has already been added', async () => {
      beforeEach(async () => {
        await aggregator.addOracle(personas.Neil, { from: personas.Carol })
      })

      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.addOracle(personas.Neil, { from: personas.Carol })
        })
      })
    })

    context('when called by anyone but the owner', async () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.addOracle(personas.Neil, { from: personas.Neil })
        })
      })
    })
  })

  describe('#removeOracle', async () => {
    beforeEach(async () => {
      await aggregator.addOracle(personas.Neil, { from: personas.Carol })
    })

    it('decreases the oracle count', async () => {
      const pastCount = await aggregator.oracleCount.call()
      await aggregator.removeOracle(personas.Neil, { from: personas.Carol })
      const currentCount = await aggregator.oracleCount.call()

      assertBigNum(currentCount, pastCount.sub(web3.utils.toBN('1')))
    })

    it('updates the answer range', async () => {
      const pastMin = await aggregator.minAnswerCount.call()
      const pastMax = await aggregator.maxAnswerCount.call()

      await aggregator.removeOracle(personas.Neil, { from: personas.Carol })

      assertBigNum(
        pastMin.sub(web3.utils.toBN('1')),
        await aggregator.minAnswerCount.call(),
      )
      assertBigNum(
        pastMax.sub(web3.utils.toBN('1')),
        await aggregator.maxAnswerCount.call(),
      )
    })

    context('when the oracle is not currently added', async () => {
      beforeEach(async () => {
        await aggregator.removeOracle(personas.Neil, { from: personas.Carol })
      })

      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.removeOracle(personas.Neil, { from: personas.Carol })
        })
      })
    })

    context('when called by anyone but the owner', async () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.removeOracle(personas.Neil, { from: personas.Ned })
        })
      })
    })
  })

  describe('#transferLINK', () => {
    context('when called by the owner', () => {
      it('succeeds', async () => {
        await aggregator.transferLINK(personas.Carol, deposit, {
          from: personas.Carol,
        })

        assertBigNum(0, await link.balanceOf.call(aggregator.address))
        assertBigNum(deposit, await link.balanceOf.call(personas.Carol))
      })

      context('with a number higher than the LINK balance', () => {
        it('fails', async () => {
          await h.assertActionThrows(async () => {
            await aggregator.transferLINK(
              personas.Carol,
              deposit.add(web3.utils.toBN('1')),
              { from: personas.Carol },
            )
          })

          assertBigNum(deposit, await link.balanceOf.call(aggregator.address))
        })
      })
    })

    context('when called by a non-owner', () => {
      it('fails', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.transferLINK(personas.Carol, deposit, {
            from: personas.Eddy,
          })
        })

        assertBigNum(deposit, await link.balanceOf.call(aggregator.address))
      })
    })
  })

  describe('#setPaymentAmount', async () => {
    const newPaymentAmount = h.toWei('2')

    it('it updates the payment amount record', async () => {
      assertBigNum(paymentAmount, await aggregator.paymentAmount.call())

      await aggregator.setPaymentAmount(newPaymentAmount, {
        from: personas.Carol,
      })

      assertBigNum(newPaymentAmount, await aggregator.paymentAmount.call())
    })

    context('when called by anyone but the owner', async () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.setPaymentAmount(newPaymentAmount, {
            from: personas.Ned,
          })
        })
      })
    })
  })

  describe('#setAnswerCountRange', async () => {
    let minAnswerCount, maxAnswerCount

    beforeEach(async () => {
      oracles = [personas.Neil, personas.Ned, personas.Nelly]
      for (const oracle of oracles) {
        await aggregator.addOracle(oracle, { from: personas.Carol })
      }
      minAnswerCount = oracles.length
      maxAnswerCount = oracles.length

      assert.equal(minAnswerCount, await aggregator.minAnswerCount.call())
      assert.equal(maxAnswerCount, await aggregator.maxAnswerCount.call())
    })

    it('it updates the min and max answer counts', async () => {
      const newMin = 1
      const newMax = 2

      await aggregator.setAnswerCountRange(newMin, newMax, {
        from: personas.Carol,
      })

      assert.equal(newMin, await aggregator.minAnswerCount.call())
      assert.equal(newMax, await aggregator.maxAnswerCount.call())
    })

    context('when it is set to higher than the number or oracles', async () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.setAnswerCountRange(minAnswerCount, 4, {
            from: personas.Carol,
          })
        })
      })
    })

    context('when it sets the min higher than the max', async () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.setAnswerCountRange(3, 2, {
            from: personas.Carol,
          })
        })
      })
    })

    context('when called by anyone but the owner', async () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.setAnswerCountRange(1, 3, {
            from: personas.Ned,
          })
        })
      })
    })
  })
})
