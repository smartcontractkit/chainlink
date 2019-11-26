import { toBuffer } from 'ethereumjs-util'
import abi from 'ethereumjs-abi'
import { checkPublicABI, decodeDietCBOR, toHex } from './support/helpers'
const Chainlink = artifacts.require('ChainlinkTestHelper.sol')

contract('Chainlink', () => {
  let cl

  beforeEach(async () => {
    cl = await Chainlink.new()
  })

  it('has a limited public interface', () => {
    checkPublicABI(Chainlink, [
      'add',
      'addBytes',
      'addInt',
      'addStringArray',
      'addUint',
      'closeEvent',
      'setBuffer',
    ])
  })

  function parseCLEvent(tx) {
    const data = toBuffer(tx.receipt.rawLogs[0].data)
    return abi.rawDecode(['bytes'], data)
  }

  describe('#close', () => {
    it('handles empty payloads', async () => {
      const tx = await cl.closeEvent()
      const [payload] = parseCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, {})
    })
  })

  describe('#setBuffer', () => {
    it('emits the buffer', async () => {
      await cl.setBuffer('0xA161616162')
      const tx = await cl.closeEvent()
      const [payload] = parseCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { a: 'b' })
    })
  })

  describe('#add', () => {
    it('stores and logs keys and values', async () => {
      await cl.add('first', 'word!!')
      const tx = await cl.closeEvent()
      const [payload] = parseCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { first: 'word!!' })
    })

    it('handles two entries', async () => {
      await cl.add('first', 'uno')
      await cl.add('second', 'dos')
      const tx = await cl.closeEvent()
      const [payload] = parseCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

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
      const [payload] = parseCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      const expected = toBuffer('0xaabbccddeeff')
      assert.deepEqual(decoded, { first: expected })
    })

    it('handles two entries', async () => {
      await cl.addBytes('first', '0x756E6F')
      await cl.addBytes('second', '0x646F73')
      const tx = await cl.closeEvent()
      const [payload] = parseCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

      const expectedFirst = toBuffer('0x756E6F')
      const expectedSecond = toBuffer('0x646F73')
      assert.deepEqual(decoded, {
        first: expectedFirst,
        second: expectedSecond,
      })
    })

    it('handles strings', async () => {
      await cl.addBytes('first', toHex('apple'))
      const tx = await cl.closeEvent()
      const [payload] = parseCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      const expected = toBuffer('apple')
      assert.deepEqual(decoded, { first: expected })
    })
  })

  describe('#addInt', () => {
    it('stores and logs keys and values', async () => {
      await cl.addInt('first', 1)
      const tx = await cl.closeEvent()
      const [payload] = parseCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { first: 1 })
    })

    it('handles two entries', async () => {
      await cl.addInt('first', 1)
      await cl.addInt('second', 2)
      const tx = await cl.closeEvent()
      const [payload] = parseCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

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
      const [payload] = parseCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { first: 1 })
    })

    it('handles two entries', async () => {
      await cl.addUint('first', 1)
      await cl.addUint('second', 2)
      const tx = await cl.closeEvent()
      const [payload] = parseCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

      assert.deepEqual(decoded, {
        first: 1,
        second: 2,
      })
    })
  })

  describe('#addStringArray', () => {
    it('stores and logs keys and values', async () => {
      await cl.addStringArray('word', [
        toHex('seinfeld'),
        toHex('"4"'),
        toHex('LIFE'),
      ])
      const tx = await cl.closeEvent()
      const [payload] = parseCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { word: ['seinfeld', '"4"', 'LIFE'] })
    })
  })
})
