import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { decodeDietCBOR, stringToBytes } from '../../../test-helpers/helpers'
import {
  getSetupFactory,
  FunctionsContracts,
  FunctionsRoles,
  anyValue,
  ids,
  createSubscription,
  getEventArg,
  parseOracleRequestEventArgs,
} from './utils'

const setup = getSetupFactory()
let contracts: FunctionsContracts
let roles: FunctionsRoles

beforeEach(async () => {
  ;({ contracts, roles } = setup())
})

describe('Functions Client', () => {
  describe('#sendSimpleRequestWithJavaScript', () => {
    it('emits events from the client and the oracle contracts', async () => {
      const subscriptionId = await createSubscription(
        roles.subOwner,
        [contracts.client.address],
        contracts.router,
        contracts.accessControl,
        contracts.linkToken,
      )
      const defaultAccountAddress = await roles.defaultAccount.getAddress()
      await expect(
        contracts.client
          .connect(roles.defaultAccount)
          .sendSimpleRequestWithJavaScript(
            'return `hello world`',
            subscriptionId,
            ids.donId,
          ),
      )
        .to.emit(contracts.client, 'RequestSent')
        .withArgs(anyValue)
        .to.emit(contracts.coordinator, 'OracleRequest')
        .withArgs(
          anyValue,
          contracts.client.address,
          defaultAccountAddress,
          subscriptionId,
          roles.subOwnerAddress,
          anyValue,
          anyValue,
        )
    })

    it('encodes user request to CBOR', async () => {
      const subscriptionId = await createSubscription(
        roles.subOwner,
        [contracts.client.address],
        contracts.router,
        contracts.accessControl,
        contracts.linkToken,
      )
      const js = 'function run(){return response}'
      const tx = await contracts.client.sendSimpleRequestWithJavaScript(
        js,
        subscriptionId,
        ids.donId,
      )
      const args = await parseOracleRequestEventArgs(tx)
      assert.equal(args.length, 5)
      const decoded = await decodeDietCBOR(args[3])
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

  describe('#fulfillRequest', () => {
    it('emits fulfillment events', async () => {
      const subscriptionId = await createSubscription(
        roles.subOwner,
        [contracts.client.address],
        contracts.router,
        contracts.accessControl,
        contracts.linkToken,
      )
      const tx = await contracts.client.sendSimpleRequestWithJavaScript(
        'function run(){return response}',
        subscriptionId,
        ids.donId,
      )
      const { events } = await tx.wait()
      const requestId = getEventArg(events, 'RequestSent', 0)
      await expect(tx)
        .to.emit(contracts.client, 'RequestSent')
        .withArgs(requestId)

      const response = stringToBytes('response')
      const error = stringToBytes('')
      const abi = ethers.utils.defaultAbiCoder

      const report = abi.encode(
        ['bytes32[]', 'bytes[]', 'bytes[]'],
        [[ethers.utils.hexZeroPad(requestId, 32)], [response], [error]],
      )

      await expect(contracts.coordinator.callReport(report))
        .to.emit(contracts.coordinator, 'OracleResponse')
        .withArgs(requestId, await roles.defaultAccount.getAddress())
        .to.emit(contracts.coordinator, 'BillingEnd')
        .to.emit(contracts.client, 'FulfillRequestInvoked')
        .withArgs(requestId, response, error)
    })
  })
})
