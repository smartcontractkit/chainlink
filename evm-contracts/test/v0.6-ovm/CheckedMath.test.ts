// SPDX-License-Identifier: MIT
// Adapted from https://github.com/OpenZeppelin/openzeppelin-contracts/blob/c9630526e24ba53d9647787588a19ffaa3dd65e1/test/math/SignedSafeMath.test.js

import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { CheckedMathTestHelperFactory } from '../../ethers/v0.6-ovm/CheckedMathTestHelperFactory'

const provider = setup.provider()
const mathFactory = new CheckedMathTestHelperFactory()
let personas: setup.Personas

beforeAll(async () => {
  personas = (await setup.users(provider)).personas
})

const int256Max = h.bigNum(2).pow(255).sub(1)
const int256Min = h.bigNum(-2).pow(255)

describe('CheckedMath', () => {
  let math: contract.Instance<CheckedMathTestHelperFactory>

  const deployment = setup.snapshot(provider, async () => {
    math = await mathFactory.connect(personas.Default).deploy()
  })

  beforeEach(async () => {
    await deployment()
  })

  describe('#add', () => {
    const a = h.bigNum('1234')
    const b = h.bigNum('5678')

    it('is commutative', async () => {
      const c1 = await math.add(a, b)
      const c2 = await math.add(b, a)

      matchers.bigNum(c1.result, c2.result)
      assert.isTrue(c1.ok)
      assert.isTrue(c2.ok)
    })

    it('is commutative with big numbers', async () => {
      const c1 = await math.add(int256Max, int256Min)
      const c2 = await math.add(int256Min, int256Max)

      matchers.bigNum(c1.result, c2.result)
      assert.isTrue(c1.ok)
      assert.isTrue(c2.ok)
    })

    it('returns false when overflowing', async () => {
      const c1 = await math.add(int256Max, 1)
      const c2 = await math.add(1, int256Max)

      matchers.bigNum(0, c1.result)
      matchers.bigNum(0, c2.result)
      assert.isFalse(c1.ok)
      assert.isFalse(c2.ok)
    })

    it('returns false when underflowing', async () => {
      const c1 = await math.add(int256Min, -1)
      const c2 = await math.add(-1, int256Min)

      matchers.bigNum(0, c1.result)
      matchers.bigNum(0, c2.result)
      assert.isFalse(c1.ok)
      assert.isFalse(c2.ok)
    })
  })

  describe('#sub', () => {
    const a = h.bigNum('1234')
    const b = h.bigNum('5678')

    it('subtracts correctly if it does not overflow and the result is negative', async () => {
      const c = await math.sub(a, b)
      const expected = a.sub(b)

      matchers.bigNum(expected, c.result)
      assert.isTrue(c.ok)
    })

    it('subtracts correctly if it does not overflow and the result is positive', async () => {
      const c = await math.sub(b, a)
      const expected = b.sub(a)

      matchers.bigNum(expected, c.result)
      assert.isTrue(c.ok)
    })

    it('returns false on overflow', async () => {
      const c = await math.sub(int256Max, -1)

      matchers.bigNum(0, c.result)
      assert.isFalse(c.ok)
    })

    it('returns false on underflow', async () => {
      const c = await math.sub(int256Min, 1)

      matchers.bigNum(0, c.result)
      assert.isFalse(c.ok)
    })
  })

  describe('#mul', () => {
    const a = h.bigNum('5678')
    const b = h.bigNum('-1234')

    it('is commutative', async () => {
      const c1 = await math.mul(a, b)
      const c2 = await math.mul(b, a)

      matchers.bigNum(c1.result, c2.result)
      assert.isTrue(c1.ok)
      assert.isTrue(c2.ok)
    })

    it('multiplies by 0 correctly', async () => {
      const c = await math.mul(a, 0)

      matchers.bigNum(0, c.result)
      assert.isTrue(c.ok)
    })

    it('returns false on multiplication overflow', async () => {
      const c = await math.mul(int256Max, 2)

      matchers.bigNum(0, c.result)
      assert.isFalse(c.ok)
    })

    it('returns false when the integer minimum is negated', async () => {
      const c = await math.mul(int256Min, -1)

      matchers.bigNum(0, c.result)
      assert.isFalse(c.ok)
    })
  })

  describe('#div', () => {
    const a = h.bigNum('5678')
    const b = h.bigNum('-5678')

    it('divides correctly', async () => {
      const c = await math.div(a, b)

      matchers.bigNum(a.div(b), c.result)
      assert.isTrue(c.ok)
    })

    it('divides a 0 numerator correctly', async () => {
      const c = await math.div(0, a)

      matchers.bigNum(0, c.result)
      assert.isTrue(c.ok)
    })

    it('returns complete number result on non-even division', async () => {
      const c = await math.div(7000, 5678)

      matchers.bigNum(1, c.result)
      assert.isTrue(c.ok)
    })

    it('reverts when 0 is the denominator', async () => {
      const c = await math.div(a, 0)

      matchers.bigNum(0, c.result)
      assert.isFalse(c.ok)
    })

    it('reverts on underflow with a negative denominator', async () => {
      const c = await math.div(int256Min, -1)

      matchers.bigNum(0, c.result)
      assert.isFalse(c.ok)
    })
  })
})
