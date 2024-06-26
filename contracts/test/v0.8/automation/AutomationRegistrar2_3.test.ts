import { ethers } from 'hardhat'
import {
  BigNumber,
  BigNumberish,
  BytesLike,
  Contract,
  ContractFactory,
  Signer,
} from 'ethers'
import { assert, expect } from 'chai'
import { evmRevert, evmRevertCustomError } from '../../test-helpers/matchers'
import { getUsers, Personas } from '../../test-helpers/setup'
import { MockV3Aggregator__factory as MockV3AggregatorFactory } from '../../../typechain/factories/MockV3Aggregator__factory'
import { UpkeepMock__factory as UpkeepMockFactory } from '../../../typechain/factories/UpkeepMock__factory'
import { ChainModuleBase__factory as ChainModuleBaseFactory } from '../../../typechain/factories/ChainModuleBase__factory'
import { MockV3Aggregator } from '../../../typechain/MockV3Aggregator'
import { UpkeepMock } from '../../../typechain/UpkeepMock'
import { randomAddress, toWei } from '../../test-helpers/helpers'
import { ChainModuleBase } from '../../../typechain/ChainModuleBase'
import { AutomationRegistrar2_3 as Registrar } from '../../../typechain/AutomationRegistrar2_3'
import { deployRegistry23 } from './helpers'
import { IAutomationRegistryMaster2_3 as IAutomationRegistry } from '../../../typechain'

// copied from AutomationRegistryBase2_3.sol
enum Trigger {
  CONDITION,
  LOG,
}
const zeroAddress = ethers.constants.AddressZero
const wrappedNativeTokenAddress = '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2'

type OnChainConfig = Parameters<IAutomationRegistry['setConfigTypeSafe']>[3]

let linkTokenFactory: ContractFactory
let mockV3AggregatorFactory: MockV3AggregatorFactory
let upkeepMockFactory: UpkeepMockFactory

let personas: Personas

before(async () => {
  personas = (await getUsers()).personas

  linkTokenFactory = await ethers.getContractFactory(
    'src/v0.8/shared/test/helpers/LinkTokenTestHelper.sol:LinkTokenTestHelper',
  )
  mockV3AggregatorFactory = (await ethers.getContractFactory(
    'src/v0.8/tests/MockV3Aggregator.sol:MockV3Aggregator',
  )) as unknown as MockV3AggregatorFactory
  upkeepMockFactory = await ethers.getContractFactory('UpkeepMock')
})

const errorMsgs = {
  onlyOwner: 'revert Only callable by owner',
  onlyAdmin: 'OnlyAdminOrOwner',
  hashPayload: 'HashMismatch',
  requestNotFound: 'RequestNotFound',
}

describe('AutomationRegistrar2_3', () => {
  const upkeepName = 'SampleUpkeep'

  const linkUSD = BigNumber.from('2000000000') // 1 LINK = $20
  const nativeUSD = BigNumber.from('400000000000') // 1 ETH = $4000
  const gasWei = BigNumber.from(100)
  const performGas = BigNumber.from(100000)
  const paymentPremiumPPB = BigNumber.from(250000000)
  const flatFeeMilliCents = BigNumber.from(0)
  const maxAllowedAutoApprove = 5
  const trigger = '0xdeadbeef'
  const offchainConfig = '0x01234567'
  const keepers = [
    randomAddress(),
    randomAddress(),
    randomAddress(),
    randomAddress(),
  ]

  const emptyBytes = '0x00'
  const stalenessSeconds = BigNumber.from(43820)
  const gasCeilingMultiplier = BigNumber.from(1)
  const checkGasLimit = BigNumber.from(20000000)
  const fallbackGasPrice = BigNumber.from(200)
  const fallbackLinkPrice = BigNumber.from(200000000)
  const fallbackNativePrice = BigNumber.from(200000000)
  const maxCheckDataSize = BigNumber.from(10000)
  const maxPerformDataSize = BigNumber.from(10000)
  const maxRevertDataSize = BigNumber.from(1000)
  const maxPerformGas = BigNumber.from(5000000)
  const minimumRegistrationAmount = BigNumber.from('1000000000000000000')
  const amount = BigNumber.from('5000000000000000000')
  const transcoder = ethers.constants.AddressZero
  const upkeepManager = ethers.Wallet.createRandom().address

  // Enum values are not auto exported in ABI so have to manually declare
  const autoApproveType_DISABLED = 0
  const autoApproveType_ENABLED_SENDER_ALLOWLIST = 1
  const autoApproveType_ENABLED_ALL = 2

  let owner: Signer
  let admin: Signer
  let someAddress: Signer
  let registrarOwner: Signer
  let stranger: Signer
  let requestSender: Signer

  let linkToken: Contract
  let linkUSDFeed: MockV3Aggregator
  let nativeUSDFeed: MockV3Aggregator
  let gasPriceFeed: MockV3Aggregator
  let mock: UpkeepMock
  let registry: IAutomationRegistry
  let registrar: Registrar
  let chainModuleBase: ChainModuleBase
  let chainModuleBaseFactory: ChainModuleBaseFactory
  let onchainConfig: OnChainConfig

  type RegistrationParams = {
    upkeepContract: string
    amount: BigNumberish
    adminAddress: string
    gasLimit: BigNumberish
    triggerType: BigNumberish
    billingToken: string
    name: string
    encryptedEmail: BytesLike
    checkData: BytesLike
    triggerConfig: BytesLike
    offchainConfig: BytesLike
  }

  function encodeRegistrationParams(params: RegistrationParams) {
    return (
      '0x' +
      registrar.interface
        .encodeFunctionData('registerUpkeep', [params])
        .slice(10)
    )
  }

  beforeEach(async () => {
    owner = personas.Default
    admin = personas.Neil
    someAddress = personas.Ned
    registrarOwner = personas.Nelly
    stranger = personas.Nancy
    requestSender = personas.Norbert

    linkToken = await linkTokenFactory.connect(owner).deploy()
    gasPriceFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(0, gasWei)
    linkUSDFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(8, linkUSD)
    nativeUSDFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(8, nativeUSD)

    chainModuleBaseFactory = await ethers.getContractFactory('ChainModuleBase')
    chainModuleBase = await chainModuleBaseFactory.connect(owner).deploy()

    registry = await deployRegistry23(
      owner,
      linkToken.address,
      linkUSDFeed.address,
      nativeUSDFeed.address,
      gasPriceFeed.address,
      zeroAddress,
      0, // onchain payout mode
      wrappedNativeTokenAddress,
    )

    mock = await upkeepMockFactory.deploy()

    const registrarFactory = await ethers.getContractFactory(
      'AutomationRegistrar2_3',
    )
    registrar = await registrarFactory.connect(registrarOwner).deploy(
      linkToken.address,
      registry.address,
      [
        {
          triggerType: Trigger.CONDITION,
          autoApproveType: autoApproveType_DISABLED,
          autoApproveMaxAllowed: 0,
        },
        {
          triggerType: Trigger.LOG,
          autoApproveType: autoApproveType_DISABLED,
          autoApproveMaxAllowed: 0,
        },
      ],
      [linkToken.address],
      [minimumRegistrationAmount],
      wrappedNativeTokenAddress,
    )

    await linkToken
      .connect(owner)
      .transfer(await requestSender.getAddress(), toWei('1000'))

    onchainConfig = {
      checkGasLimit,
      stalenessSeconds,
      gasCeilingMultiplier,
      maxCheckDataSize,
      maxPerformDataSize,
      maxRevertDataSize,
      maxPerformGas,
      fallbackGasPrice,
      fallbackLinkPrice,
      fallbackNativePrice,
      transcoder,
      registrars: [registrar.address],
      upkeepPrivilegeManager: upkeepManager,
      chainModule: chainModuleBase.address,
      reorgProtectionEnabled: true,
      financeAdmin: await admin.getAddress(),
    }
    await registry.connect(owner).setConfigTypeSafe(
      keepers,
      keepers,
      1,
      onchainConfig,
      1,
      '0x',
      [linkToken.address],
      [
        {
          gasFeePPB: paymentPremiumPPB,
          flatFeeMilliCents,
          priceFeed: await registry.getLinkUSDFeedAddress(),
          fallbackPrice: 200,
          minSpend: minimumRegistrationAmount,
          decimals: 18,
        },
      ],
    )
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await registrar.typeAndVersion()
      assert.equal(typeAndVersion, 'AutomationRegistrar 2.3.0')
    })
  })

  describe('#onTokenTransfer', () => {
    it('reverts if not called by the LINK token', async () => {
      await evmRevertCustomError(
        registrar
          .connect(someAddress)
          .onTokenTransfer(await someAddress.getAddress(), 0, '0x'),
        registrar,
        'OnlyLink',
      )
    })

    it('reverts if the admin address is 0x0000...', async () => {
      const abiEncodedBytes = encodeRegistrationParams({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: '0x0000000000000000000000000000000000000000',
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })

      await evmRevertCustomError(
        linkToken
          .connect(requestSender)
          .transferAndCall(registrar.address, amount, abiEncodedBytes),
        registrar,
        'InvalidAdminAddress',
      )
    })

    it('Auto Approve ON - registers an upkeep on KeeperRegistry instantly and emits both RegistrationRequested and RegistrationApproved events', async () => {
      //set auto approve ON with high threshold limits
      await registrar
        .connect(registrarOwner)
        .setTriggerConfig(
          Trigger.CONDITION,
          autoApproveType_ENABLED_ALL,
          maxAllowedAutoApprove,
        )

      //register with auto approve ON
      const abiEncodedBytes = encodeRegistrationParams({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })
      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)

      const [id] = await registry.getActiveUpkeepIDs(0, 1)

      //confirm if a new upkeep has been registered and the details are the same as the one just registered
      const newupkeep = await registry.getUpkeep(id)
      assert.equal(newupkeep.target, mock.address)
      assert.equal(newupkeep.admin, await admin.getAddress())
      assert.equal(newupkeep.checkData, emptyBytes)
      assert.equal(newupkeep.balance.toString(), amount.toString())
      assert.equal(newupkeep.performGas, performGas.toNumber())
      assert.equal(newupkeep.offchainConfig, offchainConfig)

      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).to.emit(registrar, 'RegistrationApproved')
    })

    it('Auto Approve ON - ignores the amount passed in and uses the actual amount sent', async () => {
      await registrar
        .connect(registrarOwner)
        .setTriggerConfig(
          Trigger.CONDITION,
          autoApproveType_ENABLED_ALL,
          maxAllowedAutoApprove,
        )

      const abiEncodedBytes = encodeRegistrationParams({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount: amount.mul(10), // muhahahaha 😈
        billingToken: linkToken.address,
      })

      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)

      const [id] = await registry.getActiveUpkeepIDs(0, 1)
      expect(await registry.getBalance(id)).to.equal(amount)
    })

    it('Auto Approve OFF - does not registers an upkeep on KeeperRegistry, emits only RegistrationRequested event', async () => {
      //get upkeep count before attempting registration
      const beforeCount = (await registry.getState()).state.numUpkeeps

      //set auto approve OFF, threshold limits dont matter in this case
      await registrar
        .connect(registrarOwner)
        .setTriggerConfig(
          Trigger.CONDITION,
          autoApproveType_DISABLED,
          maxAllowedAutoApprove,
        )

      //register with auto approve OFF
      const abiEncodedBytes = encodeRegistrationParams({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })
      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      const receipt = await tx.wait()

      //get upkeep count after attempting registration
      const afterCount = (await registry.getState()).state.numUpkeeps
      //confirm that a new upkeep has NOT been registered and upkeep count is still the same
      assert.deepEqual(beforeCount, afterCount)

      //confirm that only RegistrationRequested event is emitted and RegistrationApproved event is not
      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).not.to.emit(registrar, 'RegistrationApproved')

      const hash = receipt.logs[2].topics[1]
      const pendingRequest = await registrar.getPendingRequest(hash)
      assert.equal(await admin.getAddress(), pendingRequest[0])
      assert.ok(amount.eq(pendingRequest[1]))
    })

    it('Auto Approve ON - Throttle max approvals - does not register an upkeep on KeeperRegistry beyond the max limit, emits only RegistrationRequested event after limit is hit', async () => {
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 0)

      //set auto approve on, with max 1 allowed
      await registrar
        .connect(registrarOwner)
        .setTriggerConfig(Trigger.CONDITION, autoApproveType_ENABLED_ALL, 1)

      //set auto approve on, with max 1 allowed
      await registrar
        .connect(registrarOwner)
        .setTriggerConfig(Trigger.LOG, autoApproveType_ENABLED_ALL, 1)

      // register within threshold, new upkeep should be registered
      let abiEncodedBytes = encodeRegistrationParams({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 1) // 0 -> 1

      // try registering another one, new upkeep should not be registered
      abiEncodedBytes = encodeRegistrationParams({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas.toNumber() + 1, // make unique hash
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 1) // Still 1

      // register a second type of upkeep, different limit
      abiEncodedBytes = encodeRegistrationParams({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas, // make unique hash
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.LOG,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 2) // 1 -> 2

      // Now set new max limit to 2. One more upkeep should get auto approved
      await registrar
        .connect(registrarOwner)
        .setTriggerConfig(Trigger.CONDITION, autoApproveType_ENABLED_ALL, 2)

      abiEncodedBytes = encodeRegistrationParams({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas.toNumber() + 2, // make unique hash
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 3) // 2 -> 3

      // One more upkeep should not get registered
      abiEncodedBytes = encodeRegistrationParams({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas.toNumber() + 3, // make unique hash
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 3) // Still 3
    })

    it('Auto Approve Sender Allowlist - sender in allowlist - registers an upkeep on KeeperRegistry instantly and emits both RegistrationRequested and RegistrationApproved events', async () => {
      const senderAddress = await requestSender.getAddress()

      //set auto approve to ENABLED_SENDER_ALLOWLIST type with high threshold limits
      await registrar
        .connect(registrarOwner)
        .setTriggerConfig(
          Trigger.CONDITION,
          autoApproveType_ENABLED_SENDER_ALLOWLIST,
          maxAllowedAutoApprove,
        )

      // Add sender to allowlist
      await registrar
        .connect(registrarOwner)
        .setAutoApproveAllowedSender(senderAddress, true)

      //register with auto approve ON
      const abiEncodedBytes = encodeRegistrationParams({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })
      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)

      const [id] = await registry.getActiveUpkeepIDs(0, 1)

      //confirm if a new upkeep has been registered and the details are the same as the one just registered
      const newupkeep = await registry.getUpkeep(id)
      assert.equal(newupkeep.target, mock.address)
      assert.equal(newupkeep.admin, await admin.getAddress())
      assert.equal(newupkeep.checkData, emptyBytes)
      assert.equal(newupkeep.balance.toString(), amount.toString())
      assert.equal(newupkeep.performGas, performGas.toNumber())

      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).to.emit(registrar, 'RegistrationApproved')
    })

    it('Auto Approve Sender Allowlist - sender NOT in allowlist - does not registers an upkeep on KeeperRegistry, emits only RegistrationRequested event', async () => {
      const beforeCount = (await registry.getState()).state.numUpkeeps
      const senderAddress = await requestSender.getAddress()

      //set auto approve to ENABLED_SENDER_ALLOWLIST type with high threshold limits
      await registrar
        .connect(registrarOwner)
        .setTriggerConfig(
          Trigger.CONDITION,
          autoApproveType_ENABLED_SENDER_ALLOWLIST,
          maxAllowedAutoApprove,
        )

      // Explicitly remove sender from allowlist
      await registrar
        .connect(registrarOwner)
        .setAutoApproveAllowedSender(senderAddress, false)

      //register. auto approve shouldn't happen
      const abiEncodedBytes = encodeRegistrationParams({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })
      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      const receipt = await tx.wait()

      //get upkeep count after attempting registration
      const afterCount = (await registry.getState()).state.numUpkeeps
      //confirm that a new upkeep has NOT been registered and upkeep count is still the same
      assert.deepEqual(beforeCount, afterCount)

      //confirm that only RegistrationRequested event is emitted and RegistrationApproved event is not
      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).not.to.emit(registrar, 'RegistrationApproved')

      const hash = receipt.logs[2].topics[1]
      const pendingRequest = await registrar.getPendingRequest(hash)
      assert.equal(await admin.getAddress(), pendingRequest[0])
      assert.ok(amount.eq(pendingRequest[1]))
    })
  })

  describe('#registerUpkeep', () => {
    it('reverts with empty message if amount sent is not available in LINK allowance', async () => {
      await evmRevert(
        registrar.connect(someAddress).registerUpkeep({
          name: upkeepName,
          upkeepContract: mock.address,
          gasLimit: performGas,
          adminAddress: await admin.getAddress(),
          triggerType: Trigger.CONDITION,
          checkData: emptyBytes,
          triggerConfig: trigger,
          offchainConfig: emptyBytes,
          amount,
          encryptedEmail: emptyBytes,
          billingToken: linkToken.address,
        }),
        '',
      )
    })

    it('reverts if the amount passed in data is less than configured minimum', async () => {
      const amt = minimumRegistrationAmount.sub(1)

      await linkToken.connect(requestSender).approve(registrar.address, amt)

      await registrar
        .connect(registrarOwner)
        .setTriggerConfig(
          Trigger.CONDITION,
          autoApproveType_ENABLED_ALL,
          maxAllowedAutoApprove,
        )

      await evmRevertCustomError(
        registrar.connect(requestSender).registerUpkeep({
          name: upkeepName,
          upkeepContract: mock.address,
          gasLimit: performGas,
          adminAddress: await admin.getAddress(),
          triggerType: Trigger.CONDITION,
          checkData: emptyBytes,
          triggerConfig: trigger,
          offchainConfig: emptyBytes,
          amount: amt,
          encryptedEmail: emptyBytes,
          billingToken: linkToken.address,
        }),
        registrar,
        'InsufficientPayment',
      )
    })

    it('reverts if the billing token is not supported', async () => {
      await linkToken
        .connect(requestSender)
        .approve(registrar.address, minimumRegistrationAmount)

      await registrar
        .connect(registrarOwner)
        .setTriggerConfig(
          Trigger.CONDITION,
          autoApproveType_ENABLED_ALL,
          maxAllowedAutoApprove,
        )

      await registry
        .connect(owner)
        .setConfigTypeSafe(keepers, keepers, 1, onchainConfig, 1, '0x', [], [])

      await evmRevertCustomError(
        registrar.connect(requestSender).registerUpkeep({
          name: upkeepName,
          upkeepContract: mock.address,
          gasLimit: performGas,
          adminAddress: await admin.getAddress(),
          triggerType: Trigger.CONDITION,
          checkData: emptyBytes,
          triggerConfig: trigger,
          offchainConfig: emptyBytes,
          amount: minimumRegistrationAmount,
          encryptedEmail: emptyBytes,
          billingToken: linkToken.address,
        }),
        registrar,
        'InvalidBillingToken',
      )
    })

    it('Auto Approve ON - registers an upkeep on KeeperRegistry instantly and emits both RegistrationRequested and RegistrationApproved events', async () => {
      //set auto approve ON with high threshold limits
      await registrar
        .connect(registrarOwner)
        .setTriggerConfig(
          Trigger.CONDITION,
          autoApproveType_ENABLED_ALL,
          maxAllowedAutoApprove,
        )

      await linkToken.connect(requestSender).approve(registrar.address, amount)

      const tx = await registrar.connect(requestSender).registerUpkeep({
        name: upkeepName,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        triggerType: Trigger.CONDITION,
        checkData: emptyBytes,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        encryptedEmail: emptyBytes,
        billingToken: linkToken.address,
      })
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 1) // 0 -> 1

      //confirm if a new upkeep has been registered and the details are the same as the one just registered
      const [id] = await registry.getActiveUpkeepIDs(0, 1)
      const newupkeep = await registry.getUpkeep(id)
      assert.equal(newupkeep.target, mock.address)
      assert.equal(newupkeep.admin, await admin.getAddress())
      assert.equal(newupkeep.checkData, emptyBytes)
      assert.equal(newupkeep.balance.toString(), amount.toString())
      assert.equal(newupkeep.performGas, performGas.toNumber())
      assert.equal(newupkeep.offchainConfig, offchainConfig)

      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).to.emit(registrar, 'RegistrationApproved')
    })
  })

  describe('#setAutoApproveAllowedSender', () => {
    it('reverts if not called by the owner', async () => {
      const tx = registrar
        .connect(stranger)
        .setAutoApproveAllowedSender(await admin.getAddress(), false)
      await evmRevert(tx, 'Only callable by owner')
    })

    it('sets the allowed status correctly and emits log', async () => {
      const senderAddress = await stranger.getAddress()
      let tx = await registrar
        .connect(registrarOwner)
        .setAutoApproveAllowedSender(senderAddress, true)
      await expect(tx)
        .to.emit(registrar, 'AutoApproveAllowedSenderSet')
        .withArgs(senderAddress, true)

      let senderAllowedStatus = await registrar
        .connect(owner)
        .getAutoApproveAllowedSender(senderAddress)
      assert.isTrue(senderAllowedStatus)

      tx = await registrar
        .connect(registrarOwner)
        .setAutoApproveAllowedSender(senderAddress, false)
      await expect(tx)
        .to.emit(registrar, 'AutoApproveAllowedSenderSet')
        .withArgs(senderAddress, false)

      senderAllowedStatus = await registrar
        .connect(owner)
        .getAutoApproveAllowedSender(senderAddress)
      assert.isFalse(senderAllowedStatus)
    })
  })

  describe('#setTriggerConfig', () => {
    it('reverts if not called by the owner', async () => {
      const tx = registrar
        .connect(stranger)
        .setTriggerConfig(Trigger.LOG, autoApproveType_ENABLED_ALL, 100)
      await evmRevert(tx, 'Only callable by owner')
    })

    it('changes the config', async () => {
      const tx = await registrar
        .connect(registrarOwner)
        .setTriggerConfig(Trigger.LOG, autoApproveType_ENABLED_ALL, 100)
      await registrar.getTriggerRegistrationDetails(Trigger.LOG)
      await expect(tx)
        .to.emit(registrar, 'TriggerConfigSet')
        .withArgs(Trigger.LOG, autoApproveType_ENABLED_ALL, 100)
    })
  })

  describe('#approve', () => {
    let params: RegistrationParams

    beforeEach(async () => {
      await registrar
        .connect(registrarOwner)
        .setTriggerConfig(
          Trigger.CONDITION,
          autoApproveType_DISABLED,
          maxAllowedAutoApprove,
        )

      params = {
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      }

      //register with auto approve OFF
      const abiEncodedBytes = encodeRegistrationParams(params)

      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      await tx.wait()
    })

    it('reverts if not called by the owner', async () => {
      const tx = registrar.connect(stranger).approve({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig: emptyBytes,
        amount,
        billingToken: linkToken.address,
      })
      await evmRevert(tx, 'Only callable by owner')
    })

    it('reverts if the hash does not exist', async () => {
      const tx = registrar.connect(registrarOwner).approve({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig: emptyBytes,
        amount,
        billingToken: linkToken.address,
      })
      await evmRevertCustomError(tx, registrar, errorMsgs.requestNotFound)
    })

    it('reverts if any member of the payload changes', async () => {
      const invalidFields: any[] = [
        {
          name: 'fake',
        },
        {
          encryptedEmail: '0xdeadbeef',
        },
        {
          upkeepContract: ethers.Wallet.createRandom().address,
        },
        {
          gasLimit: performGas.add(1),
        },
        {
          adminAddress: randomAddress(),
        },
        {
          checkData: '0xdeadbeef',
        },
        {
          triggerType: Trigger.LOG,
        },
        {
          triggerConfig: '0x1234',
        },
        {
          offchainConfig: '0xdeadbeef',
        },
        {
          amount: amount.add(1),
        },
        {
          billingToken: randomAddress(),
        },
      ]
      for (let i = 0; i < invalidFields.length; i++) {
        const field = invalidFields[i]
        const badParams = Object.assign({}, params, field) as RegistrationParams
        const tx = registrar.connect(registrarOwner).approve(badParams)
        await expect(
          tx,
          `expected ${JSON.stringify(field)} to cause failure, but succeeded`,
        ).to.be.revertedWithCustomError(registrar, errorMsgs.requestNotFound)
      }
    })

    it('approves an existing registration request', async () => {
      const tx = await registrar.connect(registrarOwner).approve({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })
      await expect(tx).to.emit(registrar, 'RegistrationApproved')
    })

    it('deletes the request afterwards / reverts if the request DNE', async () => {
      await registrar.connect(registrarOwner).approve({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })
      const tx = registrar.connect(registrarOwner).approve({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })
      await evmRevertCustomError(tx, registrar, errorMsgs.requestNotFound)
    })
  })

  describe('#cancel', () => {
    let hash: string

    beforeEach(async () => {
      await registrar
        .connect(registrarOwner)
        .setTriggerConfig(
          Trigger.CONDITION,
          autoApproveType_DISABLED,
          maxAllowedAutoApprove,
        )

      //register with auto approve OFF
      const abiEncodedBytes = encodeRegistrationParams({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig,
        amount,
        billingToken: linkToken.address,
      })
      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      const receipt = await tx.wait()
      hash = receipt.logs[2].topics[1]
    })

    it('reverts if not called by the admin / owner', async () => {
      const tx = registrar.connect(stranger).cancel(hash)
      await evmRevertCustomError(tx, registrar, errorMsgs.onlyAdmin)
    })

    it('reverts if the hash does not exist', async () => {
      const tx = registrar
        .connect(registrarOwner)
        .cancel(
          '0x000000000000000000000000322813fd9a801c5507c9de605d63cea4f2ce6c44',
        )
      await evmRevertCustomError(tx, registrar, errorMsgs.requestNotFound)
    })

    it('refunds the total request balance to the admin address if owner cancels', async () => {
      const before = await linkToken.balanceOf(await admin.getAddress())
      const tx = await registrar.connect(registrarOwner).cancel(hash)
      const after = await linkToken.balanceOf(await admin.getAddress())
      assert.isTrue(after.sub(before).eq(amount.mul(BigNumber.from(1))))
      await expect(tx).to.emit(registrar, 'RegistrationRejected')
    })

    it('refunds the total request balance to the admin address if admin cancels', async () => {
      const before = await linkToken.balanceOf(await admin.getAddress())
      const tx = await registrar.connect(admin).cancel(hash)
      const after = await linkToken.balanceOf(await admin.getAddress())
      assert.isTrue(after.sub(before).eq(amount.mul(BigNumber.from(1))))
      await expect(tx).to.emit(registrar, 'RegistrationRejected')
    })

    it('deletes the request hash', async () => {
      await registrar.connect(registrarOwner).cancel(hash)
      let tx = registrar.connect(registrarOwner).cancel(hash)
      await evmRevertCustomError(tx, registrar, errorMsgs.requestNotFound)
      tx = registrar.connect(registrarOwner).approve({
        name: upkeepName,
        encryptedEmail: emptyBytes,
        upkeepContract: mock.address,
        gasLimit: performGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        triggerType: Trigger.CONDITION,
        triggerConfig: trigger,
        offchainConfig: emptyBytes,
        amount,
        billingToken: linkToken.address,
      })
      await evmRevertCustomError(tx, registrar, errorMsgs.requestNotFound)
    })
  })
})
