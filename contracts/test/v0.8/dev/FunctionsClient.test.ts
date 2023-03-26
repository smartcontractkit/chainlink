import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { Contract, ContractFactory, providers } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'
import { decodeDietCBOR, stringToBytes } from '../../test-helpers/helpers'

let concreteFunctionsClientFactory: ContractFactory
let functionsOracleFactory: ContractFactory
let functionsBillingRegistryFactory: ContractFactory
let linkTokenFactory: ContractFactory
let mockAggregatorV3Factory: ContractFactory
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
  const data = receipt.logs?.[1].data
  return ethers.utils.defaultAbiCoder.decode(
    ['address', 'address', 'uint64', 'address', 'bytes'],
    data ?? '',
  )
}

before(async () => {
  roles = (await getUsers()).roles

  concreteFunctionsClientFactory = await ethers.getContractFactory(
    'src/v0.8/tests/FunctionsClientTestHelper.sol:FunctionsClientTestHelper',
    roles.defaultAccount,
  )
  functionsOracleFactory = await ethers.getContractFactory(
    'src/v0.8/tests/FunctionsOracleHelper.sol:FunctionsOracleHelper',
    roles.defaultAccount,
  )

  functionsBillingRegistryFactory = await ethers.getContractFactory(
    'src/v0.8/tests/FunctionsBillingRegistryWithInit.sol:FunctionsBillingRegistryWithInit',
    roles.defaultAccount,
  )

  linkTokenFactory = await ethers.getContractFactory(
    'src/v0.4/LinkToken.sol:LinkToken',
    roles.consumer,
  )

  mockAggregatorV3Factory = await ethers.getContractFactory(
    'src/v0.7/tests/MockV3Aggregator.sol:MockV3Aggregator',
    roles.consumer,
  )
})

describe('FunctionsClientTestHelper', () => {
  const donPublicKey =
    '0x3804a19f2437f7bba4fcfbc194379e43e514aa98073db3528ccdbdb642e24011'
  let subscriptionId: number
  const anyValue = () => true

  let client: Contract
  let oracle: Contract
  let registry: Contract
  let linkToken: Contract
  let mockLinkEth: Contract

  beforeEach(async () => {
    // Deploy
    linkToken = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    mockLinkEth = await mockAggregatorV3Factory.deploy(
      0,
      ethers.BigNumber.from(5021530000000000),
    )
    oracle = await functionsOracleFactory.connect(roles.defaultAccount).deploy()
    registry = await functionsBillingRegistryFactory
      .connect(roles.defaultAccount)
      .deploy(linkToken.address, mockLinkEth.address, oracle.address)

    // Setup contracts
    await oracle.setRegistry(registry.address)
    await oracle.deactivateAuthorizedReceiver()
    client = await concreteFunctionsClientFactory
      .connect(roles.defaultAccount)
      .deploy(oracle.address)
    await registry.setAuthorizedSenders([oracle.address])

    await registry.setConfig(
      1_000_000,
      86_400,
      21_000 + 5_000 + 2_100 + 20_000 + 2 * 2_100 - 15_000 + 7_315,
      ethers.BigNumber.from('5000000000000000'),
      100_000,
      300,
    )

    // Setup accounts
    const createSubTx = await registry
      .connect(roles.defaultAccount)
      .createSubscription()
    const receipt = await createSubTx.wait()
    subscriptionId = receipt.events[0].args['subscriptionId'].toNumber()

    await registry
      .connect(roles.defaultAccount)
      .addConsumer(subscriptionId, client.address)

    await linkToken
      .connect(roles.defaultAccount)
      .transferAndCall(
        registry.address,
        ethers.BigNumber.from('115957983815660167'),
        ethers.utils.defaultAbiCoder.encode(['uint64'], [subscriptionId]),
      )
  })

  describe('#getDONPublicKey', () => {
    it('returns DON public key set on Oracle', async () => {
      await expect(oracle.setDONPublicKey(donPublicKey)).not.to.be.reverted
      expect(await client.callStatic.getDONPublicKey()).to.be.equal(
        donPublicKey,
      )
    })
  })

  describe('#sendSimpleRequestWithJavaScript', () => {
    it('emits events from the client and the oracle contracts', async () => {
      await expect(
        client
          .connect(roles.defaultAccount)
          .sendSimpleRequestWithJavaScript('function run() {}', subscriptionId),
      )
        .to.emit(client, 'RequestSent')
        .withArgs(anyValue)
        .to.emit(oracle, 'OracleRequest')
        .withArgs(
          anyValue,
          client.address,
          await roles.defaultAccount.getAddress(),
          subscriptionId,
          await roles.defaultAccount.getAddress(),
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
      assert.equal(5, args.length)
      const decoded = await decodeDietCBOR(args[4])
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
      const tx = await client.sendSimpleRequestWithJavaScript(
        'function run(){return response}',
        subscriptionId,
      )

      const { events } = await tx.wait()
      const requestId = getEventArg(events, 'RequestSent', 0)
      await expect(tx).to.emit(client, 'RequestSent').withArgs(requestId)

      const response = stringToBytes('response')
      const error = stringToBytes('')
      const abi = ethers.utils.defaultAbiCoder

      const report = abi.encode(
        ['bytes32[]', 'bytes[]', 'bytes[]'],
        [[ethers.utils.hexZeroPad(requestId, 32)], [response], [error]],
      )

      await expect(oracle.callReport(report))
        .to.emit(oracle, 'OracleResponse')
        .withArgs(requestId)
        .to.emit(registry, 'BillingEnd')
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId, response, error)
    })
  })
})
