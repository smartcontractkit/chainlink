import {
  contract,
  extensions,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { VRFFactory } from '../../src/generated'
import * as f from './fixtures'
extensions.ethers.BigNumber.extend()

const { bigNumberify: bn } = ethers.utils

const big1 = bn(1)
const big2 = bn(2)
const big3 = bn(3)

function assertPointsEqual(
  x: ethers.utils.BigNumber[],
  y: ethers.utils.BigNumber[],
) {
  matchers.bigNum(x[0], y[0])
  matchers.bigNum(x[1], y[1])
}

const vrfFactory = new VRFFactory()
const provider = setup.provider()

let defaultAccount: ethers.Wallet
beforeAll(async () => {
  const users = await setup.users(provider)
  defaultAccount = users.roles.defaultAccount
})

describe('VRF', () => {
  let VRF: contract.Instance<VRFFactory>
  const deployment = setup.snapshot(provider, async () => {
    VRF = await vrfFactory.connect(defaultAccount).deploy()
  })

  beforeEach(async () => {
    await deployment()
  })

  it('Accurately calculates simple and obvious bigModExp test inputs', async () => {
    const rawExp = 3 ** 2 // Appease prettier but clarify operator precedence
    matchers.bigNum(await VRF.bigModExp(3, 2, 5), rawExp % 5)
  })

  it('accurately calculates the sum of g and 2g (i.e., 3g)', async () => {
    const projectiveResult = await VRF.projectiveECAdd(
      f.generator[0],
      f.generator[1],
      f.twiceGenerator[0],
      f.twiceGenerator[1],
    )
    const zInv = projectiveResult.z3.invm(f.fieldSize)
    const affineResult = await VRF.affineECAdd(
      f.generator,
      f.twiceGenerator,
      zInv,
    )
    assertPointsEqual(f.thriceGenerator, affineResult)
  })

  it('Accurately verifies multiplication of a point by a scalar', async () => {
    assert(await VRF.ecmulVerify(f.generator, 2, f.twiceGenerator))
  })

  it('Can compute square roots', async () => {
    matchers.bigNum(2, await VRF.squareRoot(4), '4=2^2') // 4**((fieldSize-1)/2)
  })

  it('Can compute the square of the y ordinate given the x ordinate', async () => {
    matchers.bigNum(8, await VRF.ySquared(1), '8=1^3+7')
  })

  it('Hashes to the curve with the same results as the golang code', async () => {
    let result = await VRF.hashToCurve(f.generator, 1)
    matchers.bigNum(
      bn(result[0])
        .pow(big3)
        .add(bn(7))
        .umod(f.fieldSize),
      bn(result[1])
        .pow(big2)
        .umod(f.fieldSize),
      'y^2=x^3+7',
    )
    // See golang code
    result = await VRF.hashToCurve(f.generator, 1)
    matchers.bigNum(
      result[0],
      '0x530fddd863609aa12030a07c5fdb323bb392a88343cea123b7f074883d2654c4',
      'mismatch with output from services/vrf/vrf_test.go/TestVRF_HashToCurve',
    )
    matchers.bigNum(
      result[1],
      '0x6fd4ee394bf2a3de542c0e5f3c86fc8f75b278a017701a59d69bdf5134dd6b70',
      'mismatch with output from services/vrf/vrf_test.go/TestVRF_HashToCurve',
    )
  })

  it('Correctly verifies linear combinations with generator', async () => {
    assert(
      await VRF.verifyLinearCombinationWithGenerator(
        5,
        f.twiceGenerator,
        7,
        h.pubkeyToAddress(f.seventeenTimesGenerator),
      ),
      '5*(2*g)+7*g=17*g?',
    )
  })

  it('Correctly computes full linear combinations', async () => {
    const projSum = await VRF.projectiveECAdd(
      f.eightTimesGenerator[0],
      f.eightTimesGenerator[1],
      f.nineTimesGenerator[0],
      f.nineTimesGenerator[1],
    )
    const zInv = projSum[2].invm(f.fieldSize)
    assertPointsEqual(
      f.seventeenTimesGenerator,

      // '4*(2*g)+3*(3*g)=17*g?'
      await VRF.linearCombination(
        4,
        f.twiceGenerator,
        f.eightTimesGenerator,
        3,
        f.thriceGenerator,
        f.nineTimesGenerator,
        zInv,
      ),
    )
  })

  it('Computes the same hashed scalar from curve points as the golang code', async () => {
    const scalar = await VRF.scalarFromCurve(
      f.generator,
      f.generator,
      f.generator,
      h.pubkeyToAddress(f.generator),
      f.generator,
    )
    matchers.bigNum(
      '0x2b1049accb1596a24517f96761b22600a690ee5c6b6cadae3fa522e7d95ba338',
      scalar,
      'mismatch with output from services/vrf/vrf_test.go/TestVRF_ScalarFromCurve',
    )
  })

  it('Knows a good VRF proof from bad', async () => {
    const x = big1 // "secret" key in Goldberg's notation
    const pk = f.generator
    const seed = 1
    const hash = await VRF.hashToCurve(pk, seed)
    const gamma = hash // Since gamma = x * hash = hash
    const k = big1 // "Random" nonce, ha ha
    const u = f.generator // Since u = k * generator = generator
    const v = hash // Since v = k * hash = hash
    const c = await VRF.scalarFromCurve(
      hash,
      pk,
      gamma,
      h.pubkeyToAddress(u),
      v,
    )
    const s = k.sub(c.mul(x)).umod(f.groupOrder) // s = k - c * x mod group size
    const cGamma = [
      // >>> print("'0x%x',\n'0x%x'" % tuple(s.multiply(gamma, c)))
      bn('0xa2e03a05b089db7b79cd0f6655d6af3e2d06bd0129f87f9f2155612b4e2a41d8'),
      bn('0xa1dadcabf900bdfb6484e9a4390bffa6ccd666a565a991f061faf868cc9fce8'),
    ]
    const sHash = [
      // >>> print("'0x%x',\n'0x%x'" % tuple(s.multiply(hash, signature)))
      bn('0xf82b4f9161ab41ae7c11e7deb628024ef9f5e9a0bca029f0ccb5cb534c70be31'),
      bn('0xf26e7c0b4f039ca54cfa100b3457b301acb3e0b6c690d7ea5a86f8e1c481057e'),
    ]
    const projSum = await VRF.projectiveECAdd(
      cGamma[0],
      cGamma[1],
      sHash[0],
      sHash[1],
    )
    const zInv = projSum[2].invm(f.fieldSize)

    const checkOutput = async (o: ethers.utils.BigNumberish) =>
      VRF.isValidVRFOutput(
        pk,
        gamma,
        c,
        s,
        seed,
        h.pubkeyToAddress(u),
        cGamma,
        sHash,
        zInv,
        o,
      )
    assert(!(await checkOutput(0)), 'accepted a bad proof')
    const output = ethers.utils.keccak256(
      Buffer.concat(gamma.map(x => ethers.utils.arrayify(x))),
    )

    assert(await checkOutput(bn(output)), 'rejected good proof')
  })
})
