'use strict'

require('./support/helpers.js')

contract('ConcreteChainlinkLib', () => {
  const sourcePath = 'examples/ConcreteChainlinkLib.sol'
  let ccl

  beforeEach(async () => {
    ccl = await deploy(sourcePath)
  })

  it('has a limited public interface', () => {
    checkPublicABI(artifacts.require(sourcePath), [
      'add',
      'addStringArray',
      'addInt',
      'addUint',
      'closeEvent'
    ])
  })

  function parseCCLEvent (tx) {
    let data = util.toBuffer(tx.receipt.logs[0].data)
    return abi.rawDecode(['bytes'], data)
  }

  describe('#close', () => {
    it('handles empty payloads', async () => {
      let tx = await ccl.closeEvent()
      let [payload] = parseCCLEvent(tx)
      var decoded = await cbor.decodeFirst(payload)
      assert.deepEqual(decoded, {})
    })
  })

  describe('#add', () => {
    it('stores and logs keys and values', async () => {
      await ccl.add('first', 'word!!')
      let tx = await ccl.closeEvent()
      let [payload] = parseCCLEvent(tx)
      var decoded = await cbor.decodeFirst(payload)
      assert.deepEqual(decoded, { 'first': 'word!!' })
    })

    it('handles two entries', async () => {
      await ccl.add('first', 'uno')
      await ccl.add('second', 'dos')
      let tx = await ccl.closeEvent()
      let [payload] = parseCCLEvent(tx)
      var decoded = await cbor.decodeFirst(payload)

      assert.deepEqual(decoded, {
        'first': 'uno',
        'second': 'dos'
      })
    })
  })

  describe('#addInt', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addInt('first', 1)
      let tx = await ccl.closeEvent()
      let [payload] = parseCCLEvent(tx)
      var decoded = await cbor.decodeFirst(payload)
      assert.deepEqual(decoded, { 'first': 1 })
    })

    it('handles two entries', async () => {
      await ccl.addInt('first', 1)
      await ccl.addInt('second', 2)
      let tx = await ccl.closeEvent()
      let [payload] = parseCCLEvent(tx)
      var decoded = await cbor.decodeFirst(payload)

      assert.deepEqual(decoded, {
        'first': 1,
        'second': 2
      })
    })
  })

  describe('#addUint', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addUint('first', 1)
      let tx = await ccl.closeEvent()
      let [payload] = parseCCLEvent(tx)
      var decoded = await cbor.decodeFirst(payload)
      assert.deepEqual(decoded, { 'first': 1 })
    })

    it('handles two entries', async () => {
      await ccl.addUint('first', 1)
      await ccl.addUint('second', 2)
      let tx = await ccl.closeEvent()
      let [payload] = parseCCLEvent(tx)
      var decoded = await cbor.decodeFirst(payload)

      assert.deepEqual(decoded, {
        'first': 1,
        'second': 2
      })
    })
  })

  describe('#addStringArray', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addStringArray('word', ['seinfeld', '"4"', 'LIFE'])
      let tx = await ccl.closeEvent()
      let [payload] = parseCCLEvent(tx)
      var decoded = await cbor.decodeFirst(payload)

      assert.deepEqual(decoded, { 'word': ['seinfeld', '"4"', 'LIFE'] })
    })
  })
})
