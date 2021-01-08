import {
  contract,
  debug as d,
  helpers as h,
  matchers,
  setup,
} from '@chainlink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { ConcreteChainlink__factory } from '../../ethers/v0.4/factories/ConcreteChainlink__factory'

const provider = setup.provider()
const concreteChainlinkFactory = new ConcreteChainlink__factory()
const debug = d.makeDebug('ConcreteChainlink')

describe('ConcreteChainlink', () => {
  let ccl: contract.Instance<ConcreteChainlink__factory>
  let defaultAccount: ethers.Wallet
  const deployment = setup.snapshot(provider, async () => {
    defaultAccount = await setup
      .users(provider)
      .then((r) => r.roles.defaultAccount)
    ccl = await concreteChainlinkFactory.connect(defaultAccount).deploy()
  })

  beforeEach(async () => {
    await deployment()
  })

  it('has a limited public interface', () => {
    matchers.publicAbi(concreteChainlinkFactory, [
      'add',
      'addBytes',
      'addInt',
      'addStringArray',
      'addUint',
      'closeEvent',
      'setBuffer',
    ])
  })

  async function parseCCLEvent(tx: ethers.providers.TransactionResponse) {
    const receipt = await tx.wait()
    const data = receipt.logs?.[0].data
    const d = debug.extend('parseCCLEvent')
    d('data %s', data)
    return ethers.utils.defaultAbiCoder.decode(['bytes'], data ?? '')
  }

  describe('#close', () => {
    it('handles empty payloads', async () => {
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await h.decodeDietCBOR(payload)
      assert.deepEqual(decoded, {})
    })
  })

  describe('#setBuffer', () => {
    it('emits the buffer', async () => {
      await ccl.setBuffer('0xA161616162')
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await h.decodeDietCBOR(payload)
      assert.deepEqual(decoded, { a: 'b' })
    })
  })

  describe('#add', () => {
    it('stores and logs keys and values', async () => {
      await ccl.add('first', 'word!!')
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await h.decodeDietCBOR(payload)
      assert.deepEqual(decoded, { first: 'word!!' })
    })

    it('handles two entries', async () => {
      await ccl.add('first', 'uno')
      await ccl.add('second', 'dos')
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await h.decodeDietCBOR(payload)

      assert.deepEqual(decoded, {
        first: 'uno',
        second: 'dos',
      })
    })
  })

  describe('#addBytes', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addBytes('first', '0xaabbccddeeff')
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await h.decodeDietCBOR(payload)
      const expected = h.hexToBuf('0xaabbccddeeff')
      assert.deepEqual(decoded, { first: expected })
    })

    it('handles two entries', async () => {
      await ccl.addBytes('first', '0x756E6F')
      await ccl.addBytes('second', '0x646F73')
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await h.decodeDietCBOR(payload)

      const expectedFirst = h.hexToBuf('0x756E6F')
      const expectedSecond = h.hexToBuf('0x646F73')
      assert.deepEqual(decoded, {
        first: expectedFirst,
        second: expectedSecond,
      })
    })

    it('handles strings', async () => {
      await ccl.addBytes('first', ethers.utils.toUtf8Bytes('apple'))
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await h.decodeDietCBOR(payload)
      const expected = ethers.utils.toUtf8Bytes('apple')
      assert.deepEqual(decoded, { first: expected })
    })
  })

  describe('#addInt', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addInt('first', 1)
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await h.decodeDietCBOR(payload)
      assert.deepEqual(decoded, { first: 1 })
    })

    it('handles two entries', async () => {
      await ccl.addInt('first', 1)
      await ccl.addInt('second', 2)
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await h.decodeDietCBOR(payload)

      assert.deepEqual(decoded, {
        first: 1,
        second: 2,
      })
    })
  })

  describe('#addUint', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addUint('first', 1)
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await h.decodeDietCBOR(payload)
      assert.deepEqual(decoded, { first: 1 })
    })

    it('handles two entries', async () => {
      await ccl.addUint('first', 1)
      await ccl.addUint('second', 2)
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await h.decodeDietCBOR(payload)

      assert.deepEqual(decoded, {
        first: 1,
        second: 2,
      })
    })
  })

  describe('#addStringArray', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addStringArray('word', [
        ethers.utils.formatBytes32String('seinfeld'),
        ethers.utils.formatBytes32String('"4"'),
        ethers.utils.formatBytes32String('LIFE'),
      ])
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await h.decodeDietCBOR(payload)
      assert.deepEqual(decoded, { word: ['seinfeld', '"4"', 'LIFE'] })
    })
  })
})
