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
  encodeReport,
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
      const flags =
        '0x0101010101010101010101010101010101010101010101010101010101010101'
      const callbackGas = 100_000
      await contracts.router.setFlags(subscriptionId, flags)
      const defaultAccountAddress = await roles.defaultAccount.getAddress()
      await expect(
        contracts.client
          .connect(roles.defaultAccount)
          .sendSimpleRequestWithJavaScript(
            'return `hello world`',
            subscriptionId,
            ids.donId,
            callbackGas,
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
          flags,
          callbackGas,
          anyValue,
        )
    })

    it('respects gas flag setting', async () => {
      const subscriptionId = await createSubscription(
        roles.subOwner,
        [contracts.client.address],
        contracts.router,
        contracts.accessControl,
        contracts.linkToken,
      )
      const flags =
        '0x0101010101010101010101010101010101010101010101010101010101010101'
      await contracts.router.setFlags(subscriptionId, flags)
      await expect(
        contracts.client
          .connect(roles.defaultAccount)
          .sendSimpleRequestWithJavaScript(
            'return `hello world`',
            subscriptionId,
            ids.donId,
            400_000,
          ),
      )
        .to.emit(contracts.client, 'RequestSent')
        .to.emit(contracts.coordinator, 'OracleRequest')
      await expect(
        contracts.client
          .connect(roles.defaultAccount)
          .sendSimpleRequestWithJavaScript(
            'return `hello world`',
            subscriptionId,
            ids.donId,
            600_000, // limit set by gas flag == 1 is 500_000
          ),
      ).to.be.revertedWith('GasLimitTooBig(500000)')
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
        20_000,
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
        20_000,
      )
      const { events } = await tx.wait()
      const requestId = getEventArg(events, 'RequestSent', 0)
      await expect(tx)
        .to.emit(contracts.client, 'RequestSent')
        .withArgs(requestId)

      const response = stringToBytes('response')
      const error = stringToBytes('')
      const oracleRequestEvent = await contracts.coordinator.queryFilter(
        contracts.coordinator.filters.OracleRequest(),
      )
      const onchainMetadata = oracleRequestEvent[0].args?.['commitment']
      const report = await encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        response,
        error,
        onchainMetadata,
        stringToBytes(''),
      )
      await expect(contracts.coordinator.callReport(report))
        .to.emit(contracts.coordinator, 'OracleResponse')
        .withArgs(requestId, await roles.defaultAccount.getAddress())
        .to.emit(contracts.client, 'FulfillRequestInvoked')
        .withArgs(requestId, response, error)
    })
  })
})

describe('Faulty Functions Client', () => {
  it('can complete requests with an empty callback', async () => {
    const clientWithEmptyCallbackTestHelperFactory =
      await ethers.getContractFactory(
        'src/v0.8/functions/tests/v1_X/testhelpers/FunctionsClientWithEmptyCallback.sol:FunctionsClientWithEmptyCallback',
        roles.consumer,
      )

    const clientWithEmptyCallback =
      await clientWithEmptyCallbackTestHelperFactory
        .connect(roles.consumer)
        .deploy(contracts.router.address)

    const subscriptionId = await createSubscription(
      roles.subOwner,
      [clientWithEmptyCallback.address],
      contracts.router,
      contracts.accessControl,
      contracts.linkToken,
    )
    const tx = await clientWithEmptyCallback.sendSimpleRequestWithJavaScript(
      'function run(){return response}',
      subscriptionId,
      ids.donId,
      20_000,
    )
    const { events } = await tx.wait()
    const requestId = getEventArg(events, 'RequestSent', 0)
    await expect(tx)
      .to.emit(clientWithEmptyCallback, 'RequestSent')
      .withArgs(requestId)

    const response = stringToBytes('response')
    const error = stringToBytes('')
    const oracleRequestEvent = await contracts.coordinator.queryFilter(
      contracts.coordinator.filters.OracleRequest(),
    )
    const onchainMetadata = oracleRequestEvent[0].args?.['commitment']
    const report = await encodeReport(
      ethers.utils.hexZeroPad(requestId, 32),
      response,
      error,
      onchainMetadata,
      stringToBytes(''),
    )
    await expect(contracts.coordinator.callReport(report))
      .to.emit(contracts.coordinator, 'OracleResponse')
      .withArgs(requestId, await roles.defaultAccount.getAddress())
      .to.emit(contracts.router, 'RequestProcessed')
  })
})
