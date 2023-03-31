import { ethers } from 'hardhat'
import { publicAbi, decodeDietCBOR, hexToBuf } from '../../test-helpers/helpers'
import { assert, expect } from 'chai'
import { Contract, ContractFactory, providers, Signer } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'
import { makeDebug } from '../../test-helpers/debug'

const debug = makeDebug('FunctionsTestHelper')
let concreteFunctionsTestHelperFactory: ContractFactory

let roles: Roles

before(async () => {
  roles = (await getUsers()).roles
  concreteFunctionsTestHelperFactory = await ethers.getContractFactory(
    'src/v0.8/tests/FunctionsTestHelper.sol:FunctionsTestHelper',
    roles.defaultAccount,
  )
})

describe('FunctionsTestHelper', () => {
  let ctr: Contract
  let defaultAccount: Signer

  beforeEach(async () => {
    defaultAccount = roles.defaultAccount
    ctr = await concreteFunctionsTestHelperFactory
      .connect(defaultAccount)
      .deploy()
  })

  it('has a limited public interface [ @skip-coverage ]', () => {
    publicAbi(ctr, [
      'closeEvent',
      'initializeRequestForInlineJavaScript',
      'addSecrets',
      'addTwoArgs',
      'addEmptyArgs',
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
      assert.deepEqual(
        {
          ...decoded,
          language: decoded.language.toNumber(),
          codeLocation: decoded.codeLocation.toNumber(),
        },
        {
          language: 0,
          codeLocation: 0,
          source: '',
        },
      )
    })
  })

  describe('#initializeRequestForInlineJavaScript', () => {
    it('emits simple CBOR encoded request for js', async () => {
      const js = 'function run(args, responses) {}'
      await ctr.initializeRequestForInlineJavaScript(js)
      const tx = await ctr.closeEvent()
      const [payload] = await parseRequestDataEvent(tx)
      const decoded = await decodeDietCBOR(payload)
      assert.deepEqual(
        {
          ...decoded,
          language: decoded.language.toNumber(),
          codeLocation: decoded.codeLocation.toNumber(),
        },
        {
          language: 0,
          codeLocation: 0,
          source: js,
        },
      )
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
      assert.deepEqual(
        {
          ...decoded,
          language: decoded.language.toNumber(),
          codeLocation: decoded.codeLocation.toNumber(),
          secretsLocation: decoded.secretsLocation.toNumber(),
        },
        {
          language: 0,
          codeLocation: 0,
          source: js,
          secretsLocation: 1,
          secrets: hexToBuf(secrets),
        },
      )
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
      assert.deepEqual(
        {
          ...decoded,
          language: decoded.language.toNumber(),
          codeLocation: decoded.codeLocation.toNumber(),
        },
        {
          language: 0,
          codeLocation: 0,
          source: js,
          args: ['arg1', 'arg2'],
        },
      )
    })
  })

  describe('#addEmptyArgs to revert', () => {
    it('reverts with EmptyArgs() if args param is empty', async () => {
      await expect(ctr.addEmptyArgs()).to.be.revertedWith('EmptyArgs()')
    })
  })
})
