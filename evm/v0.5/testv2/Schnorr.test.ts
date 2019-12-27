import { schnorrTests, dssTests } from '../test/support/schnorr_constants'
import * as h from '../src/helpers'
import { SchnorrSECP256K1Factory } from '../src/generated/SchnorrSECP256K1Factory'
import { makeTestProvider } from '../src/provider'
import { assert } from 'chai'
import { Instance } from '../src/contract'
import BN from 'bn.js'

const schnorrSECP256K1 = new SchnorrSECP256K1Factory()

const provider = makeTestProvider()
let roles: h.Roles

beforeAll(async () => {
  const rolesAndPersonas = await h.initializeRolesAndPersonas(provider)
  roles = rolesAndPersonas.roles
})

const groupOrder = h.hexToBN(
  // Number of points in secp256k1
  '0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141',
)

describe('SchnorrSECP256K1', () => {
  let c: Instance<SchnorrSECP256K1Factory>

  const deployment = h.useSnapshot(provider, async () => {
    c = await schnorrSECP256K1.connect(roles.defaultAccount).deploy()
  })

  beforeEach(async () => {
    await deployment()
  })

  const secretKey = h.hexToBN(
    // Uniformly sampled from {0,...,groupOrder}
    '0x5d18fc9fb6494384932af3bda6fe8102c0fa7a26774e22af3993a69e2ca79565',
  )
  const publicKey = [
    // '>>>' means "computed in python"
    // >>> import py_ecc.secp256k1.secp256k1 as s
    // >>> print("'0x%x',\n'0x%x'" % tuple(s.multiply(s.G, secretKey)))
    '0x6e071bbc2060bce7bae894019d30bdf606bdc8ddc99d5023c4c73185827aeb01',
    '0x9ed10348aa5cb37be35802226259ec776119bbea355597db176c66a0f94aa183',
  ].map(h.hexToBN)
  const [msgHash, k] = [
    // Arbitrary values to test signature
    '0x18f224412c876d8efb2a3fa670837b5ad1347120363c2b310653f610d382729b',
    '0xd51e13c68bf56155a83e50fd9bc840e2a1847fb9b49cd206a577ecd1cd15e285',
  ].map(h.hexToBN)
  const kTimesG = [
    // >>> print("'0x%x',\n'0x%x'" % tuple(s.multiply(s.G, k)))
    '0x6c8644d3d376356b540e95f1727b6fd99830d53ef8af963fcc401eeb7b9f8c9f',
    '0xf142b3c0964202b45fb2f862843f75410ce07de04643b28b9ce04633b5fb225c',
  ].map(h.hexToBN)
  const kTimesGAddress = h.toAddress(kTimesG[0], kTimesG[1])
  const pubKeyYParity = publicKey[1].isEven() ? 0 : 1
  const e = h.hexToBN(
    web3.utils.soliditySha3(
      h.toPaddedHex(publicKey[0], 256),
      h.toPaddedHex(pubKeyYParity ? h.hexToBN('0x01') : h.hexToBN('0x00'), 8),
      h.toPaddedHex(msgHash, 256),
      h.toPaddedHex(kTimesGAddress, 160),
    ),
  )
  const s = k.sub(e.mul(secretKey)).umod(groupOrder) // s ≡ k - e*secretKey mod groupOrder
  it('Knows a good Schnorr signature from bad', async () => {
    assert(
      publicKey[0].lt(groupOrder.shrn(1).add(h.bigOne)),
      'x ordinate of public key must be less than half group order.',
    )
    const checkSignature = async (s: BN) =>
      c.verifySignature(
        publicKey[0].toString(),
        pubKeyYParity,
        s.toString(),
        msgHash.toString(),
        kTimesGAddress,
      )
    assert(await checkSignature(s), 'failed to verify good signature')
    assert(
      !(await checkSignature(s.add(h.bigOne))), // Corrupt signature for
      'failed to reject bad signature', //     // positive control
    )
    const gasUsed = await c.estimate.verifySignature(
      publicKey[0].toString(),
      pubKeyYParity,
      s.toString(),
      msgHash.toString(),
      kTimesGAddress,
    )
    assert.isTrue(gasUsed.lt(37500), 'burns too much gas')
  })
  it('Accepts the signatures generated on the go side', async () => {
    const testss = [dssTests, ...schnorrTests]
    for (let i = 0; i < Math.min(1, testss.length); i++) {
      const numbers = testss[i].slice(0, testss[i].length - 1)
      const [msgHash, , pX, pY, sig] = numbers.map(h.hexToBN)
      const rEIP55Address = web3.utils.toChecksumAddress(testss[i].pop())
      assert(
        await c.verifySignature(
          pX.toString(),
          pY.isEven() ? 0 : 1,
          sig.toString(),
          msgHash.toString(),
          rEIP55Address,
        ),
        'failed to verify signature constructed by golang tests',
      )
      assert(
        !(await c.verifySignature(
          pX.toString(),
          pY.isEven() ? 0 : 1,
          sig.add(h.bigOne).toString(),
          msgHash.toString(),
          rEIP55Address,
        )),
        'failed to reject bad signature',
      )
    }
  })
})
