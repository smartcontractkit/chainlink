import { ethers } from 'hardhat'
import { expect } from 'chai'
import { Signer, Contract, ContractFactory, BigNumber } from 'ethers'
import { Personas, getUsers } from '../test-helpers/setup'
import { bigNumEquals } from '../test-helpers/matchers'

let defaultAccount: Signer
let concreteSignedSafeMathFactory: ContractFactory

before(async () => {
  const personas: Personas = (await getUsers()).personas
  defaultAccount = personas.Default
  concreteSignedSafeMathFactory = await ethers.getContractFactory(
    'src/v0.6/tests/ConcreteSignedSafeMath.sol:ConcreteSignedSafeMath',
    defaultAccount,
  )
})

describe('SignedSafeMath', () => {
  // a version of the adder contract where we make all ABI exposed functions constant
  // TODO: submit upstream PR to support constant contract type generation
  let adder: Contract
  let response: BigNumber

  const INT256_MAX = BigNumber.from(
    '57896044618658097711785492504343953926634992332820282019728792003956564819967',
  )
  const INT256_MIN = BigNumber.from(
    '-57896044618658097711785492504343953926634992332820282019728792003956564819968',
  )

  beforeEach(async () => {
    adder = await concreteSignedSafeMathFactory.connect(defaultAccount).deploy()
  })

  describe('#add', () => {
    describe('given a positive and a positive', () => {
      it('works', async () => {
        response = await adder.testAdd(1, 2)
        bigNumEquals(3, response)
      })

      it('works with zero', async () => {
        response = await adder.testAdd(INT256_MAX, 0)
        bigNumEquals(INT256_MAX, response)
      })

      describe('when both are large enough to overflow', () => {
        it('throws', async () => {
          await expect(adder.testAdd(INT256_MAX, 1)).to.be.revertedWith(
            'SignedSafeMath: addition overflow',
          )
        })
      })
    })

    describe('given a negative and a negative', () => {
      it('works', async () => {
        response = await adder.testAdd(-1, -2)
        bigNumEquals(-3, response)
      })

      it('works with zero', async () => {
        response = await adder.testAdd(INT256_MIN, 0)
        bigNumEquals(INT256_MIN, response)
      })

      describe('when both are large enough to overflow', () => {
        it('throws', async () => {
          await expect(adder.testAdd(INT256_MIN, -1)).to.be.revertedWith(
            'SignedSafeMath: addition overflow',
          )
        })
      })
    })

    describe('given a positive and a negative', () => {
      it('works', async () => {
        response = await adder.testAdd(1, -2)
        bigNumEquals(-1, response)
      })
    })

    describe('given a negative and a positive', () => {
      it('works', async () => {
        response = await adder.testAdd(-1, 2)
        bigNumEquals(1, response)
      })
    })
  })

  describe('#avg', () => {
    describe('given a positive and a positive', () => {
      it('works', async () => {
        response = await adder.testAvg(2, 4)
        bigNumEquals(3, response)
      })

      it('works with zero', async () => {
        response = await adder.testAvg(0, 4)
        bigNumEquals(2, response)
        response = await adder.testAvg(4, 0)
        bigNumEquals(2, response)
      })

      it('works with large numbers', async () => {
        response = await adder.testAvg(INT256_MAX, INT256_MAX)
        bigNumEquals(INT256_MAX, response)
      })

      it('rounds towards zero', async () => {
        response = await adder.testAvg(1, 2)
        bigNumEquals(1, response)
      })
    })

    describe('given a negative and a negative', () => {
      it('works', async () => {
        response = await adder.testAvg(-2, -4)
        bigNumEquals(-3, response)
      })

      it('works with zero', async () => {
        response = await adder.testAvg(0, -4)
        bigNumEquals(-2, response)
        response = await adder.testAvg(-4, 0)
        bigNumEquals(-2, response)
      })

      it('works with large numbers', async () => {
        response = await adder.testAvg(INT256_MIN, INT256_MIN)
        bigNumEquals(INT256_MIN, response)
      })

      it('rounds towards zero', async () => {
        response = await adder.testAvg(-1, -2)
        bigNumEquals(-1, response)
      })
    })

    describe('given a positive and a negative', () => {
      it('works', async () => {
        response = await adder.testAvg(2, -4)
        bigNumEquals(-1, response)
        response = await adder.testAvg(4, -2)
        bigNumEquals(1, response)
      })

      it('works with large numbers', async () => {
        response = await adder.testAvg(INT256_MAX, -2)
        bigNumEquals(INT256_MAX.sub(2).div(2), response)
        response = await adder.testAvg(INT256_MAX, INT256_MIN)
        bigNumEquals(0, response)
      })

      it('rounds towards zero', async () => {
        response = await adder.testAvg(1, -4)
        bigNumEquals(-1, response)
        response = await adder.testAvg(4, -1)
        bigNumEquals(1, response)
      })
    })

    describe('given a negative and a positive', () => {
      it('works', async () => {
        response = await adder.testAvg(-2, 4)
        bigNumEquals(1, response)
        response = await adder.testAvg(-4, 2)
        bigNumEquals(-1, response)
      })

      it('works with large numbers', async () => {
        response = await adder.testAvg(INT256_MIN, 2)
        bigNumEquals(INT256_MIN.add(2).div(2), response)
        response = await adder.testAvg(INT256_MIN, INT256_MAX)
        bigNumEquals(0, response)
      })

      it('rounds towards zero', async () => {
        response = await adder.testAvg(-1, 4)
        bigNumEquals(1, response)
        response = await adder.testAvg(-4, 1)
        bigNumEquals(-1, response)
      })
    })
  })
})
