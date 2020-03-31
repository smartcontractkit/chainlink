import { contract, matchers, setup, wallet as w } from '@chainlink/test-helpers'
import { ethers } from 'ethers'
import { ConcreteSignedSafeMathFactory } from '../../ethers/v0.4/ConcreteSignedSafeMathFactory'

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
  let cssm: contract.Instance<ConcreteSignedSafeMathFactory>
  let response: ethers.utils.BigNumber
  const INT256_MAX = ethers.utils.bigNumberify(
    '57896044618658097711785492504343953926634992332820282019728792003956564819967',
  )
  const INT256_MIN = ethers.utils.bigNumberify(
    '-57896044618658097711785492504343953926634992332820282019728792003956564819968',
  )
  const deployment = setup.snapshot(provider, async () => {
    cssm = await concreteSignedSafeMathFactory.connect(defaultAccount).deploy()
  })

  beforeEach(async () => {
    await deployment()
  })

  describe('#add', () => {
    describe('given a positive and positive', () => {
      it('works', async () => {
        response = await cssm.testAdd(1, 2)
        matchers.bigNum(3, response)
      })

      it('works with zero', async () => {
        response = await cssm.testAdd(INT256_MAX, 0)
        matchers.bigNum(INT256_MAX, response)
      })

      describe('when both are large enough to overflow', () => {
        it('throws', async () => {
          await matchers.evmRevert(async () => {
            response = await cssm.testAdd(INT256_MAX, 1)
          })
        })
      })
    })

    describe('given a negative and negative', () => {
      it('works', async () => {
        response = await cssm.testAdd(-1, -2)
        matchers.bigNum(-3, response)
      })

      it('works with zero', async () => {
        response = await cssm.testAdd(INT256_MIN, 0)
        matchers.bigNum(INT256_MIN, response)
      })

      describe('when both are large enough to overflow', () => {
        it('throws', async () => {
          await matchers.evmRevert(async () => {
            await cssm.testAdd(INT256_MIN, -1)
          })
        })
      })
    })

    describe('given a positive and negative', () => {
      it('works', async () => {
        response = await cssm.testAdd(1, -2)
        matchers.bigNum(-1, response)
      })
    })

    describe('given a negative and positive', () => {
      it('works', async () => {
        response = await cssm.testAdd(-1, 2)
        matchers.bigNum(1, response)
      })
    })
  })

  describe('#sub', () => {
    describe('given a positive and positive', () => {
      it('works', async () => {
        response = await cssm.testSub(2, 1)
        matchers.bigNum(1, response)
      })

      it('works with zero', async () => {
        response = await cssm.testSub(INT256_MAX, 0)
        matchers.bigNum(INT256_MAX, response)
      })
    })

    describe('given a negative and negative', () => {
      it('works', async () => {
        response = await cssm.testSub(-1, -2)
        matchers.bigNum(1, response)
      })

      it('works with zero', async () => {
        response = await cssm.testSub(INT256_MIN, 0)
        matchers.bigNum(INT256_MIN, response)
      })
    })

    describe('given a positive and negative', () => {
      it('works', async () => {
        response = await cssm.testSub(1, -2)
        matchers.bigNum(3, response)
      })
    })

    describe('given a negative and positive', () => {
      it('works', async () => {
        response = await cssm.testSub(-1, 2)
        matchers.bigNum(-3, response)
      })

      describe('when both are large enough to overflow', () => {
        it('throws', async () => {
          await matchers.evmRevert(async () => {
            response = await cssm.testSub(INT256_MIN, 1)
          })
        })
      })
    })
  })

  describe('#mul', () => {
    describe('given a positive and positive', () => {
      it('works', async () => {
        response = await cssm.testMul(1, 2)
        matchers.bigNum(2, response)
      })

      it('works with zero', async () => {
        response = await cssm.testMul(INT256_MAX, 0)
        matchers.bigNum(0, response)
      })

      describe('when both are large enough to overflow', () => {
        it('throws', async () => {
          await matchers.evmRevert(async () => {
            response = await cssm.testMul(INT256_MAX, 2)
          })
        })
      })
    })

    describe('given a negative and negative', () => {
      it('works', async () => {
        response = await cssm.testMul(-1, -2)
        matchers.bigNum(2, response)
      })

      it('works with zero', async () => {
        response = await cssm.testMul(INT256_MIN, 0)
        matchers.bigNum(0, response)
      })

      describe('when both are large enough to overflow', () => {
        it('throws', async () => {
          await matchers.evmRevert(async () => {
            await cssm.testMul(INT256_MIN, 2)
          })
        })
      })
    })

    describe('given a positive and negative', () => {
      it('works', async () => {
        response = await cssm.testMul(1, -2)
        matchers.bigNum(-2, response)
      })
    })

    describe('given a negative and positive', () => {
      it('works', async () => {
        response = await cssm.testMul(-1, 2)
        matchers.bigNum(-2, response)
      })
    })
  })

  describe('#div', () => {
    describe('given a positive and positive', () => {
      it('works', async () => {
        response = await cssm.testDiv(4, 2)
        matchers.bigNum(2, response)
      })

      it('throws when dividing by zero', async () => {
        await matchers.evmRevert(async () => {
          response = await cssm.testDiv(INT256_MAX, 0)
        })
      })
    })

    describe('given a negative and negative', () => {
      it('works', async () => {
        response = await cssm.testDiv(-4, -2)
        matchers.bigNum(2, response)
      })

      it('throws when dividing by zero', async () => {
        await matchers.evmRevert(async () => {
          await cssm.testDiv(INT256_MIN, 0)
        })
      })
    })

    describe('given a positive and negative', () => {
      it('works', async () => {
        response = await cssm.testDiv(4, -2)
        matchers.bigNum(-2, response)
      })
    })

    describe('given a negative and positive', () => {
      it('works', async () => {
        response = await cssm.testDiv(-4, 2)
        matchers.bigNum(-2, response)
      })
    })
  })
})
