import { contract, matchers, setup, wallet as w } from '@chainlink/test-helpers'
import { ethers } from 'ethers'
import { ConcreteSignedSafeMathFactory } from '../../ethers/v0.6-ovm/ConcreteSignedSafeMathFactory'

const concreteSignedSafeMathFactory = new ConcreteSignedSafeMathFactory()
const provider = setup.provider()

let defaultAccount: ethers.Wallet

beforeAll(async () => {
  const { wallet } = await w.createFundedWallet(provider, 0)
  defaultAccount = wallet
})

describe('SignedSafeMath', () => {
  // a version of the adder contract where we make all ABI exposed functions constant
  // TODO: submit upstream PR to support constant contract type generation
  let adder: contract.Instance<ConcreteSignedSafeMathFactory>
  let response: ethers.utils.BigNumber
  const INT256_MAX = ethers.utils.bigNumberify(
    '57896044618658097711785492504343953926634992332820282019728792003956564819967',
  )
  const INT256_MIN = ethers.utils.bigNumberify(
    '-57896044618658097711785492504343953926634992332820282019728792003956564819968',
  )
  const deployment = setup.snapshot(provider, async () => {
    adder = await concreteSignedSafeMathFactory.connect(defaultAccount).deploy()
  })

  beforeEach(async () => {
    await deployment()
  })

  describe('#add', () => {
    describe('given a positive and a positive', () => {
      it('works', async () => {
        response = await adder.testAdd(1, 2)
        matchers.bigNum(3, response)
      })

      it('works with zero', async () => {
        response = await adder.testAdd(INT256_MAX, 0)
        matchers.bigNum(INT256_MAX, response)
      })

      describe('when both are large enough to overflow', () => {
        it('throws', async () => {
          await matchers.evmRevert(async () => {
            response = await adder.testAdd(INT256_MAX, 1)
          })
        })
      })
    })

    describe('given a negative and a negative', () => {
      it('works', async () => {
        response = await adder.testAdd(-1, -2)
        matchers.bigNum(-3, response)
      })

      it('works with zero', async () => {
        response = await adder.testAdd(INT256_MIN, 0)
        matchers.bigNum(INT256_MIN, response)
      })

      describe('when both are large enough to overflow', () => {
        it('throws', async () => {
          await matchers.evmRevert(async () => {
            await adder.testAdd(INT256_MIN, -1)
          })
        })
      })
    })

    describe('given a positive and a negative', () => {
      it('works', async () => {
        response = await adder.testAdd(1, -2)
        matchers.bigNum(-1, response)
      })
    })

    describe('given a negative and a positive', () => {
      it('works', async () => {
        response = await adder.testAdd(-1, 2)
        matchers.bigNum(1, response)
      })
    })
  })

  describe('#avg', () => {
    describe('given a positive and a positive', () => {
      it('works', async () => {
        response = await adder.testAvg(2, 4)
        matchers.bigNum(3, response)
      })

      it('works with zero', async () => {
        response = await adder.testAvg(0, 4)
        matchers.bigNum(2, response)
        response = await adder.testAvg(4, 0)
        matchers.bigNum(2, response)
      })

      it('works with large numbers', async () => {
        response = await adder.testAvg(INT256_MAX, INT256_MAX)
        matchers.bigNum(INT256_MAX, response)
      })

      it('rounds towards zero', async () => {
        response = await adder.testAvg(1, 2)
        matchers.bigNum(1, response)
      })
    })

    describe('given a negative and a negative', () => {
      it('works', async () => {
        response = await adder.testAvg(-2, -4)
        matchers.bigNum(-3, response)
      })

      it('works with zero', async () => {
        response = await adder.testAvg(0, -4)
        matchers.bigNum(-2, response)
        response = await adder.testAvg(-4, 0)
        matchers.bigNum(-2, response)
      })

      it('works with large numbers', async () => {
        response = await adder.testAvg(INT256_MIN, INT256_MIN)
        matchers.bigNum(INT256_MIN, response)
      })

      it('rounds towards zero', async () => {
        response = await adder.testAvg(-1, -2)
        matchers.bigNum(-1, response)
      })
    })

    describe('given a positive and a negative', () => {
      it('works', async () => {
        response = await adder.testAvg(2, -4)
        matchers.bigNum(-1, response)
        response = await adder.testAvg(4, -2)
        matchers.bigNum(1, response)
      })

      it('works with large numbers', async () => {
        response = await adder.testAvg(INT256_MAX, -2)
        matchers.bigNum(INT256_MAX.sub(2).div(2), response)
        response = await adder.testAvg(INT256_MAX, INT256_MIN)
        matchers.bigNum(0, response)
      })

      it('rounds towards zero', async () => {
        response = await adder.testAvg(1, -4)
        matchers.bigNum(-1, response)
        response = await adder.testAvg(4, -1)
        matchers.bigNum(1, response)
      })
    })

    describe('given a negative and a positive', () => {
      it('works', async () => {
        response = await adder.testAvg(-2, 4)
        matchers.bigNum(1, response)
        response = await adder.testAvg(-4, 2)
        matchers.bigNum(-1, response)
      })

      it('works with large numbers', async () => {
        response = await adder.testAvg(INT256_MIN, 2)
        matchers.bigNum(INT256_MIN.add(2).div(2), response)
        response = await adder.testAvg(INT256_MIN, INT256_MAX)
        matchers.bigNum(0, response)
      })

      it('rounds towards zero', async () => {
        response = await adder.testAvg(-1, 4)
        matchers.bigNum(1, response)
        response = await adder.testAvg(-4, 1)
        matchers.bigNum(-1, response)
      })
    })
  })
})
