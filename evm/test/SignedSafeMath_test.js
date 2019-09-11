import * as h from './support/helpers'
import { assertBigNum } from './support/matchers'
const ConcreteSignedSafeMath = artifacts.require('ConcreteSignedSafeMath.sol')

contract('SignedSafeMath', () => {
  let adder, response
  const INT256_MAX = web3.utils.toBN(
    '57896044618658097711785492504343953926634992332820282019728792003956564819967',
  )
  const INT256_MIN = web3.utils.toBN(
    '-57896044618658097711785492504343953926634992332820282019728792003956564819968',
  )

  beforeEach(async () => {
    adder = await ConcreteSignedSafeMath.new()
  })

  describe('#add', () => {
    context('given a positive and positive', () => {
      it('works', async () => {
        response = await adder.testAdd.call(1, 2)
        assertBigNum(3, response)
      })

      it('works with zero', async () => {
        response = await adder.testAdd.call(INT256_MAX, 0)
        assertBigNum(INT256_MAX, response)
      })

      context('when both are large enough to overflow', async () => {
        it('throws', async () => {
          await h.assertActionThrows(async () => {
            response = await adder.testAdd.call(INT256_MAX, 1)
          })
        })
      })
    })

    context('given a negative and negative', () => {
      it('works', async () => {
        response = await adder.testAdd.call(-1, -2)
        assertBigNum(-3, response)
      })

      it('works with zero', async () => {
        response = await adder.testAdd.call(INT256_MIN, 0)
        assertBigNum(INT256_MIN, response)
      })

      context('when both are large enough to overflow', async () => {
        it('throws', async () => {
          await h.assertActionThrows(async () => {
            await adder.testAdd.call(INT256_MIN, -1)
          })
        })
      })
    })

    context('given a positive and negative', () => {
      it('works', async () => {
        response = await adder.testAdd.call(1, -2)
        assertBigNum(-1, response)
      })
    })

    context('given a negative and positive', () => {
      it('works', async () => {
        response = await adder.testAdd.call(-1, 2)
        assertBigNum(1, response)
      })
    })
  })
})
