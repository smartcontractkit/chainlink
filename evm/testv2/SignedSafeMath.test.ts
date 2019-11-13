import * as h from '../src/helpersV2'
import { assertBigNum } from '../src/matchersV2'
import { ethers } from 'ethers'
import { createFundedWallet } from '../src/wallet'
import { EthersProviderWrapper } from '../src/provider'
import { ConcreteSignedSafeMathFactory } from '../src/generated/ConcreteSignedSafeMathFactory'
import { Instance } from '../src/contract'
import env from '@nomiclabs/buidler'

const concreteSignedSafeMathFactory = new ConcreteSignedSafeMathFactory()
const provider = new EthersProviderWrapper(env.ethereum)

let defaultAccount: ethers.Wallet

beforeAll(async () => {
  const { wallet } = await createFundedWallet(provider, 0)
  defaultAccount = wallet
})

describe('SignedSafeMath', () => {
  // a version of the adder contract where we make all ABI exposed functions constant
  // TODO: submit upstream PR to support constant contract type generation
  let adder: Instance<ConcreteSignedSafeMathFactory>

  let response
  const INT256_MAX = ethers.utils.bigNumberify(
    '57896044618658097711785492504343953926634992332820282019728792003956564819967',
  )
  const INT256_MIN = ethers.utils.bigNumberify(
    '-57896044618658097711785492504343953926634992332820282019728792003956564819968',
  )

  beforeEach(async () => {
    adder = await concreteSignedSafeMathFactory.connect(defaultAccount).deploy()
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
