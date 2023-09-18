import { ethers } from 'hardhat'
import { expect } from 'chai'
import { Contract, ContractFactory } from 'ethers'
import { Roles, getUsers } from '../../../test-helpers/setup'
import { stringToBytes } from '../../../test-helpers/helpers'

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

const baselineGasUsed = 641560
let currentGasUsed = 0

before(async () => {
  roles = (await getUsers()).roles

  concreteFunctionsClientFactory = await ethers.getContractFactory(
    'src/v0.8/functions/tests/v0_0_0/testhelpers/FunctionsClientTestHelper.sol:FunctionsClientTestHelper',
    roles.defaultAccount,
  )
  functionsOracleFactory = await ethers.getContractFactory(
    'src/v0.8/functions/tests/v0_0_0/testhelpers/FunctionsOracleHelper.sol:FunctionsOracleHelper',
    roles.defaultAccount,
  )

  functionsBillingRegistryFactory = await ethers.getContractFactory(
    'src/v0.8/functions/tests/v0_0_0/testhelpers/FunctionsBillingRegistryWithInit.sol:FunctionsBillingRegistryWithInit',
    roles.defaultAccount,
  )

  linkTokenFactory = await ethers.getContractFactory(
    'src/v0.8/mocks/MockLinkToken.sol:MockLinkToken',
    roles.consumer,
  )

  mockAggregatorV3Factory = await ethers.getContractFactory(
    'src/v0.8/tests/MockV3Aggregator.sol:MockV3Aggregator',
    roles.consumer,
  )
})

after(() => {
  const score = currentGasUsed - baselineGasUsed
  console.log(
    `\n               ⛳ Baseline gas used   : ${baselineGasUsed} gas`,
  )
  console.log(`\n               Current gas used   : ${currentGasUsed} gas`)
  console.log(`\n               🚩 Delta : ${score} gas`)
})

let subscriptionId: number

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
})

describe('Gas', () => {
  it('uses the expected amount of gas', async () => {
    // Setup accounts
    const createSubTx = await registry
      .connect(roles.defaultAccount)
      .createSubscription()
    const createSubscriptionTxReceipt = await createSubTx.wait()
    subscriptionId =
      createSubscriptionTxReceipt.events[0].args['subscriptionId'].toNumber()
    const createSubscriptionGasUsed = createSubscriptionTxReceipt.gasUsed

    const addConsumerTx = await registry
      .connect(roles.defaultAccount)
      .addConsumer(subscriptionId, client.address)
    const { gasUsed: addConsumerTxGasUsed } = await addConsumerTx.wait()

    const transferAndCallTx = await linkToken
      .connect(roles.defaultAccount)
      .transferAndCall(
        registry.address,
        ethers.BigNumber.from('115957983815660167'),
        ethers.utils.defaultAbiCoder.encode(['uint64'], [subscriptionId]),
      )
    const { gasUsed: transferAndCallTxGasUsed } = await transferAndCallTx.wait()

    const requestTx = await client.sendSimpleRequestWithJavaScript(
      'function run(){return response}',
      subscriptionId,
    )

    const { events, gasUsed: requestTxGasUsed } = await requestTx.wait()
    const requestId = getEventArg(events, 'RequestSent', 0)
    await expect(requestTx).to.emit(client, 'RequestSent').withArgs(requestId)

    const response = stringToBytes('response')
    const error = stringToBytes('')
    const abi = ethers.utils.defaultAbiCoder

    const report = abi.encode(
      ['bytes32[]', 'bytes[]', 'bytes[]'],
      [[ethers.utils.hexZeroPad(requestId, 32)], [response], [error]],
    )

    const fulfillmentTx = await oracle.callReport(report, {
      gasLimit: 300_000,
    })

    const { gasUsed: fulfillmentTxGasUsed } = await fulfillmentTx.wait()

    currentGasUsed = createSubscriptionGasUsed
      .add(addConsumerTxGasUsed)
      .add(transferAndCallTxGasUsed)
      .add(requestTxGasUsed)
      .add(fulfillmentTxGasUsed)
      .toNumber()
  })
})
