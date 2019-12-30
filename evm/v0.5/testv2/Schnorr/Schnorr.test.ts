import { dssTest, tests } from './fixtures'
import { SchnorrSECP256K1Factory } from '../../src/generated/SchnorrSECP256K1Factory'
import { assert } from 'chai'
import { ethers } from 'ethers'
import * as h from '../../src/helpers'
import { makeTestProvider } from '../../src/provider'
import { Instance } from '../../src/contract'
import '../../src/extensions/ethers/BigNumber'

const schnorrSECP256K1Factory = new SchnorrSECP256K1Factory()
const provider = makeTestProvider()

// Number of points in secp256k1
const groupOrder = ethers.utils.bigNumberify(
  '0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141',
)

let defaultAccount: ethers.Wallet
beforeAll(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(provider)
  defaultAccount = rolesAndPersonas.roles.defaultAccount
})

describe('SchnorrSECP256K1', () => {
  let c: Instance<SchnorrSECP256K1Factory>
  const deployment = h.useSnapshot(provider, async () => {
    c = await schnorrSECP256K1Factory.connect(defaultAccount).deploy()
  })

  beforeEach(async () => {
    await deployment()
  })

  const secretKey = ethers.utils.bigNumberify(
    // Uniformly sampled from {0,...,groupOrder}
    '0x5d18fc9fb6494384932af3bda6fe8102c0fa7a26774e22af3993a69e2ca79565',
  )
  const publicKey = [
    // '>>>' means "computed in python"
    // >>> import py_ecc.secp256k1.secp256k1 as s
    // >>> print("'0x%x',\n'0x%x'" % tuple(s.multiply(s.G, secretKey)))
    '0x6e071bbc2060bce7bae894019d30bdf606bdc8ddc99d5023c4c73185827aeb01',
    '0x9ed10348aa5cb37be35802226259ec776119bbea355597db176c66a0f94aa183',
  ].map(ethers.utils.bigNumberify)
  const [msgHash, k] = [
    // Arbitrary values to test signature
    '0x18f224412c876d8efb2a3fa670837b5ad1347120363c2b310653f610d382729b',
    '0xd51e13c68bf56155a83e50fd9bc840e2a1847fb9b49cd206a577ecd1cd15e285',
  ].map(ethers.utils.bigNumberify)
  const kTimesG = [
    // >>> print("'0x%x',\n'0x%x'" % tuple(s.multiply(s.G, k)))
    '0x046c8644d3d376356b540e95f1727b6fd99830d53ef8af963fcc401eeb7b9f8c9f',
    'f142b3c0964202b45fb2f862843f75410ce07de04643b28b9ce04633b5fb225c',
  ].join('')
  const kTimesGAddress = ethers.utils.computeAddress(kTimesG)
  const pubKeyYParity = publicKey[1].isEven() ? 0 : 1
  const e = ethers.utils.bigNumberify(
    ethers.utils.solidityKeccak256(
      ['uint256', 'uint8', 'uint256', 'uint160'],
      [publicKey[0], pubKeyYParity ? '0x01' : '0x00', msgHash, kTimesGAddress],
    ),
  )
  const s = k.sub(e.mul(secretKey)).umod(groupOrder) // s ≡ k - e*secretKey mod groupOrder

  it('Knows a good Schnorr signature from bad', async () => {
    assert(
      publicKey[0].lt(groupOrder.shrn(1).add(ethers.constants.One)),
      'x ordinate of public key must be less than half group order.',
    )

    const checkSignature = async (signature: ethers.utils.BigNumberish) =>
      c.verifySignature(
        publicKey[0],
        pubKeyYParity,
        signature,
        msgHash,
        kTimesGAddress,
      )
    assert(await checkSignature(s), 'failed to verify good signature')
    assert(
      !(await checkSignature(s.add(ethers.constants.One))), // Corrupt signature for
      'failed to reject bad signature', //     // positive control
    )

    const gasUsed = await c.estimate.verifySignature(
      publicKey[0],
      pubKeyYParity,
      s,
      msgHash,
      kTimesGAddress,
    )
    assert.isBelow(gasUsed.toNumber(), 37500, 'burns too much gas')
  })

  it('Accepts the signatures generated on the go side', async () => {
    tests.push(dssTest)
    for (let i = 0; i < Math.min(1, tests.length); i++) {
      const numbers = tests[i].slice(0, tests[i].length - 1)
      const [msgHash, , pX, pY, sig] = numbers
        .map(h.addHexPrefix)
        .map(ethers.utils.bigNumberify)
      const rEIP55Address = ethers.utils.getAddress(tests[i].pop() ?? '')
      assert(
        await c.verifySignature(
          pX,
          pY.isEven() ? 0 : 1,
          sig,
          msgHash,
          rEIP55Address,
        ),
        'failed to verify signature constructed by golang tests',
      )
      assert(
        !(await c.verifySignature(
          pX,
          pY.isEven() ? 0 : 1,
          sig.add(ethers.constants.One),
          msgHash,
          rEIP55Address,
        )),
        'failed to reject bad signature',
      )
    }
  })
})
