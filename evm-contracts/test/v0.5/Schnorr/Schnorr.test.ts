import {
  contract,
  extensions,
  helpers as h,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { SchnorrSECP256K1__factory } from '../../../ethers/v0.5/factories/SchnorrSECP256K1__factory'
import * as f from './fixtures'

const { bigNumberify: bn } = ethers.utils
extensions.ethers.BigNumber.extend(ethers.utils.BigNumber)

const schnorrSECP256K1Factory = new SchnorrSECP256K1__factory()
const provider = setup.provider()

let defaultAccount: ethers.Wallet
beforeAll(async () => {
  const users = await setup.users(provider)
  defaultAccount = users.roles.defaultAccount
})

describe('SchnorrSECP256K1', () => {
  let c: contract.Instance<SchnorrSECP256K1__factory>
  const deployment = setup.snapshot(provider, async () => {
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
