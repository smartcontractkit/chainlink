import { ethers } from 'hardhat'
import { publicAbi, decodeDietCBOR, hexToBuf } from '../../test-helpers/helpers'
import { assert, expect } from 'chai'
import { Contract, ContractFactory, providers, Signer } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'
import { makeDebug } from '../../test-helpers/debug'

const debug = makeDebug('OCR2DRTestHelper')
let concreteOCR2DRTestHelperFactory: ContractFactory

let roles: Roles

before(async () => {
  roles = (await getUsers()).roles
  concreteOCR2DRTestHelperFactory = await ethers.getContractFactory(
    'src/v0.8/tests/OCR2DRTestHelper.sol:OCR2DRTestHelper',
    roles.defaultAccount,
  )
})

describe('OCR2DRTestHelper', () => {
  let ctr: Contract
  let defaultAccount: Signer

  beforeEach(async () => {
    defaultAccount = roles.defaultAccount
    ctr = await concreteOCR2DRTestHelperFactory.connect(defaultAccount).deploy()
  })

  it('has a limited public interface [ @skip-coverage ]', () => {
    publicAbi(ctr, [
      'closeEvent',
      'initializeRequestForInlineJavaScript',
      'addSecrets',
      'addTwoArgs',
      'addEmptyArgs',
      'addQuery',
      'setTwoQueries',
      'setEmptyQueries',
      'setEmptyHeaders',
      'addQueryWithTwoHeaders',
    ])
  })

  async function parseRequestDataEvent(tx: providers.TransactionResponse) {
    const receipt = await tx.wait()
    const data = receipt.logs?.[0].data
    const d = debug.extend('parseRequestDataEvent')
    d('data %s', data)
    return ethers.utils.defaultAbiCoder.decode(['bytes'], data ?? '')
  }

  describe('#closeEvent', () => {
    it('handles empty request', async () => {
      const tx = await ctr.closeEvent()
      const [payload] = await parseRequestDataEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, {
        language: 0,
        codeLocation: 0,
        source: '',
      })
    })
  })

  describe('#initializeRequestForInlineJavaScript', () => {
    it('emits simple CBOR encoded request for js', async () => {
      const js = 'function run(args, responses) {}'
      await ctr.initializeRequestForInlineJavaScript(js)
      const tx = await ctr.closeEvent()
      const [payload] = await parseRequestDataEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, {
        language: 0,
        codeLocation: 0,
        source: js,
      })
    })
  })

  describe('#initializeRequestForInlineJavaScript to revert', () => {
    it('reverts with EmptySource() if source param is empty', async () => {
      await expect(
        ctr.initializeRequestForInlineJavaScript(''),
      ).to.be.revertedWith('EmptySource()')
    })
  })

  describe('#addSecrets', () => {
    it('emits CBOR encoded request with js and secrets', async () => {
      const js = 'function run(args, responses) {}'
      const secrets = '0xA161616162'
      await ctr.initializeRequestForInlineJavaScript(js)
      await ctr.addSecrets(secrets)
      const tx = await ctr.closeEvent()
      const [payload] = await parseRequestDataEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, {
        language: 0,
        codeLocation: 0,
        source: js,
        secretsLocation: 0,
        secrets: hexToBuf(secrets),
      })
    })
  })

  describe('#addSecrets to revert', () => {
    it('reverts with EmptySecrets() if secrets param is empty', async () => {
      const js = 'function run(args, responses) {}'
      await ctr.initializeRequestForInlineJavaScript(js)
      await expect(ctr.addSecrets('0x')).to.be.revertedWith('EmptySecrets()')
    })
  })

  describe('#addArgs', () => {
    it('emits CBOR encoded request with js and args', async () => {
      const js = 'function run(args, responses) {}'
      await ctr.initializeRequestForInlineJavaScript(js)
      await ctr.addTwoArgs('arg1', 'arg2')
      const tx = await ctr.closeEvent()
      const [payload] = await parseRequestDataEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, {
        language: 0,
        codeLocation: 0,
        source: js,
        args: ['arg1', 'arg2'],
      })
    })
  })

  describe('#addEmptyArgs to revert', () => {
    it('reverts with EmptyArgs() if args param is empty', async () => {
      await expect(ctr.addEmptyArgs()).to.be.revertedWith('EmptyArgs()')
    })
  })

  describe('#addQuery', () => {
    it('emits CBOR encoded request with js and query', async () => {
      const js = 'function run(args, responses) {}'
      const url = 'https://data.source'
      await ctr.initializeRequestForInlineJavaScript(js)
      await ctr.addQuery(url)
      const tx = await ctr.closeEvent()
      const [payload] = await parseRequestDataEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, {
        language: 0,
        codeLocation: 0,
        source: js,
        queries: [
          {
            verb: 0,
            url,
          },
        ],
      })
    })
  })

  describe('#addQuery to revert', () => {
    it('reverts with EmptyUrl() if url param is empty', async () => {
      await expect(ctr.addQuery('')).to.be.revertedWith('EmptyUrl()')
    })
  })

  describe('#setTwoQueries', () => {
    it('emits CBOR encoded request with two queries', async () => {
      const js = 'function run(args, responses) {}'
      const url1 = 'https://data.source1'
      const url2 = 'https://data.source1'
      await ctr.initializeRequestForInlineJavaScript(js)
      await ctr.setTwoQueries(url1, url2)
      const tx = await ctr.closeEvent()
      const [payload] = await parseRequestDataEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, {
        language: 0,
        codeLocation: 0,
        source: js,
        queries: [
          {
            verb: 0,
            url: url1,
          },
          {
            verb: 0,
            url: url2,
          },
        ],
      })
    })
  })

  describe('#addQueryWithTwoHeaders', () => {
    it('emits CBOR encoded request for a query with two headers', async () => {
      const js = 'function run(args, responses) {}'
      const url = 'https://data.source'
      await ctr.initializeRequestForInlineJavaScript(js)
      await ctr.addQueryWithTwoHeaders(url, 'k1', 'v1', 'k2', 'v2')
      const tx = await ctr.closeEvent()
      const [payload] = await parseRequestDataEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(decoded, {
        language: 0,
        codeLocation: 0,
        source: js,
        queries: [
          {
            verb: 0,
            url,
            headers: {
              k1: 'v1',
              k2: 'v2',
            },
          },
        ],
      })
    })
  })

  describe('#addQueryWithTwoHeaders to revert', () => {
    it('reverts with EmptyKey() if key param is empty', async () => {
      const url = 'https://data.source'
      await expect(
        ctr.addQueryWithTwoHeaders(url, 'k1', 'v1', '', 'v2'),
      ).to.be.revertedWith('EmptyKey()')
    })
    it('reverts with EmptyValue() if value param is empty', async () => {
      const url = 'https://data.source'
      await expect(
        ctr.addQueryWithTwoHeaders(url, 'k1', 'v1', 'k2', ''),
      ).to.be.revertedWith('EmptyValue()')
    })
  })

  describe('#setEmptyQueries to revert', () => {
    it('reverts with EmptyQueries() if queries param is empty', async () => {
      await expect(ctr.setEmptyQueries()).to.be.revertedWith('EmptyQueries()')
    })
  })

  describe('#setEmptyHeaders to revert', () => {
    it('reverts with EmptyHeaders() if headers param is empty', async () => {
      await expect(ctr.setEmptyHeaders()).to.be.revertedWith('EmptyHeaders()')
    })
  })
})
