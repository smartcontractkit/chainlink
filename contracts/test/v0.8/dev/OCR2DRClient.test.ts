import web3 from 'web3'
import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { Contract, ContractFactory, providers } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'
import { decodeDietCBOR } from '../../test-helpers/helpers'

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
  return ethers.utils.defaultAbiCoder.decode(
    ['address', 'bytes32', 'bytes'],
    data ?? '',
  )
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
  const subscriptionId = 1
  const anyValue = () => true

  let client: Contract
  let oracle: Contract

  beforeEach(async () => {
    oracle = await ocr2drOracleFactory.connect(roles.defaultAccount).deploy()
    client = await concreteOCR2DRClientFactory
      .connect(roles.defaultAccount)
      .deploy(oracle.address)
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
        .withArgs(
          client.address,
          anyValue,
          anyValue,
        )
    })

    it('encodes user request to CBOR', async () => {
      const js = 'function run() {}'
      const tx = await client.sendSimpleRequestWithJavaScript(
        js,
        subscriptionId,
      )
      const args = await parseOracleRequestEventArgs(tx)
      assert.equal(3, args.length)

      const decoded = await decodeDietCBOR(args[2])
      assert.deepEqual(decoded, {
        language: 0,
        codeLocation: 0,
        source: js,
      })
    })
  })

  describe('#cancelPendingRequest', () => {
    it('emits events from the client and the oracle contracts', async () => {
      const tx = await client.sendSimpleRequestWithJavaScript(
        'function run() {}',
        subscriptionId,
      )

      const { events } = await tx.wait()
      const requestId = getEventArg(events, 'RequestSent', 0)
      await expect(tx).to.emit(client, 'RequestSent').withArgs(requestId)

      await expect(client.cancelPendingRequest(requestId))
        .to.emit(client, 'RequestCancelled')
        .withArgs(requestId)
        .to.emit(oracle, 'CancelOracleRequest')
        .withArgs(requestId)
    })

    it('reverts for unknown requestId', async () => {
      await expect(
        client.cancelPendingRequest(
          '0x1bfce59c2e0d7e0f015eb02ec4e04de4e67a1fe1508a4420cfd49c650758abed',
        ),
      ).to.be.revertedWith('RequestIsNotPending')
    })
  })

  describe('#fulfillRequest', () => {
    it('emits fulfillment events', async () => {
      const tx = await client.sendSimpleRequestWithJavaScript(
        'function run() {}',
        subscriptionId,
      )

      const { events } = await tx.wait()
      const requestId = getEventArg(events, 'RequestSent', 0)
      await expect(tx).to.emit(client, 'RequestSent').withArgs(requestId)

      const response = web3.utils.asciiToHex('response')
      const error = web3.utils.asciiToHex('error')
      await expect(
        oracle.fulfillRequest(requestId, client.address, response, error),
      )
        .to.emit(oracle, 'OracleResponse')
        .withArgs(requestId)
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId, response, error)
    })
  })
})
