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
const fallbackNativePerUnitLink = 5000000000000000
export const coordinatorConfig: CoordinatorConfig = {
  maxCallbackGasLimit: 1_000_000,
  feedStalenessSeconds: 86_400,
  gasOverheadBeforeCallback: 44_615,
  gasOverheadAfterCallback: 44_615,
  requestTimeoutSeconds: 300,
  donFee: 0,
  fallbackNativePerUnitLink: BigNumber.from(fallbackNativePerUnitLink),
  maxSupportedRequestDataVersion: 1,
}
export const accessControlMockPublicKey =
  '0x32237412cC0321f56422d206e505dB4B3871AF5c'
export const accessControlMockPrivateKey =
  '2e8c8eaff4159e59711b42424c1555af1b78409e12c6f9c69a6a986d75442b20'
export type AccessControlConfig = {
  enabled: boolean
  proofSignerPublicKey: string // address
}
export const accessControlConfig: AccessControlConfig = {
  enabled: true,
  proofSignerPublicKey: accessControlMockPublicKey,
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
    roles.defaultAccount,
  )
  const accessControlFactory = await ethers.getContractFactory(
    'src/v0.8/functions/dev/1_0_0/accessControl/TermsOfServiceAllowList.sol:TermsOfServiceAllowList',
    roles.defaultAccount,
  )
  const clientTestHelperFactory = await ethers.getContractFactory(
    'src/v0.8/functions/tests/1_0_0/testhelpers/FunctionsClientTestHelper.sol:FunctionsClientTestHelper',
    roles.consumer,
  )
  const linkTokenFactory = await ethers.getContractFactory(
    'src/v0.4/LinkToken.sol:LinkToken',
    roles.defaultAccount,
  )
  const mockAggregatorV3Factory = await ethers.getContractFactory(
    'src/v0.7/tests/MockV3Aggregator.sol:MockV3Aggregator',
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
  const messageHash = await accessControl.getMessageHash(
    acceptorAddress,
    recipientAddress,
  )
  const wallet = new ethers.Wallet(accessControlMockPrivateKey)
  const proof = await wallet.signMessage(ethers.utils.arrayify(messageHash))
  return accessControl
    .connect(acceptor)
    .acceptTermsOfService(acceptorAddress, recipientAddress, proof)
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
        BigNumber.from('54666805176129187'),
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
    const routerConfigBytes = ethers.utils.defaultAbiCoder.encode(
      ['uint96', 'bytes4'],
      [...Object.values(functionsRouterConfig)],
    )
    const startingTimelockBlocks = 0
    const maxTimelockBlocks = 20
    const router = await factories.functionsRouterFactory
      .connect(roles.defaultAccount)
      .deploy(
        startingTimelockBlocks,
        maxTimelockBlocks,
        linkToken.address,
        routerConfigBytes,
      )
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
    const accessControlConfigBytes = ethers.utils.defaultAbiCoder.encode(
      ['bool', 'address'],
      [...Object.values(accessControlConfig)],
    )
    const accessControl = await factories.accessControlFactory
      .connect(roles.defaultAccount)
      .deploy(router.address, accessControlConfigBytes)
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
