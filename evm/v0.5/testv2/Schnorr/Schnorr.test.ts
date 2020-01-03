import * as f from './fixtures'
import * as h from '../../src/helpers'
import { SchnorrSECP256K1Factory } from '../../src/generated/SchnorrSECP256K1Factory'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { makeTestProvider } from '../../src/provider'
import { Instance } from '../../src/contract'
import '../../src/extensions/ethers/BigNumber'
const { bigNumberify: bn } = ethers.utils

const schnorrSECP256K1Factory = new SchnorrSECP256K1Factory()
const provider = makeTestProvider()

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

  it('Knows a good Schnorr signature from bad', async () => {
    assert(
      f.publicKey[0].lt(f.groupOrder.shrn(1).add(ethers.constants.One)),
      'x ordinate of public key must be less than half group order.',
    )

    async function checkSignature(
      signature: ethers.utils.BigNumberish,
    ): Promise<boolean> {
      return c.verifySignature(
        f.publicKey[0],
        f.pubKeyYParity,
        signature,
        f.msgHash,
        f.kTimesGAddress,
      )
    }
    assert(await checkSignature(f.s), 'failed to verify good signature')
    assert(
      !(await checkSignature(f.s.add(ethers.constants.One))), // Corrupt signature for
      'failed to reject bad signature', //     // positive control
    )

    const gasUsed = await c.estimate.verifySignature(
      f.publicKey[0],
      f.pubKeyYParity,
      f.s,
      f.msgHash,
      f.kTimesGAddress,
    )
    assert.isBelow(gasUsed.toNumber(), 37500, 'burns too much gas')
  })

  it('Accepts the signatures generated on the go side', async () => {
    f.tests.push(f.dssTest)
    for (let i = 0; i < Math.min(1, f.tests.length); i++) {
      const numbers = f.tests[i].slice(0, f.tests[i].length - 1)
      const [msgHash, , pX, pY, sig] = numbers.map(h.addHexPrefix).map(bn)
      const rEIP55Address = ethers.utils.getAddress(f.tests[i].pop() ?? '')
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
