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
      'answerRound',
      'currentAnswer',
      'oracleCount',
      'paymentAmount',
      'removeOracle',
      'transferLINK',
      'updateAnswer',
      'updatePaymentAmount',
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

    it('updates the answer', async () => {
      const answer = 100

      assert.equal(0, await aggregator.currentAnswer.call())

      await aggregator.updateAnswer(answer, nextRound, { from: personas.Neil })

      assert.equal(answer, await aggregator.currentAnswer.call())
    })

    it('increments the answer round', async () => {
      assert.equal(0, await aggregator.answerRound.call())

      const tx = await aggregator.updateAnswer(100, nextRound, { from: personas.Neil })

      assert.equal(1, await aggregator.answerRound.call())
    })

    it('pays the oracle', async () => {
      assertBigNum(0, await link.balanceOf.call(personas.Neil))

      await aggregator.updateAnswer(100, nextRound, { from: personas.Neil })

      assertBigNum(paymentAmount, await link.balanceOf.call(personas.Neil))
    })

    context('when called by a non-oracle', async () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.updateAnswer(100, nextRound, { from: personas.Carol })
        })
      })
    })

    context('when a round is passed in higher than expected', async () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.updateAnswer(100, nextRound + 1, { from: personas.Neil })
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

  describe('#updatePaymentAmount', async () => {
    const newPaymentAmount = h.toWei('2')

    it('it updates the recorded value', async () => {
      assertBigNum(paymentAmount, await aggregator.paymentAmount.call())

      await aggregator.updatePaymentAmount(newPaymentAmount, {
        from: personas.Carol,
      })

      assertBigNum(newPaymentAmount, await aggregator.paymentAmount.call())
    })

    context('when called by anyone but the owner', async () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.updatePaymentAmount(newPaymentAmount, {
            from: personas.Ned,
          })
        })
      })
    })
  })
})
