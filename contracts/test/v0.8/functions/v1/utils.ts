import { ethers } from 'hardhat'
import { BigNumber, ContractFactory, Signer, Contract, providers } from 'ethers'
import { Roles, getUsers } from '../../../test-helpers/setup'
import { EventFragment } from 'ethers/lib/utils'

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
  accessControlFactory: ContractFactory
}
export type FunctionsContracts = {
  router: Contract
  coordinator: Contract
  client: Contract
  linkToken: Contract
  mockLinkEth: Contract
  accessControl: Contract
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

export const encodeReport = async (
  requestId: string,
  result: string,
  err: string,
  onchainMetadata: any,
  offchainMetadata: string,
) => {
  const functionsResponse = await ethers.getContractFactory(
    'src/v0.8/functions/dev/v1_0_0/FunctionsCoordinator.sol:FunctionsCoordinator',
  )
  const onchainMetadataBytes = functionsResponse.interface._abiCoder.encode(
    [
      getEventInputs(
        Object.values(functionsResponse.interface.events),
        'OracleRequest',
        9,
      ),
    ],
    [[...onchainMetadata]],
  )
  const abi = ethers.utils.defaultAbiCoder
  return abi.encode(
    ['bytes32[]', 'bytes[]', 'bytes[]', 'bytes[]', 'bytes[]'],
    [[requestId], [result], [err], [onchainMetadataBytes], [offchainMetadata]],
  )
}

export type FunctionsRouterConfig = {
  maxConsumersPerSubscription: number
  adminFee: number
  handleOracleFulfillmentSelector: string
  maxCallbackGasLimits: number[]
  gasForCallExactCheck: number
  subscriptionDepositMinimumRequests: number
  subscriptionDepositJuels: BigNumber
}
export const functionsRouterConfig: FunctionsRouterConfig = {
  maxConsumersPerSubscription: 100,
  adminFee: 0,
  handleOracleFulfillmentSelector: '0x0ca76175',
  maxCallbackGasLimits: [300_000, 500_000, 1_000_000],
  gasForCallExactCheck: 5000,
  subscriptionDepositMinimumRequests: 10,
  subscriptionDepositJuels: BigNumber.from('1000000000000000000'),
}
export type CoordinatorConfig = {
  feedStalenessSeconds: number
  gasOverheadBeforeCallback: number
  gasOverheadAfterCallback: number
  requestTimeoutSeconds: number
  donFee: number
  maxSupportedRequestDataVersion: number
  fulfillmentGasPriceOverEstimationBP: number
  fallbackNativePerUnitLink: BigNumber
}
const fallbackNativePerUnitLink = 5000000000000000
export const coordinatorConfig: CoordinatorConfig = {
  feedStalenessSeconds: 86_400,
  gasOverheadBeforeCallback: 44_615,
  gasOverheadAfterCallback: 44_615,
  requestTimeoutSeconds: 300,
  donFee: 0,
  maxSupportedRequestDataVersion: 1,
  fulfillmentGasPriceOverEstimationBP: 0,
  fallbackNativePerUnitLink: BigNumber.from(fallbackNativePerUnitLink),
}
export const accessControlMockPublicKey = ethers.utils.getAddress(
  '0x32237412cC0321f56422d206e505dB4B3871AF5c',
)
export const accessControlMockPrivateKey =
  '2e8c8eaff4159e59711b42424c1555af1b78409e12c6f9c69a6a986d75442b20'
export type AccessControlConfig = {
  enabled: boolean
  signerPublicKey: string // address
}
export const accessControlConfig: AccessControlConfig = {
  enabled: true,
  signerPublicKey: accessControlMockPublicKey,
}

export async function setupRolesAndFactories(): Promise<{
  roles: FunctionsRoles
  factories: FunctionsFactories
}> {
  const roles = (await getUsers()).roles
  const functionsRouterFactory = await ethers.getContractFactory(
    'src/v0.8/functions/dev/v1_0_0/FunctionsRouter.sol:FunctionsRouter',
    roles.defaultAccount,
  )
  const functionsCoordinatorFactory = await ethers.getContractFactory(
    'src/v0.8/functions/tests/v1_0_0/testhelpers/FunctionsCoordinatorTestHelper.sol:FunctionsCoordinatorTestHelper',
    roles.defaultAccount,
  )
  const accessControlFactory = await ethers.getContractFactory(
    'src/v0.8/functions/dev/v1_0_0/accessControl/TermsOfServiceAllowList.sol:TermsOfServiceAllowList',
    roles.defaultAccount,
  )
  const clientTestHelperFactory = await ethers.getContractFactory(
    'src/v0.8/functions/tests/v1_0_0/testhelpers/FunctionsClientTestHelper.sol:FunctionsClientTestHelper',
    roles.consumer,
  )
  const linkTokenFactory = await ethers.getContractFactory(
    'src/v0.8/mocks/MockLinkToken.sol:MockLinkToken',
    roles.defaultAccount,
  )
  const mockAggregatorV3Factory = await ethers.getContractFactory(
    'src/v0.8/tests/MockV3Aggregator.sol:MockV3Aggregator',
    roles.defaultAccount,
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
      accessControlFactory,
    },
  }
}

export async function acceptTermsOfService(
  accessControl: Contract,
  acceptor: Signer,
  recipientAddress: string,
) {
  const acceptorAddress = await acceptor.getAddress()
  const message = await accessControl.getMessage(
    acceptorAddress,
    recipientAddress,
  )
  const wallet = new ethers.Wallet(accessControlMockPrivateKey)
  const flatSignature = await wallet.signMessage(ethers.utils.arrayify(message))
  const { r, s, v } = ethers.utils.splitSignature(flatSignature)
  return accessControl
    .connect(acceptor)
    .acceptTermsOfService(acceptorAddress, recipientAddress, r, s, v)
}

export async function createSubscription(
  owner: Signer,
  consumers: string[],
  router: Contract,
  accessControl: Contract,
  linkToken?: Contract,
): Promise<number> {
  const ownerAddress = await owner.getAddress()
  await acceptTermsOfService(accessControl, owner, ownerAddress)
  const tx = await router.connect(owner).createSubscription()
  const receipt = await tx.wait()
  const subId = receipt.events[0].args['subscriptionId'].toNumber()
  for (let i = 0; i < consumers.length; i++) {
    await router.connect(owner).addConsumer(subId, consumers[i])
  }
  if (linkToken) {
    await linkToken
      .connect(owner)
      .transferAndCall(
        router.address,
        BigNumber.from('1000000000000000000'),
        ethers.utils.defaultAbiCoder.encode(['uint64'], [subId]),
      )
  }
  return subId
}

export function getSetupFactory(): () => {
  contracts: FunctionsContracts
  factories: FunctionsFactories
  roles: FunctionsRoles
} {
  let contracts: FunctionsContracts
  let factories: FunctionsFactories
  let roles: FunctionsRoles

  before(async () => {
    const { roles: r, factories: f } = await setupRolesAndFactories()
    factories = f
    roles = r
  })

  beforeEach(async () => {
    const linkEthRate = BigNumber.from(5021530000000000)

    // Deploy
    const linkToken = await factories.linkTokenFactory
      .connect(roles.defaultAccount)
      .deploy()

    const mockLinkEth = await factories.mockAggregatorV3Factory.deploy(
      0,
      linkEthRate,
    )

    const router = await factories.functionsRouterFactory
      .connect(roles.defaultAccount)
      .deploy(linkToken.address, functionsRouterConfig)

    const coordinator = await factories.functionsCoordinatorFactory
      .connect(roles.defaultAccount)
      .deploy(router.address, coordinatorConfig, mockLinkEth.address)

    const accessControl = await factories.accessControlFactory
      .connect(roles.defaultAccount)
      .deploy(accessControlConfig)

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

    const allowListId = await router.getAllowListId()
    await router.proposeContractsUpdate(
      [ids.donId, allowListId],
      [coordinator.address, accessControl.address],
    )
    await router.updateContracts()

    contracts = {
      client,
      coordinator,
      router,
      linkToken,
      mockLinkEth,
      accessControl,
    }
  })

  return () => {
    return { contracts, factories, roles }
  }
}

export function getEventArg(events: any, eventName: string, argIndex: number) {
  if (Array.isArray(events)) {
    const event = events.find((e: any) => e.event === eventName)
    if (event && Array.isArray(event.args) && event.args.length > 0) {
      return event.args[argIndex]
    }
  }
  return undefined
}

export function getEventInputs(
  events: EventFragment[],
  eventName: string,
  argIndex: number,
) {
  if (Array.isArray(events)) {
    const event = events.find((e) => e.name.includes(eventName))
    if (event && Array.isArray(event.inputs) && event.inputs.length > 0) {
      return event.inputs[argIndex]
    }
  }
  throw 'Not found'
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
