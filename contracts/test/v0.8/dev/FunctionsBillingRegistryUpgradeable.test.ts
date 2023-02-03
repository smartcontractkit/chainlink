import { ethers, upgrades } from 'hardhat'
import { expect } from 'chai'
import { BigNumber, Contract, ContractFactory, Signer } from 'ethers'
import { Roles, getUsers } from '../../test-helpers/setup'

let functionsOracleFactory: ContractFactory
let clientTestHelperFactory: ContractFactory
let functionsBillingRegistryFactory: ContractFactory
let linkTokenFactory: ContractFactory
let mockAggregatorV3Factory: ContractFactory
let roles: Roles
let subOwner: Signer
let subOwnerAddress: string
let consumer: Signer
let consumerAddress: string
let stranger: Signer
let strangerAddress: string

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

const fundedLink = '1000000000000000000'
const linkEth = BigNumber.from(5021530000000000)

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

  functionsOracleFactory = await ethers.getContractFactory(
    'src/v0.8/tests/FunctionsOracleHelper.sol:FunctionsOracleHelper',
    roles.defaultAccount,
  )

  clientTestHelperFactory = await ethers.getContractFactory(
    'src/v0.8/tests/FunctionsClientTestHelper.sol:FunctionsClientTestHelper',
    roles.consumer,
  )

  functionsBillingRegistryFactory = await ethers.getContractFactory(
    'src/v0.8/dev/functions/FunctionsBillingRegistry.sol:FunctionsBillingRegistry',
    roles.consumer,
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

describe('FunctionsRegistryUpgradeable', () => {
  let registry: Contract
  let oracle: Contract
  let client: Contract
  let linkToken: Contract
  let mockLinkEth: Contract

  beforeEach(async () => {
    const { roles } = await getUsers()
    subOwner = roles.consumer
    subOwnerAddress = await subOwner.getAddress()
    consumer = roles.consumer2
    consumerAddress = await consumer.getAddress()
    stranger = roles.stranger
    strangerAddress = await stranger.getAddress()

    // Deploy
    linkToken = await linkTokenFactory.connect(roles.defaultAccount).deploy()
    mockLinkEth = await mockAggregatorV3Factory.deploy(0, linkEth)
    oracle = await functionsOracleFactory.connect(roles.defaultAccount).deploy()
    registry = await upgrades.deployProxy(
      functionsBillingRegistryFactory.connect(roles.defaultAccount),
      [linkToken.address, mockLinkEth.address, oracle.address],
    )
    client = await clientTestHelperFactory
      .connect(roles.consumer)
      .deploy(oracle.address)

    // Setup contracts
    await oracle.setRegistry(registry.address)
    await oracle.deactivateAuthorizedReceiver()

    // Setup accounts
    await linkToken.transfer(
      subOwnerAddress,
      BigNumber.from(fundedLink), // 1 LINK
    )
    await linkToken.transfer(
      strangerAddress,
      BigNumber.from(fundedLink), // 1 LINK
    )
  })

  async function createSubscription(
    owner: Signer,
    consumers: string[],
  ): Promise<number> {
    const tx = await registry.connect(owner).createSubscription()
    const receipt = await tx.wait()
    const subId = receipt.events[0].args['subscriptionId'].toNumber()
    for (let i = 0; i < consumers.length; i++) {
      await registry.connect(owner).addConsumer(subId, consumers[i])
    }
    return subId
  }

  describe('Upgrades', () => {
    let subId: number
    let requestId: string

    beforeEach(async () => {
      subId = await createSubscription(subOwner, [consumerAddress])

      await registry.setConfig(
        config.maxGasLimit,
        config.stalenessSeconds,
        config.gasAfterPaymentCalculation,
        config.weiPerUnitLink,
        config.gasOverhead,
        config.requestTimeoutSeconds,
      )

      await linkToken
        .connect(subOwner)
        .transferAndCall(
          registry.address,
          BigNumber.from('1000000000000000000'),
          ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
        )
      await registry.connect(subOwner).addConsumer(subId, client.address)
      await registry.connect(roles.defaultAccount).reg
      await registry.setAuthorizedSenders([oracle.address])

      const request = await client
        .connect(consumer)
        .sendSimpleRequestWithJavaScript(`return 'hello world'`, subId)
      requestId = (await request.wait()).events[3].args[0]
    })

    it('is successful when deployed behind a proxy', async () => {
      const report = encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('hello world'),
        stringToHex(''),
      )
      await expect(
        oracle
          .connect(roles.oracleNode)
          .callReport(report, { gasLimit: 500_000 }),
      ).to.emit(registry, 'BillingEnd')
    })

    it('can be upgraded to a new implementation', async () => {
      // Upgrade the implementation contract
      const functionsBillingRegistryMigrationFactory =
        await ethers.getContractFactory(
          'src/v0.8/tests/FunctionsBillingRegistryMigration.sol:FunctionsBillingRegistryMigration',
          roles.consumer,
        )
      const upgradedRegistry = await upgrades.upgradeProxy(
        registry.address,
        functionsBillingRegistryMigrationFactory.connect(roles.defaultAccount),
      )

      // Check that upgrade was successful
      const dummyRequest = [
        '0x',
        {
          subscriptionId: subId,
          client: client.address,
          gasLimit: 0,
          gasPrice: 0,
        },
      ]
      const registryFee = await upgradedRegistry.getRequiredFee(...dummyRequest)
      expect(registryFee).to.equal(1)

      // Check config is the same
      const currentConfig = await upgradedRegistry.getConfig()
      expect(currentConfig.maxGasLimit).to.equal(config.maxGasLimit)
      expect(currentConfig.stalenessSeconds).to.equal(config.stalenessSeconds)
      expect(currentConfig.gasAfterPaymentCalculation).to.equal(
        config.gasAfterPaymentCalculation,
      )
      expect(currentConfig.fallbackWeiPerUnitLink).to.equal(
        config.weiPerUnitLink,
      )
      expect(currentConfig.gasOverhead).to.equal(config.gasOverhead)

      // Check funds are the same
      const subscription = await upgradedRegistry.getSubscription(subId)
      expect(subscription.balance).to.equal(fundedLink)

      // Check request fulfillment still works
      const report = encodeReport(
        ethers.utils.hexZeroPad(requestId, 32),
        stringToHex('hello world'),
        stringToHex(''),
      )
      await expect(
        oracle
          .connect(roles.oracleNode)
          .callReport(report, { gasLimit: 500_000 }),
      ).to.emit(registry, 'BillingEnd')
    })
  })
})
