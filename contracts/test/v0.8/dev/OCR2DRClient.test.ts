import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { Contract, ContractFactory, providers } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'
import { decodeDietCBOR, stringToBytes } from '../../test-helpers/helpers'

let concreteOCR2DRClientFactory: ContractFactory
let ocr2drOracleFactory: ContractFactory
let roles: Roles

function getEventArg(events: any, eventName: string, argIndex: number) {
  if (Array.isArray(events)) {
    const event = events.find((e: any) => e.event == eventName)
    if (event && Array.isArray(event.args) && event.args.length > 0) {
      return event.args[argIndex]
    }
  }
  return undefined
}

async function parseOracleRequestEventArgs(tx: providers.TransactionResponse) {
  const receipt = await tx.wait()
  const data = receipt.logs?.[0].data
  return ethers.utils.defaultAbiCoder.decode(['bytes32', 'bytes'], data ?? '')
}

before(async () => {
  roles = (await getUsers()).roles

  concreteOCR2DRClientFactory = await ethers.getContractFactory(
    'src/v0.8/tests/OCR2DRClientTestHelper.sol:OCR2DRClientTestHelper',
    roles.defaultAccount,
  )
  ocr2drOracleFactory = await ethers.getContractFactory(
    'src/v0.8/dev/ocr2dr/OCR2DROracle.sol:OCR2DROracle',
    roles.defaultAccount,
  )
})

describe('OCR2DRClientTestHelper', () => {
  const donPublicKey =
    '0x3804a19f2437f7bba4fcfbc194379e43e514aa98073db3528ccdbdb642e24011'
  const subscriptionId = 1
  const anyValue = () => true

  let client: Contract
  let oracle: Contract

  beforeEach(async () => {
    const accounts = await ethers.getSigners()
    oracle = await ocr2drOracleFactory
      .connect(roles.defaultAccount)
      .deploy(accounts[0].address, donPublicKey)
    client = await concreteOCR2DRClientFactory
      .connect(roles.defaultAccount)
      .deploy(oracle.address)
  })

  describe('#getDONPublicKey', () => {
    it('returns DON public key set on Oracle', async () => {
      expect(await client.callStatic.getDONPublicKey()).to.be.equal(
        donPublicKey,
      )
    })
  })

  describe('#sendSimpleRequestWithJavaScript', () => {
    it('emits events from the client and the oracle contracts', async () => {
      await expect(
        client.sendSimpleRequestWithJavaScript(
          'function run() {}',
          subscriptionId,
        ),
      )
        .to.emit(client, 'RequestSent')
        .withArgs(anyValue)
        .to.emit(oracle, 'OracleRequest')
        .withArgs(anyValue, anyValue)
    })

    it('encodes user request to CBOR', async () => {
      const js = 'function run() {}'
      const tx = await client.sendSimpleRequestWithJavaScript(
        js,
        subscriptionId,
      )
      const args = await parseOracleRequestEventArgs(tx)
      assert.equal(2, args.length)

      const decoded = await decodeDietCBOR(args[1])
      assert.deepEqual(decoded, {
        language: 0,
        codeLocation: 0,
        source: js,
      })
    })
  })

  describe('#fulfillRequest', () => {
    it('emits fulfillment events', async () => {
      const accounts = await ethers.getSigners()
      await oracle.setAuthorizedSenders([accounts[0].address])

      const tx = await client.sendSimpleRequestWithJavaScript(
        'function run() {}',
        subscriptionId,
      )

      const { events } = await tx.wait()
      const requestId = getEventArg(events, 'RequestSent', 0)
      await expect(tx).to.emit(client, 'RequestSent').withArgs(requestId)

      const response = stringToBytes('response')
      const error = stringToBytes('error')
      await expect(oracle.fulfillRequest(requestId, response, error))
        .to.emit(oracle, 'OracleResponse')
        .withArgs(requestId)
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId, response, error)
    })
  })
})
