import { ethers } from 'hardhat'
import { publicAbi, decodeDietCBOR, hexToBuf } from '../test-helpers/helpers'
import { assert } from 'chai'
import { Contract, ContractFactory, providers, Signer } from 'ethers'
import { Roles, getUsers } from '../test-helpers/setup'
import { makeDebug } from '../test-helpers/debug'

const debug = makeDebug('ChainlinkTestHelper')
let concreteChainlinkFactory: ContractFactory

let roles: Roles

before(async () => {
  roles = (await getUsers()).roles
  concreteChainlinkFactory = await ethers.getContractFactory(
    'src/v0.6/tests/ChainlinkTestHelper.sol:ChainlinkTestHelper',
    roles.defaultAccount,
  )
})

describe('ChainlinkTestHelper', () => {
  let ccl: Contract
  let defaultAccount: Signer

  beforeEach(async () => {
    defaultAccount = roles.defaultAccount
    ccl = await concreteChainlinkFactory.connect(defaultAccount).deploy()
  })

  it('has a limited public interface [ @skip-coverage ]', () => {
    publicAbi(ccl, [
      'add',
      'addBytes',
      'addInt',
      'addStringArray',
      'addUint',
      'closeEvent',
      'setBuffer',
    ])
  })

  async function parseCCLEvent(tx: providers.TransactionResponse) {
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
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, {})
    })
  })

  describe('#setBuffer', () => {
    it('emits the buffer', async () => {
      await ccl.setBuffer('0xA161616162')
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { a: 'b' })
    })
  })

  describe('#add', () => {
    it('stores and logs keys and values', async () => {
      await ccl.add('first', 'word!!')
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { first: 'word!!' })
    })

    it('handles two entries', async () => {
      await ccl.add('first', 'uno')
      await ccl.add('second', 'dos')
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
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
      const [payload] = await parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      const expected = hexToBuf('0xaabbccddeeff')
      assert.deepEqual(decoded, { first: expected })
    })

    it('handles two entries', async () => {
      await ccl.addBytes('first', '0x756E6F')
      await ccl.addBytes('second', '0x646F73')
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)

      const expectedFirst = hexToBuf('0x756E6F')
      const expectedSecond = hexToBuf('0x646F73')
      assert.deepEqual(decoded, {
        first: expectedFirst,
        second: expectedSecond,
      })
    })

    it('handles strings', async () => {
      await ccl.addBytes('first', ethers.utils.toUtf8Bytes('apple'))
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      const expected = ethers.utils.toUtf8Bytes('apple')
      assert.deepEqual(decoded, { first: expected })
    })
  })

  describe('#addInt', () => {
    it('stores and logs keys and values', async () => {
      await ccl.addInt('first', 1)
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { first: 1 })
    })

    it('handles two entries', async () => {
      await ccl.addInt('first', 1)
      await ccl.addInt('second', 2)
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
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
      const [payload] = await parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { first: 1 })
    })

    it('handles two entries', async () => {
      await ccl.addUint('first', 1)
      await ccl.addUint('second', 2)
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
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
        ethers.utils.formatBytes32String('seinfeld'),
        ethers.utils.formatBytes32String('"4"'),
        ethers.utils.formatBytes32String('LIFE'),
      ])
      const tx = await ccl.closeEvent()
      const [payload] = await parseCCLEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, { word: ['seinfeld', '"4"', 'LIFE'] })
    })
  })
})
