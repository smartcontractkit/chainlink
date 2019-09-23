import * as h from '../src/helpersV2'
import { assertBigNum } from '../src/matchersV2'
import { ethers } from 'ethers'
import ganache from 'ganache-core'
import { AbstractContract } from '../src/contract'
import { createFundedWallet } from '../src/wallet'

const ConcreteSignedSafeMathStatic = AbstractContract.fromArtifactName(
  'ConcreteSignedSafeMath',
).toStatic()
const ganacheProvider: any = ganache.provider()
let defaultAccount: ethers.Wallet

before(async () => {
  const { wallet } = await createFundedWallet(ganacheProvider, 0)
  defaultAccount = wallet
})

describe('SignedSafeMath', () => {
  // a version of the adder contract where we make all ABI exposed functions constant
  let adderStatic: ethers.Contract

  let response
  const INT256_MAX = ethers.utils.bigNumberify(
    '57896044618658097711785492504343953926634992332820282019728792003956564819967',
  )
  const INT256_MIN = ethers.utils.bigNumberify(
    '-57896044618658097711785492504343953926634992332820282019728792003956564819968',
  )

  beforeEach(async () => {
    adderStatic = await ConcreteSignedSafeMathStatic.deploy(defaultAccount)
  })

  describe('#add', () => {
    context('given a positive and positive', () => {
      it('works', async () => {
        response = await adderStatic.testAdd(1, 2)
        assertBigNum(3, response)
      })

      it('works with zero', async () => {
        response = await adderStatic.testAdd(INT256_MAX, 0)
        assertBigNum(INT256_MAX, response)
      })

      context('when both are large enough to overflow', async () => {
        it('throws', async () => {
          await h.assertActionThrows(async () => {
            response = await adderStatic.testAdd(INT256_MAX, 1)
          })
        })
      })
    })

    context('given a negative and negative', () => {
      it('works', async () => {
        response = await adderStatic.testAdd(-1, -2)
        assertBigNum(-3, response)
      })

      it('works with zero', async () => {
        response = await adderStatic.testAdd(INT256_MIN, 0)
        assertBigNum(INT256_MIN, response)
      })

      context('when both are large enough to overflow', async () => {
        it('throws', async () => {
          await h.assertActionThrows(async () => {
            await adderStatic.testAdd(INT256_MIN, -1)
          })
        })
      })
    })

    context('given a positive and negative', () => {
      it('works', async () => {
        response = await adderStatic.testAdd(1, -2)
        assertBigNum(-1, response)
      })
    })

    context('given a negative and positive', () => {
      it('works', async () => {
        response = await adderStatic.testAdd(-1, 2)
        assertBigNum(1, response)
      })
    })
  })
})
