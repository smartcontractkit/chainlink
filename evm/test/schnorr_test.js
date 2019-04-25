import * as h from './support/helpers'

const SchnorrSECP256K1 = artifacts.require('SchnorrSECP256K1')

const BN = web3.utils.BN
const hexToBN = s => new BN(s.replace(/^0[xX]/, ''), 16) // Construct BN from hex
const groupOrder = hexToBN(
  // Number of points in secp256k1
  '0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141'
)
const bigOne = new BN(1)

// Returns d as a 0x-hex string, left-padded with zeros to l bits.
const toHex = (d, l) =>
  // Make sure there's just one 0x prefix
  web3.utils.padLeft(d.toString(16), l / 4).replace(/^(0[xX])*/, '0x')

const bottom160bits = bigOne.shln(160).sub(bigOne)
const bottom256bits = bigOne.shln(257).sub(bigOne)

// Returns the EIP55-capitalized ethereum address for this secp256k1 public key
const toAddress = (x, y) => {
  const k256 = web3.utils.soliditySha3(toHex(x, 256), toHex(y, 256))
  return web3.utils.toChecksumAddress(k256.slice(k256.length - 40))
}

contract('SchnorrSECP256K1', async accounts => {
  let c
  beforeEach(async () => {
    c = await h.deploy('SchnorrSECP256K1.sol')
  })
  const secretKey = hexToBN(
    // Uniformly sampled from {0,...,groupOrder}
    '0x5d18fc9fb6494384932af3bda6fe8102c0fa7a26774e22af3993a69e2ca79565'
  )
  const publicKey = [
    // '>>>' means "computed in python"
    // >>> import py_ecc.secp256k1.secp256k1 as s
    // >>> print("'0x%x',\n'0x%x'" % tuple(s.multiply(s.G, secretKey)))
    '0x6e071bbc2060bce7bae894019d30bdf606bdc8ddc99d5023c4c73185827aeb01',
    '0x9ed10348aa5cb37be35802226259ec776119bbea355597db176c66a0f94aa183'
  ].map(hexToBN)
  const [msgHash, k] = [
    // Arbitrary values to test signature
    '0x18f224412c876d8efb2a3fa670837b5ad1347120363c2b310653f610d382729b',
    '0xd51e13c68bf56155a83e50fd9bc840e2a1847fb9b49cd206a577ecd1cd15e285'
  ].map(hexToBN)
  const kTimesG = [
    // >>> print("'0x%x',\n'0x%x'" % tuple(s.multiply(s.G, k)))
    '0x6c8644d3d376356b540e95f1727b6fd99830d53ef8af963fcc401eeb7b9f8c9f',
    '0xf142b3c0964202b45fb2f862843f75410ce07de04643b28b9ce04633b5fb225c'
  ].map(hexToBN)
  const kTimesGAddress = toAddress(...kTimesG)
  const pubKeyYParity = publicKey[1].isEven() ? 0 : 1
  const e = hexToBN(
    web3.utils.soliditySha3(
      toHex(publicKey[0], 256),
      toHex(pubKeyYParity ? '0x01' : '0x00', 8),
      toHex(msgHash, 256),
      toHex(kTimesGAddress, 160)
    )
  )
  const s = k.sub(e.mul(secretKey)).umod(groupOrder) // s ≡ k - e*secretKey mod groupOrder
  it('Knows a good Schnorr signature from bad', async () => {
    assert(
      publicKey[0].lt(groupOrder.shrn(1).add(bigOne)),
      'x ordinate of public key must be less than half group order.'
    )
    const checkSignature = async s =>
      c.verifySignature.call(
        publicKey[0],
        pubKeyYParity,
        s,
        msgHash,
        kTimesGAddress
      )
    assert(await checkSignature(s), 'failed to verify good signature')
    assert(
      !(await checkSignature(s.add(bigOne))), // Corrupt signature for
      'failed to reject bad signature' //     // positive control
    )
    const gasUsed = await c.verifySignature.estimateGas(
      publicKey[0],
      pubKeyYParity,
      s,
      msgHash,
      kTimesGAddress
    )
    assert.isBelow(gasUsed, 37000, 'burns too much gas')
  })
  it('Accepts the signatures generated on the go side', async () => {
    const tests = require('./tests')
    const dssTest = require('./dssTest')
    tests.push(dssTest)
    for (let i = 0; i < Math.min(1, tests.length); i++) {
      const numbers = tests[i].slice(0, tests[i].length - 1)
      const [msgHash, secret, pX, pY, sig] = numbers.map(hexToBN)
      const rEIP55Address = web3.utils.toChecksumAddress(tests[i].pop())
      secret.and(bigOne) // shut linter up about unused variable
      assert(
        await c.verifySignature.call(
          pX,
          pY.isEven() ? 0 : 1,
          sig,
          msgHash,
          rEIP55Address
        ),
        'failed to verify signature constructed by golang tests'
      )
      assert(
        !(await c.verifySignature.call(
          pX,
          pY.isEven() ? 0 : 1,
          sig.add(bigOne),
          msgHash,
          rEIP55Address
        )),
        'failed to reject bad signature'
      )
    }
  })
})
