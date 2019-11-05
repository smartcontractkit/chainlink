import * as h from '../src/helpersV2'
import { assertBigNum } from '../src/matchersV2'
import { ethers } from 'ethers'
import { createFundedWallet } from '../src/wallet'
import { ConcreteSignedSafeMathFactory } from '../src/generated/ConcreteSignedSafeMathFactory'
import { Instance } from '../src/contract'
import ganache from 'ganache-core'

const concreteSignedSafeMathFactory = new ConcreteSignedSafeMathFactory()
const provider = new ethers.providers.Web3Provider(ganache.provider() as any)

let defaultAccount: ethers.Wallet

beforeAll(async () => {
  const { wallet } = await createFundedWallet(provider, 0)
  defaultAccount = wallet
})

describe('SignedSafeMath', () => {
  // a version of the adder contract where we make all ABI exposed functions constant
  // TODO: submit upstream PR to support constant contract type generation
  let adder: Instance<ConcreteSignedSafeMathFactory>
  let response: ethers.utils.BigNumber
  const INT256_MAX = ethers.utils.bigNumberify(
    '57896044618658097711785492504343953926634992332820282019728792003956564819967',
  )
  const INT256_MIN = ethers.utils.bigNumberify(
    '-57896044618658097711785492504343953926634992332820282019728792003956564819968',
  )
  const deployment = h.useSnapshot(provider, async () => {
    adder = await concreteSignedSafeMathFactory.connect(defaultAccount).deploy()
  })

  beforeEach(async () => {
    await deployment()
  })

  describe('#add', () => {
    describe('given a positive and positive', () => {
      it('works', async () => {
        response = await adder.testAdd(1, 2)
        assertBigNum(3, response)
      })

      it('works with zero', async () => {
        response = await adder.testAdd(INT256_MAX, 0)
        assertBigNum(INT256_MAX, response)
      })

      describe('when both are large enough to overflow', () => {
        it('throws', async () => {
          await h.assertActionThrows(async () => {
            response = await adder.testAdd(INT256_MAX, 1)
          })
        })
      })
    })

    describe('given a negative and negative', () => {
      it('works', async () => {
        response = await adder.testAdd(-1, -2)
        assertBigNum(-3, response)
      })

      it('works with zero', async () => {
        response = await adder.testAdd(INT256_MIN, 0)
        assertBigNum(INT256_MIN, response)
      })

      describe('when both are large enough to overflow', () => {
        it('throws', async () => {
          await h.assertActionThrows(async () => {
            await adder.testAdd(INT256_MIN, -1)
          })
        })
      })
    })

    describe('given a positive and negative', () => {
      it('works', async () => {
        response = await adder.testAdd(1, -2)
        assertBigNum(-1, response)
      })
    })

    describe('given a negative and positive', () => {
      it('works', async () => {
        response = await adder.testAdd(-1, 2)
        assertBigNum(1, response)
      })
    })
  })
})
