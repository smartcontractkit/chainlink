import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { evmRevert } from '../../../test-helpers/matchers'
import { getUsers, Personas } from '../../../test-helpers/setup'
import { BigNumber, Signer } from 'ethers'
import { LinkToken__factory as LinkTokenFactory } from '../../../../typechain/factories/LinkToken__factory'

import { MockV3Aggregator__factory as MockV3AggregatorFactory } from '../../../../typechain/factories/MockV3Aggregator__factory'
import { UpkeepMock__factory as UpkeepMockFactory } from '../../../../typechain/factories/UpkeepMock__factory'
import { KeeperRegistry20 as KeeperRegistry } from '../../../../typechain/KeeperRegistry20'
import { KeeperRegistryLogic20 as KeeperRegistryLogic } from '../../../../typechain/KeeperRegistryLogic20'
import { KeeperRegistrar20 as KeeperRegistrar } from '../../../../typechain/KeeperRegistrar20'
import { KeeperRegistry20__factory as KeeperRegistryFactory } from '../../../../typechain/factories/KeeperRegistry20__factory'
import { KeeperRegistryLogic20__factory as KeeperRegistryLogicFactory } from '../../../../typechain/factories/KeeperRegistryLogic20__factory'
import { KeeperRegistrar20__factory as KeeperRegistrarFactory } from '../../../../typechain/factories/KeeperRegistrar20__factory'

import { MockV3Aggregator } from '../../../../typechain/MockV3Aggregator'
import { LinkToken } from '../../../../typechain/LinkToken'
import { UpkeepMock } from '../../../../typechain/UpkeepMock'
import { toWei } from '../../../test-helpers/helpers'

let linkTokenFactory: LinkTokenFactory
let mockV3AggregatorFactory: MockV3AggregatorFactory
let keeperRegistryFactory: KeeperRegistryFactory
let keeperRegistryLogicFactory: KeeperRegistryLogicFactory
let keeperRegistrar: KeeperRegistrarFactory
let upkeepMockFactory: UpkeepMockFactory

let personas: Personas

before(async () => {
  personas = (await getUsers()).personas

  linkTokenFactory = await ethers.getContractFactory('LinkToken')
  mockV3AggregatorFactory = (await ethers.getContractFactory(
    'src/v0.8/tests/MockV3Aggregator.sol:MockV3Aggregator',
  )) as unknown as MockV3AggregatorFactory
  // @ts-ignore bug in autogen file
  keeperRegistryFactory = await ethers.getContractFactory('KeeperRegistry2_0')
  // @ts-ignore bug in autogen file
  keeperRegistryLogicFactory = await ethers.getContractFactory(
    'KeeperRegistryLogic2_0',
  )
  // @ts-ignore bug in autogen file
  keeperRegistrar = await ethers.getContractFactory('KeeperRegistrar2_0')
  upkeepMockFactory = await ethers.getContractFactory('UpkeepMock')
})

const errorMsgs = {
  onlyOwner: 'revert Only callable by owner',
  onlyAdmin: 'OnlyAdminOrOwner()',
  hashPayload: 'HashMismatch()',
  requestNotFound: 'RequestNotFound()',
}

describe('KeeperRegistrar2_0', () => {
  const upkeepName = 'SampleUpkeep'

  const linkEth = BigNumber.from(300000000)
  const gasWei = BigNumber.from(100)
  const executeGas = BigNumber.from(100000)
  const paymentPremiumPPB = BigNumber.from(250000000)
  const flatFeeMicroLink = BigNumber.from(0)
  const maxAllowedAutoApprove = 5

  const blockCountPerTurn = BigNumber.from(3)
  const registryGasOverhead = BigNumber.from('80000')
  const emptyBytes = '0x00'
  const stalenessSeconds = BigNumber.from(43820)
  const gasCeilingMultiplier = BigNumber.from(1)
  const checkGasLimit = BigNumber.from(20000000)
  const fallbackGasPrice = BigNumber.from(200)
  const fallbackLinkPrice = BigNumber.from(200000000)
  const maxPerformGas = BigNumber.from(5000000)
  const minUpkeepSpend = BigNumber.from('1000000000000000000')
  const amount = BigNumber.from('5000000000000000000')
  const amount1 = BigNumber.from('6000000000000000000')
  const transcoder = ethers.constants.AddressZero

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

  let linkToken: LinkToken
  let linkEthFeed: MockV3Aggregator
  let gasPriceFeed: MockV3Aggregator
  let registry: KeeperRegistry
  let registryLogic: KeeperRegistryLogic
  let mock: UpkeepMock
  let registrar: KeeperRegistrar

  beforeEach(async () => {
    owner = personas.Default
    admin = personas.Neil
    someAddress = personas.Ned
    registrarOwner = personas.Nelly
    stranger = personas.Nancy
    requestSender = personas.Norbert

    const config = {
      paymentPremiumPPB,
      flatFeeMicroLink,
      blockCountPerTurn,
      checkGasLimit,
      stalenessSeconds,
      gasCeilingMultiplier,
      minUpkeepSpend,
      maxPerformGas,
      fallbackGasPrice,
      fallbackLinkPrice,
      transcoder,
      registrar: ethers.constants.AddressZero,
    }

    linkToken = await linkTokenFactory.connect(owner).deploy()
    gasPriceFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(0, gasWei)
    linkEthFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(9, linkEth)
    registryLogic = await keeperRegistryLogicFactory
      .connect(owner)
      .deploy(
        0,
        registryGasOverhead,
        linkToken.address,
        linkEthFeed.address,
        gasPriceFeed.address,
      )

    registry = await keeperRegistryFactory
      .connect(owner)
      .deploy(
        0,
        registryGasOverhead,
        linkToken.address,
        linkEthFeed.address,
        gasPriceFeed.address,
        registryLogic.address,
        config,
      )

    mock = await upkeepMockFactory.deploy()

    registrar = await keeperRegistrar
      .connect(registrarOwner)
      .deploy(
        linkToken.address,
        autoApproveType_DISABLED,
        BigNumber.from('0'),
        registry.address,
        minUpkeepSpend,
      )

    await linkToken
      .connect(owner)
      .transfer(await requestSender.getAddress(), toWei('1000'))

    config.registrar = registrar.address
    await registry.setConfig(config)
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await registrar.typeAndVersion()
      assert.equal(typeAndVersion, 'KeeperRegistrar 2.0.0')
    })
  })

  describe('#register', () => {
    it('reverts if not called by the LINK token', async () => {
      await evmRevert(
        registrar
          .connect(someAddress)
          .register(
            upkeepName,
            emptyBytes,
            mock.address,
            executeGas,
            await admin.getAddress(),
            emptyBytes,
            amount,
            await requestSender.getAddress(),
          ),
        'OnlyLink()',
      )
    })

    it('reverts if the amount passed in data mismatches actual amount sent', async () => {
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_ENABLED_ALL,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount1,
          await requestSender.getAddress(),
        ],
      )

      await evmRevert(
        linkToken
          .connect(requestSender)
          .transferAndCall(registrar.address, amount, abiEncodedBytes),
        'AmountMismatch()',
      )
    })

    it('reverts if the sender passed in data mismatches actual sender', async () => {
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount,
          await admin.getAddress(), // Should have been requestSender.getAddress()
        ],
      )
      await evmRevert(
        linkToken
          .connect(requestSender)
          .transferAndCall(registrar.address, amount, abiEncodedBytes),
        'SenderMismatch()',
      )
    })

    it('reverts if the admin address is 0x0000...', async () => {
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          '0x0000000000000000000000000000000000000000',
          emptyBytes,
          amount,
          await requestSender.getAddress(),
        ],
      )

      await evmRevert(
        linkToken
          .connect(requestSender)
          .transferAndCall(registrar.address, amount, abiEncodedBytes),
        'RegistrationRequestFailed()',
      )
    })

    it('Auto Approve ON - registers an upkeep on KeeperRegistry instantly and emits both RegistrationRequested and RegistrationApproved events', async () => {
      //set auto approve ON with high threshold limits
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_ENABLED_ALL,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      //register with auto approve ON
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount,
          await requestSender.getAddress(),
        ],
      )
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
      assert.equal(newupkeep.executeGas, executeGas.toNumber())

      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).to.emit(registrar, 'RegistrationApproved')
    })

    it('Auto Approve OFF - does not registers an upkeep on KeeperRegistry, emits only RegistrationRequested event', async () => {
      //get upkeep count before attempting registration
      const beforeCount = (await registry.getState()).state.numUpkeeps

      //set auto approve OFF, threshold limits dont matter in this case
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_DISABLED,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      //register with auto approve OFF
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount,
          await requestSender.getAddress(),
        ],
      )
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
      await registrar.connect(registrarOwner).setRegistrationConfig(
        autoApproveType_ENABLED_ALL,
        1, // maxAllowedAutoApprove
        registry.address,
        minUpkeepSpend,
      )

      //register within threshold, new upkeep should be registered
      let abiEncodedBytes = registrar.interface.encodeFunctionData('register', [
        upkeepName,
        emptyBytes,
        mock.address,
        executeGas,
        await admin.getAddress(),
        emptyBytes,
        amount,
        await requestSender.getAddress(),
      ])
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 1) // 0 -> 1

      //try registering another one, new upkeep should not be registered
      abiEncodedBytes = registrar.interface.encodeFunctionData('register', [
        upkeepName,
        emptyBytes,
        mock.address,
        executeGas.toNumber() + 1, // make unique hash
        await admin.getAddress(),
        emptyBytes,
        amount,
        await requestSender.getAddress(),
      ])
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 1) // Still 1

      // Now set new max limit to 2. One more upkeep should get auto approved
      await registrar.connect(registrarOwner).setRegistrationConfig(
        autoApproveType_ENABLED_ALL,
        2, // maxAllowedAutoApprove
        registry.address,
        minUpkeepSpend,
      )
      abiEncodedBytes = registrar.interface.encodeFunctionData('register', [
        upkeepName,
        emptyBytes,
        mock.address,
        executeGas.toNumber() + 2, // make unique hash
        await admin.getAddress(),
        emptyBytes,
        amount,
        await requestSender.getAddress(),
      ])
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 2) // 1 -> 2

      // One more upkeep should not get registered
      abiEncodedBytes = registrar.interface.encodeFunctionData('register', [
        upkeepName,
        emptyBytes,
        mock.address,
        executeGas.toNumber() + 3, // make unique hash
        await admin.getAddress(),
        emptyBytes,
        amount,
        await requestSender.getAddress(),
      ])
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 2) // Still 2
    })

    it('Auto Approve Sender Allowlist - sender in allowlist - registers an upkeep on KeeperRegistry instantly and emits both RegistrationRequested and RegistrationApproved events', async () => {
      const senderAddress = await requestSender.getAddress()

      //set auto approve to ENABLED_SENDER_ALLOWLIST type with high threshold limits
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_ENABLED_SENDER_ALLOWLIST,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      // Add sender to allowlist
      await registrar
        .connect(registrarOwner)
        .setAutoApproveAllowedSender(senderAddress, true)

      //register with auto approve ON
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount,
          await requestSender.getAddress(),
        ],
      )
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
      assert.equal(newupkeep.executeGas, executeGas.toNumber())

      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).to.emit(registrar, 'RegistrationApproved')
    })

    it('Auto Approve Sender Allowlist - sender NOT in allowlist - does not registers an upkeep on KeeperRegistry, emits only RegistrationRequested event', async () => {
      const beforeCount = (await registry.getState()).state.numUpkeeps
      const senderAddress = await requestSender.getAddress()

      //set auto approve to ENABLED_SENDER_ALLOWLIST type with high threshold limits
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_ENABLED_SENDER_ALLOWLIST,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      // Explicitly remove sender from allowlist
      await registrar
        .connect(registrarOwner)
        .setAutoApproveAllowedSender(senderAddress, false)

      //register. auto approve shouldn't happen
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount,
          await requestSender.getAddress(),
        ],
      )
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
        registrar
          .connect(someAddress)
          .registerUpkeep(
            upkeepName,
            emptyBytes,
            mock.address,
            executeGas,
            await admin.getAddress(),
            emptyBytes,
            amount,
          ),
        '',
      )
    })

    it('reverts if the amount passed in data is less than configured minimum', async () => {
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_ENABLED_ALL,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      // amt is one order of magnitude less than minUpkeepSpend
      const amt = BigNumber.from('100000000000000000')

      await evmRevert(
        registrar
          .connect(someAddress)
          .registerUpkeep(
            upkeepName,
            emptyBytes,
            mock.address,
            executeGas,
            await admin.getAddress(),
            emptyBytes,
            amt,
          ),
        'InsufficientPayment()',
      )
    })

    it('Auto Approve ON - registers an upkeep on KeeperRegistry instantly and emits both RegistrationRequested and RegistrationApproved events', async () => {
      //set auto approve ON with high threshold limits
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_ENABLED_ALL,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      await linkToken.connect(requestSender).approve(registrar.address, amount)

      const tx = await registrar
        .connect(requestSender)
        .registerUpkeep(
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount,
        )
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 1) // 0 -> 1

      //confirm if a new upkeep has been registered and the details are the same as the one just registered
      const [id] = await registry.getActiveUpkeepIDs(0, 1)
      const newupkeep = await registry.getUpkeep(id)
      assert.equal(newupkeep.target, mock.address)
      assert.equal(newupkeep.admin, await admin.getAddress())
      assert.equal(newupkeep.checkData, emptyBytes)
      assert.equal(newupkeep.balance.toString(), amount.toString())
      assert.equal(newupkeep.executeGas, executeGas.toNumber())

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

  describe('#approve', () => {
    let hash: string

    beforeEach(async () => {
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_DISABLED,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      //register with auto approve OFF
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount,
          await requestSender.getAddress(),
        ],
      )

      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      const receipt = await tx.wait()
      hash = receipt.logs[2].topics[1]
    })

    it('reverts if not called by the owner', async () => {
      const tx = registrar
        .connect(stranger)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      await evmRevert(tx, 'Only callable by owner')
    })

    it('reverts if the hash does not exist', async () => {
      const tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          '0x000000000000000000000000322813fd9a801c5507c9de605d63cea4f2ce6c44',
        )
      await evmRevert(tx, errorMsgs.requestNotFound)
    })

    it('reverts if any member of the payload changes', async () => {
      let tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          ethers.Wallet.createRandom().address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.hashPayload)
      tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          10000,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.hashPayload)
      tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          ethers.Wallet.createRandom().address,
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.hashPayload)
      tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          '0x1234',
          hash,
        )
      await evmRevert(tx, errorMsgs.hashPayload)
    })

    it('approves an existing registration request', async () => {
      const tx = await registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      await expect(tx).to.emit(registrar, 'RegistrationApproved')
    })

    it('deletes the request afterwards / reverts if the request DNE', async () => {
      await registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      const tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.requestNotFound)
    })
  })

  describe('#cancel', () => {
    let hash: string

    beforeEach(async () => {
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_DISABLED,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      //register with auto approve OFF
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          amount,
          await requestSender.getAddress(),
        ],
      )
      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      const receipt = await tx.wait()
      hash = receipt.logs[2].topics[1]
      // submit duplicate request (increase balance)
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
    })

    it('reverts if not called by the admin / owner', async () => {
      const tx = registrar.connect(stranger).cancel(hash)
      await evmRevert(tx, errorMsgs.onlyAdmin)
    })

    it('reverts if the hash does not exist', async () => {
      const tx = registrar
        .connect(registrarOwner)
        .cancel(
          '0x000000000000000000000000322813fd9a801c5507c9de605d63cea4f2ce6c44',
        )
      await evmRevert(tx, errorMsgs.requestNotFound)
    })

    it('refunds the total request balance to the admin address', async () => {
      const before = await linkToken.balanceOf(await admin.getAddress())
      const tx = await registrar.connect(admin).cancel(hash)
      const after = await linkToken.balanceOf(await admin.getAddress())
      assert.isTrue(after.sub(before).eq(amount.mul(BigNumber.from(2))))
      await expect(tx).to.emit(registrar, 'RegistrationRejected')
    })

    it('deletes the request hash', async () => {
      await registrar.connect(registrarOwner).cancel(hash)
      let tx = registrar.connect(registrarOwner).cancel(hash)
      await evmRevert(tx, errorMsgs.requestNotFound)
      tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.requestNotFound)
    })
  })
})
