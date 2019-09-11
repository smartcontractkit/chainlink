import { toBuffer } from 'ethereumjs-util'
import abi from 'ethereumjs-abi'
import { checkPublicABI, decodeDietCBOR, toHex } from './support/helpers'
const ConcreteChainlink = artifacts.require('ConcreteChainlink.sol')

contract('ConcreteChainlink', () => {
  let ccl

  beforeEach(async () => {
    ccl = await ConcreteChainlink.new()
  })

  it('has a limited public interface', () => {
    checkPublicABI(ConcreteChainlink, [
      'add',
      'addBytes',
      'addInt',
      'addStringArray',
      'addUint',
      'closeEvent',
      'setBuffer',
    ])
  })

  function parseCCLEvent(tx) {
    const data = toBuffer(tx.receipt.rawLogs[0].data)
    return abi.rawDecode(['bytes'], data)
  }

  describe('#close', () => {
    it('handles empty payloads', async () => {
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, {})
    })
  })

  describe('#setBuffer', () => {
    it('emits the buffer', async () => {
      await ccl.setBuffer('0xA161616162')
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { a: 'b' })
    })
  })

  describe('#add', () => {
    it('stores and logs keys and values', async () => {
      await ccl.add('first', 'word!!')
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { first: 'word!!' })
    })

    it('handles two entries', async () => {
      await ccl.add('first', 'uno')
      await ccl.add('second', 'dos')
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

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
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      const expected = toBuffer('0xaabbccddeeff')
      assert.deepEqual(decoded, { first: expected })
    })

    it('handles two entries', async () => {
      await ccl.addBytes('first', '0x756E6F')
      await ccl.addBytes('second', '0x646F73')
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

      const expectedFirst = toBuffer('0x756E6F')
      const expectedSecond = toBuffer('0x646F73')
      assert.deepEqual(decoded, {
        first: expectedFirst,
        second: expectedSecond,
      })
    })

    it('handles strings', async () => {
      await ccl.addBytes('first', toHex('apple'))
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      const expected = toBuffer('apple')
      assert.deepEqual(decoded, { first: expected })
    })
  })

  describe('#addInt', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addInt('first', 1)
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { first: 1 })
    })

    it('handles two entries', async () => {
      await ccl.addInt('first', 1)
      await ccl.addInt('second', 2)
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

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
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { first: 1 })
    })

    it('handles two entries', async () => {
      await ccl.addUint('first', 1)
      await ccl.addUint('second', 2)
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

      assert.deepEqual(decoded, {
        first: 1,
        second: 2,
      })
    })
  })

  describe('#addStringArray', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addStringArray('word', [
        toHex('seinfeld'),
        toHex('"4"'),
        toHex('LIFE'),
      ])
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { word: ['seinfeld', '"4"', 'LIFE'] })
    })
  })
})
