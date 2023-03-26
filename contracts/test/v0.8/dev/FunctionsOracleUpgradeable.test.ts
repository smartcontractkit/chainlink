import { ethers, upgrades } from 'hardhat'
import { expect } from 'chai'
import { BigNumber, Contract, ContractFactory } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'

let functionsOracleOriginalFactory: ContractFactory
let clientTestHelperFactory: ContractFactory
let functionsBillingRegistryFactory: ContractFactory
let linkTokenFactory: ContractFactory
let mockAggregatorV3Factory: ContractFactory
let roles: Roles

const stringToHex = (s: string) => {
  return ethers.utils.hexlify(ethers.utils.toUtf8Bytes(s))
}

const encodeReport = (requestId: string, result: string, err: string) => {
  const abi = ethers.utils.defaultAbiCoder
  return abi.encode(
    ['bytes32[]', 'bytes[]', 'bytes[]'],
    [[requestId], [result], [err]],
  )
}

type RegistryConfig = {
  maxGasLimit: number
  stalenessSeconds: number
  gasAfterPaymentCalculation: number
  weiPerUnitLink: BigNumber
  gasOverhead: number
  requestTimeoutSeconds: number
}
const config: RegistryConfig = {
  maxGasLimit: 1_000_000,
  stalenessSeconds: 86_400,
  gasAfterPaymentCalculation:
    21_000 + 5_000 + 2_100 + 20_000 + 2 * 2_100 - 15_000 + 7_315,
  weiPerUnitLink: BigNumber.from('5000000000000000'),
  gasOverhead: 100_000,
  requestTimeoutSeconds: 300,
}

before(async () => {
  roles = (await getUsers()).roles

  functionsOracleOriginalFactory = await ethers.getContractFactory(
    'src/v0.8/tests/FunctionsOracleOriginalHelper.sol:FunctionsOracleOriginalHelper',
    roles.defaultAccount,
  )

  clientTestHelperFactory = await ethers.getContractFactory(
    'src/v0.8/tests/FunctionsClientTestHelper.sol:FunctionsClientTestHelper',
    roles.consumer,
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

describe('FunctionsOracleUpgradeable', () => {
  let subscriptionId: number
  let client: Contract
  let oracle: Contract
  let registry: Contract
  let linkToken: Contract
  let mockLinkEth: Contract
  let transmitters: string[]

  beforeEach(async () => {
    // Deploy contracts
    linkToken = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    mockLinkEth = await mockAggregatorV3Factory.deploy(
      0,
      ethers.BigNumber.from(5021530000000000),
    )
    oracle = await upgrades.deployProxy(
      functionsOracleOriginalFactory.connect(roles.defaultAccount),
    )
    registry = await functionsBillingRegistryFactory
      .connect(roles.defaultAccount)
      .deploy(linkToken.address, mockLinkEth.address, oracle.address)

    // Setup contracts
    await oracle.setRegistry(registry.address)
    await oracle.deactivateAuthorizedReceiver()
    client = await clientTestHelperFactory
      .connect(roles.defaultAccount)
      .deploy(oracle.address)
    await registry.setAuthorizedSenders([oracle.address])

    await registry.setConfig(
      config.maxGasLimit,
      config.stalenessSeconds,
      config.gasAfterPaymentCalculation,
      config.weiPerUnitLink,
      config.gasOverhead,
      config.requestTimeoutSeconds,
    )

    // Setup accounts
    const createSubTx = await registry
      .connect(roles.defaultAccount)
      .createSubscription()
    const receipt = await createSubTx.wait()
    subscriptionId = receipt.events[0].args['subscriptionId'].toNumber()

    await registry
      .connect(roles.defaultAccount)
      .addConsumer(subscriptionId, await roles.defaultAccount.getAddress())

    await registry
      .connect(roles.defaultAccount)
      .addConsumer(subscriptionId, client.address)

    await linkToken
      .connect(roles.defaultAccount)
      .transferAndCall(
        registry.address,
        ethers.BigNumber.from('300938394174049741'),
        ethers.utils.defaultAbiCoder.encode(['uint64'], [subscriptionId]),
      )

    const signers = Array.from(
      [0, 0, 0, 0],
      (_) => ethers.Wallet.createRandom().address,
    )
    transmitters = [
      await roles.oracleNode1.getAddress(),
      await roles.oracleNode2.getAddress(),
      await roles.oracleNode3.getAddress(),
      await roles.oracleNode4.getAddress(),
    ]
    await oracle.setConfig(signers, transmitters, 1, [], 1, [])
  })

  describe('Upgrades', () => {
    const placeTestRequest = async () => {
      const requestId = await client
        .connect(roles.oracleNode)
        .callStatic.sendSimpleRequestWithJavaScript(
          'function(){}',
          subscriptionId,
        )
      await expect(
        client
          .connect(roles.oracleNode)
          .sendSimpleRequestWithJavaScript('function(){}', subscriptionId),
      )
        .to.emit(client, 'RequestSent')
        .withArgs(requestId)
      return requestId
    }

    async function migrateAndCheck(factoryPath: string): Promise<Contract> {
      const functionsOracleMigrationFactory = await ethers.getContractFactory(
        factoryPath,
        roles.consumer,
      )

      // Upgrade the implementation contract
      const upgradedOracle = await upgrades.upgradeProxy(
        oracle.address,
        functionsOracleMigrationFactory.connect(roles.defaultAccount),
      )

      // Check request fulfillment still works
      const requestId = await placeTestRequest()
      const report = encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('hello world'),
        stringToHex(''),
      )
      await expect(
        upgradedOracle
          .connect(roles.oracleNode)
          .callReport(report, { gasLimit: 500_000 }),
      ).to.emit(registry, 'BillingEnd')

      return upgradedOracle
    }

    it('is successful when deployed behind a proxy', async () => {
      const requestId1 = await placeTestRequest()
      const requestId2 = await placeTestRequest()
      const result1 = stringToHex('result1')
      const result2 = stringToHex('result2')
      const err = stringToHex('')

      const abi = ethers.utils.defaultAbiCoder
      const report = abi.encode(
        ['bytes32[]', 'bytes[]', 'bytes[]'],
        [
          [requestId1, requestId2],
          [result1, result2],
          [err, err],
        ],
      )

      await expect(
        oracle
          .connect(roles.oracleNode)
          .callReport(report, { gasLimit: 300_000 }),
      )
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId1, result1, err)
        .to.emit(client, 'FulfillRequestInvoked')
        .withArgs(requestId2, result2, err)
    })

    it('can be upgraded to a new implementation', async () => {
      const upgradedOracle = await migrateAndCheck(
        'src/v0.8/tests/FunctionsOracleMigrationHelper.sol:FunctionsOracleMigrationHelper',
      )

      // Check that upgrade was successful
      const dummyRequest = [
        '0x',
        {
          subscriptionId,
          client: client.address,
          gasLimit: 0,
          gasPrice: 0,
        },
      ]
      const registryFee = await upgradedOracle.getRequiredFee(...dummyRequest)
      expect(registryFee).to.equal(1)
    })

    it('can be upgraded to the latest implementation', async () => {
      await migrateAndCheck(
        'src/v0.8/tests/FunctionsOracleUpgradeableHelper.sol:FunctionsOracleUpgradeableHelper',
      )
    })
  })
})
