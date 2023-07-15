import { ethers } from 'hardhat'
import { BigNumber, ContractFactory, Signer, Contract, providers } from 'ethers'
import { Roles, getUsers } from '../../../test-helpers/setup'

export type FunctionsRoles = Roles & {
  subOwner: Signer
  subOwnerAddress: string
  consumer: Signer
  consumerAddress: string
  stranger: Signer
  strangerAddress: string
}

export type FunctionsFactories = {
  functionsRouterFactory: ContractFactory
  functionsCoordinatorFactory: ContractFactory
  clientTestHelperFactory: ContractFactory
  linkTokenFactory: ContractFactory
  mockAggregatorV3Factory: ContractFactory
}
export type FunctionsContracts = {
  router: Contract
  coordinator: Contract
  client: Contract
  linkToken: Contract
  mockLinkEth: Contract
}

export const ids = {
  routerId: ethers.utils.formatBytes32String(''),
  donId: ethers.utils.formatBytes32String('1'),
  donId2: ethers.utils.formatBytes32String('2'),
  donId3: ethers.utils.formatBytes32String('3'),
  donId4: ethers.utils.formatBytes32String('4'),
  donId5: ethers.utils.formatBytes32String('5'),
}

export const anyValue = () => true

export const stringToHex = (s: string) => {
  return ethers.utils.hexlify(ethers.utils.toUtf8Bytes(s))
}

const linkEth = BigNumber.from(5021530000000000)

export const encodeReport = (
  requestId: string,
  result: string,
  err: string,
) => {
  const abi = ethers.utils.defaultAbiCoder
  return abi.encode(
    ['bytes32[]', 'bytes[]', 'bytes[]'],
    [[requestId], [result], [err]],
  )
}

export type FunctionsRouterConfig = {
  adminFee: number
  handleOracleFulfillmentSelector: string
}
export const functionsRouterConfig: FunctionsRouterConfig = {
  adminFee: 0,
  handleOracleFulfillmentSelector: '0x0ca76175',
}
export type CoordinatorConfig = {
  maxCallbackGasLimit: number
  feedStalenessSeconds: number
  gasOverheadBeforeCallback: number
  gasOverheadAfterCallback: number
  requestTimeoutSeconds: number
  donFee: number
  fallbackNativePerUnitLink: BigNumber
  maxSupportedRequestDataVersion: number
}
export const coordinatorConfig: CoordinatorConfig = {
  maxCallbackGasLimit: 1_000_000,
  feedStalenessSeconds: 86_400,
  gasOverheadBeforeCallback:
    21_000 + 5_000 + 2_100 + 20_000 + 2 * 2_100 - 15_000 + 7_315,
  gasOverheadAfterCallback:
    21_000 + 5_000 + 2_100 + 20_000 + 2 * 2_100 - 15_000 + 7_315,
  requestTimeoutSeconds: 300,
  donFee: 0,
  fallbackNativePerUnitLink: BigNumber.from(5000000000000000),
  maxSupportedRequestDataVersion: 1,
}

export async function setupRolesAndFactories(): Promise<{
  roles: FunctionsRoles
  factories: FunctionsFactories
}> {
  const roles = (await getUsers()).roles
  const functionsRouterFactory = await ethers.getContractFactory(
    'src/v0.8/functions/dev/1_0_0/FunctionsRouter.sol:FunctionsRouter',
    roles.defaultAccount,
  )
  const functionsCoordinatorFactory = await ethers.getContractFactory(
    'src/v0.8/functions/tests/1_0_0/testhelpers/FunctionsCoordinatorTestHelper.sol:FunctionsCoordinatorTestHelper',
    roles.consumer,
  )
  const clientTestHelperFactory = await ethers.getContractFactory(
    'src/v0.8/functions/tests/1_0_0/testhelpers/FunctionsClientTestHelper.sol:FunctionsClientTestHelper',
    roles.consumer,
  )
  const linkTokenFactory = await ethers.getContractFactory(
    'src/v0.4/LinkToken.sol:LinkToken',
    roles.consumer,
  )
  const mockAggregatorV3Factory = await ethers.getContractFactory(
    'src/v0.7/tests/MockV3Aggregator.sol:MockV3Aggregator',
    roles.consumer,
  )
  return {
    roles: {
      ...roles,
      subOwner: roles.consumer,
      subOwnerAddress: await roles.consumer.getAddress(),
      consumer: roles.consumer2,
      consumerAddress: await roles.consumer2.getAddress(),
      stranger: roles.stranger,
      strangerAddress: await roles.stranger.getAddress(),
    },
    factories: {
      functionsRouterFactory,
      functionsCoordinatorFactory,
      clientTestHelperFactory,
      linkTokenFactory,
      mockAggregatorV3Factory,
    },
  }
}

export async function createSubscription(
  owner: Signer,
  consumers: string[],
  router: Contract,
  linkToken?: Contract,
): Promise<number> {
  const tx = await router.connect(owner).createSubscription()
  const receipt = await tx.wait()
  const subId = receipt.events[0].args['subscriptionId'].toNumber()
  for (let i = 0; i < consumers.length; i++) {
    await router.connect(owner).addConsumer(subId, consumers[i])
  }
  if (linkToken)
    await linkToken
      .connect(owner)
      .transferAndCall(
        router.address,
        BigNumber.from('54666805176129187'),
        ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
      )
  return subId
}

export function getSetupFactory() {
  let contracts: FunctionsContracts
  let factories: FunctionsFactories
  let roles: FunctionsRoles

  before(async () => {
    const { roles: r, factories: f } = await setupRolesAndFactories()
    factories = f
    roles = r
  })

  beforeEach(async () => {
    // Deploy
    const linkToken = await factories.linkTokenFactory
      .connect(roles.defaultAccount)
      .deploy()
    const mockLinkEth = await factories.mockAggregatorV3Factory.deploy(
      0,
      linkEth,
    )
    const routerConfigBytes = ethers.utils.defaultAbiCoder.encode(
      ['uint96', 'bytes4'],
      [...Object.values(functionsRouterConfig)],
    )
    const router = await factories.functionsRouterFactory
      .connect(roles.defaultAccount)
      .deploy(0, 20, linkToken.address, routerConfigBytes)
    const coordinatorConfigBytes = ethers.utils.defaultAbiCoder.encode(
      [
        'uint32',
        'uint32',
        'uint32',
        'uint32',
        'int256',
        'uint32',
        'uint96',
        'uint16',
      ],
      [...Object.values(coordinatorConfig)],
    )
    const coordinator = await factories.functionsCoordinatorFactory
      .connect(roles.defaultAccount)
      .deploy(router.address, coordinatorConfigBytes, mockLinkEth.address)
    const client = await factories.clientTestHelperFactory
      .connect(roles.consumer)
      .deploy(router.address)

    // Setup accounts
    await linkToken.transfer(
      roles.subOwnerAddress,
      BigNumber.from('1000000000000000000'), // 1 LINK
    )
    await linkToken.transfer(
      roles.strangerAddress,
      BigNumber.from('1000000000000000000'), // 1 LINK
    )

    await router.proposeContractsUpdate(
      [ids.donId],
      [ethers.constants.AddressZero],
      [coordinator.address],
    )
    await router.updateContracts()

    contracts = {
      client,
      coordinator,
      router,
      linkToken,
      mockLinkEth,
    }
  })

  return () => {
    return { contracts, factories, roles }
  }
}

export function getEventArg(events: any, eventName: string, argIndex: number) {
  if (Array.isArray(events)) {
    const event = events.find((e: any) => e.event == eventName)
    if (event && Array.isArray(event.args) && event.args.length > 0) {
      return event.args[argIndex]
    }
  }
  return undefined
}

export async function parseOracleRequestEventArgs(
  tx: providers.TransactionResponse,
) {
  const receipt = await tx.wait()
  const data = receipt.logs?.[1].data
  // NOTE: indexed args are on topics, not data
  return ethers.utils.defaultAbiCoder.decode(
    ['address', 'uint64', 'address', 'bytes', 'uint16'],
    data ?? '',
  )
}
