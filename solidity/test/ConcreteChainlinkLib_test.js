import util from 'ethereumjs-util'
import abi from 'ethereumjs-abi'
import {
  checkPublicABI,
  decodeDietCBOR,
  deploy
} from './support/helpers.js'

contract('ConcreteChainlink', () => {
  const sourcePath = 'examples/ConcreteChainlink.sol'
  let ccl

  beforeEach(async () => {
    ccl = await deploy(sourcePath)
  })

  it('has a limited public interface', () => {
    checkPublicABI(artifacts.require(sourcePath), [
      'add',
      'addBytes',
      'addInt',
      'addStringArray',
      'addUint',
      'closeEvent',
      'setBuffer'
    ])
  })

  function parseCCLEvent (tx) {
    const data = util.toBuffer(tx.receipt.logs[0].data)
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
      assert.deepEqual(decoded, {'a':'b'})
    })
  })

  describe('#add', () => {
    it('stores and logs keys and values', async () => {
      await ccl.add('first', 'word!!')
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { 'first': 'word!!' })
    })

    it('handles two entries', async () => {
      await ccl.add('first', 'uno')
      await ccl.add('second', 'dos')
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

      assert.deepEqual(decoded, {
        'first': 'uno',
        'second': 'dos'
      })
    })
  })

  describe('#addBytes', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addBytes('first', '0xaabbccddeeff')
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      const expected = util.toBuffer('0xaabbccddeeff')
      assert.deepEqual(decoded, { 'first': expected })
    })

    it('handles two entries', async () => {
      await ccl.addBytes('first', '0x756E6F')
      await ccl.addBytes('second', '0x646F73')
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

      const expectedFirst = util.toBuffer('0x756E6F')
      const expectedSecond = util.toBuffer('0x646F73')
      assert.deepEqual(decoded, {
        'first': expectedFirst,
        'second': expectedSecond
      })
    })

    it('handles strings', async () => {
      await ccl.addBytes('first', 'apple')
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      const expected = util.toBuffer('apple')
      assert.deepEqual(decoded, { 'first': expected })
    })
  })

  describe('#addInt', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addInt('first', 1)
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { 'first': 1 })
    })

    it('handles two entries', async () => {
      await ccl.addInt('first', 1)
      await ccl.addInt('second', 2)
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

      assert.deepEqual(decoded, {
        'first': 1,
        'second': 2
      })
    })
  })

  describe('#addUint', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addUint('first', 1)
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { 'first': 1 })
    })

    it('handles two entries', async () => {
      await ccl.addUint('first', 1)
      await ccl.addUint('second', 2)
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

      assert.deepEqual(decoded, {
        'first': 1,
        'second': 2
      })
    })
  })

  describe('#addStringArray', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addStringArray('word', ['seinfeld', '"4"', 'LIFE'])
      const tx = await ccl.closeEvent()
      const [payload] = parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

      assert.deepEqual(decoded, { 'word': ['seinfeld', '"4"', 'LIFE'] })
    })
  })
})
