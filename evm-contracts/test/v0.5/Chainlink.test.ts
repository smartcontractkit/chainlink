import {
  contract,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { ContractReceipt } from 'ethers/contract'
import { ChainlinkTestHelper__factory } from '../../ethers/v0.5/factories/ChainlinkTestHelper__factory'

const chainlinkFactory = new ChainlinkTestHelper__factory()
const provider = setup.provider()

let defaultAccount: ethers.Wallet
beforeAll(async () => {
  defaultAccount = await setup
    .users(provider)
    .then((x) => x.roles.defaultAccount)
})

describe('Chainlink', () => {
  let cl: contract.Instance<ChainlinkTestHelper__factory>
  let clEvents: contract.Instance<
    ChainlinkTestHelper__factory
  >['interface']['events']

  const deployment = setup.snapshot(provider, async () => {
    cl = await chainlinkFactory.connect(defaultAccount).deploy()
    clEvents = cl.interface.events
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(chainlinkFactory, [
      'add',
      'addBytes',
      'addInt',
      'addStringArray',
      'addUint',
      'closeEvent',
      'setBuffer',
    ])
  })

  function getPayloadFrom(receipt: ContractReceipt) {
    const { payload } = h.eventArgs(
      h.findEventIn(receipt, clEvents.RequestData),
    )
    const decoded = h.decodeDietCBOR(payload)
    return decoded
  }

  describe('#close', () => {
    it('handles empty payloads', async () => {
      const tx = await cl.closeEvent()
      const receipt = await tx.wait()
      const decoded = getPayloadFrom(receipt)

      assert.deepEqual(decoded, {})
    })
  })

  describe('#setBuffer', () => {
    it('emits the buffer', async () => {
      await cl.setBuffer('0xA161616162')
      const tx = await cl.closeEvent()
      const receipt = await tx.wait()
      const decoded = getPayloadFrom(receipt)

      assert.deepEqual(decoded, { a: 'b' })
    })
  })

  describe('#add', () => {
    it('stores and logs keys and values', async () => {
      await cl.add('first', 'word!!')
      const tx = await cl.closeEvent()
      const receipt = await tx.wait()
      const decoded = getPayloadFrom(receipt)

      assert.deepEqual(decoded, { first: 'word!!' })
    })

    it('handles two entries', async () => {
      await cl.add('first', 'uno')
      await cl.add('second', 'dos')
      const tx = await cl.closeEvent()
      const receipt = await tx.wait()
      const decoded = getPayloadFrom(receipt)

      assert.deepEqual(decoded, {
        first: 'uno',
        second: 'dos',
      })
    })
  })

  describe('#addBytes', () => {
    it('stores and logs keys and values', async () => {
      await cl.addBytes('first', '0xaabbccddeeff')
      const tx = await cl.closeEvent()
      const receipt = await tx.wait()
      const decoded = getPayloadFrom(receipt)

      const expected = h.hexToBuf('0xaabbccddeeff')
      assert.deepEqual(decoded, { first: expected })
    })

    it('handles two entries', async () => {
      await cl.addBytes('first', '0x756E6F')
      await cl.addBytes('second', '0x646F73')
      const tx = await cl.closeEvent()
      const receipt = await tx.wait()
      const decoded = getPayloadFrom(receipt)

      const expectedFirst = h.hexToBuf('0x756E6F')
      const expectedSecond = h.hexToBuf('0x646F73')
      assert.deepEqual(decoded, {
        first: expectedFirst,
        second: expectedSecond,
      })
    })

    it('handles strings', async () => {
      await cl.addBytes('first', ethers.utils.toUtf8Bytes('apple'))
      const tx = await cl.closeEvent()
      const receipt = await tx.wait()
      const decoded = getPayloadFrom(receipt)

      const expected = ethers.utils.toUtf8Bytes('apple')
      assert.deepEqual(decoded, { first: expected })
    })
  })

  describe('#addInt', () => {
    it('stores and logs keys and values', async () => {
      await cl.addInt('first', 1)
      const tx = await cl.closeEvent()
      const receipt = await tx.wait()
      const decoded = getPayloadFrom(receipt)

      assert.deepEqual(decoded, { first: 1 })
    })

    it('handles two entries', async () => {
      await cl.addInt('first', 1)
      await cl.addInt('second', 2)
      const tx = await cl.closeEvent()
      const receipt = await tx.wait()
      const decoded = getPayloadFrom(receipt)

      assert.deepEqual(decoded, {
        first: 1,
        second: 2,
      })
    })
  })

  describe('#addUint', () => {
    it('stores and logs keys and values', async () => {
      await cl.addUint('first', 1)
      const tx = await cl.closeEvent()
      const receipt = await tx.wait()
      const decoded = getPayloadFrom(receipt)

      assert.deepEqual(decoded, { first: 1 })
    })

    it('handles two entries', async () => {
      await cl.addUint('first', 1)
      await cl.addUint('second', 2)
      const tx = await cl.closeEvent()
      const receipt = await tx.wait()
      const decoded = getPayloadFrom(receipt)

      assert.deepEqual(decoded, {
        first: 1,
        second: 2,
      })
    })
  })

  describe('#addStringArray', () => {
    it('stores and logs keys and values', async () => {
      await cl.addStringArray('word', [
        ethers.utils.formatBytes32String('seinfeld'),
        ethers.utils.formatBytes32String('"4"'),
        ethers.utils.formatBytes32String('LIFE'),
      ])
      const tx = await cl.closeEvent()
      const receipt = await tx.wait()
      const decoded = getPayloadFrom(receipt)

      assert.deepEqual(decoded, { word: ['seinfeld', '"4"', 'LIFE'] })
    })
  })
})
