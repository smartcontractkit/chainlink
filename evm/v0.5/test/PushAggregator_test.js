import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'
const Aggregator = artifacts.require('PushAggregator.sol')

const personas = h.personas
const defaultAccount = h.personas.Default
const jobId1 =
  '0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000001'
let aggregator, link

contract('PushAggregator', () => {

  beforeEach(async () => {
    link = await h.linkContract(defaultAccount)
    aggregator = await Aggregator.new(link.address, { from: personas.Carol })
  })

  it('has a limited public interface', () => {
    h.checkPublicABI(Aggregator, [
      'addOracle',
      'answerRound',
      'currentAnswer',
      'oracleCount',
      'removeOracle',
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
      await aggregator.addOracle(personas.Neil, { from: personas.Carol })
    })

    it('updates the answer', async () => {
      const answer = 100

      assert.equal(0, await aggregator.currentAnswer.call())

      await aggregator.updateAnswer(answer, { from: personas.Neil })

      assert.equal(answer, await aggregator.currentAnswer.call())
    })

    it('increments the answer round', async () => {
      assert.equal(0, await aggregator.answerRound.call())

      await aggregator.updateAnswer(100, { from: personas.Neil })

      assert.equal(1, await aggregator.answerRound.call())
    })

    context('when called by a non-oracle', async () => {
      it('reverts', async () => {
        await h.assertActionThrows(async () => {
          await aggregator.updateAnswer(100, { from: personas.Carol })
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
    const deposit = h.toWei('100')

    beforeEach(async () => {
      await link.transfer(aggregator.address, deposit)
      assertBigNum(deposit, await link.balanceOf.call(aggregator.address))
    })

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
            await aggregator.transferLINK(personas.Carol, deposit.add(web3.utils.toBN('1')), {
              from: personas.Carol,
            })
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
})
